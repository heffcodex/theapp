package tdep

import (
	"fmt"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

func NewGRPC(cfg GRPCConfig, dialOptions []grpc.DialOption, options ...Option) *D[*grpc.ClientConn] {
	resolve := func(o OptSet) (*grpc.ClientConn, error) {
		if o.IsDebug() {
			debugLog := o.DebugLogger().Named("grpc")
			debugLogDecider := func(string, error) bool { return true }
			debugLogLevelFunc := func(codes.Code) zapcore.Level { return zapcore.DebugLevel }

			dialOptions = append(dialOptions,
				grpc.WithUnaryInterceptor(grpc_zap.UnaryClientInterceptor(
					debugLog,
					grpc_zap.WithDecider(debugLogDecider),
					grpc_zap.WithLevels(debugLogLevelFunc),
				)),
				grpc.WithStreamInterceptor(grpc_zap.StreamClientInterceptor(
					debugLog,
					grpc_zap.WithDecider(debugLogDecider),
					grpc_zap.WithLevels(debugLogLevelFunc),
				)),
			)
		}

		return grpc.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), dialOptions...)
	}

	return New(resolve, options...)
}
