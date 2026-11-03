package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sikalabs/go-scripts/digitalocean-cleanup-unused-volumes/pkg/digitalocean_cleanup_unused_volumes"
	"github.com/spf13/cobra"
)

var FlagToken string
var FlagDelete bool

var Cmd = &cobra.Command{
	Use:   "digitalocean-cleanup-unused-volumes",
	Short: "Delete DigitalOcean volumes not attached to any droplet",
	Run: func(cmd *cobra.Command, args []string) {
		token := FlagToken
		if token == "" {
			token = os.Getenv("DIGITALOCEAN_TOKEN")
		}
		if token == "" {
			fmt.Fprintf(os.Stderr, "Error: --token is required (flag, DIGITALOCEAN_TOKEN env var, or .env)\n")
			os.Exit(1)
		}

		digitalocean_cleanup_unused_volumes.DigitalOceanCleanupUnusedVolumes(
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
		"Actually delete unused volumes (default is dry run)",
	)

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Warning: failed to read .env: %v\n", err)
	}
}
