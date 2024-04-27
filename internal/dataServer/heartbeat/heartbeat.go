package heartbeat

import (
	"os"
	"storage/conf"
	"storage/internal/pkg/rabbitmq"
	"time"
)

func StartHeartbeat() {
	q := rabbitmq.New(conf.RABBITMQ_SERVER)
	defer q.Close()
	for {
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
