package models

var (
	// k8s
	Cluster_Synced                  = true
	Cluster_IsSync                  = true
	Cluster_NoSync                  = false
	Pod_Container_Statue_Running    = "Running"
	Pod_Container_Statue_Terminated = "Terminated"
	Pod_Container_Statue_Waiting    = "Waiting"
	Container_Status_Run            = "Run"
	Container_Status_Pause          = "Pause"
	Container_Status_All            = "All"
	// bnech mark
	Bench_Mark_Type_Docker     = 0
	Bench_Mark_Type_Kubernetes = 1
)
