package models

var (
	// Global
	All          = "All"
	Result_Items = "items"
	Result_Total = "total"

	// k8s
	Cluster_Sync_Status_NotSynced   = "NotSynced"
	Cluster_Sync_Status_Synced      = "Synced"
	Cluster_Sync_Status_InProcess   = "InProcess"
	Cluster_Sync_Status_Clearing    = "Clearing"
	Cluster_Sync_Status_Fail        = "Fail"
	Cluster_Watch_Status_Fail       = "Fail"
	Cluster_Watch_Status_Success    = "Success"
	Cluster_IsSync                  = true
	Cluster_NoSync                  = false
	Cluster_Status_Active           = "Active"
	Cluster_Status_Unavailable      = "Unavailable"
	Cluster_Type_Kubernets          = "Kubernets"
	Cluster_Type_OpenShift          = "OpenShift"
	Cluster_Type_Rancher            = "Rancher"
	Pod_Container_Statue_Running    = "Running"
	Pod_Container_Statue_Terminated = "Terminated"
	Pod_Container_Statue_Waiting    = "Waiting"
	Api_Auth_Type_KubeConfig        = "KubeConfig"
	Api_Auth_Type_BearerToken       = "BearerToken"
	Cluster_Data_AccountName        = "Kubernetes"

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
	//镜像仓库扫描：
	TMP_Type_ImageScan = "ImageScan"
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

	Job_Status_Enable  = "Enable"
	Job_Status_Active  = "Active"
	Job_Status_Disable = "Disable"

	//安全检查类型
	SC_Type_Host      = "host"
	Sc_Type_Container = "container"

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
	Subject_Common     = "Common"
	Subject_Image_Safe = "IMAGE_SAFE"

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

	// LogConfigs
	Log_Config_SysLog_Export = "SysLogExport"

	// syslog exported log type
	SysLog_BenchScanLog      = "BenchmarkScanLog"  // 基线扫描日志
	SysLog_ImageSecLog       = "ImageSecurityLog"  // 镜像安全日志
	SysLog_ContainerVirusLog = "ContainerVirusLog" // 容器杀毒日志
	SysLog_IDSLog            = "IntrudeDetectLog"  // 入侵检测日志
	SysLog_SecAuditLog       = "SecurityAuditLog"  // 安全审计日志

	// license type
	LicType_TrialLicense    = "TrialLicense"
	LicType_StandardLicense = "StandardLicense"

	LicFile_Extension = ".lic"

	// system config
	Enable = "Enable"

	Null_Time = "0001-01-01 00:00:00 +0000 UTC"
)
