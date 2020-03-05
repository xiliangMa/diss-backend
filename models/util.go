package models

var (
	// k8s
	Cluster_Synced                  = true
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
	Bench_Mark_Type_Docker     = 0
	Bench_Mark_Type_Kubernetes = 1

	// 安全容器
	Diss_All           = -1
	Diss_Installed     = 0
	Diss_Not_Installed = 1

	//安全状态
	Diss_Status_All    = -1
	Diss_status_Safe   = 0
	Diss_Status_Unsafe = 1
)
