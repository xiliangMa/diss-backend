package sysinit

func init() {
	InitDB()
	InitSecurityLogDB()
	InitLogger()
	InitTask()
}
