package models

var (
	All = "All"

	// k8s
	Cluster_Sync_Status_OK          = 0
	Cluster_Sync_Status_IN_PROGRESS = 1
	Cluster_Sync_Status_FAIL        = 2
	Cluster_IsSync                  = true
	Cluster_NoSync                  = false
	Pod_Container_Statue_Running    = "Running"
	Pod_Container_Statue_Terminated = "Terminated"
	Pod_Container_Statue_Waiting    = "Waiting"

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

	//系统魔板类型(此处和ws resource tag 保持一致)
	TMP_Type_BM_Docker = "DockerBenchMark"
	TMP_Type_BM_K8S    = "KubernetesBenchMark"
	TMP_Type_DockerVS  = "DockerVirusScan"
	TMP_Type_HostVS    = "HostVirusScan"
	TMP_Type_LS        = "SC_LeakScan"

	TMP_Status_Enable  = "Enable"
	TMP_Status_Disable = "Disable"

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
	Subject_Common = "Common"

	// 日志级别
	Log_level_Debug = "Debug"
	Log_level_Info  = "Info"
	Log_level_Warn  = "Warn"
	Log_level_Error = "Error"

	// diss-api
	//漏洞安全等级
	Vulnerabilities_Severity_Medium     = "Medium"
	Vulnerabilities_Severity_Negligible = "Negligible"
	Vulnerabilities_Severity_Low        = "Low"
	Vulnerabilities_Severity_Critical   = "Critical"

	//exported log type
	ImageSecLog       = "ImageSecurityLog"  // 镜像安全日志
	BenchScanLog      = "BenchmarkScanLog"  // 基线扫描日志
	IDSLog            = "IntrudeDetectLog"  // 入侵检测日志
	ContainerVirusLog = "ContainerVirusLog" // 容器杀毒日志
	SecAuditLog       = "SecurityAuditLog"  // 安全审计日志
)
