package models

import (
	"net/http"
)

var (
	WSHub        *Hub
	Nats         *NatsManager
	GRM          *GoRoutineManager
	KCM          *KubernetesClientManager
	MSM          *MailServerManager
	LM           *LDAPManager
	CPM          *http.Server
	GlobalCasbin *CasbinManager
)
