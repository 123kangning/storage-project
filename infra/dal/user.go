package dal

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

// User 用户结构体，定义用户信息
type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
}

func (*User) TableName() string {
	return "users"
}

// Session 会话结构体
type Session struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"not null" json:"user_id"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
}

func (*Session) TableName() string {
	return "sessions"
}

// RegisterUser 注册用户
func RegisterUser(username, password string) error {
	user := User{
		Username: username,
		Password: password,
	}
	result := DB.Create(&user)
	if result.Error != nil {
		log.Println("RegisterUser err = ", result.Error)
		return result.Error
	}
	return nil
}

// GetUserBySession 根据会话 ID 获取用户信息
func GetUserBySession(sessionID string) (int64, error) {
	var session Session
	result := DB.Where("id = ?", sessionID).First(&session)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, result.Error
	}

	if time.Now().After(session.ExpiresAt) {
		// 会话已过期，删除会话信息
		_ = DeleteSession(sessionID)
		return 0, nil
	}

	return session.UserID, nil
}

// LoginUser 验证用户登录
func LoginUser(username, password string) (int64, error) {
	var user User
	result := DB.Where("username = ? AND password = ?", username, password).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, result.Error
	}
	return user.ID, nil
}

// CreateSession 创建新会话，删除旧会话
func CreateSession(userID int64, sessionID string, expiresAt time.Time) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		// 删除该用户的旧会话
		if err := tx.Where("user_id = ?", userID).Delete(&Session{}).Error; err != nil {
			return err
		}

		// 插入新会话
		session := Session{
			ID:        sessionID,
			UserID:    userID,
			ExpiresAt: expiresAt,
		}
		if err := tx.Create(&session).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteSession 删除会话
func DeleteSession(sessionID string) error {
	result := DB.Where("id = ?", sessionID).Delete(&Session{})
	return result.Error
}

// GetUserIDToUsernameMap 通过用户 ID 列表获取用户 ID 到用户名的映射
func GetUserIDToUsernameMap(userIDs []int64) (map[int64]string, error) {
	if len(userIDs) == 0 {
		return make(map[int64]string), nil
	}

	var users []User
	result := DB.Select("id", "username").Where("id IN ?", userIDs).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	userMap := make(map[int64]string)
	for _, user := range users {
		userMap[user.ID] = user.Username
	}
	return userMap, nil
}
