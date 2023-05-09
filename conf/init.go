package conf

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/spf13/viper"
)

var v *viper.Viper

func Init(configPath string) {
	v = viper.New()
	v.AddConfigPath(configPath)
	//v.AddConfigPath(path + "/../conf")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		logger.Fatal("config init: ", err)
	}
}
