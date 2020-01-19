package utils

import (
	"testing"
)

var (
	kubeconfig string
	configName = "config"
	path       = "../kubeconfig"
	namespaces = "default"
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
