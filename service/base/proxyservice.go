package base

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/plugins/proxy"
)

type ProxyService struct {
	TargetUrl string
	Close     bool
}

func (this *ProxyService) ScopeProxyOperator() error {
	// 默认清空原有代理服务对象
	if models.CPM != nil {
		models.CPM.Close()
	}
	// 创建新的代理对象
	proxy := proxy.ProxyServer{TargetUrl: this.TargetUrl}
	srv := proxy.StartServer()
	models.CPM = srv
	logs.Info("Stop proxy server success")

	//if this.Close && models.CPM != nil {
	//	err := models.CPM.Close()
	//	if err != nil {
	//		return err
	//	}
	//	models.CPM = nil
	//	return nil
	//} else {
	//	proxy := proxy.ProxyServer{TargetUrl: this.TargetUrl}
	//	srv := proxy.StartServer()
	//	models.CPM = srv
	//	logs.Info("Stop proxy server success")
	//}
	return nil
}
