package cmd

import (
	"github.com/sikalabs/go-scripts/yaml-to-keepass/pkg/yaml_to_keepass"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "yaml-to-keepass <input.yaml> <output.kdbx>",
	Short: "Create a KeePass database from a YAML file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return yaml_to_keepass.YamlToKeepass(args[0], args[1])
	},
}
