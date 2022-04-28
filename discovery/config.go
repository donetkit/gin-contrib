package discovery

type Config struct {
	Id                  string
	ServiceName         string
	ServiceRegisterAddr string
	ServiceRegisterPort int
	ServiceCheckAddr    string
	ServiceCheckPort    int
	Tags                []string
	IntervalTime        int
	DeregisterTime      int
	TimeOut             int
	CheckHTTP           string
}
