package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sikalabs/go-scripts/vault-secret-sync/pkg/vault_secret_sync"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var FlagSourceAddr string
var FlagSourceToken string
var FlagTargetAddr string
var FlagTargetToken string
var FlagSyncPath string

var Cmd = &cobra.Command{
	Use:   "vault-secret-sync",
	Short: "Sync Vault secrets from one Vault to another",
	Run: func(cmd *cobra.Command, args []string) {
		values := map[string]string{
			"source-addr":  viper.GetString("source-addr"),
			"source-token": viper.GetString("source-token"),
			"target-addr":  viper.GetString("target-addr"),
			"target-token": viper.GetString("target-token"),
			"sync-path":    viper.GetString("sync-path"),
		}

		for _, name := range []string{"source-addr", "source-token", "target-addr", "target-token", "sync-path"} {
			if values[name] == "" {
				envName := "VAULT_SECRET_SYNC_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
				fmt.Fprintf(os.Stderr, "Error: --%s is required (flag, %s env var, or .env)\n", name, envName)
				os.Exit(1)
			}
		}

		vault_secret_sync.VaultSecretSync(
			values["source-addr"],
			values["source-token"],
			values["target-addr"],
			values["target-token"],
			values["sync-path"],
		)
	},
}

func init() {
	Cmd.Flags().StringVar(
		&FlagSourceAddr,
		"source-addr",
		"",
		"Source Vault address (e.g. https://vault-source.example.com)",
	)
	Cmd.Flags().StringVar(
		&FlagSourceToken,
		"source-token",
		"",
		"Source Vault token",
	)
	Cmd.Flags().StringVar(
		&FlagTargetAddr,
		"target-addr",
		"",
		"Target Vault address (e.g. https://vault-target.example.com)",
	)
	Cmd.Flags().StringVar(
		&FlagTargetToken,
		"target-token",
		"",
		"Target Vault token",
	)
	Cmd.Flags().StringVar(
		&FlagSyncPath,
		"sync-path",
		"",
		"Path to sync, including mount (e.g. secret/my-app)",
	)

	for _, name := range []string{"source-addr", "source-token", "target-addr", "target-token", "sync-path"} {
		viper.BindPFlag(name, Cmd.Flags().Lookup(name))
	}

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Warning: failed to read .env: %v\n", err)
	}

	viper.SetEnvPrefix("VAULT_SECRET_SYNC")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}
