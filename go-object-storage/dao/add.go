package dao

import (
	"errors"
	"log"
	"time"
)

func Add(name, hash string, size int64) error {
	hs := Get(name)
	if len(hs) > 0 {
		return errors.New("命名冲突，请先删除重名文件")
	}
	// 准备预编译语句
	stmt, err := DB.Prepare("INSERT INTO file(name,size,hash,is_delete,updated_at) VALUES (?,?,?,false,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// 执行预编译语句插入数据
	_, err = stmt.Exec(name, size, hash, time.Now()) // 将 value1, value2 替换为实际的数据
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
