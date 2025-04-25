package test

import (
	"encoding/json"
	"log"
	"storage/infra/myes"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	go myes.Init()
	time.Sleep(time.Second)
	testGet("文档")
	testGet("speak")
	testGet("一")
	testGet("二")
}
func testGet(name string) {
	f1s, total, err := myes.GetFile(name, 0, 10)
	if err != nil {
		log.Println("err1 = ", err)
	}
	fs, _ := json.Marshal(f1s)
	log.Println(total, string(fs))
}
