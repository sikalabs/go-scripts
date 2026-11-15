package digitalocean_cleanup_empty_loadbalancers

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

// productionTeamUUIDs is the allowlist of production teams. Deleting load
// balancers in these teams requires an explicit interactive confirmation.
var productionTeamUUIDs = map[string]string{
	// "<team-uuid>": "SikaLabs PROD",
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
	fmt.Println("!!! WARNING: YOU ARE ABOUT TO DELETE LOAD BALANCERS IN PRODUCTION !!!")
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

func listLoadBalancers(ctx context.Context, client *godo.Client) ([]godo.LoadBalancer, error) {
	var all []godo.LoadBalancer
	opt := &godo.ListOptions{PerPage: 200}
	for {
		lbs, resp, err := client.LoadBalancers.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		all = append(all, lbs...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}
		opt.Page = page + 1
	}
	return all, nil
}

// isEmpty reports whether the load balancer has no backends: no droplets
// (regional LBs, incl. tag based ones) and no target LBs (global LBs).
func isEmpty(lb godo.LoadBalancer) bool {
	return len(lb.DropletIDs) == 0 && len(lb.TargetLoadBalancerIDs) == 0
}

func DigitalOceanCleanupEmptyLoadBalancers(token string, delete bool) {
	ctx := context.Background()
	client := godo.NewFromToken(token)

	production := checkTeam(ctx, client)

	lbs, err := listLoadBalancers(ctx, client)
	if err != nil {
		log.Fatalf("Failed to list load balancers: %v", err)
	}

	var empty []godo.LoadBalancer
	for _, lb := range lbs {
		if isEmpty(lb) {
			empty = append(empty, lb)
		}
	}

	if len(empty) == 0 {
		fmt.Println("No empty load balancers found.")
		fmt.Println("Done!")
		return
	}

	for _, lb := range empty {
		fmt.Printf("Empty load balancer: %s (%s, %s, %s)\n", lb.Name, lb.ID, regionSlug(lb), lb.IP)
	}
	fmt.Printf("Found %d empty load balancer(s).\n", len(empty))

	if !delete {
		fmt.Println("Run with --delete to delete them.")
		fmt.Println("Done!")
		return
	}

	if production {
		confirmProduction()
	}

	for _, lb := range empty {
		fmt.Printf("Deleting load balancer: %s (%s, %s, %s)\n", lb.Name, lb.ID, regionSlug(lb), lb.IP)
		if _, err := client.LoadBalancers.Delete(ctx, lb.ID); err != nil {
			log.Fatalf("Failed to delete load balancer %s: %v", lb.ID, err)
		}
	}

	fmt.Println("Done!")
}

func regionSlug(lb godo.LoadBalancer) string {
	if lb.Region == nil {
		return ""
	}
	return lb.Region.Slug
}
