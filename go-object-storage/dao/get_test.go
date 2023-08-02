package dao

import (
	"log"
	"testing"
)

func TestGet(t *testing.T) {
	ans := Get("十号文档.txt")
	log.Println(ans)
	if ans.Name != "十号文档.txt" {
		t.Error("Get error")
	}
}
