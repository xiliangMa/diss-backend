package k8s

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type DeploymentService struct {
	ClientGo models.ClientGo
	Cluster  *models.Cluster
	Close    chan bool
}

func (this *DeploymentService) List() (*v1.DeploymentList, error) {
	return this.ClientGo.ClientSet.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
}

func (this *DeploymentService) Wtach() {
	deployWatch, err := this.ClientGo.ClientSet.AppsV1().Deployments("").Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		logs.Error("Wtach deployment error: %s  ", err)
		return
	}
	// 开启 watch 事件
Retry:
	for {
		select {
		case <-this.Close:
			logs.Info("Close deploymentWatch, cluster: %s", this.Cluster.Name)
			deployWatch.Stop()
			return
		case event, ok := <-deployWatch.ResultChan():
			if event.Object != nil || ok {
				object := event.Object.(*v1.Deployment)
				id := string(object.UID)
				name := object.Name
				clusterId := this.Cluster.Id
				clusterName := this.Cluster.Name

				KMetaData, _ := json.Marshal(object.ObjectMeta)
				KSpec, _ := json.Marshal(object.Spec)
				KStatus, _ := json.Marshal(object.Status)

				deploy := models.Deployment{}
				deploy.Id = id
				deploy.Name = name
				deploy.AccountName = models.Account_Admin
				deploy.ClusterName = clusterName
				deploy.KMetaData = string(KMetaData)
				deploy.KSpec = string(KSpec)
				deploy.KStatus = string(KStatus)

				logs.Info("Watch >>> Deployment: %s <<<, >>> ClusterOBJ: %s <<<,  >>> NameSpace: %s <<<, >>> EventType: %s <<<", id, clusterId, name, event.Type)
				switch event.Type {
				case watch.Added:
					deploy.Add()
				case watch.Modified:
					deploy.Delete()
					deploy.Add()
				case watch.Deleted:
					deploy.Delete()
				case watch.Bookmark:
					//todo
				case watch.Error:
					//todo
				}
			} else {
				// 如果 watch 异常退回重新 watch
				logs.Warn("DeploymentWatch chan has been close!!!!, cluster: %s", this.Cluster.Name)

				// 清除全局 GRM（携程对象）
				watchType := this.Cluster.Id + `_` + utils.Deployment
				delete(models.GRM.GoRoutineMap, watchType)
				logs.Info("Remove DeploymentWatch from global GRM object, cluster: %s", this.Cluster.Name)

				// 清除数据库数据
				deploy := models.Deployment{}
				deploy.ClusterName = this.Cluster.Name
				deploy.Delete()

				// 重启 watch 携程
				deployService := DeploymentService{Cluster: this.Cluster, ClientGo: this.ClientGo, Close: make(chan bool)}
				models.GRM.GoRoutineMap[watchType] = deployService

				logs.Info("Retry DeploymentWatch, cluster: %s", this.Cluster.Name)
				go deployService.Wtach()
				break Retry
			}
		}
	}
}
