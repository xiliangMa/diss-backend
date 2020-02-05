package utils

var (
	Fail               = 1
	GetHostListErr     = 2
	AddHostErr         = 3
	DeleteHostErr      = 4
	GetHostZero        = 5
	HostExistError     = 6
	GetHostMetricError = 7
	EditHostErr        = 8
	SiginErr           = 17
	AuthorizeErr       = 18
	GetUserInfoErr     = 19
	ElasticConnErr     = 21
	ElasticSearchErr   = 22

	//system
	CheckK8sFilePostErr         = 60
	CheckK8sFileCreateClientErr = 61
	CheckK8sFileIsExistErr      = 62
	UploadK8sFileErr            = 63
	CheckK8sFileTestErr         = 64
)
