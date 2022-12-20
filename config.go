package theapp

import "time"

type IConfig interface {
	AppName() string
	AppKey() string
	LogLevel() string
	ShutdownTimeout() time.Duration
	FrontendURL() string

	Load() error
}
