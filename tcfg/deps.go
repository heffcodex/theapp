package tcfg

import "time"

type BunPostgres struct {
	DSN            string `mapstructure:"dsn"`
	MaxConnections int    `mapstructure:"maxConnections"`
	MaxIdleTime    int    `mapstructure:"maxIdleTime"`
}

func (c *BunPostgres) MaxIdleTimeSeconds() time.Duration {
	if c.MaxIdleTime < 1 {
		return 0
	}

	return time.Duration(c.MaxIdleTime) * time.Second
}

type GRPCClient struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}
