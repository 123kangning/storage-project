package main

import (
	"storage/myes"
	"time"
)

func main() {
	go myes.Init()
	time.Sleep(time.Second) //客户端启动需要一定的时间，所以要睡1s,否则client来不及赋值，就为nil
	myes.Run()
}
