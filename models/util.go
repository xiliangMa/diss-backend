package models

import "time"

var (
	// Global
	All          = "All"
	Result_Items = "items"
	Result_Total = "total"

	// k8s
	Cluster_Sync_Status_NotSynced             = "NotSynced"
	Cluster_Sync_Status_Synced                = "Synced"
	Cluster_Sync_Status_InProcess             = "InProcess"
	Cluster_Sync_Status_Clearing              = "Clearing"
	Cluster_Sync_Status_Fail                  = "Fail"
	Cluster_Watch_Status_Fail                 = "Fail"
	Cluster_Watch_Status_Success              = "Success"
	Cluster_IsSync                            = true
	Cluster_NoSync                            = false
	Cluster_Status_Active                     = "Active"
	Cluster_Status_Unavailable                = "Unavailable"
	Cluster_Type_Kubernets                    = "Kubernets"
	Cluster_Type_OpenShift                    = "OpenShift"
	Cluster_Type_Rancher                      = "Rancher"
	Pod_Container_Statue_Running              = "Running"
	Pod_Container_Statue_Terminated           = "Terminated"
	Pod_Container_Statue_Waiting              = "Waiting"
	Api_Auth_Type_KubeConfig                  = "KubeConfig"
	Api_Auth_Type_BearerToken                 = "BearerToken"
	Cluster_Data_AccountName                  = "Kubernetes"
	Clster_Node_Label_Control_Rancher         = "node-role.kubernetes.io/controlplane"
	Clster_Node_Label_Master                  = "node-role.kubernetes.io/master"
	Clster_Node_Label_Worker                  = "node-role.kubernetes.io/worker"
	Clster_Node_Roler_All                     = "All"
	Clster_Node_Roler_Worker                  = "Worker"
	Clster_Node_Roler_Master                  = "Master"
	Network_Policy_Type_Value_Allow           = "Allow"
	Network_Policy_Type_Value_AllowAll        = "AllowAll"
	Network_Policy_Type_Value_Refuse          = "Refuse"
	Network_Policy_Type_Ingress               = "Ingress"
	Network_Policy_Type_Egress                = "Egress"
	Cluster_Scope_InternalUrl                 = "InternalUrl"
	Cluster_Scope_PublicUrl                   = "PublicUrl"
	Cluster_Scope_UrlPort                     = "32666"
	Cluster_Scope_Operator_Status_ActiveFail  = "ActiveFail"
	Cluster_Scope_Operator_Status_Actived     = "Actived"
	Cluster_Scope_Operator_Status_Activing    = "Activing"
	Cluster_Scope_Operator_Status_DisableFail = "DisableFail"
	Cluster_Scope_Operator_Status_Disabled    = "Disabled"
	Cluster_Scope_Operator_Status_Disableing  = "Disabling"
	Cluster_Scope_Operator_Status_Null        = ""

	Cmd_History_Type_Host      = "Host"
	Cmd_History_Type_Container = "Container"

	//容器状态
	Container_Status_Run   = "Run"
	Container_Status_Pause = "Pause"
	Container_Status_All   = "All"

	// 主机相关状态
	Host_Docker_Status_Nornal   = "Normal"
	Host_Docker_Status_Abnormal = "Abnormal"

	Host_Type_Server = "Server"
	Host_Type_Vm     = "Vm"

	Host_Status_Normal   = "Normal"
	Host_Status_Abnormal = "Abnormal"

	// 主机-安全容器状态
	Diss_Installed    = "Installed"
	Diss_NotInstalled = "NotInstalled"

	// 主机-安全状态
	Diss_status_Safe   = "Safe"
	Diss_Status_Unsafe = "Unsafe"

	//// 系统魔板类型(此处和ws resource tag 保持一致，和授权模块关联)
	// 镜像扫描：
	TMP_Type_ImageScan         = "ImageScan"
	TMP_Type_HostImageVulnScan = "HostImageVulnScan"
	//基线扫描：
	TMP_Type_BM_Docker = "DockerBenchMark"
	TMP_Type_BM_K8S    = "KubernetesBenchMark"
	//入侵扫描
	TMP_Type_IDS_Docker = "DockerIntrudeDetectScan"
	TMP_Type_IDS_Host   = "HostIntrudeDetectScan"
	//安全审计
	TMP_Type_DockerSecurityAudit     = "DockerSecurityAudit"
	TMP_Type_KubernetesSecurityAudit = "KubernetesSecurityAudit"
	TMP_Type_CommandSecurityAudit    = "CommandSecurityAudit"
	//病毒扫描：
	TMP_Type_DockerVS = "DockerVirusScan"
	TMP_Type_HostVS   = "HostVirusScan"
	//漏洞扫描
	TMP_Type_LS = "SC_LeakScan"

	TMP_Status_Enable  = "Enable"
	TMP_Status_Disable = "Disable"

	Job_Status_Active     = "Active"
	Job_Status_Deactiving = "Deactiving"
	Job_Status_Deactived  = "Deactived"

	//安全检查类型
	SC_Type_Host      = "host"
	Sc_Type_Container = "container"
	Sc_Type_Image     = "image"

	// 任务状态
	Task_Status_Pending        = "Pending"
	Task_Status_Running        = "Running"
	Task_Status_Removing       = "Removing"
	Task_Status_Pause          = "Pause"
	Task_Status_Finished       = "Finished"
	Task_Status_Failed         = "Failed"
	Task_Status_Deliver_Failed = "DeliverFailed"
	Task_Status_Received       = "Received"
	Task_Status_Receive_Failed = "ReceiveFailed"
	Task_Status_Unavailable    = "Unavailable"

	//任务操作
	Task_Action_Deactive = "Deactive"

	//任务类型
	Job_Type_Once     = "Once"
	Job_Type_Periodic = "Periodical"

	//任务级别
	Job_Level_System = "System" // 系统级
	Job_Level_User   = "User"   // 用户级

	//基线日志类型
	BMLT_Host_All    = "host"
	BMLT_Docker      = "docker"
	BMLT_K8s         = "k8s"
	BML_Template_ALL = "All"
	BML_Level_ALL    = "All"
	BML_Level_High   = "High"
	BML_Level_Medium = "Medium"
	BML_Level_Low    = "Low"
	BML_Result_ALL   = "All"
	BML_Result_Pass  = "Pass"
	BML_Result_Fail  = "Fail"

	// 入侵检测类型
	IDLT_Docker = "container"
	IDLT_Host   = "host"

	// 分组类型
	Group_Type_Host      = "Host"
	Group_Type_Container = "Container"

	// 租户
	Account_Admin = "admin"

	// nats subject (diss-backend 主题， agent 则以主机id区分标识)
	Subject_Common        = "Common"
	Subject_Image_Safe    = "IMAGE_SAFE"
	Subject_IntrudeDetect = "INTRUDE_DETECT"

	// 告警信息类型
	WarningInfo_File      = "ALERT_TYPE_FILE"
	WarningInfo_Other     = "ALERT_TYPE_OTHER"
	WarningInfo_Process   = "ALERT_TYPE_PROCESS"
	WarningInfo_Container = "ALERT_TYPE_CONTAINER"
	WarningInfo_Image     = "ALERT_TYPE_IMAGE"
	WarningInfo_MailError = "ALERT_MAIL_ERROR"

	// 告警信息级别
	WarningLevel_High   = "ALERT_SEVERITY_HIGH"
	WarningLevel_Medium = "ALERT_SEVERITY_MEDIUM"
	WarningLevel_Low    = "ALERT_SEVERITY_LOW"

	//告警信息状态
	WarningStatus_Not_Dealed = "未处理"

	// 日志级别
	Log_level_Debug = "Debug"
	Log_level_Info  = "Info"
	Log_level_Warn  = "Warn"
	Log_level_Error = "Error"

	// kubernetes 相关
	Kubernetes_Object_Spec     = "Spec"
	Kubernetes_Object_MetaData = "MetaData"

	// diss-api
	//漏洞安全等级
	Vulnerabilities_Severity_Medium     = "Medium"
	Vulnerabilities_Severity_Negligible = "Negligible"
	Vulnerabilities_Severity_Low        = "Low"
	Vulnerabilities_Severity_Critical   = "Critical"

	// log type tag
	LogToEmail_Prefix            = "日志导出"
	LogType_BenchMarkLog         = "BenchmarkLog"              // 基线扫描日志
	LogType_ImageSecLog          = "ImageSecurityLog"          // 镜像安全日志
	LogType_ContainerVirusLog    = "ContainerVirusLog"         // 容器杀毒日志
	LogType_IntrudeDetectLog     = "IntrudeDetectLog"          // 入侵检测日志
	LogType_ContainerSecAuditLog = "ContainerSecurityAuditLog" // 容器安全审计日志
	LogType_CommandSecAuditLog   = "CommandSecurityAuditLog"   // 命令安全审计日志
	LogSubType_DockerEvent       = "DockerEvent"               // DockerEvent子类型（容器安全审计下）

	// license type
	LicType_TrialLicense    = "TrialLicense"
	LicType_StandardLicense = "StandardLicense"
	LicFile_Extension       = ".lic"

	LicModuleType_BenchMark = "BenchMark"

	EncryptedFileType_License     = "license"
	EncryptedFileType_FeatureCode = "featureCode"

	// 系统配置类型
	FeatureCode       = "FeatureCode"
	EmailServerConfig = "EmailServerConfig"
	LDAPClientConfig  = "LDAPClientConfig"
	Login_Type_LDAP   = "LDAP"

	// LogConfigs
	Log_Config_SysLog_Export = "SysLogExport"
	Log_Config_To_Mail       = "LogToMail"

	MailServer_Not_Available_Msg = "邮箱服务器不可用"
	Mail_CanNotSend_Msg          = "邮件发送不成功"

	// 告警邮件字段定义
	MailField_Subject     = "Subject"
	MailField_LogType     = "LogType"
	MailField_InfoSubType = "InfoSubType"
	MailField_Body        = "Body"
	MailField_From        = "From"
	MailField_To          = "To"

	MailServerStatus_Normal   = "Normal"
	MailServerStatus_Abnormal = "Abnormal"

	// system config
	Enable = "Enable"

	Null_Time = "0001-01-01 00:00:00 +0000 UTC"

	// time and zone
	CstZone = time.FixedZone("CST", 8*3600)

	// db image type
	DB_Image_type_Mysql    = "Mysql"
	DB_Image_type_Oracle   = "Oracle"
	DB_Image_type_Redis    = "Redis"
	DB_Image_type_Postgres = "Postgres"
	DB_Image_type_Mongodb  = "Mongodb"
	DB_Image_type_Memcache = "Memcache"
	DB_Image_type_DB2      = "DB2"
	DB_Image_type_Hbase    = "Hbase"

	// ==============  anchore engine ==============
	// 镜像漏洞等级
	Image_vuln_Severity_Low    = "Low"
	Image_vuln_Severity_Medium = "Medium"

	// 镜像 content 类型
	Image_Content_Type_OS        = "os"
	Image_Vuln_Type_All          = "all"
	Image_Metadata_Type_Manifest = "manifest"

	// white list
	WarnWhiteListConfigKey             = "WarnWhiteList"
	WarnWhiteListCnTrans_Node          = []string{"节点", "node"}
	WarnWhiteListCnTrans_ContainerId   = []string{"容器ID", "container\\.id"}
	WarnWhiteListCnTrans_ContainerName = []string{"容器名称", "container\\.name"}
	WarnWhiteListCnTrans_CmdLine       = []string{"命令包含", "proc\\.cmdline"}

	// png type
	PictureType = "image/png"

	WarnInfoStatus = "已处理"

	FailStatus = "处理失败"

	// status
	Status = "Processed Success"

	// resp center operation
	Container = "Container"

	//const key
	StatusKey = "Status"

	//const key
	WarningInfoId = "WarningInfoId"
)
