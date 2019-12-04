package utils

import (
	"crypto/tls"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/elastic/go-elasticsearch"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	cfg elasticsearch.Config
)

func initCfg() elasticsearch.Config {
	var Adress []string
	var UserName string
	var Password string
	if os.Getenv("RunMode") != "prod" {
		Adress = strings.Split(beego.AppConfig.String("es::Address"), ",")
		UserName = beego.AppConfig.String("es::UserName")
		Password = beego.AppConfig.String("es::Password")
	} else {
		Adress = strings.Split(os.Getenv("Address"), ",")
		UserName = os.Getenv("UserName")
		Password = os.Getenv("Password")
	}

	cfg = elasticsearch.Config{
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

	return cfg
}

func GetESClient() (*elasticsearch.Client, error) {
	client, err := elasticsearch.NewClient(initCfg())

	if err != nil {
		logs.Error("Error creating the ES client: %s", err)
		err = err
	}
	return client, err
}
