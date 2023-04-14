package tdep

import (
	"fmt"

	"github.com/heffcodex/redix"
)

func NewRedix(config redix.Config, options ...Option) *D[*redix.Client] {
	resolve := func(o OptSet) (*redix.Client, error) {
		if env := o.Env(); !env.IsEmpty() {
			config.AppendNamespace(env.String())
		}

		client, err := redix.NewClient(config)
		if err != nil {
			return nil, fmt.Errorf("new client: %w", err)
		}

		return client, nil
	}

	return New(resolve, options...)
}
