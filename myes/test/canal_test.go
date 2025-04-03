package test

import (
	"encoding/json"
	"log"
	"storage/myes"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	go myes.Init()
	time.Sleep(time.Second) //客户端启动需要一定的时间，所以要睡1s,否则client来不及赋值，就为nil
	myes.Run()
}

func TestGet(t *testing.T) {
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
