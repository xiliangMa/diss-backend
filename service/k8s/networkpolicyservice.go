package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type NetworkPolicyService struct {
	ClientGo      models.ClientGo
	Cluster       *models.Cluster
	Close         chan bool
	NetworkPolicy *models.NetworkPolicy
}

func (this *NetworkPolicyService) List() (*v1.NetworkPolicyList, error) {
	return this.ClientGo.ClientSet.NetworkingV1().NetworkPolicies("").List(metav1.ListOptions{})
}

func (this *NetworkPolicyService) Wtach() {
	ns := ""
	if this.NetworkPolicy != nil {
		ns = this.NetworkPolicy.NameSpaceName
	}
	netpolWatch, err := this.ClientGo.ClientSet.NetworkingV1().NetworkPolicies(ns).Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// 开启 watch 事件
Retry:
	for {
		select {
		case <-this.Close:
			logs.Info("Close NetworkPolicyWatch, cluster: %s", this.Cluster.Name)
			netpolWatch.Stop()
			return
		case event, ok := <-netpolWatch.ResultChan():
			if event.Object != nil || ok {
				object := event.Object.(*v1.NetworkPolicy)
				id := string(object.UID)
				name := object.Name
				clusterId := this.Cluster.Id
				clusterName := this.Cluster.Name

				KMetaData, _ := json.Marshal(object.ObjectMeta)
				KSpec, _ := json.Marshal(object.Spec)

				netpol := models.NetworkPolicy{}
				netpol.Id = id
				netpol.Name = name
				netpol.AccountName = models.Account_Admin
				netpol.ClusterName = clusterName
				netpol.KMetaData = string(KMetaData)
				netpol.KSpec = string(KSpec)

				logs.Info("Watch >>> NetworkPolicy: %s <<<, >>> Cluster: %s <<<,  >>> NameSpace: %s <<<, >>> EventType: %s <<<", id, clusterId, name, event.Type)
				switch event.Type {
				case watch.Added:
					netpol.Add()
				case watch.Modified:
					netpol.Delete()
					netpol.Add()
				case watch.Deleted:
					netpol.Delete()
				case watch.Bookmark:
					//todo
				case watch.Error:
					//todo
				}
			} else {
				// 如果 watch 异常退回重新 watch
				logs.Warn("NetworkPolicyWatch chan has been close!!!!, cluster: %s", this.Cluster.Name)

				// 清除全局 GRM（携程对象）
				watchType := this.Cluster.Id + `_` + utils.NetworkPolicy
				delete(models.GRM.GoRoutineMap, watchType)
				logs.Info("Remove NetworkPolicyWatch from global GRM object, cluster: %s", this.Cluster.Name)

				// 清除数据库数据
				netpol := models.NetworkPolicy{}
				netpol.ClusterName = this.Cluster.Name
				netpol.Delete()

				// 重启 watch 携程
				deployService := NetworkPolicyService{Cluster: this.Cluster, ClientGo: this.ClientGo, Close: make(chan bool)}
				models.GRM.GoRoutineMap[watchType] = deployService

				logs.Info("Retry NetworkPolicyWatch, cluster: %s", this.Cluster.Name)
				go deployService.Wtach()
				break Retry
			}
		}
	}
}

func (this *NetworkPolicyService) Create() (*v1.NetworkPolicy, error) {
	netpol := v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "nginx"}},
			Ingress: []v1.NetworkPolicyIngressRule{
				v1.NetworkPolicyIngressRule{
					From: []v1.NetworkPolicyPeer{
						v1.NetworkPolicyPeer{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"access": "true"}}},
					},
				},
			},
		},
	}
	return this.ClientGo.ClientSet.NetworkingV1().NetworkPolicies(this.NetworkPolicy.NameSpaceName).Create(&netpol)
}

func (this *NetworkPolicyService) Delete() error {
	return this.ClientGo.ClientSet.NetworkingV1().NetworkPolicies(this.NetworkPolicy.NameSpaceName).Delete(this.NetworkPolicy.Name, nil)
}
