package tcmd

import (
	"fmt"

	"github.com/heffcodex/zapex"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/heffcodex/theapp"
	"github.com/heffcodex/theapp/tcfg"
)

type NewApp[C tcfg.Config, A theapp.App[C]] func() (A, error)

type Cmd[C tcfg.Config, A theapp.App[C]] struct {
	newApp   NewApp[C, A]
	opts     []CmdOption
	commands []*cobra.Command
}

func New[C tcfg.Config, A theapp.App[C]](newApp NewApp[C, A], opts ...CmdOption) *Cmd[C, A] {
	return &Cmd[C, A]{
		newApp: newApp,
		opts:   opts,
	}
}

func (c *Cmd[C, A]) Add(commands ...*cobra.Command) {
	c.commands = append(c.commands, commands...)
}

func (c *Cmd[C, A]) Execute() error {
	defer zapex.OnRecover(func(err error) { zapex.Default().Fatal("panic", zap.Error(err)) })() //nolint: revive // it's ok

	shut := newShutter()
	root := c.makeRoot(shut)

	if err := root.Execute(); err != nil {
		shut.down()
		return err
	}

	return nil
}

func (c *Cmd[C, A]) makeRoot(shut *shutter) *cobra.Command {
	root := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			app, err := c.newApp()
			if err != nil {
				return fmt.Errorf("new app: %w", err)
			}

			cancelFn := cmdInject[C, A](cmd, app, shut)
			timeout := app.Config().ShutdownTimeout()

			shut.setup(app.L(), cancelFn, app.Close, timeout)
			go func() {
				shut.waitInterrupt()
				shut.down()
			}()

			return nil
		},
		PersistentPostRun: func(*cobra.Command, []string) {
			shut.down()
		},
	}

	// override cobra global defaults:
	cobra.EnableCommandSorting = false //nolint: reassign // it's ok

	for _, opt := range c.opts {
		opt(root)
	}

	for _, cmd := range c.commands {
		root.AddCommand(cmd)
	}

	return root
}
