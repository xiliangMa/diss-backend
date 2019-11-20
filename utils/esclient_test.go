package utils

import (
	"context"
	"crypto/tls"
	"github.com/elastic/go-elasticsearch"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_ESClient(t *testing.T) {
	Adress := []string{"http://122.51.240.195:9200"}
	UserName := "elastic"
	Password := "changeme"
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

	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		t.Log("ESClient create fail, ", err)
	} else {
		info, _ := es.Info()
		t.Log("ESClient create success, info: ", info)

		res, _ := es.API.Search(es.Search.WithContext(context.Background()),
			es.Search.WithIndex("metric*"),
			es.Search.WithBody(strings.NewReader(`{"query" : { "match_all":{} }}`)),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty())
		t.Log("the metric index is ok ", res)
	}


}
