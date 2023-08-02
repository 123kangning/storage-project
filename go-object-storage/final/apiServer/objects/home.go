package objects

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"storage/dao"
)

type home struct {
	BaseResp BaseResp
	Files    []dao.File
}

func Home(c *gin.Context) { //获取所有的Metadata集合
	resp := home{}
	files := dao.GetAll()
	resp.Files = files
	resp.BaseResp.Set(0, "success")
	c.JSON(http.StatusOK, resp)
}
