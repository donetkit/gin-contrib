package discovery

// Option for queue system
type Option func(*config)

// WithId set id function
func WithId(id string) Option {
	return func(cfg *config) {
		cfg.Id = id
	}
}

// WithServiceName set serviceName function
func WithServiceName(serviceName string) Option {
	return func(cfg *config) {
		cfg.ServiceName = serviceName
	}
}

// WithServiceRegisterAddr set serviceRegisterAddr function
func WithServiceRegisterAddr(serviceRegisterAddr string) Option {
	return func(cfg *config) {
		cfg.ServiceRegisterAddr = serviceRegisterAddr
	}
}

// WithServiceRegisterPort set serviceRegisterPort function
func WithServiceRegisterPort(serviceRegisterPort int) Option {
	return func(cfg *config) {
		cfg.ServiceRegisterPort = serviceRegisterPort
	}
}

// WithServiceCheckAddr set serviceCheckAddr function
func WithServiceCheckAddr(serviceCheckAddr string) Option {
	return func(cfg *config) {
		cfg.ServiceCheckAddr = serviceCheckAddr
	}
}

// WithServiceCheckPort set serviceCheckPort function
func WithServiceCheckPort(serviceCheckPort int) Option {
	return func(cfg *config) {
		cfg.ServiceCheckPort = serviceCheckPort
	}
}

// WithTags set tags function
func WithTags(tags []string) Option {
	return func(cfg *config) {
		cfg.Tags = tags
	}
}

// WithIntervalTime set intervalTime function
func WithIntervalTime(intervalTime int) Option {
	return func(cfg *config) {
		cfg.IntervalTime = intervalTime
	}
}

// WithDeregisterTime set deregisterTime function
func WithDeregisterTime(deregisterTime int) Option {
	return func(cfg *config) {
		cfg.DeregisterTime = deregisterTime
	}
}

// WithTimeOut set timeOut function
func WithTimeOut(timeOut int) Option {
	return func(cfg *config) {
		cfg.TimeOut = timeOut
	}
}

// WithCheckHTTP set checkHTTP function
func WithCheckHTTP(checkHTTP string) Option {
	return func(cfg *config) {
		cfg.CheckHTTP = checkHTTP
	}
}
