package tdep

import (
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strings"
)

type DRedis *D[*Redis]

type RedisConfig struct {
	DSN      string `mapstructure:"dsn"`
	Cert     string `mapstructure:"cert"`
	KeyGroup string `mapstructure:"keyGroup"`
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

func NewRedis(cfg RedisConfig, options ...Option) DRedis {
	resolve := func(o OptSet) (*Redis, error) {
		opts, err := redis.ParseURL(cfg.DSN)
		if err != nil {
			return nil, fmt.Errorf("parse DSN: %w", err)
		}

		if opts.TLSConfig != nil {
			opts.TLSConfig.InsecureSkipVerify = true

			if cfg.Cert != "" {
				ca, err := os.ReadFile(cfg.Cert)
				if err != nil {
					return nil, fmt.Errorf("read cert: %w", err)
				}

				rootCAs := x509.NewCertPool()
				if !rootCAs.AppendCertsFromPEM(ca) {
					return nil, errors.New("append cert")
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
