package middleware

import (
	"github.com/siddontang/go-log/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"storage/infra/dal"
)

func SessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.GetHeader("X-Session-ID")
		if sessionID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "message": "未提供会话 ID"})
			log.Error("未提供会话 ID")
			c.Abort()
			return
		}

		userID, err := dal.GetUserBySession(sessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "验证会话时出错"})
			log.Error("验证会话时出错")
			c.Abort()
			return
		}

		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "message": "无效的会话 ID"})
			log.Error("无效的会话 ID")
			c.Abort()
			return
		}

		// 将用户 ID 存储到上下文中，供后续处理函数使用
		c.Set("userID", userID)
		c.Next()
	}
}
