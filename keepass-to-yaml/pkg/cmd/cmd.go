package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sikalabs/go-scripts/keepass-to-yaml/pkg/keepass_to_yaml"
	"github.com/spf13/cobra"
)

var FlagPassword string
var FlagPasswordFile string

var Cmd = &cobra.Command{
	Use:   "keepass-to-yaml <input.kdbx> <output.yaml>",
	Short: "Export a KeePass database to a YAML file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		password := FlagPassword
		if password == "" && FlagPasswordFile != "" {
			data, err := os.ReadFile(FlagPasswordFile)
			if err != nil {
				return fmt.Errorf("failed to read password file: %w", err)
			}
			password = strings.TrimRight(string(data), "\r\n")
		}
		if password == "" {
			password = os.Getenv("KEEPASS_MASTER_PASSWORD")
		}
		if password == "" {
			return fmt.Errorf("master password is required (--password, --password-file or KEEPASS_MASTER_PASSWORD env var)")
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
	Cmd.Flags().StringVar(
		&FlagPasswordFile,
		"password-file",
		"",
		"File containing the KeePass master password",
	)
	Cmd.MarkFlagsMutuallyExclusive("password", "password-file")
}
