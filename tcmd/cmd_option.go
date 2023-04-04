package tcmd

import (
	"github.com/spf13/cobra"
)

type CmdOption func(cmd *cobra.Command)

func WithCmdSilenceErrors(v bool) CmdOption {
	return func(cmd *cobra.Command) {
		cmd.SilenceErrors = v
	}
}

func WithCmdSilenceUsage(v bool) CmdOption {
	return func(cmd *cobra.Command) {
		cmd.SilenceUsage = v
	}
}

func WithGlobCmdSorting(v bool) CmdOption {
	return func(_ *cobra.Command) {
		cobra.EnableCommandSorting = v
	}
}
