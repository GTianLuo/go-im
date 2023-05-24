package http

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	r := gin.New()
	r.POST("/user/register", UserRegister)
	r.POST("/user/login", UserLogin)
	return r
}
