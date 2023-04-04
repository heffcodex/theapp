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

func GetApp[C tcfg.IConfig](cmd *cobra.Command) IApp[C] {
	return getApp[C](cmd)
}

func getApp[C tcfg.IConfig](cmd *cobra.Command) IApp[C] {
	return cmd.Context().Value(AppCmdContextKey).(IApp[C])
}

func WaitInterrupt(cmd *cobra.Command) {
	getShutter(cmd).waitInterrupt()
}

func getShutter(cmd *cobra.Command) *shutter {
	return cmd.Context().Value(ShutterCmdContextKey).(*shutter)
}

func cmdInject[C tcfg.IConfig](cmd *cobra.Command, app IApp[C], shut *shutter) (cancel context.CancelFunc) {
	ctx := cmd.Context()

	ctx = context.WithValue(ctx, AppCmdContextKey, app)
	ctx = context.WithValue(ctx, ShutterCmdContextKey, shut)

	ctx, cancel = context.WithCancel(ctx)
	cmd.SetContext(ctx)

	return cancel
}
