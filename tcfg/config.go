package tcfg

import (
	"encoding/base64"
	"fmt"
	"time"
)

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

type Key string

func (k Key) String() string {
	return string(k)
}

func (k Key) Bytes() ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(string(k))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	} else if len(b) != 32 {
		return nil, fmt.Errorf("invalid len (must be 32): %d", len(b))
	}

	return b, nil
}

type IConfig interface {
	AppName() string
	AppKey() Key
	AppEnv() Env
	LogLevel() string
	ShutdownTimeout() time.Duration

	Load() error
}
