package tcfg

import "time"

type Env string

func (e Env) String() string {
	return string(e)
}

const (
	EnvDev   Env = "dev"
	EnvStage Env = "stage"
	EnvProd  Env = "prod"
)

type IConfig interface {
	AppName() string
	AppKey() string
	AppEnv() Env
	LogLevel() string
	ShutdownTimeout() time.Duration
	FrontendURL() string

	Load() error
}
