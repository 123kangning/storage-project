package myes

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"storage/conf"
)

type BinlogSync struct {
	canal.DummyEventHandler
}

// OnRotate 获取 binlog 下个日志文件名字和位置
func (h *BinlogSync) OnRotate(header *replication.EventHeader, r *replication.RotateEvent) error {
	fmt.Printf("下一个日志为 %s 位置为 %d \n", string(r.NextLogName), r.Position)
	return nil
}

// OnTableChanged 在 OnDDL 之前执行
func (h *BinlogSync) OnTableChanged(header *replication.EventHeader, schema string, table string) error {
	result := fmt.Sprintf("修改了数据库%s中表%s的结构", schema, table)
	fmt.Println(result)
	return nil
}

// OnDDL query 事件中的一些信息，如执行的 sql 语句
func (h *BinlogSync) OnDDL(header *replication.EventHeader, nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	fmt.Println(string(queryEvent.Query))
	return nil
}

// OnXID 打印事件 Xid 的结束为止
func (h *BinlogSync) OnXID(header *replication.EventHeader, m mysql.Position) error {
	fmt.Println("XID", m.Pos)
	return nil
}

// OnRow 获取 event_type 为 write_rows, update_rows, delete_rows 的数据
func (h *BinlogSync) OnRow(ev *canal.RowsEvent) error {
	fmt.Println("原始数据：", ev.Rows)
	fmt.Printf("sql的操作行为：%s\t", ev.Action)

	for idx := range ev.Table.PKColumns {
		fmt.Printf("主键为：%s\n", ev.Table.Columns[ev.Table.PKColumns[idx]].Name)
	}
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
	cfg := &canal.Config{
		Addr:     conf.MysqlAddr,
		User:     "canal",
		Password: "canal",
		Charset:  "utf8mb4",
		ServerID: conf.ServerID,
		Dump: canal.DumpConfig{
			TableDB: "file",
		},
	}
	// 表名
	//cfg.Dump.Tables = []string{"file.file"}
	//cfg.Dump.ExecutionPath = ""

	c, err := canal.NewCanal(cfg)
	if err != nil {
		panic(err)
	}
	// Register a handler to handler Events
	c.SetEventHandler(&BinlogSync{})
	<-Start
	err = c.Run()
	if err != nil {
		panic(err)
	}
}
