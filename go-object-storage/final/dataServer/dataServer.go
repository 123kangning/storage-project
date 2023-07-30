package main

import (
	"log"
	"net/http"
	"os"
	"storage/final/dataServer/heartbeat"
	"storage/final/dataServer/locate"
	"storage/final/dataServer/objects"
	"storage/final/dataServer/temp"
)

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler) //GetStream对象调用
	http.HandleFunc("/temp/", temp.Handler)       //apiServer中的TempPutStream对象调用
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
