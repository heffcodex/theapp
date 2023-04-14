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

func (e Env) IsEmpty() bool {
	return e == ""
}

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

func (k Key) Validate() error {
	_, err := k.bytes()
	return err
}

func (k Key) String() string {
	return string(k)
}

func (k Key) Bytes() []byte {
	b, err := k.bytes()
	if err != nil {
		panic(err)
	}

	return b
}

func (k Key) bytes() ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(string(k))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	} else if len(b) != 32 {
		return nil, fmt.Errorf("invalid len (must be 32): %d", len(b))
	}

	return b, nil
}
