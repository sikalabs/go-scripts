package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sikalabs/go-scripts/digitalocean-cleanup-empty-loadbalancers/pkg/digitalocean_cleanup_empty_loadbalancers"
	"github.com/spf13/cobra"
)

var FlagToken string
var FlagDelete bool

var Cmd = &cobra.Command{
	Use:   "digitalocean-cleanup-empty-loadbalancers",
	Short: "Delete DigitalOcean load balancers without any droplets attached",
	Run: func(cmd *cobra.Command, args []string) {
		token := FlagToken
		if token == "" {
			token = os.Getenv("DIGITALOCEAN_TOKEN")
		}
		if token == "" {
			fmt.Fprintf(os.Stderr, "Error: --token is required (flag, DIGITALOCEAN_TOKEN env var, or .env)\n")
			os.Exit(1)
		}

		digitalocean_cleanup_empty_loadbalancers.DigitalOceanCleanupEmptyLoadBalancers(
			token,
			FlagDelete,
		)
	},
}

func init() {
	Cmd.Flags().StringVar(
		&FlagToken,
		"token",
		"",
		"DigitalOcean API token (default from DIGITALOCEAN_TOKEN)",
	)
	Cmd.Flags().BoolVar(
		&FlagDelete,
		"delete",
		false,
		"Actually delete empty load balancers (default is dry run)",
	)

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Warning: failed to read .env: %v\n", err)
	}
}
