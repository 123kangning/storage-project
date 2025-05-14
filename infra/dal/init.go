package dal

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"storage/conf"
)

var DB *gorm.DB

func init() {
	dsn := conf.MySQLDefaultDSN
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
