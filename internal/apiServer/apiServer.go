package main

import (
	"storage/infra/myes"
	"storage/internal/apiServer/heartbeat"
)

/**
 * @Description:	起点，处理各个请求
 */
func main() {
	//os.Setenv("LISTEN_ADDRESS", "10.29.2.1:12345")
	go heartbeat.ListenHeartbeat()
	myes.Init()
	r := InitRouter()
	err := r.Run(":8081")
	if err != nil {
		panic(err)
	}
}
