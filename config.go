package theapp

type IConfig interface {
	AppName() string
	AppKey() string
	LogLevel() string
	ShutdownGracePeriod() int
	FrontendURL() string

	Load() error
}
