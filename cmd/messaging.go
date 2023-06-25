package cmd

import (
	"github.com/spf13/cobra"
	"go-im/messaging"
)

func init() {
	rootCmd.AddCommand(messagingCmd)
}

var messagingCmd = &cobra.Command{
	Use:   "messaging",
	Short: "start messaging service",
	Run:   messagingCmdHandle,
}

func messagingCmdHandle(cmd *cobra.Command, args []string) {
	messaging.RunMain()
}
