package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ConfigPath string
)

func init() {
	cobra.OnInitialize(initConfig)
	path, _ := os.Getwd()
	ConfigPath = path + "/../conf/"
}

var rootCmd = &cobra.Command{
	Use:   "im",
	Short: "这是一个超牛逼的IM系统",
	Run:   IM,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func IM(cmd *cobra.Command, args []string) {

}

type I interface {
	//io.Closer
	Close() error
}

type X struct {
}

func (x X) Close() error { //TODO implement me
	panic("implement me")
}

func initConfig() {
}
