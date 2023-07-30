package dao

import (
	"log"
	"testing"
)

func TestGet(t *testing.T) {
	ans := Get("speak.txt")
	log.Println(ans)
	if len(ans) != 2 {
		t.Error("Get error")
	}
}
