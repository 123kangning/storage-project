package temp

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// 文件在数据层真正的put过程
func put(w http.ResponseWriter, r *http.Request) {
	//uuid为[hash.writer分片索引]
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 读取临时文件信息反序列化为结构体
	tempInfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	//获取对应.dat文件
	datFile := infoFile + ".dat"
	f, e := os.Open(datFile)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	//移除/temp目录下 临时文件信息
	os.Remove(infoFile)
	if actual != tempInfo.Size { //出现意外情况，不是预期行为
		// 删除/temp目录下 临时数据文件
		os.Remove(datFile)
		log.Println("actual size mismatch, expect", tempInfo.Size, "actual", actual)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//临时文件转正式文件 tempInfo->原文件属性，datFile->之后.dat文件的名称
	commitTempObject(datFile, tempInfo)
}
