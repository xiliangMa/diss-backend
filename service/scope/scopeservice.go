package scope

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"log"
	"net/http"
	"os"
)

type ScopeService struct {
	Cluster  *models.Cluster
	IsActive bool
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
		log.Fatal(err)
	}

	// 创建 scope app
	kubernetesHandler := models.KubernetesHandler{ClientGo: &client, DymaicClient: &dymaicClient, IsActive: this.IsActive, File: f}
	return kubernetesHandler.CreateOrDeleteResourceByYml()

}
