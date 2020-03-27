package models

var (
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

	// 分组类型
	Group_Host      = 0
	Group_Container = 1

	// 安全容器
	Diss_All           = -1
	Diss_Installed     = 0
	Diss_Not_Installed = 1

	//安全状态
	Diss_Status_All    = -1
	Diss_status_Safe   = 0
	Diss_Status_Unsafe = 1

	//基线日志类型
	BMLT_Host_All  = "host"
	BMLT_Docker    = "docker"
	BMLT_K8s       = "k8s"
	BML_Level_Info = "INFO"
	BML_Level_Warn = "WARN"
	BML_Level_Pass = "PASS"
	BML_Level_fail = "FAIL"

	// 入侵检测类型
	IDLT_Docker = "docker"
	IDLT_Host   = "host"

	// 租户
	Account_Admin = "admin"
)
