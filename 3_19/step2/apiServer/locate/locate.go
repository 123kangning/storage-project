package locate

import (
	"os"
	"project/go-object-storage/src/lib/rabbitmq"
	"strconv"
	"time"
)

// Locate 通过rabbitmq向dataServers发送文件名，返回存储该文件的节点地址
func Locate(name string) string {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

func Exist(name string) bool {
	return Locate(name) != ""
}
