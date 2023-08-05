package main

import (
	"net/http"
	"os"
	"storage/final/apiServer/heartbeat"
	"storage/final/apiServer/locate"
	"storage/myes"
)

/**
 * @Description:	起点，处理各个请求
 */
func main() {
	os.Setenv("LISTEN_ADDRESS", "10.29.2.1:12345")
	go heartbeat.ListenHeartbeat()
	go myes.Init()
	go myes.Run()
	r := InitRouter()
	err := r.Run(os.Getenv("LISTEN_ADDRESS"))
	if err != nil {
		panic(err)
	}
	//http.HandleFunc("/objects/", objects.Handler) //webServer中的uploadHandler、downloadHandler调用
	http.HandleFunc("/locate/", locate.Handler)
	//http.HandleFunc("/versions/", versions.Handler) //webServer中的listHandler调用
	//log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
