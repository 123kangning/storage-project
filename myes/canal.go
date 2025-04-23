package myes

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"log"
	"storage/conf"
	//"storage/conf"
)

type BinlogSync struct {
	canal.DummyEventHandler
}

// OnRow 获取 event_type 为 write_rows, update_rows, delete_rows 的数据
func (h *BinlogSync) OnRow(ev *canal.RowsEvent) error {
	log.Println("原始数据：", ev.Rows)
	if ev.Table.Name != "file" {
		log.Println("表名不匹配,不处理")
		return nil
	}
	log.Printf("sql的操作行为：%s\t", ev.Action)

	switch ev.Action {
	case canal.InsertAction:
		{
			insertFunc(ev)
		}
	case canal.DeleteAction:
		{
			deleteFunc(ev)
		}
	case canal.UpdateAction:
		{
			updateFunc(ev)
		}
	}

	return nil
}
func insertFunc(ev *canal.RowsEvent) {
	after := make(map[string]interface{})
	for columnIndex, currColumn := range ev.Table.Columns {
		after[currColumn.Name] = ev.Rows[0][columnIndex]
	}
	fmt.Printf("%+v\n", ev)
	fmt.Println("after = ", after)
	fmt.Println(after["size"].(int32), after["id"].(int64), after["name"].(string), after["hash"].(string))
	AddFile(after["size"].(int32), after["id"].(int64), after["name"].(string), after["hash"].(string))
}

// 暂无实际用途
func deleteFunc(ev *canal.RowsEvent) {
	before := make(map[string]interface{})
	for columnIndex, currColumn := range ev.Table.Columns {
		before[currColumn.Name] = ev.Rows[0][columnIndex]
	}
	fmt.Println("before = ", before)
	DeleteFile(before["id"].(int64))
}

// 实际用于假删除和从回收站回收，更改is_delete字段
func updateFunc(ev *canal.RowsEvent) {
	before := make(map[string]interface{})
	after := make(map[string]interface{})
	for columnIndex, currColumn := range ev.Table.Columns {
		before[currColumn.Name] = ev.Rows[0][columnIndex]
	}
	for columnIndex, currColumn := range ev.Table.Columns {
		after[currColumn.Name] = ev.Rows[1][columnIndex]
	}
	if before["is_delete"].(int8) == 0 && after["is_delete"].(int8) == 1 { //删除
		DeleteFile(before["id"].(int64))
	}
	if before["is_delete"].(int8) == 1 && after["is_delete"].(int8) == 0 { //回收
		AddFile(after["size"].(int32), after["id"].(int64), after["name"].(string), after["hash"].(string))
	}
}

func Run() {
	c, err := canal.NewCanal(&canal.Config{
		Addr:     conf.MysqlAddr,
		User:     "canal",
		Password: "canal",
		ServerID: conf.ServerID,
	})
	if err != nil {
		log.Fatal(err)
	}
	h := &BinlogSync{}
	c.SetEventHandler(h)
	go Init()
	<-Start
	err = c.Run()
	if err != nil {
		log.Fatal(err)
	}
}
