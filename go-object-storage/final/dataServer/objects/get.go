package objects

import (
	"net/http"
	"strings"
)

// 返回压缩后的文件
func get(w http.ResponseWriter, r *http.Request) {
	file := getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	SendFile(w, file)
}
