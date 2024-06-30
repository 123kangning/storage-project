package dal

import (
	"log"
)

func Del(hash string) {
	// 准备预编译语句
	stmt, err := DB.Prepare("update file set is_delete=1 where hash=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// 执行预编译语句插入数据
	result, err := stmt.Exec(hash)
	log.Println("result = ", result)
	if err != nil {
		log.Fatal(err)
	}
}
