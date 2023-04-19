package main

import (
	"log"
	"net/http"
	"os"
	"project/3_19/step2/dataServer/heartbeat"
	"project/3_19/step2/dataServer/locate"
	"project/3_19/step2/dataServer/objects"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
