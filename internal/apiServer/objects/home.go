package objects

import (
	"github.com/gin-gonic/gin"
	"net/http"
	dal2 "storage/infra/dal"
)

type home struct {
	BaseResp BaseResp    `json:"baseResp"`
	Files    []dal2.File `json:"files"`
}

func Home(c *gin.Context) { //获取所有的Metadata集合
	resp := home{}
	files := dal2.GetAll()
	resp.Files = files
	resp.BaseResp.Set(0, "success")
	c.JSON(http.StatusOK, resp)
}
