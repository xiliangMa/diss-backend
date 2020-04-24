package global

import (
	modelws "github.com/xiliangMa/diss-backend/models/ws"
	"github.com/xiliangMa/diss-backend/service/nats"
)

var (
	WSHub       *modelws.Hub
	NatsManager *nats.NatsManager
)
