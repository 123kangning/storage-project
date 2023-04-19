package objects

import (
	"io"
	"net/http"
)

// 将r中的数据写入
func storeObject(r io.Reader, object string) (int, error) {
	//创建一个io通道，并返回这个通道的写入流（putStream）
	stream, e := putStream(object)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}
	//将r中的数据写入putStream中，即写入io通道中，
	io.Copy(stream, r)
	e = stream.Close()
	if e != nil {
		return http.StatusInternalServerError, e
	}
	return http.StatusOK, nil
}
