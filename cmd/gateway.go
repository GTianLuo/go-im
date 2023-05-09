package cmd

import (
	"github.com/spf13/cobra"
	"go-im/gateway"
)

func init() {
	rootCmd.AddCommand(gatewayCmd)
}

var gatewayCmd = &cobra.Command{
	Short: "start gateway service",
	Use:   "gateway",
	Run:   handleGateway,
}

func handleGateway(cmd *cobra.Command, args []string) {
	gateway.RunMain()
}
