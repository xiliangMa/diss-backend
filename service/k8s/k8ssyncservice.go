package k8s

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/synccheck"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strings"
	"time"
)

type K8sSyncService struct {
	Cluster        *models.Cluster
	SyncCheckPoint int64
	Clientgo       ClientGo
}

func NewK8sSyncService(syncCheckPoint int64, cluster *models.Cluster) *K8sSyncService {
	return &K8sSyncService{
		SyncCheckPoint: syncCheckPoint,
		Cluster:        cluster,
	}
}

func (this *K8sSyncService) Sync() models.Result {
	if result := this.check(); result.Code != http.StatusOK {
		return result
	}
	return this.SyncCluster()
}

func (this *K8sSyncService) CreateK8sClient() ClientGo {
	// 构建k8s客户端
	params := new(ApiParams)
	params.AuthType = this.Cluster.AuthType
	params.BearerToken = this.Cluster.BearerToken
	params.MasterUrl = this.Cluster.MasterUrls
	params.KubeConfigPath = this.Cluster.FileName
	return CreateK8sClient(params)
}

func (this *K8sSyncService) check() models.Result {
	var resultData models.Result

	if this.Cluster == nil {
		resultData.Code = utils.ClusterNoExistErr
		resultData.Message = string(utils.ClusterNoExistErr)
		return resultData
	}

	data := this.Cluster.List(0, 0).Data.(map[string]interface{})
	total := data["total"]
	clusterList := data["items"]
	if data == nil || total == 0 {
		resultData.Code = utils.ClusterNoExistErr
		resultData.Message = string(utils.ClusterNoExistErr)
		return resultData
	} else {
		this.Cluster = clusterList.([]*models.Cluster)[0]
	}

	if this.Cluster.Status != models.Cluster_Status_Active {
		resultData.Code = utils.ClusterNotActiveOrNotvailableErr
		resultData.Message = string(utils.ClusterNotActiveOrNotvailableErr)
		return resultData
	}

	if this.Cluster.SyncStatus == models.Cluster_Sync_Status_Clearing || this.Cluster.SyncStatus == models.Cluster_Sync_Status_InProcess {
		resultData.Code = utils.ClusterIsSyncingOrIsClearingErr
		resultData.Message = string(utils.ClusterIsSyncingOrIsClearingErr)
		return resultData
	}

	resultData.Code = http.StatusOK
	return resultData
}

func (this *K8sSyncService) SyncCluster() models.Result {
	var resultData models.Result
	clusterName := this.Cluster.Name
	clusterId := this.Cluster.Id
	this.Clientgo = this.CreateK8sClient()

	if this.Clientgo.ErrMessage == "" {
		this.Cluster.SyncStatus = models.Cluster_Sync_Status_InProcess
		this.Cluster.Update()
		logs.Info("########################################## cluster:  %s, Sync start.", clusterName)
		defer func() {
			if err := recover(); err != nil {
				// 更新集群的同步状态
				this.Cluster.SyncStatus = models.Cluster_Sync_Status_Fail
				this.Cluster.Update()
				logs.Error("########################################## cluster:  %s id: %s , Sync fail. err: %s", clusterName, clusterId, err)
			}
		}()

		this.SyncCheckPoint = this.SyncCheckPoint
		// 同步 namespace
		this.SyncNameSpace()

		// 同步 集群内的 hostconfig && hostInfo
		this.SyncHostConfigAndInfo()

		// 单独同步 hostinfo
		//this.SyncHostInfo(c.Name)

		// 同步HostImageConfig（无法通过k8s采集镜像的详细信息 imageconfig & imageinfo 均由agent采集）
		//this.SyncHostImageConfig()

		// 同步 namespace 下的 pod
		this.SyncNamespacePod()

		// 同步 pod 内的 containerconfig && containerinfo
		this.SyncPodContainerConfigAndInfo()

		logs.Info("Sync end...., cluster name: %s ", clusterName)
		// 更新同步时间、状态
		this.Cluster.SyncStatus = models.Cluster_Sync_Status_Synced
		this.Cluster.Update()
		logs.Info("########################################## cluster:  %s, Sync end.", clusterName)
	} else {
		// 更新集群的同步状态
		this.Cluster.SyncStatus = models.Cluster_Sync_Status_Fail
		this.Cluster.Update()
		logs.Error("########################################## cluster:  %s, Sync fail. err: %s ", clusterName, this.Clientgo.ErrMessage)
		return resultData
	}

	resultData.Code = http.StatusOK
	return resultData
}

func (this *K8sSyncService) SyncHostConfigAndInfo() {
	nodeService := NodeService{NodeInterface: this.Clientgo.ClientSet.CoreV1().Nodes()}
	nodes, err := nodeService.List()
	clusterName := this.Cluster.Name
	clusterId := this.Cluster.Id

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
			config.ClusterName = clusterName
			config.IsInK8s = true
			config.Diss = models.Diss_NotInstalled
			config.DissStatus = models.Diss_Status_Unsafe
			config.Status = models.Host_Status_Normal

			// 同步 hostinfo
			info := new(models.HostInfo)
			info.HostName = name
			info.Id = nodeId
			if n.Status.Addresses[0].Type == "InternalIP" {
				info.InternalAddr = n.Status.Addresses[0].Address
				config.InternalAddr = n.Status.Addresses[0].Address
			} else {
				info.InternalAddr = n.Status.Addresses[1].Address
				config.InternalAddr = n.Status.Addresses[1].Address
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
			info.DockerStatus = models.Host_Docker_Status_Nornal
			info.ClusterId = clusterId
			info.ClusterName = clusterName

			// 初始化原始数据
			metaData, _ := json.Marshal(n.ObjectMeta)
			config.KMetaData = string(metaData)
			info.KMetaData = string(metaData)
			spec, _ := json.Marshal(n.Spec)
			config.KSpec = string(spec)
			info.KSpec = string(spec)
			status, _ := json.Marshal(n.Status)
			config.KStatus = string(status)
			info.KStatus = string(status)

			config.Add() // 添加 hostconfig
			info.Add()   // 添加 hostinfo

		}
		logs.Info("########## Sync HostConfig && HostInfo, cluster: %s >>> end <<< ##########, size: %d", clusterName, len(nodes.Items))
	}
}

func (this *K8sSyncService) SyncHostImageConfig() {
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

func (this *K8sSyncService) SyncNameSpace() {
	clusterName := this.Cluster.Name
	clusterId := this.Cluster.Id
	nameSpaceService := NameSpaceService{NameSpaceInterface: this.Clientgo.ClientSet.CoreV1().Namespaces()}
	nameSpaces, err := nameSpaceService.List()
	if err != nil {
		logs.Error("Sync namspace err: %s", err.Error())
	} else {
		CheckObject := new(models.NameSpace)
		CheckObject.SyncCheckPoint = this.SyncCheckPoint
		CheckObject.ClusterId = clusterId
		logs.Info("########## Sync NameSpace, Cluster: %s >>> strat <<< ##########, size: %d", clusterName, len(nameSpaces.Items))
		for _, o := range nameSpaces.Items {
			nsName := o.ObjectMeta.Name
			ns := new(models.NameSpace)
			nId := string(o.UID)
			nsId := nId
			ns.Id = nsId
			ns.Name = nsName
			ns.ClusterName = clusterName
			ns.AccountName = models.Account_Admin
			ns.ClusterId = clusterId
			ns.SyncCheckPoint = this.SyncCheckPoint
			// 初始化原始数据
			metaData, _ := json.Marshal(o.ObjectMeta)
			ns.KMetaData = string(metaData)
			spec, _ := json.Marshal(o.Spec)
			ns.KSpec = string(spec)
			status, _ := json.Marshal(o.Status)
			ns.KStatus = string(status)

			ns.Add(true)
		}
		// 清除脏数据
		size := len(nameSpaces.Items)
		if size != 0 {
			k8sCheckHandler := synccheck.K8SCheckHadler{nil, nil, CheckObject, nil}
			k8sCheckHandler.Check(models.Resource_NameSpace)
			logs.Info("########## Empty Dirty Data, Model: %s ##########", models.Resource_NameSpace)
		}
		logs.Info("########## Sync NameSpace, Cluster: %s >>> end <<< ##########, size: %d", clusterName, len(nameSpaces.Items))
	}
}

func (this *K8sSyncService) SyncNamespacePod() {
	clusterName := this.Cluster.Name
	nameSpaceService := NameSpaceService{NameSpaceInterface: this.Clientgo.ClientSet.CoreV1().Namespaces()}
	nameSpaces, err := nameSpaceService.List()
	if err != nil {
		logs.Error("Sync namspace err: %s", err.Error())
	} else {
		CheckObject := new(models.Pod)
		CheckObject.SyncCheckPoint = this.SyncCheckPoint
		CheckObject.ClusterName = clusterName
		for _, ns := range nameSpaces.Items {
			nsName := ns.Name
			podService := PodService{PodInterface: this.Clientgo.ClientSet.CoreV1().Pods(nsName)}
			pods, err := podService.List()
			if err != nil {
				logs.Error("Sync namespace: %s pods err: %s", nsName, err.Error())
			} else {
				logs.Info("########## Sync Pod, NameSpace: %s >>> strat <<< ##########, size: %d, NSName %s", nsName, len(pods.Items), nsName)
				for _, o := range pods.Items {
					// 同步 pod
					pod := new(models.Pod)
					pod.Id = string(o.UID)
					pod.Name = o.ObjectMeta.Name
					pod.NameSpaceName = nsName
					pod.HostName = o.Spec.NodeName
					pod.ClusterName = clusterName
					pod.SyncCheckPoint = this.SyncCheckPoint
					// 初始化原始数据
					metaData, _ := json.Marshal(ns.ObjectMeta)
					pod.KMetaData = string(metaData)
					spec, _ := json.Marshal(ns.Spec)
					pod.KSpec = string(spec)
					status, _ := json.Marshal(ns.Status)
					pod.KStatus = string(status)

					pod.Add()
				}
				logs.Info("########## Sync Pod, NameSpace: %s >>> end <<< ##########, size: %d, NSName %s", nsName, len(pods.Items), nsName)
			}
		}
		// 清除脏数据
		size := len(nameSpaces.Items)
		if size != 0 {
			k8sCheckHandler := synccheck.K8SCheckHadler{nil, nil, nil, CheckObject}
			k8sCheckHandler.Check(models.Resource_Pod)
			logs.Info("########## Empty Dirty Data, Model: %s ##########", models.Resource_Pod)
		}
	}
}

func (this *K8sSyncService) SyncPodContainerConfigAndInfo() {
	clusterName := this.Cluster.Name
	nameSpaceService := NameSpaceService{NameSpaceInterface: this.Clientgo.ClientSet.CoreV1().Namespaces()}
	nameSpaces, err := nameSpaceService.List()
	if err != nil {
		logs.Error("Sync namspace err: %s", err.Error())
	} else {
		CheckObject1 := new(models.ContainerConfig)
		CheckObject1.SyncCheckPoint = this.SyncCheckPoint
		CheckObject1.ClusterName = clusterName

		CheckObject2 := new(models.ContainerInfo)
		CheckObject2.SyncCheckPoint = this.SyncCheckPoint
		CheckObject2.ClusterName = clusterName

		for _, ns := range nameSpaces.Items {
			nsName := ns.Name
			podService := PodService{PodInterface: this.Clientgo.ClientSet.CoreV1().Pods(nsName)}
			pods, err := podService.List()
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
						ccob.SyncCheckPoint = this.SyncCheckPoint
						ccob.AccountName = models.Account_Admin

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
						ciob.ClusterName = clusterName
						ciob.SyncCheckPoint = this.SyncCheckPoint

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
		// 清除脏数据
		size := len(nameSpaces.Items)
		if size != 0 {
			k8sCheckHandler := synccheck.K8SCheckHadler{CheckObject1, CheckObject2, nil, nil}

			k8sCheckHandler.Check(models.Resource_ContainerConfig)
			logs.Info("########## Empty Dirty Data, Model: %s ##########", models.Resource_ContainerConfig)

			k8sCheckHandler.Check(models.Resource_ContainerInfo)
			logs.Info("########## Empty Dirty Data, Model: %s ##########", models.Resource_ContainerInfo)
		}
	}
}
