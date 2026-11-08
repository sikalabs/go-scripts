package cmd

import (
	"fmt"
	"os"

	"github.com/sikalabs/go-scripts/keepass-to-yaml/pkg/keepass_to_yaml"
	"github.com/spf13/cobra"
)

var FlagPassword string

var Cmd = &cobra.Command{
	Use:   "keepass-to-yaml <input.kdbx> <output.yaml>",
	Short: "Export a KeePass database to a YAML file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		password := FlagPassword
		if password == "" {
			password = os.Getenv("KEEPASS_MASTER_PASSWORD")
		}
		if password == "" {
			return fmt.Errorf("master password is required (--password or KEEPASS_MASTER_PASSWORD env var)")
		}
		return keepass_to_yaml.KeepassToYaml(args[0], args[1], password)
	},
}

func init() {
	Cmd.Flags().StringVarP(
		&FlagPassword,
		"password",
		"p",
		"",
		"KeePass master password (default from KEEPASS_MASTER_PASSWORD)",
	)
}
