package proxy

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ProxyServer struct {
	TargetUrl string
	Token     string
	Method    string
	Body      interface{}
	Err       chan string
}

func (this *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remote, err := url.Parse(this.TargetUrl)
	if err != nil {
		logs.Error("Parse TargetUrl failed, Err: %s.", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
	logs.Info("Reverse proxy success, TargetUrl: %s.", this.TargetUrl)
}

func (this *ProxyServer) StartServer() *http.Server {
	port := beego.AppConfig.String("proxy::Port")
	srv := &http.Server{Addr: "0.0.0.0:" + port}
	srv.Handler = this
	// todo get ListenAndServeTLS err
	go srv.ListenAndServeTLS("conf/ca.crt", "conf/server.key")
	logs.Info("Start proxy server success, Port: %s.", port)
	return srv
}

func (this *ProxyServer) Request(user string, pwd string) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli := &http.Client{Transport: tr}

	if this.Method == "" {
		this.Method = "GET"
	}

	if !strings.HasPrefix(this.TargetUrl, "http://") && !strings.HasPrefix(this.TargetUrl, "https://") {
		return nil, errors.New("Invalid Scheme")
	}

	var requestBody bytes.Buffer
	if this.Body != nil {
		json.NewEncoder(&requestBody).Encode(this.Body)
	}

	req, err := http.NewRequest(this.Method, this.TargetUrl, &requestBody)

	req.Header.Add("content-type", "application/json")

	if this.Token != "" {
		req.Header.Add("x-auth-token", this.Token)
		req.Header.Add("Authorization", this.Token)
	} else {
		req.SetBasicAuth(user, pwd)
	}
	resp, err := cli.Do(req)
	return resp, err
}
