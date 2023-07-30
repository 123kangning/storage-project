package objects

import (
	"fmt"
	"storage/final/apiServer/heartbeat"
	"storage/final/apiServer/locate"
	"storage/src/lib/rs"
)

//Get @Description:	取出对象，注意按照版本取，URL里没指定就取最新的版本
//
//func Get(c *gin.Context) {
//	name := c.Query("name")
//	meta, e := es.GetMetadata(name, version) //	从ES中取出来
//	if e != nil {
//		log.Println(e)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	if meta.Hash == "" { //空的就是没找到咯，del里面删除不就是置空吗
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	hash := url.PathEscape(meta.Hash) //对字符串转义，保证传输过程中的安全
//	stream, e := GetStream(hash, meta.Size)
//	if e != nil {
//		log.Println(e)
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	offset := utils.GetOffsetFromHeader(r.Header)
//	if offset != 0 {
//		stream.Seek(offset, io.SeekCurrent)
//		w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", offset, meta.Size-1, meta.Size))
//		w.WriteHeader(http.StatusPartialContent)
//	}
//	//可选 gzip进行压缩
//	acceptGzip := false
//	encoding := r.Header["Accept-Encoding"]
//	for i := range encoding {
//		if encoding[i] == "gzip" {
//			acceptGzip = true
//			break
//		}
//	}
//	if acceptGzip {
//		w.Header().Set("content-encoding", "gzip")
//		w2 := gzip.NewWriter(w)
//		io.Copy(w2, stream)
//		w2.Close()
//	} else {
//		io.Copy(w, stream)
//	}
//	stream.Close()
//}

func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	locateInfo := locate.Locate(hash)
	if len(locateInfo) < rs.DATA_SHARDS { // 如果定位的数据服务节点数少于数据恢复的最低数量，则无法恢复完整数据，返回定位失败错误
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)
	}
	// 选择数据服务节点，用于数据恢复
	dataServers := make([]string, 0)
	if len(locateInfo) != rs.ALL_SHARDS { //从其他数据节点中选择一些节点来进行数据恢复
		dataServers = heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}
