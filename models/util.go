package models

const (
	// k8s
	Cluster_Status_Run = iota
	Cluster_Status_NoRun
	Cluster_Synced   = true
	Cluster_NoSynced = false
	Cluster_IsSync   = true
	Cluster_NoSync   = false

	// bnech mark
	Bench_Mark_Type_Docker     = 0
	Bench_Mark_Type_Kubernetes = 1
)
