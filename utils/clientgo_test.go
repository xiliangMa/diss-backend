package utils

import (
	"io/ioutil"
	"testing"
)

var (
	kubeconfig string
	configName = "config"
	path       = "../kubeconfig"
	jobfile    = "../conf/kube-bench/kube-bench-job.yml"
	namespaces = "default"
	pod        = "kube-bench"
	clientgo   ClientGo
)

func Test_CreateK8sClient(t *testing.T) {
	clientgo = CreateK8sClient(kubeconfig, path, configName)
	if clientgo.err == nil {
		t.Log("K8S Client create Success")
	} else {
		t.Error("K8S Client create Fail")
	}
}

func Test_GetNodes(t *testing.T) {
	if clientgo.err == nil {
		nodes, err := clientgo.GetNodes()
		if err == nil {
			t.Logf("集群节点个数 %d", len(nodes.Items))
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}

func Test_GetPodsByNameSpace(t *testing.T) {
	if clientgo.err == nil {
		pods, err := clientgo.GetPodsByNameSpace(namespaces)
		if err == nil {
			t.Logf("default 命名空间下的 pod个数 %d", len(pods.Items))
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}

func Test_GetPodLogsByNameSpace(t *testing.T) {
	if clientgo.err == nil {
		request := clientgo.GetPodLogs(namespaces, pod)
		if body, _ := request.Stream(); body != nil {
			log, _ := ioutil.ReadAll(body)
			t.Logf("Pod logs： %s", log)
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}

func Test_GetJob(t *testing.T) {
	if clientgo.err == nil {
		job, err := clientgo.GetJob(namespaces, pod)
		if err != nil && job == nil {
			t.Logf("Get job err, %s", err)
		} else {
			t.Logf("Job status %d", job.Status.Succeeded)
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}

func Test_CreateJobByYml(t *testing.T) {
	if clientgo.err == nil {
		job, err := clientgo.CreateJobByYml(jobfile, namespaces)
		if err != nil {
			t.Logf("Create job err, %s", err)
		} else {
			t.Logf("Job status %d", job.Status.Succeeded)
		}
	} else {
		t.Error("K8S Client create Fail")
	}
}
