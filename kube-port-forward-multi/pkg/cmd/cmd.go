package cmd

import (
	"fmt"
	"os"

	"github.com/sikalabs/go-scripts/kube-port-forward-multi/pkg/port_forward"
	"github.com/spf13/cobra"
)

var FlagNamespace string
var FlagContext string
var FlagKubeconfig string

var Cmd = &cobra.Command{
	Use:   "kube-port-forward-multi <local-port>:<svc/name|pod/name>:<remote-port> [...]",
	Short: "Run multiple kubectl port-forwards at once",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := port_forward.Run(port_forward.Options{
			Targets:    args,
			Namespace:  FlagNamespace,
			Context:    FlagContext,
			Kubeconfig: FlagKubeconfig,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	Cmd.Flags().StringVarP(
		&FlagNamespace,
		"namespace",
		"n",
		"",
		"Kubernetes namespace",
	)

	Cmd.Flags().StringVar(
		&FlagContext,
		"context",
		"",
		"Kubeconfig context to use",
	)

	Cmd.Flags().StringVar(
		&FlagKubeconfig,
		"kubeconfig",
		"",
		"Path to kubeconfig file",
	)
}
