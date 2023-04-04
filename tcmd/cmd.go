package tcmd

import (
	"fmt"

	"github.com/heffcodex/zapex"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/heffcodex/theapp"
	"github.com/heffcodex/theapp/tcfg"
)

type NewAppFn[C tcfg.IConfig, A theapp.IApp[C]] func() (A, error)

type Cmd[C tcfg.IConfig, A theapp.IApp[C]] struct {
	newAppFn NewAppFn[C, A]
	opts     []CmdOption
	commands []*cobra.Command
}

func NewCmd[C tcfg.IConfig, A theapp.IApp[C]](newAppFn NewAppFn[C, A], opts ...CmdOption) *Cmd[C, A] {
	return &Cmd[C, A]{
		newAppFn: newAppFn,
		opts:     opts,
	}
}

func (c *Cmd[C, A]) Add(commands ...*cobra.Command) {
	c.commands = append(c.commands, commands...)
}

func (c *Cmd[C, A]) Execute() error {
	defer zapex.OnRecover(func(err error) { zapex.Default().Fatal("panic", zap.Error(err)) })()

	shut := newShutter()
	root := c.makeRoot(shut)

	if err := root.Execute(); err != nil {
		shut.shutdown()
		return err
	}

	return nil
}

func (c *Cmd[C, A]) makeRoot(shut *shutter) *cobra.Command {
	root := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			app, err := c.newAppFn()
			if err != nil {
				return fmt.Errorf("create app: %w", err)
			}

			cancelFn := cmdInject[C, A](cmd, app, shut)
			timeout := app.Config().ShutdownTimeout()

			shut.setup(app.L(), cancelFn, app.Close, timeout)

			return nil
		},
		PersistentPostRun: func(*cobra.Command, []string) {
			shut.shutdown()
		},
	}

	// override cobra global defaults:
	cobra.EnableCommandSorting = false

	for _, opt := range c.opts {
		opt(root)
	}

	for _, cmd := range c.commands {
		root.AddCommand(cmd)
	}

	return root
}
