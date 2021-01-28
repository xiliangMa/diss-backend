package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"testing"
)

func Test_NatsV2(t *testing.T) {
	var servers = "nats://diss:dissP@ssw0rd@127.0.0.1:4222,nats://diss:dissP@ssw0rd@127.0.0.1:4223,nats://diss:dissP@ssw0rd@127.0.0.1:4224"
	nc, err := nats.Connect(servers)
	if err != nil {
		log.Fatal(err)
	}
	sc, err := stan.Connect("diss-cluster", "1111", stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
	}
	sc.Subscribe("test", func(m *stan.Msg) {
		log.Printf("[Received] %+v", m)
	})
	sc.Publish("test", []byte("nats cluster test111111111"))
	select {}
}
