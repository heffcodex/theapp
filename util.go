package theapp

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/heffcodex/theapp/tcfg"
)

const (
	AppCmdContextKey     = "theapp.app"
	ShutterCmdContextKey = "theapp.shutter"
)

func GetApp[C tcfg.IConfig, A IApp[C]](cmd *cobra.Command) A {
	return getApp[C, A](cmd)
}

func getApp[C tcfg.IConfig, A IApp[C]](cmd *cobra.Command) A {
	return cmd.Context().Value(AppCmdContextKey).(A)
}

func WaitInterrupt(cmd *cobra.Command) {
	getShutter(cmd).waitInterrupt()
}

func getShutter(cmd *cobra.Command) *shutter {
	return cmd.Context().Value(ShutterCmdContextKey).(*shutter)
}

func cmdInject[C tcfg.IConfig, A IApp[C]](cmd *cobra.Command, app A, shut *shutter) (cancel context.CancelFunc) {
	ctx := cmd.Context()

	ctx = context.WithValue(ctx, AppCmdContextKey, app)
	ctx = context.WithValue(ctx, ShutterCmdContextKey, shut)

	ctx, cancel = context.WithCancel(ctx)
	cmd.SetContext(ctx)

	return cancel
}
