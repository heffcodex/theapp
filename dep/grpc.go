package dep

import (
	"fmt"
	"google.golang.org/grpc"
)

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

func NewGRPC(cfg GRPCConfig, dialOptions ...grpc.DialOption) *D[*grpc.ClientConn] {
	resolve := func() (*grpc.ClientConn, error) {
		return grpc.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), dialOptions...)
	}

	return NewDep(true, resolve)
}
