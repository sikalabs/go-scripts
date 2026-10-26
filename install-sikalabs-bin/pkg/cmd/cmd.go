package cmd

import (
	"fmt"

	"github.com/sikalabs/go-scripts/install-sikalabs-bin/pkg/install_sikalabs_bin"
	"github.com/spf13/cobra"
)

var FlagDir string

var Cmd = &cobra.Command{
	Use:   "install-sikalabs-bin <name>",
	Short: "Install a binary from Sikalabs S3 (on DigitalOcean)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dst, err := install_sikalabs_bin.InstallBin(args[0], FlagDir)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Installed %s\n", dst)
	},
}

func init() {
	Cmd.Flags().StringVarP(
		&FlagDir,
		"dir",
		"d",
		"/usr/local/bin",
		"Directory to install the binary to",
	)
}
