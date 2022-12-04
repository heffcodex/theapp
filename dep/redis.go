package dep

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"net/url"
	"os"
)

type RedisConfig struct {
	DSN        string `mapstructure:"dsn"`
	Cert       string `mapstructure:"cert"`
	KeysPrefix string `mapstructure:"keys_prefix"`
}

type RedisClient struct {
	*redis.Client
	keysPrefix string
}

func (rc *RedisClient) KeysPrefix() string {
	return rc.keysPrefix
}

type Redis Dep[*RedisClient]

func NewRedis(cfg RedisConfig) *Redis {
	resolve := func() (*RedisClient, error) {
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

		return &RedisClient{
			Client:     redis.NewClient(opts),
			keysPrefix: cfg.KeysPrefix,
		}, nil
	}

	return (*Redis)(NewDep(true, resolve))
}
