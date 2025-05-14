package objects

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"storage/infra/dal"
)

// Del 删除对象的操作
func Del(c *gin.Context) {
	resp := &BaseResp{}
	// 获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		resp.Set(1, "未获取到用户 ID")
		c.JSON(http.StatusUnauthorized, resp)
		return
	}
	userIDInt64, ok := userID.(int64)
	if !ok {
		resp.Set(1, "用户 ID 类型错误")
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	hash := c.Query("hash")

	file := dal.Get(hash) //找出最近的版本
	if file.Hash == "" {
		log.Println("没有该文件")
		resp.Set(1, "没有该文件")
		c.JSON(http.StatusOK, resp)
		return
	}
	err := dal.Delete(hash, userIDInt64)
	if err != nil {
		log.Println(err)
		resp.Set(1, err.Error())
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Set(0, "success")
	c.JSON(http.StatusOK, resp)
	return
}
