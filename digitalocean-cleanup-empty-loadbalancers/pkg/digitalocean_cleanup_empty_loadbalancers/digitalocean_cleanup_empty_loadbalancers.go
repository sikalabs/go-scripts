package digitalocean_cleanup_empty_loadbalancers

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

	checkTeam(ctx, client)

	lbs, err := listLoadBalancers(ctx, client)
	if err != nil {
		log.Fatalf("Failed to list load balancers: %v", err)
	}

	found := 0
	for _, lb := range lbs {
		if !isEmpty(lb) {
			continue
		}
		found++

		region := ""
		if lb.Region != nil {
			region = lb.Region.Slug
		}

		if !delete {
			fmt.Printf("Would delete load balancer: %s (%s, %s, %s)\n", lb.Name, lb.ID, region, lb.IP)
			continue
		}

		fmt.Printf("Deleting load balancer: %s (%s, %s, %s)\n", lb.Name, lb.ID, region, lb.IP)
		if _, err := client.LoadBalancers.Delete(ctx, lb.ID); err != nil {
			log.Fatalf("Failed to delete load balancer %s: %v", lb.ID, err)
		}
	}

	if found == 0 {
		fmt.Println("No empty load balancers found.")
	} else if !delete {
		fmt.Printf("Found %d empty load balancer(s). Run with --delete to delete them.\n", found)
	}

	fmt.Println("Done!")
}
