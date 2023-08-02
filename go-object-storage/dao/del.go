package dao

import (
	"log"
	"time"
)

func Del(name string) {
	// 准备预编译语句
	stmt, err := DB.Prepare("update file set is_delete=1,update_at=? where name=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// 执行预编译语句插入数据
	result, err := stmt.Exec(time.Now().Add(8*time.Hour), name)
	log.Println("result = ", result)
	if err != nil {
		log.Fatal(err)
	}
}
