package utils

import (
	"crypto/tls"
	"github.com/astaxie/beego"
	"github.com/elastic/go-elasticsearch"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_ESClient(t *testing.T) {
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
		t.Log("ESClient create fail, ", err)
	}
	t.Log("ESClient create success, info: ", client.Info)
}

