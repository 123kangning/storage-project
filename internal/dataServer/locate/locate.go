package locate

import (
	"os"
	"path/filepath"
	"storage/conf"
	"storage/internal/pkg/rabbitmq"
	"storage/internal/pkg/types"
	"strconv"
	"strings"
	"sync"
)

/*
存储hash-id键值对
对于每一个文件而言，都有唯一确定的hash值，以及多个分片id
*/
var objects = make(map[string]int)
var mutex sync.Mutex

func Locate(hash string) int {
	mutex.Lock()
	id, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return id
}

func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
}

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

// StartLocate 从消息队列中获取msg 查找资源是否存在
func StartLocate() {
	q := rabbitmq.New(conf.RabbitmqServer)
	defer q.Close()
	q.Bind("dataServer")
	c := q.Consume()
	for msg := range c {
		//提取出要查找的文件hash值
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		id := Locate(hash) //从hash-id的键值对中取出id
		if id != -1 {
			//根据hash找到文件存储地址(即该节点本身地址)以及hash对应的id
			q.Send(msg.ReplyTo, types.LocateMessage{Addr: os.Getenv("LISTEN_ADDRESS"), Id: id})
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		//file->[hash.分片id.dat]
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files[i])
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			panic(e)
		}
		objects[hash] = id
	}
}
