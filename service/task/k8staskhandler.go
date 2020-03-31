package task

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/k8s"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strings"
	"time"
)

type K8STaskHandler struct {
	Clientgo utils.ClientGo
}

func NewK8STaskHandler(path string) *K8STaskHandler {
	return &K8STaskHandler{
		Clientgo: utils.CreateK8sClient(path),
	}
}

func (this *K8STaskHandler) SyncHostConfigAndInfo(clusterName, clusterId string) {
	nodes, err := this.Clientgo.GetNodes()
	if err != nil {
		logs.Error("Sync node err: %s", err.Error())
	} else {
		logs.Info("########## Sync HostConfig && HostInfo, cluster: %s >>> strat <<< ##########, size: %d", clusterName, len(nodes.Items))
		for _, n := range nodes.Items {
			name := n.ObjectMeta.Name
			// 同步 hostconfig
			nodeId := strings.ToLower(n.Status.NodeInfo.SystemUUID)
			config := new(models.HostConfig)
			config.HostName = name
			config.Id = nodeId
			config.OS = n.Status.NodeInfo.OSImage
			config.IsInK8s = true
			config.ClusterId = clusterId
			config.IsInK8s = true

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
			info.Mem = fmt.Sprintf("%.2f", float64(m)/1024/1024/1024)
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
			info.DockerStatus = "Running"
			config.Inner_AddHostConfig() // 添加 hostconfig
			info.Inner_AddHostInfo()     // 添加 hostinfo

		}
		logs.Info("########## Sync HostConfig && HostInfo, cluster: %s >>> end <<< ##########, size: %d", clusterName, len(nodes.Items))
	}
}

func (this *K8STaskHandler) SyncHostInfo(clusterName string) {
	nodes, err := this.Clientgo.GetNodes()
	if err != nil {
		logs.Error("Sync node err: %s", err.Error())
	} else {
		logs.Info("########## Sync cluster: %s HostInfo >>> strat <<< ##########, size: %d", clusterName, len(nodes.Items))
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
			info.Mem = fmt.Sprintf("%.2f", m/1024/1024/1024)
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
		logs.Info("########## Sync cluster: %s HostInfo >>> end <<< ##########, size: %d", clusterName, len(nodes.Items))
	}
}

func (this *K8STaskHandler) SyncHostImageConfig() {
	nodes, err := this.Clientgo.GetNodes()
	if err != nil {
		logs.Error("Sync node err: %s", err.Error())
	} else {
		for _, n := range nodes.Items {
			logs.Info("########## Sync ImageConfig, Host: %s >>> strat <<< ########## , size: %d", n.Name, len(n.Status.Images))
			//同步主机images
			nodeId := n.Status.NodeInfo.SystemUUID
			for _, o := range n.Status.Images {
				for _, name := range o.Names {
					// 同步imageinfo
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
			logs.Info("########## Sync ImageConfig, Host: %s  >>> end <<< ########## , size: %d", n.Name, len(n.Status.Images))
		}
	}
}

func (this *K8STaskHandler) SyncNameSpace(clusetrName, clusterId string) {
	nameSpaces, err := this.Clientgo.GetNameSpaces()
	if err != nil {
		logs.Error("Sync namspace err: %s", err.Error())
	} else {
		logs.Info("########## Sync NameSpace, Cluster: %s >>> strat <<< ##########, size: %d", clusetrName, len(nameSpaces.Items))
		for _, ns := range nameSpaces.Items {
			nsName := ns.ObjectMeta.Name
			ob := new(k8s.NameSpace)
			nId := string(ns.UID)
			nsId := nId
			ob.Id = nsId
			ob.Name = nsName
			ob.ClusterId = clusterId
			ob.Add(true)
		}
		logs.Info("########## Sync NameSpace, Cluster: %s >>> end <<< ##########, size: %d", clusetrName, len(nameSpaces.Items))
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
				logs.Info("########## Sync Pod, NameSpace: %s >>> strat <<< ##########, size: %d, NSName %s", nsName, len(pods.Items), nsName)
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
				logs.Info("########## Sync Pod, NameSpace: %s >>> end <<< ##########, size: %d, NSName %s", nsName, len(pods.Items), nsName)
			}
		}

	}
}

func (this *K8STaskHandler) SyncPodContainerConfigAndInfo(clusterName string) {
	nameSpaces, err := this.Clientgo.GetNameSpaces()
	SyncCheckPoint := time.Now().Unix()
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
					// pod 相关数据
					podIp := pod.Status.PodIP
					hostName := pod.Spec.NodeName

					labelsByte, _ := json.Marshal(pod.Labels)
					labels := string(labelsByte)

					volumesByte, _ := json.Marshal(pod.Labels)
					volumes := string(volumesByte)

					podName := pod.ObjectMeta.Name

					logs.Info("########## Sync ContainerConfig && ContainerInfo, Pod: %s >>> strat <<< ##########, size: %d, NSName: %s", podName, len(pods.Items), nsName)

					for _, c := range pod.Status.ContainerStatuses {
						//公用变量
						cid := strings.Replace(c.ContainerID, "docker://", "", 1)
						cname := c.Name
						podId := string(pod.UID)
						imageId := c.ImageID
						imageName := c.Image
						hostName := hostName
						startTime := pod.Status.StartTime.String()

						//计算运行时间
						now := time.Now()
						createTime, _ := time.Parse(time.RFC3339Nano, startTime)
						created := now.Sub(createTime)

						//动态的回去容器状态
						status := models.Pod_Container_Statue_Running
						if c.State.Terminated != nil {
							status = models.Pod_Container_Statue_Terminated
						}
						if c.State.Waiting != nil {
							status = models.Pod_Container_Statue_Waiting
						}

						//同步 containerconfig
						ccob := new(models.ContainerConfig)
						ccob.Id = cid
						ccob.Name = cname
						ccob.PodId = podId
						ccob.PodName = podName
						ccob.NameSpaceName = nsName
						ccob.ImageName = imageName
						ccob.HostName = hostName
						ccob.Status = status
						ccob.ClusterName = clusterName
						ccob.Age = "Up " + created.String()
						ccob.CreateTime = startTime
						ccob.UpdateTime = startTime
						ccob.SyncCheckPoint = SyncCheckPoint

						//同步 containerinfo
						ciob := new(models.ContainerInfo)
						ciob.Id = cid
						ciob.Name = cname
						ciob.HostName = pod.Spec.NodeName
						ciob.NameSpaceName = nsName
						ciob.PodId = podId
						ciob.PodName = podName
						ciob.ImageId = imageId
						ciob.ImageName = imageName
						//ciob.HostId = ""
						ciob.HostName = hostName
						ciob.StartedAt = createTime.String()
						ciob.CreatedAt = createTime.String()
						ciob.Status = status
						ciob.Ip = podIp
						ciob.Labels = labels
						ciob.Volumes = volumes
						ciob.SyncCheckPoint = SyncCheckPoint

						// 通过 containers获取的数据
						for _, cs := range pod.Spec.Containers {
							if cs.Name == cname {
								commandByte, _ := json.Marshal(cs.Command)
								command := string(commandByte)
								ccob.Command = command

								portsByte, _ := json.Marshal(cs.Ports)
								ports := string(portsByte)
								ciob.Ports = ports

								volumeMountsByte, _ := json.Marshal(cs.VolumeMounts)
								volumeMounts := string(volumeMountsByte)
								ciob.Mounts = volumeMounts
								ciob.Command = command
							}
						}

						ccob.Add() // 添加 containerconfig
						ciob.Add() // 添加contain儿info

					}

					logs.Info("########## Sync ContainerConfig && ContainerInfo, Pod: %s >>> end <<< ##########, size: %d, NSName: %s", podName, len(pods.Items), nsName)

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
	result := cluster.List(0, 0)

	if result.Code == http.StatusOK && result.Data != nil {
		data := result.Data.(map[string]interface{})
		for _, c := range data["items"].([]*k8s.Cluster) {
			if c.IsSync == models.Cluster_IsSync {
				clusterName := c.Name
				c.SyncStatus = models.Cluster_Sync_Status_IN_PROGRESS
				c.Update()
				logs.Info("########################################## cluster:  %s, Sync start.", c.Name)
				// 创建k8s客户端
				this := NewK8STaskHandler(c.FileName)

				// 同步 namespace
				this.SyncNameSpace(c.Name, c.Id)

				// 同步 集群内的 hostconfig && hostInfo
				this.SyncHostConfigAndInfo(c.Name, c.Id)

				// 单独同步 hostinfo
				//this.SyncHostInfo(c.Name)

				// 同步HostImageConfig（无法通过k8s采集镜像的详细信息 imageconfig & imageinfo 均由agent采集）
				//this.SyncHostImageConfig()

				// 同步 namespace 下的 pod
				this.SyncNamespacePod()

				// 同步 pod 内的 containerconfig && containerinfo
				this.SyncPodContainerConfigAndInfo(clusterName)

				logs.Info("Sync end...., cluster name: %s ", c.Name)
				// 更新同步时间、状态
				c.SyncStatus = models.Cluster_Sync_Status_OK
				c.Update()
				logs.Info("########################################## cluster:  %s, Sync end.", c.Name)
			}
		}

	}
}
