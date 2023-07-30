package main

import (
	"github.com/gin-gonic/gin"
	"storage/final/apiServer/objects"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	baseGroup := r.Group("/file")
	baseGroup.PUT("/put/", objects.Put)
	return r
}
