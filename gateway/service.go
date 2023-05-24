package gateway

import (
	"go-im/common/conf"
	"go-im/common/log"
	"os"
)

func RunMain() {
	//读取配置文件
	path, _ := os.Getwd()
	conf.Init(path + "/conf/")
	if err := initManager(); err != nil {
		log.Fatal("failed run gateway:", err.Error())
	}
	select {}
}
