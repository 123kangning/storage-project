package main

import (
	"log"
	"net/http"
	"os"
	"project/3_19/step2/apiServer/heartbeat"
	"project/3_19/step2/apiServer/objects"
)

func main() {
	go heartbeat.ListenHeartbeat()                //心跳保活机制
	http.HandleFunc("/objects/", objects.Handler) //将put和get请求转发到dataServer层，本身不对数据进行操作
	//http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
