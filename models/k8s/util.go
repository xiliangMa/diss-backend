package k8s

const (
	Cluster_Status_Run = iota
	Cluster_Status_NoRun
	Cluster_Synced   = true
	Cluster_NoSynced = false
	Cluster_IsSync   = true
	Cluster_NoSync   = false
)
