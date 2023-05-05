package conf

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/spf13/viper"
	"os"
)

var v *viper.Viper

func Init() {
	v = viper.New()
	path, _ := os.Getwd()
	v.AddConfigPath(path + "/conf/")
	//v.AddConfigPath(path + "/../conf")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		logger.Fatal("config init: ", err)
	}
}
