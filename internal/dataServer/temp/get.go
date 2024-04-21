package temp

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".dat")
	if len(files) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	file := files[0]
	f, err := os.Open(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, f)
}
