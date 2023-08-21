package tcfg

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config interface {
	AppName() string
	AppKey() Key
	AppEnv() Env
	LogLevel() string
	ShutdownTimeout() time.Duration

	BeforeRead(v *viper.Viper) error
	AfterRead(v *viper.Viper, c Config) error
}

var _ Config = (*Base)(nil)

type Base struct {
	App App `mapstructure:"app"`
}

func (c *Base) AppName() string {
	return c.App.Name
}

func (c *Base) AppKey() Key {
	return c.App.Key
}

func (c *Base) AppEnv() Env {
	return c.App.Env
}

func (c *Base) LogLevel() string {
	return c.App.LogLevel
}

func (c *Base) ShutdownTimeout() time.Duration {
	return time.Duration(c.App.ShutdownTimeout) * time.Second
}

func (*Base) BeforeRead(*viper.Viper) error {
	return nil
}

func (c *Base) AfterRead(*viper.Viper, Config) error {
	if err := c.AppKey().Validate(); err != nil {
		return fmt.Errorf("validate app key: %w", err)
	}

	return nil
}

func LoadConfig[C Config]() (C, error) {
	v := viper.New()

	v.AddConfigPath(".")
	v.AddConfigPath("./.data")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvPrefix("CFG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config C

	if err := config.BeforeRead(v); err != nil {
		return config, fmt.Errorf("before read: %w", err)
	}

	if err := v.ReadInConfig(); err != nil {
		return *new(C), fmt.Errorf("read: %w", err)
	}

	if err := v.UnmarshalExact(&config); err != nil {
		return *new(C), fmt.Errorf("unmarshal exact: %w", err)
	}

	if err := config.AfterRead(v, config); err != nil {
		return config, fmt.Errorf("after read: %w", err)
	}

	return config, nil
}
