package tdep

import (
	"strconv"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
)

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

func NewGRPC(cfg GRPCConfig, dialOptions []grpc.DialOption, options ...Option) *D[*grpc.ClientConn] {
	resolve := func(o OptSet) (*grpc.ClientConn, error) {
		target := cfg.Host + ":" + strconv.FormatInt(int64(cfg.Port), 10)

		logDecider := func(_ string, err error) bool { return o.IsDebug() || err != nil }
		withLogDecider := grpc_zap.WithDecider(logDecider)

		dialOptions = append(dialOptions,
			grpc.WithUserAgent(o.Name()),
			grpc.WithUnaryInterceptor(grpc_zap.UnaryClientInterceptor(o.Log(), withLogDecider)),
			grpc.WithStreamInterceptor(grpc_zap.StreamClientInterceptor(o.Log(), withLogDecider)),
		)

		return grpc.Dial(target, dialOptions...)
	}

	return New(resolve, options...)
}
