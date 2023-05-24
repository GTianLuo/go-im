package conf

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/spf13/viper"
)

var V *viper.Viper

func Init(configPath string) {
	V = viper.New()
	V.AddConfigPath(configPath)
	//V.AddConfigPath(path + "/../conf")
	V.SetConfigName("config")
	V.SetConfigType("yaml")
	if err := V.ReadInConfig(); err != nil {
		logger.Fatal("config init: ", err)
	}
}
