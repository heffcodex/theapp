package tcfg

import "time"

type Env string

func (e Env) String() string {
	return string(e)
}

const (
	EnvDev   Env = "dev"
	EnvTest  Env = "test"
	EnvStage Env = "staging"
	EnvProd  Env = "production"
)

type IConfig interface {
	AppName() string
	AppKey() string
	AppEnv() Env
	LogLevel() string
	ShutdownTimeout() time.Duration

	Load() error
}
