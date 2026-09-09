package cmd

import (
	"github.com/sikalabs/go-scripts/example-hello-world/pkg/example_hello_world"
	"github.com/spf13/cobra"
)

var FlagName string

var Cmd = &cobra.Command{
	Use:   "example-hello-world",
	Short: "example hello world",
	Run: func(cmd *cobra.Command, args []string) {
		example_hello_world.ExampleHelloWorld(FlagName)
	},
}

func init() {
	Cmd.Flags().StringVarP(
		&FlagName,
		"name",
		"n",
		"World",
		"Name to greet",
	)
}
