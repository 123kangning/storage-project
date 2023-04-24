package objects

import (
	"fmt"
	"log"
	"project/go-object-storage/final/apiServer/heartbeat"
	"project/go-object-storage/src/lib/rs"
)

// 调用dataServer生成文件，暂时还未写入，返回writer
func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	log.Println("api.objects.putStream")
	// 获取全部的数据服务节点，无需排除任何节点
	servers := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	// 如果数据服务节点长度不等于分片长度，则无法完整保存数据，提示报错
	if len(servers) != rs.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}
	return rs.NewRSPutStream(servers, hash, size) //调用dataServer生成文件，暂时还未写入
}
