package objects

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"storage/dao"
)

// Del 删除对象的操作
func Del(c *gin.Context) {
	name := c.Query("name")
	resp := &BaseResp{}
	file := dao.Get(name) //找出最近的版本
	if file.Hash == "" {
		resp.Set(1, "没有该文件")
		c.JSON(http.StatusOK, resp)
		return
	}
	dao.Del(name)
	resp.Set(0, "success")
	c.JSON(http.StatusOK, resp)
	return
}
