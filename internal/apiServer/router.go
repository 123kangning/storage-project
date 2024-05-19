package main

import (
	"github.com/gin-gonic/gin"
	"storage/internal/apiServer/objects"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	baseGroup := r.Group("/file")
	baseGroup.GET("/search/", objects.Search)
	baseGroup.PUT("/put/", objects.Put)
	baseGroup.GET("/get/", objects.Get)
	baseGroup.DELETE("/del/", objects.Del)
	baseGroup.GET("/", objects.Home)
	return r
}
