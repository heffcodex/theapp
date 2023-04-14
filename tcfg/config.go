package tcfg

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type IConfig interface {
	AppName() string
	AppKey() Key
	AppEnv() Env
	LogLevel() string
	ShutdownTimeout() time.Duration
}

var _ IConfig = (*Config)(nil)

type Config struct {
	App App `mapstructure:"app"`
}

func (c Config) AppName() string {
	return c.App.Name
}

func (c Config) AppKey() Key {
	return c.App.Key
}

func (c Config) AppEnv() Env {
	return c.App.Env
}

func (c Config) LogLevel() string {
	return c.App.LogLevel
}

func (c Config) ShutdownTimeout() time.Duration {
	return time.Duration(c.App.ShutdownTimeout) * time.Second
}

func LoadConfig[C IConfig]() (C, error) {
	v := viper.New()

	v.AddConfigPath(".")
	v.AddConfigPath("./.data")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvPrefix("CFG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return *new(C), fmt.Errorf("read: %w", err)
	}

	var config C

	if err := v.UnmarshalExact(&config); err != nil {
		return *new(C), fmt.Errorf("unmarshal exact: %w", err)
	}

	if err := config.AppKey().Validate(); err != nil {
		return *new(C), fmt.Errorf("validate app key: %w", err)
	}

	return config, nil
}
