package tcfg

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/hkdf"
)

var (
	KeyEncoding = base64.StdEncoding
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
	_, err := k.getBytes()
	return err
}

func (k Key) String() string {
	return string(k)
}

func (k Key) Bytes() []byte {
	b, err := k.getBytes()
	if err != nil {
		panic(err)
	}

	return b
}

func (k Key) Extract(salt string) Key {
	derived := hkdf.Extract(sha256.New, k.Bytes(), []byte(salt))
	b64 := KeyEncoding.EncodeToString(derived)

	return Key(b64)
}

func (k Key) getBytes() ([]byte, error) {
	if len(k) == 32 {
		return []byte(k), nil
	}

	b, err := KeyEncoding.DecodeString(string(k))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	} else if len(b) != 32 {
		return nil, fmt.Errorf("invalid len (must be 32): %d", len(b))
	}

	return b, nil
}
