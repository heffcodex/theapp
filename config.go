package theapp

type IConfig interface {
	AppName() string
	LogLevel() string
	ShutdownGracePeriod() int

	Load() error
}
