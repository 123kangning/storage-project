package dal

import (
	"database/sql"
	"errors"
	"log"
	"storage/conf"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 打开数据库连接
func openDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", conf.MySQLDefaultDSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// RegisterUser 注册用户
func RegisterUser(username, password string) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO user (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		log.Println("RegisterUser err = ", err)
		return err
	}
	return nil
}

// GetUserBySession 根据会话 ID 获取用户信息
func GetUserBySession(sessionID string) (int, error) {
	db, err := openDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var userID int
	var expiresAt time.Time
	err = db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE id = ?", sessionID).Scan(&userID, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	if time.Now().After(expiresAt) {
		// 会话已过期，删除会话信息
		_ = DeleteSession(sessionID)
		return 0, nil
	}

	return userID, nil
}

// LoginUser 验证用户登录
func LoginUser(username, password string) (int, error) {
	db, err := openDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var userID int
	err = db.QueryRow("SELECT id FROM user WHERE username = ? AND password = ?", username, password).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return userID, nil
}

// CreateSession 创建新会话，删除旧会话
func CreateSession(userID int, sessionID string, expiresAt time.Time) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 删除该用户的旧会话
	_, err = tx.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 插入新会话
	_, err = tx.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, userID, expiresAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// DeleteSession 删除会话
func DeleteSession(sessionID string) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	return err
}
