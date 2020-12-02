package router

import (
	"github.com/gin-gonic/gin"

	"errortest/api"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode("debug")

	group := r.Group("/api")
	group.GET("GetUser/:id", api.GetUser)

	return r
}
