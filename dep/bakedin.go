package dep

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"google.golang.org/grpc"
	"net/url"
	"os"
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
	options ...Option[*bun.DB],
) *D[*bun.DB] {
	resolve := func() (*bun.DB, error) {
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

		return bunDB, nil
	}

	return New(resolve, options...)
}

// GRPC Client

type GRPCConfig struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

func NewGRPC(cfg GRPCConfig, dialOptions []grpc.DialOption, options ...Option[*grpc.ClientConn]) *D[*grpc.ClientConn] {
	resolve := func() (*grpc.ClientConn, error) {
		return grpc.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), dialOptions...)
	}

	return New(resolve, options...)
}

// Redis Client

type RedisConfig struct {
	DSN        string `mapstructure:"dsn"`
	Cert       string `mapstructure:"cert"`
	KeysPrefix string `mapstructure:"keys_prefix"`
}

type Redis struct {
	*redis.Client
	keysPrefix string
}

func (r *Redis) KeysPrefix() string {
	return r.keysPrefix
}

func NewRedis(cfg RedisConfig, options ...Option[*Redis]) *D[*Redis] {
	resolve := func() (*Redis, error) {
		dsnURL, err := url.Parse(cfg.DSN)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse DSN as URL")
		}

		opts, err := redis.ParseURL(cfg.DSN)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse DSN as options")
		}

		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
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
			opts.TLSConfig.ServerName = dsnURL.Hostname()
		}

		return &Redis{
			Client:     redis.NewClient(opts),
			keysPrefix: cfg.KeysPrefix,
		}, nil
	}

	return New(resolve, options...)
}
