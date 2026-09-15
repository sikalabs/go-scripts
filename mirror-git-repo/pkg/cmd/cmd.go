package cmd

import (
	"github.com/sikalabs/go-scripts/mirror-git-repo/pkg/mirror_git_repo"
	"github.com/spf13/cobra"
)

var FlagSourceUrl string
var FlagTargetUrl string
var FlagSourceToken string
var FlagTargetToken string

var Cmd = &cobra.Command{
	Use:   "mirror-git-repo",
	Short: "Mirror git repo",
	Run: func(cmd *cobra.Command, args []string) {
		mirror_git_repo.MirrorGitRepo(
			FlagSourceUrl,
			FlagTargetUrl,
			FlagSourceToken,
			FlagTargetToken,
		)
	},
}

func init() {
	Cmd.Flags().StringVar(
		&FlagSourceUrl,
		"source-url",
		"",
		"Source repo URL (e.g. https://gitlab.example.com/source-group/source-repo.git)",
	)
	Cmd.MarkFlagRequired("source-url")
	Cmd.Flags().StringVar(
		&FlagTargetUrl,
		"target-url",
		"",
		"Target repo URL (e.g. https://gitlab.example.com/target-group/target-repo.git)",
	)
	Cmd.MarkFlagRequired("target-url")
	Cmd.Flags().StringVar(
		&FlagSourceToken,
		"source-token",
		"",
		"Source repo access token (optional)",
	)
	Cmd.Flags().StringVar(
		&FlagTargetToken,
		"target-token",
		"",
		"Target repo access token (optional)",
	)
}
