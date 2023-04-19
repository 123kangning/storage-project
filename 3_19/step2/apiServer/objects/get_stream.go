package objects

import (
	"fmt"
	"io"
	"project/3_19/step2/apiServer/locate"
	"project/go-object-storage/src/lib/objectstream"
)

func getStream(object string) (io.Reader, error) {
	//得到存储该文件的节点地址
	server := locate.Locate(object)
	if server == "" { //定位失败，可能是该文件不存在
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return objectstream.NewGetStream(server, object)
}
