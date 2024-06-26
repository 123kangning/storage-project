package heartbeat

import (
	"storage/conf"
	"storage/internal/pkg/rabbitmq"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time) //保存所有的数据缓存节点
var mutex sync.RWMutex                       //处理并发的锁

// ListenHeartbeat
/**
 * @Description: 监听心跳信号
*		通过创建一个消息队列来绑定apiServer exchange，通过开启一个go的channel来监听来自每一个节点的心跳信号
*		心跳信号是有内容的，包括发自于哪个节点的地址和发送的时间，要把发送心跳信号的时间更新到dataServers
*/
func ListenHeartbeat() {
	q := rabbitmq.New(conf.RabbitmqServer)
	defer q.Close()
	q.Bind("apiServer")
	c := q.Consume()
	// 移除过期节点
	go removeExpiredDataServer()
	// 监听数据服务节点心跳，将心跳信息写入全局缓存
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock() //记得加锁
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

/*
*
  - @Description: 移除过期的节点
    每五秒都把那些超过10秒没有收到心跳消息的节点给删掉
    这个函数因为是死循环所以不能放在主线程执行，一会单开一个goroutine来执行
    涉及到并发所以同样加锁
*/
func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

// GetDataServers
/**
 * @Description: 返回当前的数据节点
 * @return []string
 */
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}
