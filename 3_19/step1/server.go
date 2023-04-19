package main

import (
	"log"
	"net/http"
	"os"
	"project/3_19/step1/objects"
)

func main() {
	http.HandleFunc("/objects/", objects.Handler)                    //处理以"/object/"开头的URL，那么就交给objects.Handler来具体实现
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil)) //就是记录日志 这个不是重点 但服务器必须要记录操作日志
}
