package main

import (
	"log"
	"net/http"
	"os"
	"storage/internal/dataServer/heartbeat"
	"storage/internal/dataServer/locate"
	"storage/internal/dataServer/objects"
	"storage/internal/dataServer/temp"
)

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler) //GetStream对象调用
	http.HandleFunc("/temp/", temp.Handler)       //apiServer中的TempPutStream对象调用
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
