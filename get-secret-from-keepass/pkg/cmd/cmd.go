package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sikalabs/go-scripts/get-secret-from-keepass/pkg/get_secret_from_keepass"
	"github.com/spf13/cobra"
)

var FlagPassword string
var FlagPasswordFile string

var Cmd = &cobra.Command{
	Use:     "get-secret-from-keepass <database.kdbx> <path>",
	Short:   "Print password of a KeePass entry by its path",
	Example: "get-secret-from-keepass secrets.kdbx /Root/group1/bar_entry",
	Args:    cobra.ExactArgs(2),
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
		secret, err := get_secret_from_keepass.GetSecretFromKeepass(args[0], args[1], password)
		if err != nil {
			return err
		}
		fmt.Println(secret)
		return nil
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
