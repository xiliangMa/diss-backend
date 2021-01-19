package models

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/xiliangMa/diss-backend/utils"
	"io"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"math/rand"
	"net/http"
	"os"
)

type KubernetesClientManager struct {
	ClientHub       map[string]ClientGo
	DymaicClientHub map[string]DymaicClient
}

func NewKubernetesClientManager() *KubernetesClientManager {
	return &KubernetesClientManager{
		ClientHub:       make(map[string]ClientGo),
		DymaicClientHub: make(map[string]DymaicClient),
	}
}

type ApiParams struct {
	KubeConfigPath string
	BearerToken    string
	MasterUrl      string
	AuthType       string
}

type ClientGo struct {
	ClientSet  *kubernetes.Clientset
	ErrMessage string
}

type DymaicClient struct {
	Interface  dynamic.Interface
	ErrMessage string
}

type KubernetesHandler struct {
	*ClientGo
	*DymaicClient
	IsActive bool
	File     *os.File
}

func CreateK8sClient(params *ApiParams) ClientGo {
	var config *rest.Config
	var err error
	clientgo := ClientGo{nil, ""}
	if params.AuthType == "BearerToken" {
		// 这里是使用用户名和密码调用APIserver，所以kubeconfig为空
		kubeconfig := flag.String(string(rand.Intn(1000)), "", "BearerToken")
		flag.Parse()
		config, err = clientcmd.BuildConfigFromFlags(params.MasterUrl, *kubeconfig)
		config.Insecure = true
		config.BearerToken = params.BearerToken
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", params.KubeConfigPath)
	}
	if err != nil {
		clientgo.ErrMessage = "BuildConfigFromFlags"
		return clientgo
	}
	// 根据指定的 config 创建 clientset
	clientSet, err := kubernetes.NewForConfig(config)
	clientgo.ClientSet = clientSet
	if err != nil {
		clientgo.ErrMessage = err.Error()
	}

	// 检测是否可用
	if _, err := clientgo.GetNodes(); err != nil {
		clientgo.ErrMessage = err.Error()
	}
	return clientgo
}

func (clientgo *ClientGo) GetNodes() (*v1.NodeList, error) {
	return clientgo.ClientSet.CoreV1().Nodes().List(nil, v12.ListOptions{})
}

func InitClientHub() {
	// 获取所有集群
	cluster := new(Cluster)
	result := cluster.List(0, 0)

	if result.Code == http.StatusOK && result.Data != nil {
		data := result.Data.(map[string]interface{})
		clusterList := data["items"].([]*Cluster)
		for _, c := range clusterList {
			// 构建k8s客户端
			clientGo := CreateK8sClient(BuildApiParams(c))
			if clientGo.ErrMessage == "" {
				// 加入缓存
				KCM.ClientHub[c.Id] = clientGo
				logs.Info("Init ClientHub add cluster： %s success.", c.Name)
			}
		}
	}
}

func CreateK8sDymaicClient(params *ApiParams) DymaicClient {
	var config *rest.Config
	var err error
	dymaicClient := DymaicClient{nil, ""}
	if params.AuthType == "BearerToken" {
		// 这里是使用用户名和密码调用APIserver，所以kubeconfig为空
		kubeconfig := flag.String(string(rand.Intn(1000)), "", "BearerToken")
		flag.Parse()
		config, err = clientcmd.BuildConfigFromFlags(params.MasterUrl, *kubeconfig)
		config.Insecure = true
		config.BearerToken = params.BearerToken
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", params.KubeConfigPath)
	}
	if err != nil {
		dymaicClient.ErrMessage = "BuildConfigFromFlags"
		return dymaicClient
	}
	// 根据指定的 config 创建 clientset
	dyClientInface, err := dynamic.NewForConfig(config)
	dymaicClient.Interface = dyClientInface
	if err != nil {
		dymaicClient.ErrMessage = err.Error()
	}
	// todo 检测是否可用
	return dymaicClient
}

func InitDymaicClientHub() {
	// 获取所有集群
	cluster := new(Cluster)
	result := cluster.List(0, 0)

	if result.Code == http.StatusOK && result.Data != nil {
		data := result.Data.(map[string]interface{})
		clusterList := data["items"].([]*Cluster)
		for _, c := range clusterList {
			// 构建k8s客户端
			dymaicClient := CreateK8sDymaicClient(BuildApiParams(c))
			if dymaicClient.ErrMessage == "" {
				// 加入缓存
				KCM.DymaicClientHub[c.Id] = dymaicClient
				logs.Info("Init DymaicClientHub add cluster： %s success.", c.Name)
			}
		}
	}
}

func BuildApiParams(c *Cluster) *ApiParams {
	params := new(ApiParams)
	params.AuthType = c.AuthType
	params.BearerToken = c.BearerToken
	params.MasterUrl = c.MasterUrls
	params.KubeConfigPath = c.FileName
	return params
}

func GetDymaicClient(c *Cluster) DymaicClient {
	clusterId := c.Id
	if _, ok := KCM.DymaicClientHub[clusterId]; !ok {
		KCM.DymaicClientHub[clusterId] = CreateK8sDymaicClient(BuildApiParams(c))
	}
	return KCM.DymaicClientHub[clusterId]
}

func GetClient(c *Cluster) ClientGo {
	clusterId := c.Id
	if _, ok := KCM.ClientHub[clusterId]; !ok {
		KCM.ClientHub[clusterId] = CreateK8sClient(BuildApiParams(c))
	}
	return KCM.ClientHub[clusterId]
}

func (this *KubernetesHandler) CreateOrDeleteResourceByYml() Result {
	result := Result{Code: http.StatusOK}
	namespace := "default"
	operatorType := "Create"

	// 动态解析参数
	d := yaml.NewYAMLOrJSONDecoder(this.File, 4096)

	restMapperRes, err := restmapper.GetAPIGroupResources(this.ClientGo.ClientSet.Discovery())
	if err != nil {
		logs.Error("RestMapperRes yml fail. Err: %s", err)
	}
	restMapper := restmapper.NewDiscoveryRESTMapper(restMapperRes)

	ext := runtime.RawExtension{}
	if err := d.Decode(&ext); err != nil {
		if err == io.EOF {
			logs.Error("Decode scope yml fail. Err: %s", err)
		}
	}
	list := new(v1.List)
	json.Unmarshal(ext.Raw, list)
	for _, object := range list.Items {
		obj, gvk, err := unstructured.UnstructuredJSONScheme.Decode(object.Raw, nil, nil)
		if err != nil {
			logs.Error("UnstructuredJSONScheme scope yml fail. Err: %s", err)
		}
		mapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			logs.Error("RESTMapping scope yml fail. Err: %s", err)
		}

		// runtime.Object转换为unstructed
		unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			logs.Error("DefaultUnstructuredConverter scope yml fail. Err: %s", err)
		}

		var unstruct unstructured.Unstructured

		unstruct.Object = unstructuredObj

		if md, ok := unstruct.Object["metadata"]; ok {
			metadata := md.(map[string]interface{})
			if internalns, ok := metadata["namespace"]; ok {
				namespace = internalns.(string)
			}
		}
		// 创建 Resource
		if this.IsActive {
			if mapping.Resource.Resource == "namespaces" || mapping.Resource.Resource == "clusterroles" || mapping.Resource.Resource == "clusterrolebindings" {
				_, err = this.DymaicClient.Interface.Resource(mapping.Resource).Create(nil, &unstruct, metav1.CreateOptions{})
			} else {
				_, err = this.DymaicClient.Interface.Resource(mapping.Resource).Namespace(namespace).Create(nil, &unstruct, metav1.CreateOptions{})
			}
		} else {
			operatorType = "Delete"
			if mapping.Resource.Resource == "namespaces" || mapping.Resource.Resource == "clusterroles" || mapping.Resource.Resource == "clusterrolebindings" {
				err = this.DymaicClient.Interface.Resource(mapping.Resource).Delete(nil, unstruct.GetName(), metav1.DeleteOptions{})
			} else {
				err = this.DymaicClient.Interface.Resource(mapping.Resource).Namespace(namespace).Delete(nil, unstruct.GetName(), metav1.DeleteOptions{})
			}
		}
		if err != nil {
			msg := fmt.Sprintf("%s kind: %s, name: %s  fail. Err: %s", operatorType, mapping.Resource.Resource, unstruct.GetName(), err)
			logs.Error(msg)
			result.Code = utils.CretaeResourceError
			result.Message = msg
			return result
		} else {
			msg := fmt.Sprintf("%s kind: %s, name: %s  success.", operatorType, mapping.Resource.Resource, unstruct.GetName())
			logs.Info(msg)
		}
	}
	return result
}
