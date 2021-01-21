package models

import "net/http"

var (
	WSHub *Hub
	Nats  *NatsManager
	GRM   *GoRoutineManager
	KCM   *KubernetesClientManager
	MSM   *MailServerManager
	LM    *LDAPManager
	AEM   *AnchoreEngineManager
	CPM   *http.Server
)
