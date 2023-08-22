package tcmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/heffcodex/theapp"
	"github.com/heffcodex/theapp/tcfg"
)

type (
	appKey struct{}
)

func App[C tcfg.Config, A theapp.App[C]](cmd *cobra.Command) A {
	return getApp[C, A](cmd)
}

func getApp[C tcfg.Config, A theapp.App[C]](cmd *cobra.Command) A {
	return cmd.Context().Value(appKey{}).(A)
}

func cmdInject[C tcfg.Config, A theapp.App[C]](cmd *cobra.Command, app A) (cancel context.CancelFunc) {
	ctx := context.WithValue(cmd.Context(), appKey{}, app)

	ctx, cancel = context.WithCancel(ctx)
	cmd.SetContext(ctx)

	return cancel
}
