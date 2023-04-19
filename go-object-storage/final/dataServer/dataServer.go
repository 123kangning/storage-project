package main

import (
	"log"
	"net/http"
	"os"
	"project/go-object-storage/final/dataServer/heartbeat"
	"project/go-object-storage/final/dataServer/locate"
	"project/go-object-storage/final/dataServer/objects"
	"project/go-object-storage/final/dataServer/temp"
)

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
