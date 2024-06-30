package dal

import (
	"log"
	"time"
)

type File struct {
	Hash     string `json:"-"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	UpdateAt string `json:"updateAt"`
}

type FileDO struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Size      int       `json:"size"`
	Hash      string    `json:"hash"`
	IsDelete  bool      `json:"is_delete"`
	UpdatedAt time.Time `json:"updated_at"`
}

func Add(name, hash string, size int64) error {
	stmt, err := DB.Prepare("select id,is_delete,hash from file where hash = ?")
	if err != nil {
		log.Fatal("Prepare failed,err = ", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(hash)
	if err != nil {
		log.Fatal("add check failed,err = ", err)
	}
	var file FileDO
	if rows.Next() {
		err = rows.Scan(&file.ID, &file.IsDelete, &file.Hash)
		if err != nil {
			log.Println("check file failed,err = ", err)
		}
	}
	if file.ID == 0 {
		//无记录，直接插入
		// 准备预编译语句
		stmt, err = DB.Prepare("INSERT INTO file(name,size,hash,is_delete,update_at) VALUES (?,?,?,false,?)")
		if err != nil {
			log.Fatal("insert prepare fail, err = ", err)
		}
		defer stmt.Close()

		// 执行预编译语句插入数据
		_, err = stmt.Exec(name, size, hash, time.Now().Add(8*time.Hour))
		if err != nil {
			log.Fatal("insert Exec fail, err = ", err)
		}
	} else {
		//有记录，但是已经删除，直接更新
		// 准备预编译语句
		stmt, err = DB.Prepare("update file set is_delete=false where id=?")
		if err != nil {
			log.Fatal("update prepare fail, err = ", err)
		}
		defer stmt.Close()

		// 执行预编译语句插入数据
		_, err = stmt.Exec(file.ID)
		if err != nil {
			log.Fatal("update Exec fail, err = ", err)
		}
	}

	return nil
}
