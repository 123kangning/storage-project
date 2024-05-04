package heartbeat

import (
	"os"
	"storage/conf"
	"storage/internal/pkg/rabbitmq"
	"time"
)

func StartHeartbeat() {
	q := rabbitmq.New(conf.RabbitmqServer)
	defer q.Close()
	for {
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
