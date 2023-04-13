package tcfg

import (
	"encoding/base64"
	"fmt"
)

type App struct {
	Name            string `mapstructure:"name"`
	Key             Key    `mapstructure:"key"`
	Env             Env    `mapstructure:"env"`
	LogLevel        string `mapstructure:"logLevel"`
	ShutdownTimeout int    `mapstructure:"shutdownTimeout"`
}

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

func (k Key) MustBytes() []byte {
	b, err := k.Bytes()
	if err != nil {
		panic(err)
	}

	return b
}
