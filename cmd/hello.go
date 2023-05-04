package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(helloCmd)
}

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "test command",
	Run:   helloHandle,
}

func helloHandle(cmd *cobra.Command, args []string) {
	fmt.Println("world!!!")
}
