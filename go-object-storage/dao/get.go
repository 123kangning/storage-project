package dao

import (
	"log"
)

func Get(name string) []string {
	// 准备预编译语句
	stmt, err := DB.Prepare("select hash from file where name = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// 执行预编译语句插入数据
	rows, err := stmt.Query(name)
	if err != nil {
		log.Fatal(err)
	}
	ans := make([]string, 0)
	for rows.Next() {
		var hash string
		err = rows.Scan(&hash)
		if err != nil {
			log.Println(err)
			continue
		}
		ans = append(ans, hash)
	}
	return ans
}
