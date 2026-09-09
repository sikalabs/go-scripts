package cmd

import (
	"github.com/sikalabs/go-scripts/mirror-container-registry/pkg/mirror_container_registry"
	"github.com/spf13/cobra"
)

var FlagSourceRegistry string
var FlagTargetRegistry string
var FlagSourceGroup string
var FlagTargetGroup string
var FlagSourceToken string
var FlagTargetToken string

var Cmd = &cobra.Command{
	Use:   "mirror-container-registry",
	Short: "Mirror container registry",
	Run: func(cmd *cobra.Command, args []string) {
		mirror_container_registry.MirrorContainerRegistry(
			FlagSourceRegistry,
			FlagTargetRegistry,
			FlagSourceGroup,
			FlagTargetGroup,
			FlagSourceToken,
			FlagTargetToken,
		)
	},
}

func init() {
	Cmd.Flags().StringVar(
		&FlagSourceRegistry,
		"source-registry",
		"",
		"Source registry (e.g. registry-source.example.com)",
	)
	Cmd.MarkFlagRequired("source-registry")
	Cmd.Flags().StringVar(
		&FlagTargetRegistry,
		"target-registry",
		"",
		"Target registry (e.g. registry-target.example.com)",
	)
	Cmd.MarkFlagRequired("target-registry")
	Cmd.Flags().StringVar(
		&FlagSourceGroup,
		"source-group",
		"",
		"Source group to filter (e.g. source-group)",
	)
	Cmd.MarkFlagRequired("source-group")
	Cmd.Flags().StringVar(
		&FlagTargetGroup,
		"target-group",
		"",
		"Target group, replaces source group (e.g. target-group)",
	)
	Cmd.MarkFlagRequired("target-group")
	Cmd.Flags().StringVar(
		&FlagSourceToken,
		"source-token",
		"",
		"Source registry deploy token (optional)",
	)
	Cmd.Flags().StringVar(
		&FlagTargetToken,
		"target-token",
		"",
		"Target registry deploy token (optional)",
	)
}
