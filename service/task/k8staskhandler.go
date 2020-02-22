package task

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/k8s"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
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

func (this *K8STaskHandler) SyncHostConfig(clusterName, clusterId string) {
	nodes, err := this.Clientgo.GetNodes()
	if err != nil {
		logs.Error("Sync node err: %s", err.Error())
	} else {
		logs.Info("Sync cluster: %s HostConfig, size: %d", clusterName, len(nodes.Items))
		for _, n := range nodes.Items {
			name := n.ObjectMeta.Name
			// 同步 hostconfig
			nodeId := n.Status.NodeInfo.SystemUUID
			config := new(models.HostConfig)
			config.HostName = name
			config.Id = nodeId
			config.OS = n.Status.NodeInfo.OSImage
			config.IsInK8s = true
			config.ClusterId = clusterId
			config.Inner_AddHostConfig()
		}
	}
}

func (this *K8STaskHandler) SyncHostInfo(clusterName string) {
	nodes, err := this.Clientgo.GetNodes()
	if err != nil {
		logs.Error("Sync node err: %s", err.Error())
	} else {
		logs.Info("Sync cluster: %s HostInfo, size: %d", clusterName, len(nodes.Items))
		for _, n := range nodes.Items {
			name := n.ObjectMeta.Name
			nodeId := n.Status.NodeInfo.SystemUUID
			// 同步 hostinfo
			info := new(models.HostInfo)
			info.HostName = name
			info.Id = nodeId
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
			info.Disk = utils.UnitConvert(d)
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

func (this *K8STaskHandler) SyncHostImageConfig() {
	nodes, err := this.Clientgo.GetNodes()
	if err != nil {
		logs.Error("Sync node err: %s", err.Error())
	} else {
		for _, n := range nodes.Items {
			logs.Info("Sync host: %s ImageConfig, size: %d", n.Name, len(n.Status.Images))
			//同步主机images
			nodeId := n.Status.NodeInfo.SystemUUID
			for _, o := range n.Status.Images {
				for _, name := range o.Names {
					image := new(models.ImageConfig)
					image.HostId = nodeId
					image.Size = utils.UnitConvert(o.SizeBytes)
					// to do image create time
					if !strings.Contains(name, "@sha256:") {
						image.Name = name
						image.ImageId = nodeId + "---" + image.Name
						image.Id = nodeId + "---" + image.Name
						image.Add()
					}
				}
			}
		}
	}
}

func (this *K8STaskHandler) SyncNameSpace(clusetrName, clusterId string) {
	nameSpaces, err := this.Clientgo.GetNameSpaces()
	if err != nil {
		logs.Error("Sync namspace err: %s", err.Error())
	} else {
		logs.Info("Sync Cluster: %s NameSpace, size: %d", clusetrName, len(nameSpaces.Items))
		for _, ns := range nameSpaces.Items {
			nsName := ns.ObjectMeta.Name
			ob := new(k8s.NameSpace)
			nId := string(ns.UID)
			nsId := nId
			ob.Id = nsId
			ob.Name = nsName
			ob.ClusterId = clusterId
			ob.Add()
		}
	}
}

func (this *K8STaskHandler) SyncNamespacePod() {
	nameSpaces, err := this.Clientgo.GetNameSpaces()
	if err != nil {
		logs.Error("Sync namspace err: %s", err.Error())
	} else {
		for _, ns := range nameSpaces.Items {
			nsName := ns.Name
			pods, err := this.Clientgo.GetPodsByNameSpace(nsName)
			if err != nil {
				logs.Error("Sync namespace: %s pods err: %s", nsName, err.Error())
			} else {
				logs.Info("Sync NameSpace: %s Pod, size: %d, NSName %s", nsName, len(pods.Items), nsName)
				for _, pod := range pods.Items {
					// 同步 pod
					podob := new(k8s.Pod)
					podob.Id = string(pod.UID)
					podob.Name = pod.ObjectMeta.Name
					podob.NameSpaceName = nsName
					podob.HostIp = pod.Status.HostIP
					podob.HostName = pod.Spec.NodeName
					podob.PodIp = pod.Status.PodIP
					podob.Status = string(pod.Status.Phase)
					podob.Add()
				}
			}
		}

	}
}

func (this *K8STaskHandler) SyncPodContainerConfig() {
	nameSpaces, err := this.Clientgo.GetNameSpaces()
	if err != nil {
		logs.Error("Sync namspace err: %s", err.Error())
	} else {
		for _, ns := range nameSpaces.Items {
			nsName := ns.Name
			pods, err := this.Clientgo.GetPodsByNameSpace(nsName)
			if err != nil {
				logs.Error("Sync namespace: %s pods err: %s", nsName, err.Error())
			} else {
				for _, pod := range pods.Items {
					logs.Info("Sync Pod: %s ContainerConfig, size: %d, NSName: %s", pod.ObjectMeta.Name, len(pods.Items), nsName)
					//同步 containerconfig
					for _, c := range pod.Status.ContainerStatuses {
						ccob := new(models.ContainerConfig)
						ccob.Id = strings.Replace(c.ContainerID, "docker://", "", 1)
						ccob.Name = c.Name
						ccob.PodId = string(pod.UID)
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
	}
}

type Data struct {
	items []*k8s.Cluster
	total int
}

func SyncAll() {
	// cluster
	var cluster k8s.Cluster
	result := cluster.List(0, 100000)

	if result.Code == http.StatusOK && result.Data != nil {
		data := result.Data.(map[string]interface{})
		for _, c := range data["items"].([]*k8s.Cluster) {
			if c.IsSync == models.Cluster_IsSync && c.Synced == models.Cluster_NoSync {
				logs.Info("########################################## cluster:  %s, Sync start.", c.Name)
				// 创建k8s客户端
				this := NewK8STaskHandler(c.FileName)

				// 同步 namespace
				this.SyncNameSpace(c.Name, c.Id)

				// 同步 hostconfig
				this.SyncHostConfig(c.Name, c.Id)

				// 同步 hostinfo
				this.SyncHostInfo(c.Name)

				// 同步HostImageConfig
				this.SyncHostImageConfig()

				// 同步 namespace 下的 pod
				this.SyncNamespacePod()

				// 同步 pod 下的 container
				this.SyncPodContainerConfig()

				logs.Info("Sync end...., cluster name: %s ", c.Name)
				// 更新初始化状态
				c.Synced = models.Cluster_Synced
				c.Update()
				logs.Info("########################################## cluster:  %s, Sync end.", c.Name)
			}
		}

	}
}
