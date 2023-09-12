package tdep

import (
	"slices"

	"github.com/heffcodex/redix"
)

func NewRedix(config redix.Config, options ...Option) *D[*redix.Client] { //nolint: gocritic // hugeParam: it's ok
	resolve := func(o OptSet) (*redix.Client, error) {
		_config := redix.Config{
			Name:      config.Name,
			Namespace: config.Namespace,
			DSN:       config.DSN,
			Cert: redix.ConfigCert{
				Env:  config.Cert.Env,
				File: config.Cert.File,
				Data: slices.Clone(config.Cert.Data),
			},
		}

		if _config.Name == "" {
			_config.Name = o.Name()
		}

		if _config.Namespace == "" {
			_config.Namespace = redix.Namespace(o.Name())

			if env := o.Env(); !env.IsEmpty() {
				_config.Namespace = _config.Namespace.Append(env.String())
			}
		}

		return redix.NewClient(_config)
	}

	return New(resolve, options...)
}
