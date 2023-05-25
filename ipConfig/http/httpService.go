package http

import (
	"github.com/gin-gonic/gin"
	"log"
)

func InitHttpService() {
	router := NewRouter()
	go func() {
		if err := router.Run(":9999"); err != nil {
			log.Fatal("failed to start http service:", err)
		}
	}()
}

func NewRouter() *gin.Engine {
	r := gin.New()
	r.POST("/user/register", UserRegister)
	r.POST("/user/login", UserLogin)
	return r
}
