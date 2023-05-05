package ipConfig

import (
	"github.com/gin-gonic/gin"
	"go-im/conf"
	"go-im/ipConfig/domain"
	"go-im/ipConfig/source"
	"net/http"
)

func RunMain() {
	//读取配置文件
	conf.Init()
	//初始化服务发现
	source.Init()
	//初始化服务调度
	domain.Init()
	e := gin.Default()
	e.GET("/ip/list", func(context *gin.Context) {
		context.JSON(http.StatusOK, domain.DisPatch())
	})
	e.Run(":4567")
}
