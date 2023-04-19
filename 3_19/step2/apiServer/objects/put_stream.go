package objects

import (
	"fmt"
	"project/3_19/step2/apiServer/heartbeat"
	"project/go-object-storage/src/lib/objectstream"
)

func putStream(object string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer() //随机返回一个当前存活节点的地址
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewPutStream(server, object), nil
}
