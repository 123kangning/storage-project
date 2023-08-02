package dao

import (
	"log"
	"time"
)

// Get 精确获取单个文件
func Get(name string) File {
	// 准备预编译语句
	stmt, err := DB.Prepare("select name,size,hash,update_at from file where name = ? and not is_delete order by 'update_at' desc")
	if err != nil {
		log.Fatal("err1 = ", err)
	}
	defer stmt.Close()

	// 执行预编译语句插入数据
	rows, err := stmt.Query(name)
	if err != nil {
		log.Fatal("err2 = ", err)
	}
	var file File
	if rows.Next() {
		err = rows.Scan(&file.Name, &file.Size, &file.Hash, &file.UpdateAt)
		if err != nil {
			log.Println("err3 = ", err)
		}
	}
	return file
}

// GetAll 获取全部文件
func GetAll() []File {
	// 准备预编译语句
	stmt, err := DB.Prepare("select name,size,hash,update_at from file order by 'update_at' desc limit 10")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// 执行预编译语句插入数据
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	files := make([]File, 0)
	for rows.Next() {
		var file File
		var t time.Time
		err = rows.Scan(&file.Name, &file.Size, &file.Hash, &t)
		file.UpdateAt = t.Format("2006-01-02 15:04:05")
		if err != nil {
			log.Println(err)
		}
		files = append(files, file)
	}
	return files
}
