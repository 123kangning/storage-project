package main

import (
	"github.com/gin-gonic/gin"
	"storage/internal/apiServer/middleware"
	"storage/internal/apiServer/objects"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")

	baseGroup := v1.Group("/file")
	{
		// 添加 SessionAuth 中间件
		baseGroup.Use(middleware.SessionAuth())
		baseGroup.GET("/search/", objects.Search)
		baseGroup.POST("/post/", objects.Post)
		baseGroup.GET("/get/", objects.Get)
		baseGroup.DELETE("/del/", objects.Del)
		baseGroup.GET("/", objects.Home)
	}

	userGroup := v1.Group("/user")
	{
		userGroup.POST("/login", objects.Login)
		userGroup.POST("/register", objects.Register)
		userGroup.POST("/logout", objects.Logout)
	}
	return r
}
