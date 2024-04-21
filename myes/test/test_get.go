package main

import (
	"encoding/json"
	"log"
	"storage/myes"
	"time"
)

func main() {
	go myes.Init()
	time.Sleep(time.Second)
	testGet("文档")
	testGet("speak")
	testGet("一")
	testGet("二")
}
func testGet(name string) {
	f1s, err := myes.GetFile(name)
	if err != nil {
		log.Println("err1 = ", err)
	}
	fs, _ := json.Marshal(f1s)
	log.Println(string(fs))
}
