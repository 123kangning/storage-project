package locate

import (
	"encoding/json"
	"os"
	"storage/internal/pkg/rabbitmq"
	"storage/internal/pkg/rs"
	"storage/internal/pkg/types"
	"time"
)

// Locate
/*
  - @Description: 用来通过hash定位一个object对象的地址
    新建一个消息队列，群发广播去找这个对象，如果一秒后没有响应就关闭，返回没找到
  - @param name
  - @return locateInfo	对象所在的地址,Id与Addr的键值对
*/
func Locate(name string) (locateInfo map[int]string) {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	locateInfo = make(map[int]string)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return
		}
		var info types.LocateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.Id] = info.Addr
	}
	return
}

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DATA_SHARDS
}
