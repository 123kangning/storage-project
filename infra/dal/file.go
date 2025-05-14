package dal

import (
	"errors"
	"github.com/siddontang/go-log/log"
	"gorm.io/gorm"
	"time"
)

// File 文件结构体
type File struct {
	Name     string    `gorm:"not null" json:"name"`
	Size     int       `gorm:"not null" json:"size"`
	Hash     string    `gorm:"size:64;not null;unique" json:"hash"`
	IsDelete bool      `gorm:"not null;default:false" json:"is_delete"`
	UpdateAt time.Time `gorm:"not null" json:"update_at"`
	Source   int64     `gorm:"not null" json:"source"`
}

func (*File) TableName() string {
	return "files"
}

// Add 添加文件记录
func Add(name, hash string, size int, source int64) error {
	file := File{
		Name:     name,
		Size:     size,
		Hash:     hash,
		IsDelete: false,
		UpdateAt: time.Now(),
		Source:   source,
	}
	result := DB.Create(&file)
	return result.Error
}

// Delete 删除文件记录，需验证用户 ID
func Delete(hash string, userID int64) error {
	query := DB.Model(&File{}).Where("hash =?", hash)
	// 提供给外界用户的删除接口，进行鉴权。若用户 ID 为 0 则不进行鉴权
	if userID != 0 {
		query = query.Where("source =?", userID)
	}
	result := query.Update("is_delete", true)
	if result.RowsAffected == 0 {
		return nil
	}
	return result.Error
}

// Get 精确获取单个文件
func Get(hash string) File {
	var file File
	result := DB.Where("hash = ? AND is_delete = ?", hash, false).First(&file)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Error("Get file failed, err = ", result.Error)
	}
	return file
}
