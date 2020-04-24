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

	// bnech mark
	//Bench_Mark_Type_Docker     = 0
	//Bench_Mark_Type_Kubernetes = 1

	// 安全容器
	Diss_All           = -1
	Diss_Installed     = 0
	Diss_Not_Installed = 1

	//安全状态
	Diss_Status_All    = -1
	Diss_status_Safe   = 0
	Diss_Status_Unsafe = 1

	//系统魔板类型(此处和ws resource tag 保持一致)
	TMP_Type_BM_Docker = "DockerBenchMark"
	TMP_Type_BM_K8S    = "KubernetesBenchMark"
	TMP_Type_VS        = "SC_VirusScan"
	TMP_Type_LS        = "SC_LeakScan"
	TMP_Status_ALl     = -1
	SC_Type_Host       = "host"
	Sc_Type_Container  = "container"

	// 任务状态
	Task_Status_Pending        = "Pending"
	Task_Status_Running        = "Running"
	Task_Status_Pause          = "Pause"
	Task_Status_Finshed        = "Finshed"
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
	Group_All       = -1
	Group_Host      = 0
	Group_Container = 1

	// 租户
	Account_Admin = "admin"

	// nats topic
	Topic_Task    = "Task"
	Topic_Metrics = "Metrics"
)
