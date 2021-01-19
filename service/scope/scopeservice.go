package scope

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	v1 "k8s.io/api/core/v1"
	"net/http"
	"os"
	"strings"
)

type ScopeService struct {
	Cluster  *models.Cluster
	IsActive bool
	*v1.Pod
	*v1.Namespace
}

/**
 ** install or active scope
 ** default port 32666
 */
func (this *ScopeService) ActiveOrDisableScope() models.Result {
	result := models.Result{Code: http.StatusOK}
	clusterName := this.Cluster.Name

	// 检测集群客户端
	dymaicClient := models.GetDymaicClient(this.Cluster)
	client := models.GetClient(this.Cluster)
	if dymaicClient.ErrMessage != "" {
		result.Code = utils.CreateK8sClientErr
		result.Message = dymaicClient.ErrMessage
		logs.Error("Scope install fail, ClusterName: %s, Err: %s", clusterName, result.Message)
		return result
	}

	// 读取yml
	f, err := os.Open(utils.GetScopeYml())
	if err != nil {
		result.Code = utils.OpenFileErr
		result.Message = err.Error()
		logs.Error("Open scope yml fail, ClusterName: %s, Err: %s", clusterName, result.Message)
		return result
	}

	// 创建 scope app
	kubernetesHandler := models.KubernetesHandler{ClientGo: &client, DymaicClient: &dymaicClient, IsActive: this.IsActive, File: f}
	return kubernetesHandler.CreateOrDeleteResourceByYml()
}

func (this *ScopeService) UpdatetClusterScopeUrlAndStatus(isDelete bool) models.Result {
	result := models.Result{Code: http.StatusOK}
	ns := this.Pod.Namespace
	name := this.Pod.Name
	node := this.Pod.Spec.NodeName
	if ns == utils.GetScopeNameSpace() && strings.Contains(name, utils.GetScopeAppName()) {
		podStatus := string(this.Pod.Status.Phase)
		hi := models.HostInfo{}
		hi.HostName = node
		hi.ClusterId = this.Cluster.Id
		// 获取主机信息
		result := hi.List()
		if result.Code != http.StatusOK {
			return result
		}
		if result.Code == http.StatusOK && result.Data != nil {
			// 更新url
			scopeUrlJson := make(map[string]string)
			data := result.Data.(map[string]interface{})
			hiList := data[models.Result_Items].([]*models.HostInfo)
			if hiList != nil && len(hiList) > 0 {
				if hiList[0].InternalAddr != "" {
					scopeUrlJson[models.Cluster_Scope_InternalUrl] = `http://` + hiList[0].InternalAddr + `:` + models.Cluster_Scope_UrlPort
				}
				if hiList[0].PublicAddr != "" {
					scopeUrlJson[models.Cluster_Scope_PublicUrl] = `http://` + hiList[0].PublicAddr + `:` + models.Cluster_Scope_UrlPort
				}
				scopeUrlStr, _ := json.Marshal(scopeUrlJson)
				this.Cluster.ScopeUrl = string(scopeUrlStr)

				// 更新状态
				// 注意： 删除scope时由于只判断 pod 状态不能确定scope是否被删除完成。
				// 需要判断 ns 是否存在，具体看 nsWatch 的代码
				if podStatus == models.Pod_Container_Statue_Running {
					this.Cluster.SocpeStatus = models.Cluster_Scope_Operator_Status_Actived
					if isDelete {
						this.Cluster.SocpeStatus = models.Cluster_Scope_Operator_Status_Disableing
					}
					result = this.Cluster.Update()
				}
			}

		}
	}
	return result
}

func (this *ScopeService) CheckScopeIsDisable() {
	ns := this.Namespace.Name
	if ns == utils.GetScopeNameSpace() {
		this.Cluster.SocpeStatus = models.Cluster_Scope_Operator_Status_Disabled
		this.Cluster.ScopeUrl = ""
		this.Cluster.Update()
	}

}
