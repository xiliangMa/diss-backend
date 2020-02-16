package utils

var (
	Fail                    = 1
	GetHostListErr          = 2
	AddHostErr              = 3
	DeleteHostErr           = 4
	GetHostZero             = 5
	HostExistError          = 6
	GetHostMetricError      = 7
	EditHostErr             = 8
	SiginErr                = 17
	AuthorizeErr            = 18
	GetUserInfoErr          = 19
	ElasticConnErr          = 21
	ElasticSearchErr        = 22
	AddBenchMarkTemplateErr = 23
	GetBenchMarkTemplateErr = 24
	AddBenchMarkLogErr      = 25
	GetBenchMarkLogErr      = 26

	//system
	CheckK8sFilePostErr         = 60
	CheckK8sFileCreateClientErr = 61
	CheckK8sFileIsExistErr      = 62
	UploadK8sFileErr            = 63
	CheckK8sFileTestErr         = 64

	//k8s
	AddNameSpaceErr       = 100
	AddPodErr             = 101
	AddClusterErr         = 102
	GetClusterErr         = 103
	AddHostConfigErr      = 104
	GetHostConfigErr      = 105
	AddImageConfigErr     = 106
	GetImageConfigErr     = 107
	GetPodErr             = 108
	AddContainerConfigErr = 109
	GetContainerConfigErr = 110

	//task
	AddTaskErr = 130
	GetTaskErr = 131
)
