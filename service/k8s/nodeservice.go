package k8s

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"strings"
)

type NodeService struct {
	NodeInterface corev1.NodeInterface
	Cluster       *models.Cluster
	Close         chan bool
}

func (this *NodeService) List() (*v1.NodeList, error) {
	return this.NodeInterface.List(metav1.ListOptions{})
}

func (this *NodeService) Wtach() {
	nodeWatch, err := this.NodeInterface.Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// 开启 watch 事件
Retry:
	for {
		select {
		case <-this.Close:
			logs.Info("Close nodeWatch, cluster: %s", this.Cluster.Name)
			nodeWatch.Stop()
			return
		case event, ok := <-nodeWatch.ResultChan():
			if event.Object != nil || ok {
				object := event.Object.(*v1.Node)
				id := strings.ToLower(object.Status.NodeInfo.SystemUUID)
				name := object.ObjectMeta.Name
				os := object.Status.NodeInfo.OSImage
				clusterId := this.Cluster.Id
				clusterName := this.Cluster.Name

				KMetaData, _ := json.Marshal(object.ObjectMeta)
				KSpec, _ := json.Marshal(object.Spec)
				KStatus, _ := json.Marshal(object.Status)

				// HostConfig
				hc := new(models.HostConfig)
				hc.HostName = name
				hc.Id = id
				hc.OS = os
				hc.IsInK8s = true
				hc.ClusterId = clusterId
				hc.ClusterName = clusterName
				hc.IsInK8s = true
				hc.Diss = models.Diss_NotInstalled
				hc.DissStatus = models.Diss_Status_Unsafe
				hc.Status = models.Host_Status_Normal

				// HostInfo
				hi := new(models.HostInfo)
				hi.HostName = name
				hi.Id = id
				if object.Status.Addresses[0].Type == "InternalIP" {
					hi.InternalAddr = object.Status.Addresses[0].Address
					hc.InternalAddr = object.Status.Addresses[0].Address
				} else {
					hi.InternalAddr = object.Status.Addresses[1].Address
					hc.InternalAddr = object.Status.Addresses[1].Address
				}
				capacity := object.Status.Capacity
				c, _ := capacity.Cpu().AsInt64()
				hi.CpuCore = c
				m, _ := capacity.Memory().AsInt64()
				hi.Mem = fmt.Sprintf("%.2f", float64(m)/1024/1024/1024)
				d, _ := capacity.StorageEphemeral().AsInt64()
				hi.Disk = utils.UnitConvert(d)
				nStatusNodeinfo := object.Status.NodeInfo
				hi.OS = os
				hi.Kernel = nStatusNodeinfo.KernelVersion
				hi.Architecture = nStatusNodeinfo.Architecture
				hi.DockerRuntime = nStatusNodeinfo.ContainerRuntimeVersion
				hi.KubeletVer = nStatusNodeinfo.KubeletVersion
				hi.Kubeproxy = nStatusNodeinfo.KubeProxyVersion
				hi.KubernetesVer = nStatusNodeinfo.KubeletVersion
				hi.DockerStatus = models.Host_Docker_Status_Nornal
				hi.ClusterId = clusterId
				hi.ClusterName = clusterName

				// 初始化原始数据
				hc.KMetaData = string(KMetaData)
				hi.KMetaData = string(KMetaData)
				hc.KSpec = string(KSpec)
				hi.KSpec = string(KSpec)
				hc.KStatus = string(KStatus)
				hi.KStatus = string(KStatus)

				logs.Info("Watch >>> Node: %s <<<, >>> Cluster: %s <<<, >>> EventType: %s <<<", id, clusterId, event.Type)
				switch event.Type {
				case watch.Added:
					hc.Add()
					hi.Add()
				case watch.Modified:
					hc.Add()
					hi.Add()
				case watch.Deleted:
					hc.Delete()
					hi.Delete()
				case watch.Bookmark:
					//todo
				case watch.Error:
					//todo
				}
			} else {
				// 如果 watch 异常退回重新 watch
				logs.Warn("nodeWatch chan has been close!!!!, cluster: %s", this.Cluster.Name)

				watchType := this.Cluster.Id + `_` + utils.Host
				delete(models.GRM.GoRoutineMap, watchType)
				logs.Info("Remove nodeWatch from global GRM object, cluster: %s", this.Cluster.Name)

				k8sWatchService := K8sWatchService{Cluster: this.Cluster}
				clientGo := k8sWatchService.CreateK8sClient()

				nodeService := NodeService{Cluster: this.Cluster, NodeInterface: clientGo.ClientSet.CoreV1().Nodes(), Close: make(chan bool)}
				models.GRM.GoRoutineMap[watchType] = nodeService

				logs.Info("Retry node watch, cluster: %s", this.Cluster.Name)
				go nodeService.Wtach()
				break Retry
			}
		}
	}
}
