package main

import (
	"log"
	"myes"
	"net/http"
	"os"
	"storage/final/apiServer/heartbeat"
	"storage/final/apiServer/locate"
	"storage/final/apiServer/objects"
	"storage/final/apiServer/temp"
	"storage/final/apiServer/versions"
)

/**
 * @Description:	起点，处理各个请求
 */
func main() {
	myes.Init()
	go heartbeat.ListenHeartbeat()
	r := InitRouter()
	err := r.Run(os.Getenv("LISTEN_ADDRESS"))
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/objects/", objects.Handler) //webServer中的uploadHandler、downloadHandler调用
	http.HandleFunc("/temp/", temp.Handler)       //和/objects中的post一起看,head只查看，主要看put
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler) //webServer中的listHandler调用
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
