package theapp

import (
	"context"
	"github.com/spf13/cobra"
)

const (
	AppCmdContextKey     = "theapp.app"
	ShutterCmdContextKey = "theapp.shutter"
)

func GetApp[A IApp](cmd *cobra.Command) A {
	return getApp[A](cmd)
}

func getApp[A IApp](cmd *cobra.Command) A {
	return cmd.Context().Value(AppCmdContextKey).(A)
}

func WaitShutdown(cmd *cobra.Command) {
	getShutter(cmd).waitShutdownComplete(cmd.Context())
}

func getShutter(cmd *cobra.Command) *shutter {
	return cmd.Context().Value(ShutterCmdContextKey).(*shutter)
}

func cmdInject[A IApp](cmd *cobra.Command, app A, shutter *shutter) context.CancelFunc {
	ctx, cancel := context.WithCancel(cmd.Context())

	ctx = context.WithValue(ctx, AppCmdContextKey, app)
	ctx = context.WithValue(ctx, ShutterCmdContextKey, shutter)

	cmd.SetContext(ctx)

	return cancel
}
