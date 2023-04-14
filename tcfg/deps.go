package tcfg

type BunPostgres struct {
	DSN            string `mapstructure:"dsn"`
	MaxConnections int    `mapstructure:"maxConnections"`
}

type GRPCClient struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}
