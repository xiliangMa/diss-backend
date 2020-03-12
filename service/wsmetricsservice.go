package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/securitylog"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type WSMetricsService struct {
	Message []byte
}

func (wsmh *WSMetricsService) Save() error {
	ms := models.MetricsResult{}
	if err := json.Unmarshal(wsmh.Message, &ms); err != nil {
		logs.Error("Paraces MetricsResult error %s", err)
		return err
	}
	switch ms.ResTag {
	case models.Tag_HostConfig:
		hostConfig := models.HostConfig{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &hostConfig); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		logs.Info("############################ Sync agent data, >>> HostId: %s, Type: %s <<<", hostConfig.Id, models.Tag_HostConfig)
		if err := hostConfig.Inner_AddHostConfig(); err != nil {
			return err
		}
	case models.Tag_HostInfo:
		hostInfo := models.HostInfo{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &hostInfo); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", hostInfo.Id, models.Tag_HostInfo)
		if err := hostInfo.Inner_AddHostInfo(); err != nil {
			return err
		}
	case models.Tag_ContainerConfig:
		containerConfigList := []models.ContainerConfig{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &containerConfigList); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		size := len(containerConfigList)
		if size != 0 {
			logs.Info("############################ Sync agent data, >>>  HostName: %s, Type: %s, Size: %d <<<", containerConfigList[0].HostName, models.Tag_ContainerConfig, len(containerConfigList))
		}
		for _, containerConfig := range containerConfigList {
			if result := containerConfig.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
		}
	case models.Tag_ContainerInfo:
		containerInfoList := []models.ContainerInfo{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &containerInfoList); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		size := len(containerInfoList)
		if size != 0 {
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", containerInfoList[0].HostId, models.Tag_ContainerInfo, len(containerInfoList))
		}
		for _, containerInfo := range containerInfoList {
			if result := containerInfo.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
		}
	case models.Tag_ImageConfig:
		imageConfigList := []models.ImageConfig{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &imageConfigList); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		size := len(imageConfigList)
		if size != 0 {
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", imageConfigList[0].HostId, models.Tag_ImageConfig, len(imageConfigList))
			//删除主机下的所有imageconfig
			//imageConfigList[0].Delete()
		}

		for _, imageConfig := range imageConfigList {
			if result := imageConfig.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
		}
	case models.Tag_ImageInfo:
		imageInfoList := []models.ImageInfo{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &imageInfoList); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		size := len(imageInfoList)
		if size != 0 {
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", imageInfoList[0].HostId, models.Tag_ImageInfo, len(imageInfoList))
			//删除主机下的所有imageginfo
			imageInfoList[0].Delete()
		}
		for _, imageInfo := range imageInfoList {
			if result := imageInfo.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
		}
	case models.Tag_HostPs:
		hostPsList := []models.HostPs{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &hostPsList); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		size := len(hostPsList)
		if size != 0 {
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", hostPsList[0].HostId, models.Tag_HostPs, size)
			// 删除该主机下所有的进程
			hostPsList[0].Delete()

		}
		for _, hostPs := range hostPsList {
			if result := hostPs.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
		}
	case models.Tag_ContainerPs:
		containerPsList := []models.ContainerPs{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &containerPsList); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		size := len(containerPsList)

		if size != 0 {
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", containerPsList[0].HostId, models.Tag_ContainerPs, len(containerPsList))
			containerPsList[0].Delete()
		}

		for _, containerTop := range containerPsList {
			if result := containerTop.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
		}
	case models.Tag_DockerBenchMarkLog:
		index := beego.AppConfig.String("security_log::BenchMarkIndex")
		benchMarkLog := securitylog.BenchMarkLog{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &benchMarkLog); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", benchMarkLog.HostId, models.Tag_DockerBenchMarkLog)
		if result := benchMarkLog.Add(); result.Code != http.StatusOK {
			return errors.New(result.Message)
		}

		// 上报es
		esClient, err := utils.GetESClient()
		if err != nil {
			logs.Error("Get ESClient err: %s", err)
		}

		respones, err := esClient.Create(index, benchMarkLog.Id, bytes.NewReader(s))
		if err != nil {
			logs.Error("Add security_log to es fail, benchMarkLog.Id: %s", benchMarkLog.Id)
		} else {
			logs.Info("Add security_log to es success, benchMarkLog.Id: %s", benchMarkLog.Id)
		}
		defer respones.Body.Close()
	case models.Tag_KubernetesBenchMarkLog:
		index := beego.AppConfig.String("security_log::BenchMarkIndex")
		benchMarkLog := securitylog.BenchMarkLog{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &benchMarkLog); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", benchMarkLog.HostId, models.Tag_KubernetesBenchMarkLog)
		if result := benchMarkLog.Add(); result.Code != http.StatusOK {
			return errors.New(result.Message)
		}
		// 上报es
		esClient, err := utils.GetESClient()
		if err != nil {
			logs.Error("Get ESClient err: %s", err)
		}
		respones, err := esClient.Create(index, benchMarkLog.Id, bytes.NewReader(s))
		if err != nil {
			logs.Error("Add security_log to es fail, benchMarkLog.Id: %s", benchMarkLog.Id)
		} else {
			logs.Info("Add security_log to es success, benchMarkLog.Id: %s", benchMarkLog.Id)
		}
		defer respones.Body.Close()
	case models.Tag_HostCmdHistory:
		cmdHistoryList := models.CmdHistoryList{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &cmdHistoryList.List); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		size := len(cmdHistoryList.List)
		if size != 0 {
			logs.Info("############################ Sync agent data, >>> HostId: %s, ype: %s, Size: %d <<<", cmdHistoryList.List[0].HostId, models.Tag_HostCmdHistory, size)
			// 删除该主机下所有的记录 type = 0
			cmdHistoryList.List[0].Delete()
		}

		for _, cmdHistory := range cmdHistoryList.List {
			if result := cmdHistory.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
		}
	case models.Tag_ContainerCmdHistory:
		cmdHistoryList := models.CmdHistoryList{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &cmdHistoryList.List); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		size := len(cmdHistoryList.List)
		if size != 0 {
			logs.Info("############################ Sync agent data, >>> HostId: %s, Type: %s, Size: %d <<<", cmdHistoryList.List[0].HostId, models.Tag_ContainerCmdHistory, size)
			// 删除该主机下所有的记录 type = 1
			cmdHistoryList.List[0].Delete()
		}
		for _, cmdHistory := range cmdHistoryList.List {
			if result := cmdHistory.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
		}
	}

	return nil
}
