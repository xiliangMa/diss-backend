package task

import (
	"github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/k8s"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
)

type K8STaskHandler struct {
	Clientgo utils.ClientGo
}

func NewK8STaskHandler(path string) *K8STaskHandler {
	return &K8STaskHandler{
		Clientgo: utils.CreateK8sClient(path),
	}
}

func (this *K8STaskHandler) SyncHost(clusterId string) {
	nodes, err := this.Clientgo.GetNodes()
	if err != nil {
		logs.Error("Sync node err: %s", err.Error())
	} else {
		for _, n := range nodes.Items {
			name := n.ObjectMeta.Name
			// 同步 hostconfig
			nodeId := n.Status.NodeInfo.MachineID
			config := new(models.HostConfig)
			config.HostName = name
			config.Id = nodeId
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
			info.Disk = strconv.FormatInt(d/1024/1024/1024, 10)
			info.Id = nodeId
			nStatusNodeinfo := n.Status.NodeInfo
			info.OS = nStatusNodeinfo.OSImage
			info.Kernel = nStatusNodeinfo.KernelVersion
			info.Architecture = nStatusNodeinfo.Architecture
			info.DockerRuntime = nStatusNodeinfo.ContainerRuntimeVersion
			info.KubeletVer = nStatusNodeinfo.KubeletVersion
			info.Kubeproxy = nStatusNodeinfo.KubeProxyVersion
			info.KubernetesVer = nStatusNodeinfo.KubeletVersion
			info.Inner_AddHostInfo()

			// 同步主机images
			for _, o := range n.Status.Images {
				var imageName string
				imageId, _ := uuid.NewV4()
				image := new(models.ImageConfig)
				image.Id = imageId.String()
				image.HostId = nodeId
				if len(o.Names) == 1 {
					image.Name = o.Names[0]
				} else {
					for _, name := range o.Names {
						if !strings.Contains(name, "@sha256:") {
							imageName = imageName + name + ","
						}
					}
					image.Name = imageName
				}
				image.Size = string(o.SizeBytes)
				// to do image create time

				image.Add()
			}
		}
	}
}

func (this *K8STaskHandler) SyncHostInfo() {
	nodes, err := this.Clientgo.GetNodes()
	if err != nil {
		logs.Error("Sync node err: %s", err.Error())
	} else {
		for _, n := range nodes.Items {
			nodeId := n.Status.NodeInfo.MachineID
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
			ob.Id = nodeId
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
			nId := string(ns.UID)
			nsId := nId
			ob.Id = nsId
			ob.Name = nsName
			ob.ClusterId = clusterId
			ob.Add()

			//同步pod
			this.SyncPod(nsName)
		}
	}

}

func (this *K8STaskHandler) SyncPod(nsName string) {
	pods, err := this.Clientgo.GetPodsByNameSpace(nsName)
	if err != nil {
		logs.Error("Sync namespace: %s pods err: %s", nsName, err.Error())
	} else {
		for _, pod := range pods.Items {
			// 同步 pod
			podob := new(k8s.Pod)
			pId := string(pod.UID)
			podob.Id = pId
			podob.Name = pod.ObjectMeta.Name
			podob.NamSpaceName = nsName
			podob.HostIp = pod.Status.HostIP
			podob.HostName = pod.Spec.NodeName
			podob.PodIp = pod.Status.PodIP
			podob.Status = string(pod.Status.Phase)
			podob.Add()

			//同步 containerconfig
			for _, c := range pod.Status.ContainerStatuses {
				ccob := new(models.ContainerConfig)
				ccob.Id = strings.Replace(c.ContainerID, "docker://", "", 1)
				ccob.Name = c.Name
				ccob.PodId = pId
				ccob.PodName = pod.ObjectMeta.Name
				ccob.NameSpaceName = nsName
				//var commandArr string
				//for _, commad := range c.Command {
				//	commandArr = commandArr + commad + ","
				//}
				//ccob.Command = commandArr
				ccob.ImageName = c.Image
				ccob.HostName = pod.Spec.NodeName
				//ccob.Status = c.State.Running.StartedAt.String()
				ccob.UpdateTime = pod.Status.StartTime.Time
				ccob.Add()
			}

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
			if c.IsSync == models.Cluster_IsSync && c.Synced == models.Cluster_NoSynced {
				logs.Info("Sync start...., cluster name: %s ", c.Name)
				// 创建k8s客户端
				this := NewK8STaskHandler(c.FileName)

				// 同步 hostconfig & hostinfo & imageconfig
				this.SyncHost(c.Id)

				// 同步 ns &&  ns 下的 pod
				this.SyncNameSpace(c.Id)

				logs.Info("Sync end...., cluster name: %s ", c.Name)
				// 更新初始化状态
				c.Synced = models.Cluster_Synced
				c.Update()
				logs.Info("Update cluster: %s  Synced status: true ", c.Name)
			}
		}

	}
}
