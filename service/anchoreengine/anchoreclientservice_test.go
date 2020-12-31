package anchoreengine

import (
	"encoding/base64"
	"github.com/anchore/client-go/pkg/external"
	"testing"
)

func Test_GetImageContent(t *testing.T) {
	cfg := external.NewConfiguration()
	cfg.AddDefaultHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:foobar")))
	cfg.AddDefaultHeader("accept", "application/json")
	cfg.BasePath = "http://49.232.153.63:8080/v1"
	anchoreClient := external.NewAPIClient(cfg)
	list, _, err := anchoreClient.ImagesApi.GetImageContentByTypeImageId(nil, "f70734b6a266dcb5f44c383274821207885b549b75c8e119404917a61335981a", "files", nil)
	if err != nil {
		t.Errorf("Create anchore client failed: err: %s.", err.Error())
	}
	t.Log(list)
	//for _, obj := range list.Content {
	//	t.Log(obj)
	//}
}

func Test_GetImageVuln(t *testing.T) {
	cfg := external.NewConfiguration()
	cfg.AddDefaultHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:foobar")))
	cfg.AddDefaultHeader("accept", "application/json")
	cfg.BasePath = "http://49.232.153.63:8080/v1"
	anchoreClient := external.NewAPIClient(cfg)
	list, _, err := anchoreClient.ImagesApi.GetImageVulnerabilitiesByTypeImageId(nil, "f70734b6a266dcb5f44c383274821207885b549b75c8e119404917a61335981a", "all", nil)
	if err != nil {
		t.Errorf("Create anchore client failed: err: %s.", err.Error())
	}
	for _, obj := range list.Vulnerabilities {
		t.Log(obj)
	}
}

func Test_GetImageVulnStatistics(t *testing.T) {
	cfg := external.NewConfiguration()
	cfg.AddDefaultHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:foobar")))
	cfg.AddDefaultHeader("accept", "application/json")
	cfg.BasePath = "http://49.232.153.63:8080/v1"
	anchoreClient := external.NewAPIClient(cfg)
	list, _, err := anchoreClient.ImagesApi.GetImageVulnerabilitiesByTypeImageId(nil, "f70734b6a266dcb5f44c383274821207885b549b75c8e119404917a61335981a", "all", nil)
	if err != nil {
		t.Errorf("Create anchore client failed: err: %s.", err.Error())
	}
	vulnStatistics := make(map[string]interface{})
	vulnLowCount := 0
	vulnMediumCount := 0
	if &list != nil && len(list.Vulnerabilities) > 0 {
		for _, vuln := range list.Vulnerabilities {
			switch vuln.Severity {
			case "Low":
				vulnLowCount = vulnLowCount + 1
			case "Medium":
				vulnMediumCount = vulnMediumCount + 1
			}
		}
	}
	vulnStatistics["Low"] = vulnLowCount
	vulnStatistics["Medium"] = vulnMediumCount
	t.Log(vulnStatistics)
}

func Test_GetImageMetadata(t *testing.T) {
	cfg := external.NewConfiguration()
	cfg.AddDefaultHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:foobar")))
	cfg.AddDefaultHeader("accept", "application/json")
	cfg.BasePath = "http://49.232.153.63:8080/v1"
	anchoreClient := external.NewAPIClient(cfg)
	list, _, err := anchoreClient.ImagesApi.GetImageMetadataByType(nil, "sha256:39eda93d15866957feaee28f8fc5adb545276a64147445c64992ef69804dbf01", "manifest", nil)
	if err != nil {
		t.Errorf("Create anchore client failed: err: %s.", err.Error())
	}
	t.Log(list)
}
