package main

import (
	"net/http"
	"storage/internal/apiServer/heartbeat"
	"storage/internal/apiServer/locate"
	"storage/myes"
)

/**
 * @Description:	起点，处理各个请求
 */
func main() {
	//os.Setenv("LISTEN_ADDRESS", "10.29.2.1:12345")
	go heartbeat.ListenHeartbeat()
	go myes.Init()
	go myes.Run()
	r := InitRouter()
	err := r.Run(":8081")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/locate/", locate.Handler)
}
