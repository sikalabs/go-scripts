package cmd

import (
	"github.com/sikalabs/go-scripts/vault-secret-sync/pkg/vault_secret_sync"
	"github.com/spf13/cobra"
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
		vault_secret_sync.VaultSecretSync(
			FlagSourceAddr,
			FlagSourceToken,
			FlagTargetAddr,
			FlagTargetToken,
			FlagSyncPath,
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
	Cmd.MarkFlagRequired("source-addr")
	Cmd.Flags().StringVar(
		&FlagSourceToken,
		"source-token",
		"",
		"Source Vault token",
	)
	Cmd.MarkFlagRequired("source-token")
	Cmd.Flags().StringVar(
		&FlagTargetAddr,
		"target-addr",
		"",
		"Target Vault address (e.g. https://vault-target.example.com)",
	)
	Cmd.MarkFlagRequired("target-addr")
	Cmd.Flags().StringVar(
		&FlagTargetToken,
		"target-token",
		"",
		"Target Vault token",
	)
	Cmd.MarkFlagRequired("target-token")
	Cmd.Flags().StringVar(
		&FlagSyncPath,
		"sync-path",
		"",
		"Path to sync, including mount (e.g. secret/my-app)",
	)
	Cmd.MarkFlagRequired("sync-path")
}
