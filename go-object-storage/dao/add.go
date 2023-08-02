package dao

import (
	"errors"
	"log"
	"time"
)

type File struct {
	Hash     string `json:"-"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	UpdateAt string `json:"updateAt"`
}

func Add(name, hash string, size int64) error {
	file := Get(name)
	if file.Hash != "" {
		return errors.New("命名冲突，请先删除重名文件")
	}
	// 准备预编译语句
	stmt, err := DB.Prepare("INSERT INTO file(name,size,hash,is_delete,update_at) VALUES (?,?,?,false,?)")
	if err != nil {
		log.Fatal("err1 = ", err)
	}
	defer stmt.Close()

	// 执行预编译语句插入数据
	_, err = stmt.Exec(name, size, hash, time.Now().Add(8*time.Hour)) // 将 value1, value2 替换为实际的数据
	if err != nil {
		log.Fatal("err2 = ", err)
	}
	return nil
}
