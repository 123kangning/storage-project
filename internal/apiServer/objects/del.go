package objects

import (
	"github.com/gin-gonic/gin"
	"net/http"
	dal2 "storage/infra/dal"
)

// Del 删除对象的操作
func Del(c *gin.Context) {
	hash := c.Query("hash")
	resp := &BaseResp{}
	file := dal2.Get(hash) //找出最近的版本
	if file.Hash == "" {
		resp.Set(1, "没有该文件")
		c.JSON(http.StatusOK, resp)
		return
	}
	dal2.Del(hash)
	resp.Set(0, "success")
	c.JSON(http.StatusOK, resp)
	return
}
