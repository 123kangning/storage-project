package dal

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"storage/conf"
)

var DB *sql.DB

func init() {
	var err error
	// 使用Open方法创建数据库连接
	DB, err = sql.Open("mysql", conf.MySQLDefaultDSN)
	if err != nil {
		log.Fatal(err)
	}

	// 尝试与数据库建立连接
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("sql ", DB, " open success...")
}
