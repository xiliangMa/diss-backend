package models

import (
	"github.com/casbin/casbin"
	"net/http"
)

var (
	WSHub    *Hub
	Nats     *NatsManager
	GRM      *GoRoutineManager
	KCM      *KubernetesClientManager
	MSM      *MailServerManager
	LM       *LDAPManager
	CPM      *http.Server
	Enforcer *casbin.Enforcer
)
