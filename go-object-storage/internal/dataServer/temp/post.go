package temp

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*
文件 "STORAGE_ROOT" + "/temp/" + uuid 负责保存指定uuid 的tempInfo文件信息对象
文件 "STORAGE_ROOT" + "/temp/" + t.Uuid + ".dat" 为指定uuid的暂存文件
文件 "STORAGE_ROOT" + "/objects/" + hash.id.shardHash 为真正存储的zip压缩文件
*/

type tempInfo struct {
	Uuid string
	Name string //[hash.writer分片索引]
	Size int64
}

// 创建/temp中的两个文件，一个配置，一个副本
func post(w http.ResponseWriter, r *http.Request) {
	//生成uuid
	log.Println("post")
	output, _ := exec.Command("uuidgen").Output()
	uuid := strings.TrimSuffix(string(output), "\n")
	//name解析出来为[hash.writer分片索引]
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := tempInfo{uuid, name, size}
	//写入关联uuid文件的tempinfo对象
	e = t.writeToFile()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//创建.dat文件
	os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid + ".dat")
	w.Write([]byte(uuid))
}

func (t *tempInfo) writeToFile() error {
	f, e := os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid)
	if e != nil {
		return e
	}
	defer f.Close()
	b, _ := json.Marshal(t)
	f.Write(b)
	return nil
}
