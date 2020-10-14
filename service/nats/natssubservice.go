package nats

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/nats-io/stan.go"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/synccheck"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type NatsSubService struct {
	Message       []byte
	Conn          stan.Conn
	ClientSubject string
	DelTask       *models.Task
}

func (this *NatsSubService) Save() error {
	ms := models.WsData{}
	if err := json.Unmarshal(this.Message, &ms); err != nil {
		logs.Error("Paraces WsData error %s", err)
		return err
	}

	if ms.Data == nil {
		return nil
	}
	switch ms.Type {
	case models.Type_Metric: // 数据上报
		switch ms.Tag {
		case models.Resource_HeartBeat:
			// to do 可以暂时采用点对点的 websocket 上报
			//heartBeat := models.HeartBeat{}
			//s, _ := json.Marshal(ms.Data)
			//if err := json.Unmarshal(s, &heartBeat); err != nil {
			//	logs.Error("Paraces %s error %s", ms.Tag, err)
			//	return err
			//}
			//ip := strings.Split(this.Conn.RemoteAddr().String(), ":")
			//client := &models.Client{Hub: models.WSHub, Conn: this.Conn, Send: make(chan []byte, 256), ClientIp: ip[0], SystemId: heartBeat.SystemId}
			//client.Hub.Register <- client
			//logs.Info("Nats ############################ Agent Heater Beat data, >>> HostId: %s, Type: %s <<<", heartBeat.SystemId, models.Resource_HeartBeat)
			//metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			//this.ReceiveData(metricsResult)
		case models.Resource_CmdHistory_LatestTime:
			cmdHistory := models.CmdHistory{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &cmdHistory); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			msg := fmt.Sprintf("Nats ############################ Agent Fetch LatestTime data, >>> HostId: %s, Type: %s <<<", cmdHistory.HostId, models.Resource_CmdHistory_LatestTime)
			if cmdHistory.Type == models.Cmd_History_Type_Container {
				msg = fmt.Sprintf("Nats ############################ Agent Fetch LatestTime data, >>> HostId: %s, ContainerId: %s, Type: %s <<<", cmdHistory.HostId, cmdHistory.ContainerId, models.Resource_CmdHistory_LatestTime)
			}
			logs.Info(msg)
			result := cmdHistory.GetLatestTime()
			metricsResult := models.WsData{Code: result.Code, Type: models.Type_Control, Tag: models.Resource_CmdHistory_LatestTime, RCType: models.Resource_Control_Type_Get, Data: result.Data}
			this.ReceiveData(metricsResult)
		case models.Resource_DockerEvent_LatestTime:
			dockerEvent := models.DockerEvent{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &dockerEvent); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			msg := fmt.Sprintf("Nats ############################ Agent Fetch LatestTime data, >>> HostId: %s, Type: %s <<<", dockerEvent.HostId, models.Resource_DockerEvent_LatestTime)
			logs.Info(msg)
			result := dockerEvent.GetLatestTime()
			metricsResult := models.WsData{Code: result.Code, Type: models.Type_Control, Tag: models.Resource_DockerEvent_LatestTime, RCType: models.Resource_Control_Type_Get, Data: result.Data}
			this.ReceiveData(metricsResult)
		case models.Resource_HostInfoDynamic:
			// k8s 主机更新主机的外网ip
			hostInfo := models.HostInfo{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostInfo); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("Nats ############################ Sync agent data, >>> HostId: %s, Type: %s <<<", hostInfo.Id, models.Resource_HostInfoDynamic)
			if result := hostInfo.UpdateDynamic(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case models.Resource_HostConfigDynamic:
			// k8s 主机更新主机的外网ip
			hostConfig := models.HostConfig{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostConfig); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}

			logs.Info("Nats ############################ Sync agent data, >>> HostId: %s, Type: %s <<<", hostConfig.Id, models.Resource_HostConfigDynamic)
			if result := hostConfig.UpdateDynamic(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case models.Resource_HostConfig:
			hostConfig := models.HostConfig{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostConfig); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("Nats ############################ Sync agent data, >>> HostId: %s, Type: %s <<<", hostConfig.Id, models.Resource_HostConfig)
			if err := hostConfig.Add(); err != nil {
				return err
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case models.Resource_HostInfo:
			hostInfo := models.HostInfo{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostInfo); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", hostInfo.Id, models.Resource_HostInfo)
			if err := hostInfo.Add(); err != nil {
				return err
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case models.Resource_ContainerConfig:
			containerConfigList := []models.ContainerConfig{}
			CheckObject := new(models.ContainerConfig)
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &containerConfigList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(containerConfigList)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>>  HostName: %s, Type: %s, Size: %d <<<", containerConfigList[0].HostName, models.Resource_ContainerConfig, size)
			}
			for _, containerConfig := range containerConfigList {
				//if result := containerConfig.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				containerConfig.AccountName = models.Account_Admin
				containerConfig.Add()
			}
			// 清除脏数据
			if size != 0 {
				CheckObject.HostName = containerConfigList[0].HostName
				CheckObject.SyncCheckPoint = containerConfigList[0].SyncCheckPoint
				agentCheckHandler := synccheck.AgentCheckHadler{CheckObject, nil}
				agentCheckHandler.Check(models.Resource_ContainerConfig)
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_ContainerInfo:
			containerInfoList := []models.ContainerInfo{}
			CheckObject := new(models.ContainerInfo)
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &containerInfoList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(containerInfoList)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", containerInfoList[0].HostId, models.Resource_ContainerInfo, size)
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
				agentCheckHandler.Check(models.Resource_ContainerInfo)
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_ImageConfig:
			imageConfigList := []models.ImageConfig{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &imageConfigList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(imageConfigList)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", imageConfigList[0].HostId, models.Resource_ImageConfig, len(imageConfigList))
				//删除主机下的所有imageconfig
				imageConfigList[0].Delete()
			}

			for _, imageConfig := range imageConfigList {
				//if result := imageConfig.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				imageConfig.Add()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_ImageInfo:
			imageInfoList := []models.ImageInfo{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &imageInfoList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(imageInfoList)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", imageInfoList[0].HostId, models.Resource_ImageInfo, len(imageInfoList))
				//删除主机下的所有imageginfo
				imageInfoList[0].Delete()
			}
			for _, imageInfo := range imageInfoList {
				//if result := imageInfo.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				imageInfo.Add()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case models.Resource_HostPs:
			hostPsList := []models.HostPs{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostPsList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(hostPsList)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", hostPsList[0].HostId, models.Resource_HostPs, size)
				// 删除该主机下所有的进程
				hostPsList[0].Delete()

			}
			for _, hostPs := range hostPsList {
				//if result := hostPs.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				hostPs.Add()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_ContainerPs:
			containerPsList := []models.ContainerPs{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &containerPsList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(containerPsList)

			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", containerPsList[0].HostId, models.Resource_ContainerPs, len(containerPsList))
				containerPsList[0].Delete()
			}

			for _, containerTop := range containerPsList {
				//if result := containerTop.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				containerTop.Add()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_DockerBenchMark:
			//index := beego.AppConfig.String("security_log::BenchMarkIndex")
			benchMarkLog := models.BenchMarkLog{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &benchMarkLog); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", benchMarkLog.HostId, models.Resource_DockerBenchMark)
			if result := benchMarkLog.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
			// hostconfig填充最新的基线结果个数
			host := models.HostConfig{}
			host.Id = benchMarkLog.HostId
			hostConfig := host.Get()
			if hostConfig != nil {
				dockerCISCount := map[string]int{}
				dockerCISCount["NoteCount"] = benchMarkLog.NoteCount
				dockerCISCount["InfoCount"] = benchMarkLog.InfoCount
				dockerCISCount["WarnCount"] = benchMarkLog.WarnCount
				dockerCISCount["PassCount"] = benchMarkLog.PassCount
				benchCountJson, _ := json.Marshal(dockerCISCount)
				hostConfig.DockerCISCount = string(benchCountJson)

				hostConfig.Update()
			}

			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
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
		case models.Resource_KubernetesBenchMark:
			//index := beego.AppConfig.String("security_log::BenchMarkIndex")
			benchMarkLog := models.BenchMarkLog{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &benchMarkLog); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", benchMarkLog.HostId, models.Resource_KubernetesBenchMark)
			if result := benchMarkLog.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
			// hostconfig填充最新的基线结果个数
			host := models.HostConfig{}
			host.Id = benchMarkLog.HostId
			hostConfig := host.Get()
			if hostConfig != nil {
				kubeCISCount := map[string]int{}
				kubeCISCount["FailCount"] = benchMarkLog.FailCount
				kubeCISCount["InfoCount"] = benchMarkLog.InfoCount
				kubeCISCount["WarnCount"] = benchMarkLog.WarnCount
				kubeCISCount["PassCount"] = benchMarkLog.PassCount
				benchCountJson, _ := json.Marshal(kubeCISCount)
				hostConfig.KubeCISCount = string(benchCountJson)

				hostConfig.Update()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
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
		case models.Resource_HostCmdHistory:
			cmdHistoryList := models.CmdHistoryList{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &cmdHistoryList.List); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(cmdHistoryList.List)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>> HostId: %s, ype: %s, Size: %d <<<", cmdHistoryList.List[0].HostId, models.Resource_HostCmdHistory, size)
				// 删除该主机下所有的记录 type = 0
				//cmdHistoryList.List[0].Delete()
			}

			for _, cmdHistory := range cmdHistoryList.List {
				cmdHistory.Add()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_ContainerCmdHistory:
			cmdHistoryList := models.CmdHistoryList{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &cmdHistoryList.List); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(cmdHistoryList.List)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>> HostId: %s, Type: %s, Size: %d <<<", cmdHistoryList.List[0].HostId, models.Resource_ContainerCmdHistory, size)
				// 删除该主机下所有的记录 type = 1
				//cmdHistoryList.List[0].Delete()
			}
			for _, cmdHistory := range cmdHistoryList.List {
				cmdHistory.Add()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_DockerEvent:
			dockerEvent := models.DockerEvent{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &dockerEvent); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("Nats ############################ Sync agent data, >>> HostId: %s, Type: %s<<<", dockerEvent.HostId, models.Resource_DockerEvent)
			dockerEvent.Add()
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_WarningInfo:
			sysConfig := models.SysConfig{}
			sysConfig.Key = models.Resource_WarningInfo
			warningInfoConfig := sysConfig.Get()

			if warningInfoConfig != nil {
				// WarningInfo 配置项开启时执行数据的上报
				if warningInfoConfig.Value != "" {
					sysConfigJson := map[string]bool{}
					err := json.Unmarshal([]byte(warningInfoConfig.Value), &sysConfigJson)
					if err != nil {
						logs.Error("Parse %s Config error %s", models.Resource_WarningInfo, err)
						return err
					}
					if sysConfigJson[models.Enable] == true {
						warningInfo := models.WarningInfo{}
						s, _ := json.Marshal(ms.Data)
						if err := json.Unmarshal(s, &warningInfo); err != nil {
							logs.Error("Parse %s error %s", ms.Tag, err)
							return err
						}
						logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", warningInfo.HostId, models.Resource_WarningInfo+"-"+warningInfo.Type)
						if warningInfo.HostId != "" {
							if warningInfo.Type == models.WarningInfo_File || warningInfo.Type == models.WarningInfo_Other || warningInfo.Type == models.WarningInfo_Process || warningInfo.Type == models.WarningInfo_Container {
								hostParam := models.HostConfig{Id: warningInfo.Id}
								hostconfig := hostParam.Get()
								warningInfo.Cluster = hostconfig.ClusterName
								warningInfo.Account = "admin"
							}
						}

						if result := warningInfo.Add(); result.Code != http.StatusOK {
							return errors.New(result.Message)
						}
						metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
						this.ReceiveData(metricsResult)
					}
				}
			} else {
				// 建立配置项，并设置默认值 Enable: false
				warninginfoConfig := map[string]bool{models.Enable: false}
				warninginfoConfigJson, _ := json.Marshal(warninginfoConfig)
				sysConfig.Value = string(warninginfoConfigJson)
				sysConfig.Add()
			}
		}
	case models.Type_Control:
		switch ms.Tag {
		case models.Resource_Task:
			// 获取任务列表接口
			switch ms.RCType {
			case models.Resource_Control_Type_Get:
				metricsResult := models.WsData{Code: http.StatusOK, Type: models.Type_Control, Tag: models.Resource_Task, RCType: models.Resource_Control_Type_Get}
				task := models.Task{}
				s, _ := json.Marshal(ms.Data)
				if err := json.Unmarshal(s, &task); err != nil {
					logs.Error("Paraces: %s type: %s error: %s  ", ms.Tag, ms.RCType, err)
					return err
				}
				if result := task.GetUnFinishedTaskList(); result.Code != http.StatusOK {
					metricsResult.Code = result.Code
					metricsResult.Msg = result.Message
					logs.Error("Nats ############################ Get un finished task list  fail, >>> HostId: %s, error: <<<", task.Host.Id, result.Message)
					return errors.New(result.Message)
				} else {
					metricsResult.Code = result.Code
					metricsResult.Msg = result.Message
					metricsResult.Data = result.Data
					if result.Data != nil {
						data := result.Data.(map[string]interface{})
						total := data["total"]
						if total != 0 {
							logs.Info("Nats ############################  Get un finished task list, >>> HostId: %s, Type: %s, task size:  %v <<<", task.Host.Id, models.Resource_Task, total)
							this.ReceiveData(metricsResult)
						} else {
							logs.Info("Nats ############################  Get un finished task list, >>> HostId: %s, Type: %s, task size:  %v <<<", task.Host.Id, models.Resource_Task, 0)
						}
					}
				}
			case models.Resource_Control_Type_Put:
				//更新任务状态
				metricsResult := models.WsData{Code: http.StatusOK, Type: models.Type_Control, Tag: models.Resource_Task, RCType: models.Resource_Control_Type_Put}
				taskList := []models.Task{}
				s, _ := json.Marshal(ms.Data)
				if err := json.Unmarshal(s, &taskList); err != nil {
					logs.Error("Paraces: %s type: %s error: %s  ", ms.Tag, ms.RCType, err)
					return err
				}
				logTag := "Nats ############################ "
				for _, task := range taskList {
					metricsResult.Data = task
					if result := task.Update(); result.Code != http.StatusOK {
						metricsResult.Code = result.Code
						metricsResult.Msg = result.Message
						msg := ""
						if task.Host == nil {
							msg = fmt.Sprintf("Update task KStatus: %s, fail, >>> HostName: %s, task id: %s, error: %s <<<", task.Status, task.Container.HostName, task.Id, result.Message)
						} else {
							msg = fmt.Sprintf("Update task KStatus: %s, fail, >>> HostId: %s, task id: %s, error: %s <<<", task.Status, task.Host.Id, task.Id, result.Message)
						}
						logs.Error(logTag + msg)
						taskRawInfo, _ := json.Marshal(task)
						taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Account: task.Account, Level: models.Log_level_Error}
						taskLog.Add()
						return errors.New(result.Message)
					} else {
						msg := ""
						if task.Host == nil {
							msg = fmt.Sprintf("Update task Success, KStatus: %s >>> HostName: %s, Type: %s, task id: %s <<<", task.Status, task.Container.HostName, models.Resource_Task, task.Id)
						} else {
							msg = fmt.Sprintf("Update task Success, KStatus: %s >>> HostId: %s, Type: %s, task id: %s <<<", task.Status, task.Host.Id, models.Resource_Task, task.Id)
						}
						logs.Info(logTag + msg)
						taskRawInfo, _ := json.Marshal(task)
						taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Account: task.Account, Level: models.Log_level_Info}
						taskLog.Add()
					}
					this.ReceiveData(metricsResult)
				}
			case models.Resource_Control_Type_Delete:
				logTag := "Nats ############################ "
				metricsResult := models.WsData{Code: http.StatusOK, Type: models.Type_Control, Tag: models.Resource_Task, RCType: models.Resource_Control_Type_Delete}
				task := models.Task{}
				s, _ := json.Marshal(ms.Data)
				if err := json.Unmarshal(s, &task); err != nil {
					logs.Error("Paraces: %s type: %s error: %s  ", ms.Tag, ms.RCType, err)
					return err
				}
				if result := task.Delete(); result.Code != http.StatusOK {
					metricsResult.Code = result.Code
					metricsResult.Msg = result.Message
					msg := ""
					if task.Host == nil {
						msg = fmt.Sprintf("Delete task fail, >>> HostName: %s, , task id: %s, error: %s <<<", task.Container.HostName, task.Id, result.Message)
					} else {
						msg = fmt.Sprintf("Delete task fail, >>> HostId: %s, , task id: %s, error: %s <<<", task.Host.Id, task.Id, result.Message)
					}
					logs.Error(msg)
					taskRawInfo, _ := json.Marshal(task)
					taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Account: task.Account, Level: models.Log_level_Error}
					taskLog.Add()
					return errors.New(result.Message)
				} else {
					msg := ""
					if task.Host == nil {
						msg = fmt.Sprintf("Delete task success, >>> HostName: %s, Type: %s, task id: %s<<<", task.Container.HostName, models.Resource_Task, task.Id)
					} else {
						msg = fmt.Sprintf("Delete task success, >>> HostId: %s, Type: %s, task id: %s<<<", task.Host.Id, models.Resource_Task, task.Id)
					}
					logs.Info(logTag + msg)
					taskRawInfo, _ := json.Marshal(task)
					taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Account: task.Account, Level: models.Log_level_Info}
					taskLog.Add()
				}
				this.ReceiveData(metricsResult)
			}
		}
	}

	return nil
}

func (this *NatsSubService) DeleteTask() error {
	logs.Info("################ Delete Task <<<start>>> ################")
	task := this.DelTask
	result := models.WsData{Type: models.Type_Control, Tag: models.Resource_Task, RCType: models.Resource_Control_Type_Delete, Data: task}
	data, err := json.Marshal(result)

	// 下发删除任务
	subject := ""
	// nats
	if task.Host != nil {
		subject = task.Host.Id
	}
	if task.Container != nil {
		subject = task.Container.HostName
	}
	if subject != "" {
		err = models.Nats.Conn.Publish(subject, data)
		if err == nil {
			logs.Info("Deliver Task to Nats Success, Subject: %s Id: %s, RCType: %s, data: %v", subject, task.Id, models.Resource_Control_Type_Delete, result)
		}
	} else {
		return errors.New(string(utils.ResourceNotFoundErr))
	}
	if err != nil {
		msg := fmt.Sprintf("Delete Task Fail, Id: %s, err: %s", task.Id, err.Error())
		logs.Error(msg)
		taskRawInfo, _ := json.Marshal(task)
		taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Account: task.Account, Level: models.Log_level_Error}
		taskLog.Add()
		return err
	}
	logs.Info("################ Delete Task <<<end>>> ################")
	return nil
}

func (this *NatsSubService) ReceiveData(result models.WsData) {
	data, _ := json.Marshal(result)
	err := this.Conn.Publish(this.ClientSubject, data)
	if err != nil {
		logs.Error("Nats ############################ Received data from agent fail ############################", err)
	}
}

func RunClientSub(clientSubject string) {
	natsManager := models.Nats
	if natsManager != nil && natsManager.Conn != nil {
		serverSubject := models.Subject_Common + `-` + clientSubject
		natsManager.Conn.Subscribe(serverSubject, func(m *stan.Msg) {
			natsSubService := NatsSubService{Conn: natsManager.Conn, Message: m.Data, ClientSubject: clientSubject}
			natsSubService.Save()
		}, stan.DurableName(serverSubject))
	}
}

func RunClientSub_Image_Safe() {
	natsManager := models.Nats
	if natsManager != nil && natsManager.Conn != nil {
		registries := models.Registries{}
		registryList := registries.List()

		for _, registry := range registryList {
			libname := registry.Registry
			imageSafeSubject := libname + `_` + models.Subject_Image_Safe
			logs.Info("Subscribe image registry :", imageSafeSubject)
			natsManager.Conn.Subscribe(imageSafeSubject, func(m *stan.Msg) {
				natsSubService := NatsSubService{Conn: natsManager.Conn, Message: m.Data, ClientSubject: libname}
				natsSubService.Save()
			}, stan.DurableName(imageSafeSubject))
		}
	}
}

func RunClientSub_IDL() {
	natsManager := models.Nats
	if natsManager != nil && natsManager.Conn != nil {
		host := models.HostConfig{}
		hostList := host.List(0, 0)

		hostListData := hostList.Data.(map[string]interface{})["items"]

		if hostListData != nil {
			for _, host := range hostListData.([]*models.HostConfig) {
				hostid := host.Id
				IDLSubject := hostid + `_` + models.Subject_IntrudeDetect
				logs.Info("Subscribe Intrude Detect :", IDLSubject)
				natsManager.Conn.Subscribe(IDLSubject, func(m *stan.Msg) {
					natsSubService := NatsSubService{Conn: natsManager.Conn, Message: m.Data, ClientSubject: hostid}
					natsSubService.Save()
				}, stan.DurableName(IDLSubject))
			}
		}

	}
}
