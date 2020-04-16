package ws

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/bean"
	"github.com/xiliangMa/diss-backend/models/job"
	"github.com/xiliangMa/diss-backend/models/securitylog"
	"github.com/xiliangMa/diss-backend/models/ws"
	"github.com/xiliangMa/diss-backend/service/synccheck"
	"net/http"
	"strings"
)

type WSMetricsService struct {
	Message []byte
	Conn    *websocket.Conn
}

func (this *WSMetricsService) Save() error {
	ms := ws.MetricsResult{}
	if err := json.Unmarshal(this.Message, &ms); err != nil {
		logs.Error("Paraces MetricsResult error %s", err)
		return err
	}
	switch ms.ResType {
	case ws.Type_Metric: // 数据上报
		switch ms.ResTag {
		case ws.Tag_HeartBeat:
			heartBeat := bean.HeartBeat{}
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &heartBeat); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			ip := strings.Split(this.Conn.RemoteAddr().String(), ":")
			client := &Client{Hub: WSHub, Conn: this.Conn, Send: make(chan []byte, 256), ClientIp: ip[0], SystemId: heartBeat.SystemId}
			client.Hub.Register <- client
			logs.Info("############################ Agent Heater Beat data, >>> HostId: %s, Type: %s <<<", heartBeat.SystemId, ws.Tag_HeartBeat)
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case ws.Tag_HostConfig:
			hostConfig := models.HostConfig{}
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &hostConfig); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			logs.Info("############################ Sync agent data, >>> HostId: %s, Type: %s <<<", hostConfig.Id, ws.Tag_HostConfig)
			if err := hostConfig.Inner_AddHostConfig(); err != nil {
				return err
			}
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case ws.Tag_HostInfo:
			hostInfo := models.HostInfo{}
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &hostInfo); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", hostInfo.Id, ws.Tag_HostInfo)
			if err := hostInfo.Inner_AddHostInfo(); err != nil {
				return err
			}
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case ws.Tag_ContainerConfig:
			containerConfigList := []models.ContainerConfig{}
			CheckObject := new(models.ContainerConfig)
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &containerConfigList); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			size := len(containerConfigList)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>>  HostName: %s, Type: %s, Size: %d <<<", containerConfigList[0].HostName, ws.Tag_ContainerConfig, size)
			}
			for _, containerConfig := range containerConfigList {
				//if result := containerConfig.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				containerConfig.Add()
			}
			// 清除脏数据
			if size != 0 {
				CheckObject.HostName = containerConfigList[0].HostName
				CheckObject.SyncCheckPoint = containerConfigList[0].SyncCheckPoint
				agentCheckHandler := synccheck.AgentCheckHadler{CheckObject, nil}
				agentCheckHandler.Check(ws.Tag_ContainerConfig)
			}
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case ws.Tag_ContainerInfo:
			containerInfoList := []models.ContainerInfo{}
			CheckObject := new(models.ContainerInfo)
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &containerInfoList); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			size := len(containerInfoList)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", containerInfoList[0].HostId, ws.Tag_ContainerInfo, size)
			}
			for _, containerInfo := range containerInfoList {
				//if result := containerInfo.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				containerInfo.Add()
			}

			// 清除脏数据
			if size != 0 {
				CheckObject.HostName = containerInfoList[0].HostName
				CheckObject.SyncCheckPoint = containerInfoList[0].SyncCheckPoint
				agentCheckHandler := synccheck.AgentCheckHadler{nil, CheckObject}
				agentCheckHandler.Check(ws.Tag_ContainerInfo)
			}
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case ws.Tag_ImageConfig:
			imageConfigList := []models.ImageConfig{}
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &imageConfigList); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			size := len(imageConfigList)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", imageConfigList[0].HostId, ws.Tag_ImageConfig, len(imageConfigList))
				//删除主机下的所有imageconfig
				imageConfigList[0].Delete()
			}

			for _, imageConfig := range imageConfigList {
				//if result := imageConfig.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				imageConfig.Add()
			}
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case ws.Tag_ImageInfo:
			imageInfoList := []models.ImageInfo{}
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &imageInfoList); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			size := len(imageInfoList)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", imageInfoList[0].HostId, ws.Tag_ImageInfo, len(imageInfoList))
				//删除主机下的所有imageginfo
				imageInfoList[0].Delete()
			}
			for _, imageInfo := range imageInfoList {
				//if result := imageInfo.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				imageInfo.Add()
			}
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case ws.Tag_HostPs:
			hostPsList := []models.HostPs{}
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &hostPsList); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			size := len(hostPsList)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", hostPsList[0].HostId, ws.Tag_HostPs, size)
				// 删除该主机下所有的进程
				hostPsList[0].Delete()

			}
			for _, hostPs := range hostPsList {
				//if result := hostPs.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				hostPs.Add()
			}
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case ws.Tag_ContainerPs:
			containerPsList := []models.ContainerPs{}
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &containerPsList); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			size := len(containerPsList)

			if size != 0 {
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", containerPsList[0].HostId, ws.Tag_ContainerPs, len(containerPsList))
				containerPsList[0].Delete()
			}

			for _, containerTop := range containerPsList {
				//if result := containerTop.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				containerTop.Add()
			}
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case ws.Tag_DockerBenchMarkLog:
			//index := beego.AppConfig.String("security_log::BenchMarkIndex")
			benchMarkLog := securitylog.BenchMarkLog{}
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &benchMarkLog); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", benchMarkLog.HostId, ws.Tag_DockerBenchMarkLog)
			if result := benchMarkLog.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
			// 上报es
			//esClient, err := utils.GetESClient()
			//if err != nil {
			//	logs.Error("Get ESClient err: %s", err)
			//}
			//
			//respones, err := esClient.Create(index, benchMarkLog.Id, bytes.NewReader(s))
			//if err != nil {
			//	logs.Error("Add security_log to es fail, benchMarkLog.Id: %s", benchMarkLog.Id)
			//} else {
			//	logs.Info("Add security_log to es success, benchMarkLog.Id: %s", benchMarkLog.Id)
			//}
			//defer respones.Body.Close()
		case ws.Tag_KubernetesBenchMarkLog:
			//index := beego.AppConfig.String("security_log::BenchMarkIndex")
			benchMarkLog := securitylog.BenchMarkLog{}
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &benchMarkLog); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", benchMarkLog.HostId, ws.Tag_KubernetesBenchMarkLog)
			if result := benchMarkLog.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
			// 上报es
			//esClient, err := utils.GetESClient()
			//if err != nil {
			//	logs.Error("Get ESClient err: %s", err)
			//}
			//respones, err := esClient.Create(index, benchMarkLog.Id, bytes.NewReader(s))
			//if err != nil {
			//	logs.Error("Add security_log to es fail, benchMarkLog.Id: %s", benchMarkLog.Id)
			//} else {
			//	logs.Info("Add security_log to es success, benchMarkLog.Id: %s", benchMarkLog.Id)
			//}
			//defer respones.Body.Close()
		case ws.Tag_HostCmdHistory:
			cmdHistoryList := models.CmdHistoryList{}
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &cmdHistoryList.List); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			size := len(cmdHistoryList.List)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>> HostId: %s, ype: %s, Size: %d <<<", cmdHistoryList.List[0].HostId, ws.Tag_HostCmdHistory, size)
				// 删除该主机下所有的记录 type = 0
				cmdHistoryList.List[0].Delete()
			}

			for _, cmdHistory := range cmdHistoryList.List {
				//if result := cmdHistory.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				cmdHistory.Add()
			}
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case ws.Tag_ContainerCmdHistory:
			cmdHistoryList := models.CmdHistoryList{}
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &cmdHistoryList.List); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			size := len(cmdHistoryList.List)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>> HostId: %s, Type: %s, Size: %d <<<", cmdHistoryList.List[0].HostId, ws.Tag_ContainerCmdHistory, size)
				// 删除该主机下所有的记录 type = 1
				cmdHistoryList.List[0].Delete()
			}
			for _, cmdHistory := range cmdHistoryList.List {
				//if result := cmdHistory.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				cmdHistory.Add()
			}
			metricsResult := ws.MetricsResult{ResType: ws.Type_ReceiveState, ResTag: ws.Tag_Received, Metric: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		}
	case ws.Type_SyncData: // 数据下发
		switch ms.ResTag {
		case ws.Tag_Task:
			task := job.Task{}
			s, _ := json.Marshal(ms.Metric)
			if err := json.Unmarshal(s, &task); err != nil {
				logs.Error("Paraces %s error %s", ms.ResTag, err)
				return err
			}
			if result := task.GetUnFinishedTaskList(); result.Code != http.StatusOK {
				logs.Info("############################ Get un finished task list  fail, >>> HostId: %s, error: <<<", task.Host.Id, result.Message)
				return errors.New(result.Message)
			} else {
				metricsResult := ws.MetricsResult{ResType: ws.Type_Response, ResTag: ws.Tag_Task, Metric: result.Data, Config: ""}
				this.ReceiveData(metricsResult)
				jsonStr, _ := json.Marshal(result.Data)
				logs.Info("############################  Get un finished task list, >>> HostId: %s, Type: %s, task data:  %v <<<", task.Host.Id, ws.Tag_Task, string(jsonStr))
			}
		}
	}

	return nil
}

func (this *WSMetricsService) ReceiveData(result ws.MetricsResult) {
	data, _ := json.Marshal(result)
	err := this.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logs.Info("############################ Received data from agent fail ############################", err)
	}
}
