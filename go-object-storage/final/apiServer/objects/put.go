package objects

import (
	"log"
	"net/http"
	"net/url"
	"project/go-object-storage/src/lib/es"
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
