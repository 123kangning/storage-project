package temp

import (
	"compress/gzip"
	"io"
	"net/url"
	"os"
	"project/go-object-storage/final/dataServer/locate"
	"project/go-object-storage/src/lib/utils"
	"strconv"
	"strings"
)

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}

func commitTempObject(datFile string, tempinfo *tempInfo) {
	f, _ := os.Open(datFile)
	defer f.Close()
	d := url.PathEscape(utils.CalculateHash(f))
	f.Seek(0, io.SeekStart)
	w, _ := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + tempinfo.Name + "." + d)
	w2 := gzip.NewWriter(w) //将gzip压缩后的文件写入w中
	io.Copy(w2, f)
	w2.Close()
	//删除副本文件
	os.Remove(datFile)
	//加入hash-id键值对
	locate.Add(tempinfo.hash(), tempinfo.id())
}
