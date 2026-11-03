package digitalocean_cleanup_unused_volumes

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
)

// eligibleTeamUUIDs is the allowlist of teams this script is allowed to
// clean up, to avoid running it against a wrong account by accident.
var eligibleTeamUUIDs = map[string]string{
	"531319e3-3988-4578-a9fd-ec94e22df877": "SikaLabs DEV",
}

func checkTeam(ctx context.Context, client *godo.Client) {
	account, _, err := client.Account.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}
	if account.Team == nil {
		log.Fatalf("Token is not associated with any team")
	}
	if _, ok := eligibleTeamUUIDs[account.Team.UUID]; !ok {
		log.Fatalf("Team %s (%s) is not eligible for cleanup", account.Team.Name, account.Team.UUID)
	}
	fmt.Printf("Team: %s (%s)\n", account.Team.Name, account.Team.UUID)
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

	checkTeam(ctx, client)

	volumes, err := listVolumes(ctx, client)
	if err != nil {
		log.Fatalf("Failed to list volumes: %v", err)
	}

	found := 0
	var totalSize int64
	for _, v := range volumes {
		if !isUnused(v) {
			continue
		}
		found++
		totalSize += v.SizeGigaBytes

		region := ""
		if v.Region != nil {
			region = v.Region.Slug
		}

		if !delete {
			fmt.Printf("Would delete volume: %s (%s, %s, %d GB)\n", v.Name, v.ID, region, v.SizeGigaBytes)
			continue
		}

		fmt.Printf("Deleting volume: %s (%s, %s, %d GB)\n", v.Name, v.ID, region, v.SizeGigaBytes)
		if _, err := client.Storage.DeleteVolume(ctx, v.ID); err != nil {
			log.Fatalf("Failed to delete volume %s: %v", v.ID, err)
		}
	}

	if found == 0 {
		fmt.Println("No unused volumes found.")
	} else if !delete {
		fmt.Printf("Found %d unused volume(s) (%d GB total). Run with --delete to delete them.\n", found, totalSize)
	}

	fmt.Println("Done!")
}
