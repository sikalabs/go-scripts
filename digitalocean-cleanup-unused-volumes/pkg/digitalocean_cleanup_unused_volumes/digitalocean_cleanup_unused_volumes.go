package digitalocean_cleanup_unused_volumes

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/digitalocean/godo"
)

// eligibleTeamUUIDs is the allowlist of teams this script is allowed to
// clean up, to avoid running it against a wrong account by accident.
var eligibleTeamUUIDs = map[string]string{
	"531319e3-3988-4578-a9fd-ec94e22df877": "SikaLabs DEV",
}

// productionTeamUUIDs is the allowlist of production teams. Deleting volumes
// in these teams requires an explicit interactive confirmation.
var productionTeamUUIDs = map[string]string{
	"928fd8dd-5d60-43b3-a253-f54aea41e7f0": "SikaLabs PROD",
}

const productionConfirmText = "THIS IS PRODUCTION"
const productionWait = 10 * time.Second

// checkTeam verifies the token's team is eligible and reports whether it is
// a production team.
func checkTeam(ctx context.Context, client *godo.Client) bool {
	account, _, err := client.Account.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}
	if account.Team == nil {
		log.Fatalf("Token is not associated with any team")
	}
	if _, ok := productionTeamUUIDs[account.Team.UUID]; ok {
		fmt.Printf("Team: %s (%s) [PRODUCTION]\n", account.Team.Name, account.Team.UUID)
		return true
	}
	if _, ok := eligibleTeamUUIDs[account.Team.UUID]; !ok {
		log.Fatalf("Team %s (%s) is not eligible for cleanup", account.Team.Name, account.Team.UUID)
	}
	fmt.Printf("Team: %s (%s)\n", account.Team.Name, account.Team.UUID)
	return false
}

func countdown(d time.Duration) {
	for i := int(d.Seconds()); i > 0; i-- {
		fmt.Printf("\r%d ", i)
		time.Sleep(time.Second)
	}
	fmt.Print("\r   \r")
}

// confirmProduction waits, asks the user to type the confirmation text and
// waits again. Exits the program if the confirmation does not match.
func confirmProduction() {
	fmt.Println()
	fmt.Println("!!! WARNING: YOU ARE ABOUT TO DELETE VOLUMES IN PRODUCTION !!!")
	fmt.Printf("Waiting %s before confirmation (Ctrl+C to abort)...\n", productionWait)
	countdown(productionWait)

	fmt.Printf("Type \"%s\" to continue: ", productionConfirmText)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read confirmation: %v", err)
	}
	if strings.TrimSpace(line) != productionConfirmText {
		log.Fatalf("Confirmation does not match, aborting")
	}

	fmt.Printf("Deleting in %s (Ctrl+C to abort)...\n", productionWait)
	countdown(productionWait)
}

func listVolumes(ctx context.Context, client *godo.Client) ([]godo.Volume, error) {
	var all []godo.Volume
	params := &godo.ListVolumeParams{
		ListOptions: &godo.ListOptions{PerPage: 200},
	}
	for {
		volumes, resp, err := client.Storage.ListVolumes(ctx, params)
		if err != nil {
			return nil, err
		}
		all = append(all, volumes...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}
		params.ListOptions.Page = page + 1
	}
	return all, nil
}

// isUnused reports whether the volume is not attached to any droplet.
func isUnused(v godo.Volume) bool {
	return len(v.DropletIDs) == 0
}

func DigitalOceanCleanupUnusedVolumes(token string, delete bool) {
	ctx := context.Background()
	client := godo.NewFromToken(token)

	production := checkTeam(ctx, client)

	volumes, err := listVolumes(ctx, client)
	if err != nil {
		log.Fatalf("Failed to list volumes: %v", err)
	}

	var unused []godo.Volume
	var totalSize int64
	for _, v := range volumes {
		if !isUnused(v) {
			continue
		}
		unused = append(unused, v)
		totalSize += v.SizeGigaBytes
	}

	if len(unused) == 0 {
		fmt.Println("No unused volumes found.")
		fmt.Println("Done!")
		return
	}

	for _, v := range unused {
		fmt.Printf("Unused volume: %s (%s, %s, %d GB)\n", v.Name, v.ID, regionSlug(v), v.SizeGigaBytes)
	}
	fmt.Printf("Found %d unused volume(s) (%d GB total).\n", len(unused), totalSize)

	if !delete {
		fmt.Println("Run with --delete to delete them.")
		fmt.Println("Done!")
		return
	}

	if production {
		confirmProduction()
	}

	for _, v := range unused {
		fmt.Printf("Deleting volume: %s (%s, %s, %d GB)\n", v.Name, v.ID, regionSlug(v), v.SizeGigaBytes)
		if _, err := client.Storage.DeleteVolume(ctx, v.ID); err != nil {
			log.Fatalf("Failed to delete volume %s: %v", v.ID, err)
		}
	}

	fmt.Println("Done!")
}

func regionSlug(v godo.Volume) string {
	if v.Region == nil {
		return ""
	}
	return v.Region.Slug
}
