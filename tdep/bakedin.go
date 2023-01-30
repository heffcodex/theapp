package tdep

import (
	"crypto/x509"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"os"
	"strings"
)

// Bun

type BunConfigPostgres struct {
	DSN            string `mapstructure:"dsn"`
	MaxConnections int    `mapstructure:"max_connections"`
}

func NewBunPostgres(
	cfg BunConfigPostgres,
	onTuneConnector func(conn *pgdriver.Connector),
	onTuneSQLDB func(db *sql.DB),
	onTuneBunDB func(db *bun.DB),
	options ...Option,
) *D[*bun.DB] {
	resolve := func(o OptSet) (*bun.DB, error) {
		conn := pgdriver.NewConnector(pgdriver.WithDSN(cfg.DSN))
		if onTuneConnector != nil {
			onTuneConnector(conn)
		}

		sqlDB := sql.OpenDB(conn)
		sqlDB.SetMaxOpenConns(cfg.MaxConnections)
		if onTuneSQLDB != nil {
			onTuneSQLDB(sqlDB)
		}

		bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())
		if onTuneBunDB != nil {
			onTuneBunDB(bunDB)
		}

		bunLogger, _ := zap.NewStdLogAt(o.debugLog.Named("bun"), zap.DebugLevel)
		bunDB.AddQueryHook(
			bundebug.NewQueryHook(
				bundebug.WithEnabled(o.IsDebug()),
				bundebug.WithVerbose(o.IsDebug()),
				bundebug.WithWriter(bunLogger.Writer()),
			),
		)

		return bunDB, nil
	}

	return New(resolve, options...)
}

// GRPC Client

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

// Redis Client

type RedisConfig struct {
	DSN      string `mapstructure:"dsn"`
	Cert     string `mapstructure:"cert"`
	KeyGroup string `mapstructure:"key_group"`
}

type Redis struct {
	*redis.Client
	keyGroup string
}

func (r *Redis) KeyPrefix() string {
	return r.key("")
}

func (r *Redis) Key(parts ...string) string {
	return r.key(parts...)
}

func (r *Redis) key(parts ...string) string {
	return strings.Join(append([]string{r.keyGroup}, parts...), ":")
}

func NewRedis(cfg RedisConfig, options ...Option) *D[*Redis] {
	resolve := func(o OptSet) (*Redis, error) {
		opts, err := redis.ParseURL(cfg.DSN)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse DSN as options")
		}

		if opts.TLSConfig != nil {
			opts.TLSConfig.InsecureSkipVerify = true

			if cfg.Cert != "" {
				ca, err := os.ReadFile(cfg.Cert)
				if err != nil {
					return nil, errors.Wrap(err, "can't read root CA")
				}

				rootCAs := x509.NewCertPool()
				if !rootCAs.AppendCertsFromPEM(ca) {
					return nil, errors.New("can't append root CA")
				}

				opts.TLSConfig.InsecureSkipVerify = false
				opts.TLSConfig.RootCAs = rootCAs
			}
		}

		keyGroup := cfg.KeyGroup
		if o.env != "" {
			keyGroup = keyGroup + ":" + o.env.String()
		}

		return &Redis{
			Client:   redis.NewClient(opts),
			keyGroup: keyGroup,
		}, nil
	}

	return New(resolve, options...)
}
