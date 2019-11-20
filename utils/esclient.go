package utils

import (
	"crypto/tls"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/elastic/go-elasticsearch"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	ESClient *elasticsearch.Client
)

func InitEsClient() {
	Adress := strings.Split(beego.AppConfig.String("es::Address"), ",")
	UserName := beego.AppConfig.String("es::UserName")
	Password := beego.AppConfig.String("es::Password")

	cfg := elasticsearch.Config{
		Addresses: Adress,
		Username:  UserName,
		Password:  Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
	}


	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logs.Error("Error creating the ES client: %s", err)
	}
	ESClient = client
}

func GetESClient() *elasticsearch.Client {
	if ESClient == nil {
		InitEsClient()
	}
	return ESClient
}
