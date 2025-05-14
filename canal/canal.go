package main

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"log"
	"storage/conf"
	"storage/infra/myes"
)

type BinlogSync struct {
	canal.DummyEventHandler
}

// OnRow 获取 event_type 为 write_rows, update_rows, delete_rows 的数据
func (h *BinlogSync) OnRow(ev *canal.RowsEvent) error {
	log.Println("原始数据：", ev.Rows)
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
	for _, row := range ev.Rows {
		after := make(map[string]interface{})
		for columnIndex, currColumn := range ev.Table.Columns {
			after[currColumn.Name] = row[columnIndex]
		}
		fmt.Println(after["size"], after["id"], after["name"], after["hash"], after["source"])
		fmt.Println(after["size"].(int32), after["id"].(int64), after["name"].(string), after["hash"].(string), after["source"].(int64))
		myes.AddFile(int(after["size"].(int32)), after["id"].(int64), after["name"].(string), after["hash"].(string), after["source"].(int64))
	}
}

// 暂无实际用途
func deleteFunc(ev *canal.RowsEvent) {
	before := make(map[string]interface{})
	for columnIndex, currColumn := range ev.Table.Columns {
		before[currColumn.Name] = ev.Rows[0][columnIndex]
	}
	myes.DeleteFile(before["id"].(int64))
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
		myes.DeleteFile(before["id"].(int64))
	} else if before["is_delete"].(int8) == 1 && after["is_delete"].(int8) == 0 { //回收
		myes.AddFile(int(after["size"].(int32)), after["id"].(int64), after["name"].(string), after["hash"].(string), after["source"].(int64))
	} else {
		myes.UpdateFile(before["id"].(int64), after["source"].(int64), int(after["size"].(int32)), after["name"].(string), after["hash"].(string))
	}
}

func main() {
	cfg := &canal.Config{}
	cfg.Addr = conf.MysqlAddr
	cfg.User = "canal"
	cfg.Password = "canal"
	cfg.ServerID = conf.ServerID
	cfg.IncludeTableRegex = []string{"file\\.files"}

	c, err := canal.NewCanal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	h := &BinlogSync{}
	c.SetEventHandler(h)
	myes.Init()
	err = c.Run()
	if err != nil {
		log.Fatal(err)
	}
}
