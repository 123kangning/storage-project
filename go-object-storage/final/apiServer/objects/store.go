package objects

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"project/go-object-storage/final/apiServer/locate"
	"project/go-object-storage/src/lib/utils"
)

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.Exist(url.PathEscape(hash)) { //如果该对象已经存在，直接返回
		return http.StatusOK, nil
	}
	//调用dataServer POST创建文件，但不写入 获取接口服务节点存储对象的流
	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		return http.StatusInternalServerError, e
	}
	/*
		   	当reader从r中读数据时，同时写入stream
			写入stream时，执行TempPutStream的Write方法（向dataServer发送patch请求，写入数据）
	*/
	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader) //计算该文件内容的hash值
	if d != hash {
		// hash值不一致，删除临时文件
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
	//hash值一致，顺利写入
	stream.Commit(true)
	return http.StatusOK, nil
}
