package dal

import (
	"log"
	"time"
)

// Get 精确获取单个文件
func Get(hash string) File {
	// 准备预编译语句
	stmt, err := DB.Prepare("select name,size,hash,update_at from file where hash = ? and is_delete=0 order by 'update_at' desc")
	if err != nil {
		log.Fatal("Prepare failed,err = ", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(hash)
	if err != nil {
		log.Fatal("Query failed,err = ", err)
	}
	var file File
	if rows.Next() {
		err = rows.Scan(&file.Name, &file.Size, &file.Hash, &file.UpdateAt)
		if err != nil {
			log.Println("get file failed,err = ", err)
		}
	}
	return file
}

// GetAll 获取全部文件
func GetAll() []File {
	// 准备预编译语句
	stmt, err := DB.Prepare("select name,size,hash,update_at from file where is_delete=0 order by 'update_at' desc limit 10")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	files := make([]File, 0)
	for rows.Next() {
		var file File
		var t time.Time
		err = rows.Scan(&file.Name, &file.Size, &file.Hash, &t)
		file.UpdateAt = t.Format(time.DateTime)
		if err != nil {
			log.Println(err)
		}
		files = append(files, file)
	}
	return files
}
