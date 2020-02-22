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
	GetHostPsErr            = 27
	AddHostPsErr            = 28
	EditHostInfoErr         = 29
	GetHostInfoErr          = 30
	AddHostConfigErr        = 31
	EditHostConfigErr       = 32
	GetHostConfigErr        = 33
	EditContainerConfigErr  = 34
	AddImageConfigErr       = 35
	GetImageConfigErr       = 36
	AddContainerConfigErr   = 37
	GetContainerConfigErr   = 38
	EditImageConfigErr      = 39
	EditHostPsErr           = 40
	GetContainerInfoErr     = 41
	AddContainerInfoErr     = 42
	EditContainerInfoErr    = 43
	DeleteHostPsErr         = 44
	DeleteContainerInfoErr  = 45

	//system
	CheckK8sFilePostErr         = 60
	CheckK8sFileCreateClientErr = 61
	CheckK8sFileIsExistErr      = 62
	UploadK8sFileErr            = 63
	CheckK8sFileTestErr         = 64

	//k8s
	AddNameSpaceErr  = 100
	AddPodErr        = 101
	AddClusterErr    = 102
	GetClusterErr    = 103
	GetPodErr        = 104
	EditNameSpaceErr = 105
	EditPodErr       = 106
	GetNameSpaceErr  = 107

	//task
	AddTaskErr = 130
	GetTaskErr = 131
)
