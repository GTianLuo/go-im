package cmd

import (
	"github.com/spf13/cobra"
	"go-im/ipConfig"
)

func init() {
	rootCmd.AddCommand(ipConfigCmd)
}

var ipConfigCmd = &cobra.Command{
	Use:   "ipConfig",
	Short: "start ipConfig service",
	Run:   ipConfigHandle,
}

func ipConfigHandle(cmd *cobra.Command, args []string) {
	ipConfig.RunMain()
}
