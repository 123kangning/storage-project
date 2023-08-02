package main

import (
	"github.com/gin-gonic/gin"
	"storage/final/apiServer/objects"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	baseGroup := r.Group("/file")
	baseGroup.PUT("/put/", objects.Put)
	baseGroup.GET("/get/", objects.Get)
	baseGroup.DELETE("/del/", objects.Del)
	baseGroup.GET("/", objects.Home)
	return r
}
