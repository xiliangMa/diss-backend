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
		job, err := clientgo.CreateJobByYml(jobFile, namespaces)
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

func Test_DeleteJob(t *testing.T) {
	if clientgo.ErrMessage == "" {
		err := clientgo.DeleteJob(namespaces, jobName)
		if err != nil {
			t.Logf("Delete job err, %s", err)
		} else {
			t.Logf("Delete job %s success", jobName)
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}
