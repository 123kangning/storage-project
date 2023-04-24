package main

import (
	"log"
	"net/http"
	"os"
	"project/go-object-storage/final/apiServer/heartbeat"
	"project/go-object-storage/final/apiServer/locate"
	"project/go-object-storage/final/apiServer/objects"
	"project/go-object-storage/final/apiServer/temp"
	"project/go-object-storage/final/apiServer/versions"
)

/**
 * @Description:	起点，处理各个请求
 */
func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler) //webServer中的uploadHandler、downloadHandler调用
	http.HandleFunc("/temp/", temp.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler) //webServer中的listHandler调用
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
