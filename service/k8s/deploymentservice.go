package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type DeploymentService struct {
	DeploymentInterface appsv1.DeploymentInterface
	Cluster             *models.Cluster
	Close               chan bool
}

func (this *DeploymentService) List() (*v1.DeploymentList, error) {
	return this.DeploymentInterface.List(metav1.ListOptions{})
}

func (this *DeploymentService) Wtach() {
	deployWatch, err := this.DeploymentInterface.Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// 开启 watch 事件
Done:
	for {
		select {
		case <-this.Close:
			logs.Info("Close deploymentWatch, cluster: %s", this.Cluster.Name)
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

				ns := models.Deployment{}
				ns.Id = id
				ns.Name = name
				ns.AccountName = models.Account_Admin
				ns.ClusterName = clusterName
				ns.KMetaData = string(KMetaData)
				ns.KSpec = string(KSpec)
				ns.KStatus = string(KStatus)

				logs.Info("Watch >>> Deployment: %s <<<, >>> Cluster: %s <<<,  >>> NameSpace: %s <<<, >>> EventType: %s <<<", id, clusterId, name, event.Type)
				switch event.Type {
				case watch.Added:
					ns.Add()
				case watch.Modified:
					ns.Delete()
					ns.Add()
				case watch.Deleted:
					ns.Delete()
				case watch.Bookmark:
					//todo
				case watch.Error:
					//todo
				}
			} else {
				// 如果 watch 异常退回重新 watch
				logs.Error("deploymentWatch chan has been close!!!!")
				break Done
			}
		}
	}
}
