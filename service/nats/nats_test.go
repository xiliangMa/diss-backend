package nats

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"testing"
)

var (
	CLuster_Id = "diss-cluster"
	Client_Id  = "diss-agent22222"
	Nats_Url   = "nats://diss:dissP@ssw0rd@localhost:4222"
)

//func Test_Nats_Pub(t *testing.T) {
//	nc, err := stan.Connect(CLuster_Id, Client_Id, stan.NatsURL(Nats_Url))
//	err = nc.Publish("Client1-Task", []byte("task----11----- "))
//	err = nc.Publish("Client2-Task", []byte("task----22----- "))
//
//	if err != nil {
//		t.Logf("Pub......... fail, err: %s", err)
//	}
//}

func Test_Nats_Sub(t *testing.T) {
	subject := "89e38324-e2bd-4c4b-abf3-d132f268187c"
	nc, err := stan.Connect(CLuster_Id, Client_Id, stan.NatsURL(Nats_Url))
	if err != nil {
		t.Errorf("Client1 connet Nats server fail. err: %s", err)
	}
	_, err = nc.Subscribe(subject, func(m *stan.Msg) {
		fmt.Println("22222222", string(m.Data))
		fmt.Println("Client1 get a message:", string(m.Data))
	}, stan.DurableName(subject))
	if err != nil {
		t.Errorf("Client1 get mesage fail. err: %s", err)
	}
}

//
//func Test_Nats_Sub2(t *testing.T) {
//	subject := "bogon_Task"
//	nc, err := stan.Connect(CLuster_Id, Client_Id, stan.NatsURL(Nats_Url))
//	if err != nil {
//		t.Errorf("Client2 connet Nats server fail. err: %s", err)
//	}
//	_, err = nc.Subscribe(subject, func(m *stan.Msg) {
//		fmt.Println("Client2 get a message:", string(m.Data))
//	}, stan.DurableName(subject))
//	if err != nil {
//		t.Errorf("Client2 get mesage fail. err: %s", err)
//	}
//}
