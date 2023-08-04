package myes

import (
	"encoding/json"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"log"
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
	rowData := make(map[string]interface{})
	rowList := make([]interface{}, len(ev.Rows))

	fmt.Println("原始数据：", ev.Rows)
	fmt.Printf("sql的操作行为：%s\t", ev.Action)

	for idx, _ := range ev.Table.PKColumns {
		fmt.Printf("主键为：%s\n", ev.Table.Columns[ev.Table.PKColumns[idx]].Name)
	}
	for idxRow, _ := range ev.Rows {
		for columnIndex, currColumn := range ev.Table.Columns {
			// 字段名和对应的值
			row := fmt.Sprintf("%v:%v", currColumn.Name, ev.Rows[idxRow][columnIndex])
			fmt.Println(row)
			rowData[currColumn.Name] = ev.Rows[idxRow][columnIndex]
			rowList[idxRow] = rowData
		}
	}

	rowJson, err := json.Marshal(rowList)
	if err != nil {
		return fmt.Errorf("序列化错误：%s", err)
	}

	fmt.Printf("序列化为json格式：%s\n\n", string(rowJson))
	return nil
}

func Run() {
	cfg := canal.NewDefaultConfig()
	//cfg.Addr = "127.0.0.1:3306"
	cfg.User = "file"
	cfg.Password = "file"
	// 数据库名
	cfg.Dump.TableDB = "file"
	cfg.ServerID = 1
	// 表名
	cfg.Dump.Tables = []string{"file.file"}
	cfg.Dump.ExecutionPath = ""

	c, err := canal.NewCanal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Register a handler to handler Events
	c.SetEventHandler(&BinlogSync{})

	err = c.Run()
	if err != nil {
		log.Fatal(err)
	}
}
