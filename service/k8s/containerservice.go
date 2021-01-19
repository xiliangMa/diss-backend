package k8s

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/xiliangMa/diss-backend/models"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"strings"
	"time"
)

type ContainerService struct {
	Cluster *models.Cluster
	Pod     *v1.Pod
}

func (this *ContainerService) InitContainer(EventType watch.EventType) {
	switch EventType {
	case watch.Added:
		this.AddContainer()
	case watch.Modified:
		this.DeleteContainer()
		this.AddContainer()
	case watch.Deleted:
		this.DeleteContainer()
	case watch.Bookmark:
	}
}

func (this *ContainerService) DeleteContainer() {
	pod := this.Pod
	podId := string(pod.UID)
	cc := models.ContainerConfig{PodId: podId}
	ci := models.ContainerInfo{PodId: podId}
	cc.Delete()
	ci.Delete()
}

func (this *ContainerService) AddContainer() {
	clusterName := this.Cluster.Name
	nsName := this.Pod.Namespace
	pod := this.Pod
	// pod 相关数据
	podIp := pod.Status.PodIP
	hostName := pod.Spec.NodeName

	labelsByte, _ := json.Marshal(pod.Labels)
	labels := string(labelsByte)

	volumesByte, _ := json.Marshal(pod.Labels)
	volumes := string(volumesByte)

	podName := pod.ObjectMeta.Name

	logs.Info("########## Sync ContainerConfig && ContainerInfo, Pod: %s >>> strat <<< ##########, size: %d, NSName: %s", podName, len(pod.Spec.Containers), nsName)

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
		ciob.Add() // 添加containerinfo

		logs.Info("########## Sync ContainerConfig && ContainerInfo, Pod: %s >>> end <<< ##########, size: %d, NSName: %s", podName, len(pod.Spec.Containers), nsName)
	}
}
