package amis

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Index(c *gin.Context) {
	//从文件test.json读取json数据
	body, err := os.ReadFile(schema + "app.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "读取文件失败," + schema,
		})
	}
	c.HTML(200, "home.html", gin.H{
		"App": string(body),
	})
}

func Schema(c *gin.Context) {
	pageKey := c.Query("page_key") + ".json"
	//从文件test.json读取json数据
	body, err := os.ReadFile(schema + pageKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "读取文件失败," + schema + pageKey,
		})
	}
	c.JSON(http.StatusOK, json.RawMessage(body))
}
