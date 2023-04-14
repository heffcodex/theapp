package tdep

import (
	"fmt"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"

	"github.com/heffcodex/theapp/tcfg"
)

func NewGRPC(cfg tcfg.GRPCClient, dialOptions []grpc.DialOption, options ...Option) *D[*grpc.ClientConn] {
	resolve := func(o OptSet) (*grpc.ClientConn, error) {
		log := o.Log().Named("grpc")
		logDecider := func(_ string, err error) bool { return o.IsDebug() || err != nil }

		dialOptions = append(dialOptions,
			grpc.WithUnaryInterceptor(grpc_zap.UnaryClientInterceptor(log, grpc_zap.WithDecider(logDecider))),
			grpc.WithStreamInterceptor(grpc_zap.StreamClientInterceptor(log, grpc_zap.WithDecider(logDecider))),
		)

		return grpc.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), dialOptions...)
	}

	return New(resolve, options...)
}
