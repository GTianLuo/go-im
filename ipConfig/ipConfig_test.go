package ipConfig

import (
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"go-im/common/discovery"
	"go-im/conf"
	"go-im/ipConfig/serviceManage"
	"go-im/ipConfig/source"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestData(t *testing.T) {
	source.Init()
	time.Sleep(100000 * time.Minute)
}

func TestWatch(t *testing.T) {
	//读取配置文件
	path, _ := os.Getwd()
	conf.Init(path + "/conf")
	//初始化服务发现
	source.Init()
	//初始化服务调度
	serviceManage.Init()
	e := gin.Default()
	e.GET("/ip/list", func(context *gin.Context) {
		context.JSON(http.StatusOK, serviceManage.DisPatch())
	})
	e.Run(":4567")
}

func TestServiceRegister(t *testing.T) {

	for i := 0; i < 10; i++ {
		s, err := (&discovery.EndpointInfo{
			IP:   "localhost",
			Port: "9999" + strconv.Itoa(i),
			Metadata: map[string]interface{}{
				"connect_num":   99999,
				"message_bytes": 993882,
			},
		}).Marshal()
		if err != nil {
			logger.Fatal(err)
		}
		_ = discovery.NewServerRegister(
			context.Background(),
			[]string{"http://localhost:2379"},
			time.Second*5,
			10,
			fmt.Sprint("im/ipConfig/"+"node-1"+strconv.Itoa(i)),
			s,
		)
	}

	time.Sleep(time.Second * 100000000)
}

func TestSelect(t *testing.T) {
	ch := make(chan int)
	select {
	case <-ch:
		fmt.Println("Read")
	case ch <- 1:
		fmt.Println("write")
	}
}