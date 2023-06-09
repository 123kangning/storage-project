package versions

import (
	"encoding/json"
	"log"
	"net/http"
	"project/go-object-storage/src/lib/es"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) { //获取所有的Metadata集合
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0
	size := 1000
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	for { //获取该name对应的所有 Metadata 的集合
		metas, e := es.SearchAllVersions(name, from, size)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range metas {
			b, _ := json.Marshal(metas[i])
			w.Write(b)
			w.Write([]byte("\n"))
		}
		if len(metas) != size { //读完了
			return
		}
		from += size
	}
}
