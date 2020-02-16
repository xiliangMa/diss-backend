package utils

import (
	"io/ioutil"
	"testing"
)

var (
	kubeconfig string
	path       = "../kubeconfig/config"
	jobFile    = "../conf/kube-bench/kube-bench-job.yml"
	namespaces = "default"
	jobName    = "kube-bench"
	podName    = "kube-bench"
	clientgo   ClientGo
	kubeBenchJobCommand    = []string{"kube-bench", "node", "--benchmark", "cis-1.3"}
)

func Test_CreateK8sClient(t *testing.T) {
	clientgo = CreateK8sClient(path)
	if clientgo.ErrMessage == "" {
		t.Log("K8S Client create Success")
	} else {
		t.Error("K8S Client create Fail")
	}
}

func Test_GetNodes(t *testing.T) {
	if clientgo.ErrMessage == "" {
		nodes, err := clientgo.GetNodes()
		if err == nil {
			t.Logf("集群节点个数 %d", len(nodes.Items))
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}

func Test_GetPodsByNameSpace(t *testing.T) {
	if clientgo.ErrMessage == "" {
		pods, err := clientgo.GetPodsByNameSpace(namespaces)
		if err == nil {
			t.Logf("default 命名空间下的 pod 个数 %d", len(pods.Items))
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}

func Test_GetPodLogsByNameSpace(t *testing.T) {
	if clientgo.ErrMessage == "" {
		request := clientgo.GetPodLogs(namespaces, podName)
		if body, _ := request.Stream(); body != nil {
			log, _ := ioutil.ReadAll(body)
			t.Logf("Pod logs： %s", log)
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}

func Test_CreateJobByYml(t *testing.T) {
	if clientgo.ErrMessage == "" {
		job, err := clientgo.CreateJobByYml(kubeBenchJobCommand, jobFile, namespaces)
		if err != nil {
			t.Logf("Create job err, %s", err)
		} else {
			t.Logf("Job status %d", job.Status.Succeeded)
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}

func Test_GetJob(t *testing.T) {
	if clientgo.ErrMessage == "" {
		job, err := clientgo.GetJob(namespaces, podName)
		if err != nil && job == nil {
			t.Logf("Get job err, %s", err)
		} else {
			t.Logf("Job status %d", job.Status.Succeeded)
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}

//func Test_DeleteJob(t *testing.T) {
//	if clientgo.ErrMessage == "" {
//		err := clientgo.DeleteJob(namespaces, jobName)
//		if err != nil {
//			t.Logf("Delete job err, %s", err)
//		} else {
//			t.Logf("Delete job %s success", jobName)
//		}
//	} else {
//		t.Error("K8S Client create Fail")
//	}
//}

func Test_GetNameSpaces(t *testing.T) {
	if clientgo.ErrMessage == "" {
		ns, err := clientgo.GetNameSpaces()
		if err != nil {
			t.Logf("Get namespaces err: %s", err)
		} else {
			t.Logf("集群内的 namespace 个数为: %d ", len(ns.Items))
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}

func Test_GetPodMetricsAndList(t *testing.T) {
	if clientgo.ErrMessage == "" {
		podMetrics, err := clientgo.GetPodMetrics("default", "nginx-5b48b4bf7c-d7flv")
		if err != nil {
			t.Logf("Get PodMetrics err: %s", err)
		} else {
			t.Log("Get PodMetrics", podMetrics)
		}
		podMetricsList, err := clientgo.GetPodMetricsList("default")
		if err != nil {
			t.Logf("Get PodMetricsList err: %s", err)
		} else {
			t.Log("Get PodMetricsList", podMetricsList)
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}

func Test_GetNodeMetricsAndList(t *testing.T) {
	if clientgo.ErrMessage == "" {
		nodeMetricsList, err := clientgo.GetNodeMetricsList()
		if err != nil {
			t.Logf("Get GetNodeMetricsList err: %s", err)
		} else {
			t.Log("Get NodeMetricsList", nodeMetricsList)
		}
		nodeMetrics, err := clientgo.GetNodeMetrics("izm5e3cntl0pztm4dj3phrz")
		if err != nil {
			t.Logf("Get nodeMetrics err: %s", err)
		} else {
			t.Log("Get nodeMetrics", nodeMetrics)
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}
