package objectstream

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type TempPutStream struct {
	Server string
	Uuid   string
}

// NewTempPutStream 向dataServer的/temp发送POST请求，创建文件（不写入内容），object为[hash.writer分片索引]
func NewTempPutStream(server, object string, size int64) (*TempPutStream, error) {
	request, e := http.NewRequest("POST", "http://"+server+"/temp/"+object, nil)
	if e != nil {
		return nil, e
	}
	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	response, e := client.Do(request)
	if e != nil {
		return nil, e
	}
	uuid, e := io.ReadAll(response.Body) //uuid由dataServer中的post方法生成
	if e != nil {
		return nil, e
	}
	return &TempPutStream{server, string(uuid)}, nil
}

// 向dataServer发送patch请求，将传入的[]byte作为该请求的消息体
func (w *TempPutStream) Write(p []byte) (n int, err error) {
	request, e := http.NewRequest("PATCH", "http://"+w.Server+"/temp/"+w.Uuid, strings.NewReader(string(p)))
	if e != nil {
		return 0, e
	}
	client := http.Client{}
	r, e := client.Do(request)
	if e != nil {
		return 0, e
	}
	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}
	return len(p), nil
}

// Commit 成功时，将/temp中的文件写入到/objects中;失败时，直接将该文件从/temp中删除
func (w *TempPutStream) Commit(good bool) {
	method := "DELETE"
	if good {
		method = "PUT"
	}
	request, _ := http.NewRequest(method, "http://"+w.Server+"/temp/"+w.Uuid, nil)
	client := http.Client{}
	client.Do(request)
}

func NewTempGetStream(server, uuid string) (*GetStream, error) {
	return newGetStream("http://" + server + "/temp/" + uuid)
}
