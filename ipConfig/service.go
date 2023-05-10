package ipConfig

import (
	"github.com/gin-gonic/gin"
	"go-im/conf"
	"go-im/ipConfig/serviceManage"
	"go-im/ipConfig/source"
	"net/http"
	"os"
)

func RunMain() {
	//读取配置文件
	path, _ := os.Getwd()
	conf.Init(path + "/conf/")
	//初始化服务发现
	source.Init()
	//初始化服务调度
	serviceManage.Init()
	e := gin.Default()
	e.GET("/ip/list", func(context *gin.Context) {
		context.JSON(http.StatusOK, serviceManage.DisPatch())
	})
	e.Run(":9999")
}
