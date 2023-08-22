package tdep

import (
	"fmt"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
)

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

func NewGRPC(cfg GRPCConfig, dialOptions []grpc.DialOption, options ...Option) *D[*grpc.ClientConn] {
	resolve := func(o OptSet) (*grpc.ClientConn, error) {
		logDecider := func(_ string, err error) bool { return o.IsDebug() || err != nil }

		dialOptions = append(dialOptions,
			grpc.WithUnaryInterceptor(grpc_zap.UnaryClientInterceptor(o.Log(), grpc_zap.WithDecider(logDecider))),
			grpc.WithStreamInterceptor(grpc_zap.StreamClientInterceptor(o.Log(), grpc_zap.WithDecider(logDecider))),
		)

		return grpc.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), dialOptions...)
	}

	return New(resolve, options...)
}
