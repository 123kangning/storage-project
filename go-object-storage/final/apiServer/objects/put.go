package objects

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"project/go-object-storage/final/apiServer/heartbeat"
	"project/go-object-storage/src/lib/es"
	"project/go-object-storage/src/lib/rs"
	"project/go-object-storage/src/lib/utils"
	"strings"
)

/**
 * @Description: 核心函数PUT
 * @param w
 * @param r
 */
func put(w http.ResponseWriter, r *http.Request) {
	log.Println("api.object.put")

	hash := utils.GetHashFromHeader(r.Header) //先获取hash值
	if hash == "" {                           //hash值为空，记得返回问题
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	size := utils.GetSizeFromHeader(r.Header) //从头部获取size信息
	/*
		不同的编程语言和系统可能对数据的处理方式有所不同，因此在存储请求数据时，需要对散列值进行转义，以确保数据在不同的环境中都能被正确处理
	*/
	c, e := storeObject(r.Body, url.PathEscape(hash), size) //存储文件到/objects,返回状态码以及error
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2] //组成名字
	e = es.AddVersion(name, hash, size)                //更新版本
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

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
