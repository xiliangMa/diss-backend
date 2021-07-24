package k8s

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type ClusterService struct {
	Cluster *models.Cluster
}

func (this *ClusterService) UpdateCluster() *models.Result {
	result := models.Result{Code: http.StatusOK}
	// 对比集群名是否变更
	dbCluster := &models.Cluster{Id: this.Cluster.Id}
	dbCluster = dbCluster.Get()
	if dbCluster == nil {
		result.Code = utils.ClusterNoExistErr
		return &result
	}
	result = this.Cluster.Update()

	if dbCluster.Name == this.Cluster.Name {
		return &result
	}

	// 停止watch 携程、清除全局集群客户端、清除集群数据
	var list []*models.Cluster
	list = append(list, this.Cluster)
	k8sClearService := K8sClearService{ClusterList: list, DropCluster: false}
	k8sClearService.ClearAll()

	// 更新全局集群客户端、重新watch集群
	k8sWatchService := K8sWatchService{Cluster: this.Cluster}
	k8sWatchService.WatchCluster()
	return &result
}

func (this *ClusterService) ListCluster(from, limit int) *models.Result {
	result := models.Result{Code: http.StatusOK}
	err, _, list := this.Cluster.BaseList(from, limit)
	if err != nil {
		result.Code = utils.GetClusterErr
		result.Message = err.Error()
		return &result
	}
	var kubeScan models.KubeScan
	for _, cluster := range list {
		kubeScan.ClusterId = cluster.Id
		cluster.KubeVulnCount = len(kubeScan.Vulnerabilities)
		cluster.Update()
	}
	result = this.Cluster.List(from, limit)
	return &result
}
