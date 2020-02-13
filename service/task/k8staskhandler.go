package task

import (
	"github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/k8s"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
)

type K8STaskHandler struct {
	Clientgo utils.ClientGo
}

func NewK8STaskHandler(path string) *K8STaskHandler {
	return &K8STaskHandler{
		Clientgo: utils.CreateK8sClient(path),
	}
}

func (this *K8STaskHandler) SyncCluster() {
}

func (this *K8STaskHandler) SyncHostConfig(clusterId string) {
	nodes, err := this.Clientgo.GetNodes()
	if err != nil {
		logs.Error("Sync node err: %s", err.Error())
	} else {
		for _, n := range nodes.Items {
			name := n.ObjectMeta.Name
			// 同步 hostconfig
			uid, _ := uuid.NewV4()
			config := new(models.HostConfig)
			config.HostName = name
			config.Id = uid.String()
			config.OS = n.Status.NodeInfo.OSImage
			config.IsInK8s = true
			config.ClusterId = clusterId
			config.Inner_AddHostConfig()

			// 同步 hostinfo
			info := new(models.HostInfo)
			info.HostName = name
			if n.Status.Addresses[0].Type == "InternalIP" {
				info.InternalAddr = n.Status.Addresses[0].Address
			} else {
				info.InternalAddr = n.Status.Addresses[1].Address
			}
			capacity := n.Status.Capacity
			c, _ := capacity.Cpu().AsInt64()
			info.CpuCore = c
			m, _ := capacity.Memory().AsInt64()
			info.Mem = m / 1024 / 1024 / 1024
			d, _ := capacity.StorageEphemeral().AsInt64()
			info.Disk = strconv.FormatInt(d / 1024 / 1024 / 1024, 10)
			info.Id = uid.String()
			nStatusNodeinfo := n.Status.NodeInfo
			info.OS = nStatusNodeinfo.OSImage
			info.Kernel = nStatusNodeinfo.KernelVersion
			info.Architecture = nStatusNodeinfo.Architecture
			info.DockerRuntime = nStatusNodeinfo.ContainerRuntimeVersion
			info.KubeletVer = nStatusNodeinfo.KubeletVersion
			info.Kubeproxy = nStatusNodeinfo.KubeProxyVersion
			info.KubernetesVer = nStatusNodeinfo.KubeletVersion
			info.Inner_AddHostInfo()
		}
	}
}

func (this *K8STaskHandler) SyncHostInfo() {
	nodes, err := this.Clientgo.GetNodes()
	if err != nil {
		logs.Error("Sync node err: %s", err.Error())
	} else {
		for _, n := range nodes.Items {
			uid, _ := uuid.NewV4()
			ob := new(models.HostInfo)
			ob.HostName = n.ObjectMeta.Name
			if n.Status.Addresses[0].Type == "InternalIP" {
				ob.InternalAddr = n.Status.Addresses[0].Address
			} else {
				ob.InternalAddr = n.Status.Addresses[1].Address
			}
			c, _ := n.Status.Capacity.Cpu().AsInt64()
			ob.CpuCore = c
			m, _ := n.Status.Capacity.Memory().AsInt64()
			ob.Mem = m
			ob.Id = uid.String()
			nStatusNodeinfo := n.Status.NodeInfo
			ob.OS = nStatusNodeinfo.OSImage
			ob.Kernel = nStatusNodeinfo.KernelVersion
			ob.Architecture = nStatusNodeinfo.Architecture
			ob.DockerRuntime = nStatusNodeinfo.ContainerRuntimeVersion
			ob.KubeletVer = nStatusNodeinfo.KubeletVersion
			ob.Kubeproxy = nStatusNodeinfo.KubeProxyVersion
			ob.KubernetesVer = nStatusNodeinfo.KubeletVersion
			ob.Inner_AddHostInfo()
		}
	}
}

func (this *K8STaskHandler) SyncNameSpace(clusterId string) {
	nameSpaces, err := this.Clientgo.GetNameSpaces()
	if err != nil {
		logs.Error("Sync namspace err: %s", err.Error())
	} else {
		for _, ns := range nameSpaces.Items {
			nsName := ns.ObjectMeta.Name
			ob := new(k8s.NameSpace)
			NID, _ := uuid.NewV4()
			nsId := NID.String()
			ob.Id = nsId
			ob.Name = nsName
			ob.ClusterId = clusterId
			ob.Add()
			this.SyncPod(nsName, nsId, clusterId)
		}
	}

}

func (this *K8STaskHandler) SyncPod(nsName, clusterId, nsId string) {
	pods, err := this.Clientgo.GetPodsByNameSpace(nsName)
	if err != nil {
		logs.Error("Sync namespace: %s pods err: %s", nsName, err.Error())
	} else {
		for _, pod := range pods.Items {
			ob := new(k8s.Pod)
			NID, _ := uuid.NewV4()
			ob.Id = NID.String()
			ob.Name = pod.ObjectMeta.Name
			ob.ClusterId = clusterId
			ob.NameSpaceId = nsId
			ob.Add()
		}
	}
}

type Data struct {
	items []*k8s.Cluster
	total int
}

func SyncAll() {
	// cluster
	var cluster k8s.Cluster
	result := cluster.List()

	if result.Code == http.StatusOK && result.Data != nil {
		data := result.Data.(map[string]interface{})
		for _, c := range data["items"].([]*k8s.Cluster) {
			if c.IsSync == k8s.Cluster_IsSync && c.Synced == k8s.Cluster_NoSynced {
				logs.Info("Sync start...., cluster name: %s ", c.Name)
				// 创建k8s客户端
				this := NewK8STaskHandler(c.FileName)

				// 同步主机
				this.SyncHostConfig(c.Id)

				// 同步 ns &&  ns 下的 pod
				this.SyncNameSpace(c.Id)
				logs.Info("Sync end...., cluster name: %s ", c.Name)
				// 更新初始化状态
				c.Synced = k8s.Cluster_Synced
				c.Update()
				logs.Info("Update cluster: %s  Synced status: true ", c.Name)
			}
		}

	}
}
