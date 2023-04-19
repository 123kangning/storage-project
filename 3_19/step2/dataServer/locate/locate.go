package locate

import (
	"os"
	"project/go-object-storage/src/lib/rabbitmq"
	"strconv"
)

func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// StartLocate 从消息队列中获取msg 查找资源是否存在
func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		//提取出要查找的文件名称
		object, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + object) { //资源存在
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}
