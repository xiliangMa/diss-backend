package k8s

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	v1 "k8s.io/api/networking/v1"
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
	return this.ClientGo.ClientSet.NetworkingV1().NetworkPolicies("").List(context.Background(), metav1.ListOptions{})
}

func (this *NetworkPolicyService) Wtach() {
	ns := ""
	if this.NetworkPolicy != nil {
		ns = this.NetworkPolicy.NameSpaceName
	}
	netpolWatch, err := this.ClientGo.ClientSet.NetworkingV1().NetworkPolicies(ns).Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		logs.Error("Wtach networkPolicy error: %s  ", err)
		return
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
				if event.Type == watch.Error {
					break
				}
				networkPolicy := v1.NetworkPolicy{}
				if !reflect.DeepEqual(event.Object.GetObjectKind(), networkPolicy.GetObjectKind()) {
					break
				}
				object := event.Object.(*v1.NetworkPolicy)
				id := string(object.UID)
				name := object.Name
				ns := object.Namespace
				clusterId := this.Cluster.Id
				clusterName := this.Cluster.Name

				KMetaData, _ := json.Marshal(object.ObjectMeta)
				KSpec, _ := json.Marshal(object.Spec)

				netpol := new(models.NetworkPolicy)
				netpol.Id = id
				netpol.Name = name
				netpol.ClusterName = clusterName
				netpol.AccountName = models.Cluster_Data_AccountName
				netpol.NameSpaceName = ns
				netpol.ClusterId = clusterId
				netpol.KMetaData = string(KMetaData)
				netpol.KSpec = string(KSpec)

				netpol.IngressPolicy = models.Network_Policy_Type_Value_Allow
				netpol.EgressPolicy = models.Network_Policy_Type_Value_Allow
				if object.Spec.Ingress == nil {
					netpol.IngressPolicy = models.Network_Policy_Type_Value_Refuse
				} else {
					if object.Spec.Ingress[0].From == nil && object.Spec.Ingress[0].Ports == nil {
						netpol.IngressPolicy = models.Network_Policy_Type_Value_AllowAll
					}
				}
				if object.Spec.Egress == nil {
					netpol.EgressPolicy = models.Network_Policy_Type_Value_Refuse
				} else {
					if object.Spec.Egress[0].To == nil && object.Spec.Egress[0].Ports == nil {
						netpol.EgressPolicy = models.Network_Policy_Type_Value_AllowAll
					}
				}

				logs.Info("Watch >>> NetworkPolicy: %s <<<, >>> ClusterOBJ: %s <<<,  >>> NameSpace: %s <<<, >>> EventType: %s <<<", id, clusterId, ns, event.Type)
				switch event.Type {
				case watch.Added:
					netpol.Add()
				case watch.Modified:
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

	// ============ ObjectMeta ============
	objectMeta := metav1.ObjectMeta{}
	err := json.Unmarshal([]byte(this.NetworkPolicy.KMetaData), &objectMeta)
	if err != nil {
		logs.Error("Paraces: %s, error: %s  ", models.Kubernetes_Object_MetaData, err)
		return nil, err
	}

	// ============ Spec ============
	spec := v1.NetworkPolicySpec{}
	err = json.Unmarshal([]byte(this.NetworkPolicy.KSpec), &spec)
	if err != nil {
		logs.Error("Paraces: %s, error: %s  ", models.Kubernetes_Object_Spec, err)
		return nil, err
	}

	netpol := v1.NetworkPolicy{ObjectMeta: objectMeta, Spec: spec}
	netpol.ObjectMeta.ResourceVersion = ""

	return this.ClientGo.ClientSet.NetworkingV1().NetworkPolicies(this.NetworkPolicy.NameSpaceName).Create(context.Background(), &netpol, metav1.CreateOptions{})
}

func (this *NetworkPolicyService) Update() (*v1.NetworkPolicy, error) {

	// ============ ObjectMeta ============
	objectMeta := metav1.ObjectMeta{}
	err := json.Unmarshal([]byte(this.NetworkPolicy.KMetaData), &objectMeta)
	if err != nil {
		logs.Error("Paraces: %s, error: %s  ", models.Kubernetes_Object_MetaData, err)
		return nil, err
	}

	// ============ Spec ============
	spec := v1.NetworkPolicySpec{}
	err = json.Unmarshal([]byte(this.NetworkPolicy.KSpec), &spec)
	if err != nil {
		logs.Error("Paraces: %s, error: %s  ", models.Kubernetes_Object_Spec, err)
		return nil, err
	}

	netpol := v1.NetworkPolicy{ObjectMeta: objectMeta, Spec: spec}
	netpol.ObjectMeta.ResourceVersion = ""

	return this.ClientGo.ClientSet.NetworkingV1().NetworkPolicies(this.NetworkPolicy.NameSpaceName).Update(context.Background(), &netpol, metav1.UpdateOptions{})
}

func (this *NetworkPolicyService) Delete() error {
	return this.ClientGo.ClientSet.NetworkingV1().NetworkPolicies(this.NetworkPolicy.NameSpaceName).Delete(context.Background(), this.NetworkPolicy.Name, metav1.DeleteOptions{})
}
