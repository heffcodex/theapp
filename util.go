package theapp

import (
	"context"
	"github.com/spf13/cobra"
)

const (
	AppCmdContextKey = "app"
)

func GetApp[A IApp](cmd *cobra.Command) A {
	return getApp[A](cmd)
}

func getApp[A IApp](cmd *cobra.Command) A {
	return cmd.Context().Value(AppCmdContextKey).(A)
}

func injectAppAndCancelIntoCmd[A IApp](cmd *cobra.Command, app A) context.CancelFunc {
	ctx, cancel := context.WithCancel(cmd.Context())
	ctx = context.WithValue(ctx, AppCmdContextKey, app)

	cmd.SetContext(ctx)

	return cancel
}
