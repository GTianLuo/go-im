package ipConfig

import (
	"go-im/common/conf"
	"go-im/common/conf/dbConf"
	"go-im/ipConfig/http"
	"go-im/ipConfig/rpc/service"
	"go-im/ipConfig/serviceManage"
	"go-im/ipConfig/source"
	"os"
)

func RunMain() {
	//读取配置文件
	path, _ := os.Getwd()
	conf.Init(path + "/common/conf/")
	//初始化数据库服务
	dbConf.InitDbService()
	//初始化服务发现
	source.Init()
	//初始化服务调度
	serviceManage.Init()
	//初始化http服务
	http.InitHttpService()
	//初始化rpc服务
	service.InitRpcService()
}
