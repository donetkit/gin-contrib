package glog

// Option for queue system
type Option func(*Config)

// WithFile set file function
func WithFile(log2File bool) Option {
	return func(cfg *Config) {
		cfg.log2File = log2File
	}
}

// WithHostInfo set hostName, ip function
func WithHostInfo(hostName, ip string) Option {
	return func(cfg *Config) {
		cfg.hostName = hostName
		cfg.ip = ip
	}
}

// WithLogLevel set logLevel function
func WithLogLevel(logLevel LogLevel) Option {
	return func(cfg *Config) {
		cfg.logLevel = logLevel
	}
}

// WithLogColor set logColor function
func WithLogColor(logColor bool) Option {
	return func(cfg *Config) {
		cfg.logColor = logColor
	}
}
