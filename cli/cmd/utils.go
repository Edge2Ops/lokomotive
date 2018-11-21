package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// doesKubeconfigExist checks if the kubeconfig provided by user exists
func doesKubeconfigExist(*cobra.Command, []string) error {
	var err error
	kubeconfig := viper.GetString("kubeconfig")
	if _, err = os.Stat(kubeconfig); os.IsNotExist(err) {
		return fmt.Errorf("Kubeconfig %q not found", kubeconfig)
	}
	return err
}

// any `lokoctl <subcommand>` that needs kubeconfig flag should use this utility
// function to add this flag.
func addKubeConfigFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(
		"kubeconfig",
		os.ExpandEnv("$HOME/.kube/config"),
		"Path to kubeconfig file (required)")
	viper.BindPFlag("kubeconfig", cmd.PersistentFlags().Lookup("kubeconfig"))
}
