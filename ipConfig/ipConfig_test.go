package ipConfig

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-im/common/conf"
	"go-im/common/discovery"
	"go-im/ipConfig/serviceManage"
	"go-im/ipConfig/source"
	"io"
	"net/http"
	"path/filepath"
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
	conf.Init("../conf")
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
		s := &discovery.EndpointInfo{
			IP:   "localhost",
			Port: "9999" + strconv.Itoa(i),
			Metadata: map[string]interface{}{
				"connect_num":   99999,
				"message_bytes": 993882,
			},
		}
		_ = discovery.NewServerRegister(
			context.Background(),
			[]string{"http://localhost:2379"},
			time.Second*5,
			10,
			fmt.Sprint("im/gatewayServer/"+"node-1"+strconv.Itoa(i)),
			s,
		)
	}

	time.Sleep(time.Second * 100000000)
}

func TestSelect(t *testing.T) {
	resp, err := http.Get("http://localhost:9999/ip/list")
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
	ipListBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	ipList := make([]string, 0)
	println(string(ipListBytes))
	if err = json.Unmarshal(ipListBytes, &ipList); err != nil {
		panic(err)
	}
	fmt.Println(ipList)
}

func TestPath(t *testing.T) {
	fmt.Println(filepath.Join("hello\\hello.c", ".."))
}
