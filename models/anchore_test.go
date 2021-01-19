package models

import (
	"github.com/anchore/client-go/pkg/external"
	"github.com/beego/beego/v2/core/logs"
	"testing"
)

func Test_NewAnchoreEngineManager(t *testing.T) {
	cfg := external.NewConfiguration()
	cfg.AddDefaultHeader("Authorization", "Basic YWRtaW46Zm9vYmFy")
	cfg.AddDefaultHeader("accept", "application/json")
	cfg.BasePath = "http://49.232.153.63:8080/v1"
	anchoreClient := external.NewAPIClient(cfg)
	list, _, err := anchoreClient.ImagesApi.ListImageContent(nil, "f70734b6a266dcb5f44c383274821207885b549b75c8e11", nil)
	if err != nil {
		logs.Error("Get %s failed from anchore engine, err: %s ")
	}
	//imageList, _, err := anchoreClient.ImagesApi.ListImages(nil, nil)
	if err != nil {
		t.Errorf("Create anchore client failed: err: %s.", err.Error())
	}
	for _, obj := range list {
		t.Log(obj)
	}
}
