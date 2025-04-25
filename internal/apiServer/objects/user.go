package objects

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"storage/infra/dal"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register 处理用户注册请求
func Register(c *gin.Context) {
	resp := &BaseResp{}
	req := &User{}
	if err := c.Bind(req); err != nil {
		log.Println("err = ", err, " req = ", req)
		c.JSON(http.StatusBadRequest, "bind error")
		return
	}

	if req.Username == "" || req.Password == "" {
		resp.Set(1, "用户名和密码不能为空")
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err := dal.RegisterUser(req.Username, req.Password)
	if err != nil {
		resp.Set(1, fmt.Sprintf("该名称已被占用"))
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Set(0, "注册成功")
	c.JSON(http.StatusOK, resp)
}

// Login 处理用户登录请求
func Login(c *gin.Context) {
	resp := &BaseResp{}
	req := &User{}
	if err := c.Bind(req); err != nil {
		log.Println("err = ", err, " req = ", req)
		c.JSON(http.StatusBadRequest, "bind error")
		return
	}

	if req.Username == "" || req.Password == "" {
		resp.Set(1, "用户名和密码不能为空")
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	userID, err := dal.LoginUser(req.Username, req.Password)
	if err != nil {
		resp.Set(1, fmt.Sprintf("登录失败: %v", err))
		c.JSON(http.StatusOK, resp)
		return
	}

	if userID == 0 {
		resp.Set(1, "用户名或密码错误")
		c.JSON(http.StatusOK, resp)
		return
	}

	// 生成会话 ID
	sessionID := generateSessionID()
	expiresAt := time.Now().Add(24 * time.Hour)

	// 创建会话，删除旧会话
	err = dal.CreateSession(userID, sessionID, expiresAt)
	if err != nil {
		resp.Set(1, fmt.Sprintf("创建会话失败: %v", err))
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Set(0, "登录成功")
	resp.Data = map[string]interface{}{
		"session_id": sessionID,
	}
	c.JSON(http.StatusOK, resp)
}

// Logout 处理用户注销请求
func Logout(c *gin.Context) {
	resp := &BaseResp{}

	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		resp.Set(1, "未找到会话信息")
		c.JSON(http.StatusOK, resp)
		return
	}

	// 删除会话
	err := dal.DeleteSession(sessionID)
	if err != nil {
		resp.Set(1, fmt.Sprintf("退出失败: %v", err))
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Set(0, "退出成功")
	c.JSON(http.StatusOK, resp)
}

// generateSessionID 生成随机会话 ID
func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
