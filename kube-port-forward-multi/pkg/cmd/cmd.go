package cmd

import (
	"fmt"
	"os"

	"github.com/sikalabs/go-scripts/kube-port-forward-multi/pkg/port_forward"
	"github.com/spf13/cobra"
)

var FlagContext string
var FlagKubeconfig string

var Cmd = &cobra.Command{
	Use:   "kube-port-forward-multi <local-port>:<svc/name|pod/name>:<remote-port> [...]",
	Short: "Run multiple kubectl port-forwards at once",
	Long: "Run multiple kubectl port-forwards at once.\n\n" +
		"Each target is <local-port>:<svc/name|pod/name>:<remote-port> for the current\n" +
		"namespace, or <local-port>:<namespace/svc/name|namespace/pod/name>:<remote-port>\n" +
		"to target a specific namespace.",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := port_forward.Run(port_forward.Options{
			Targets:    args,
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
