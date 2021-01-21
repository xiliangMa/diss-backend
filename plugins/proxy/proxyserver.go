package proxy

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyServer struct {
	TargetUrl string
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
