package theapp

import (
	"fmt"

	"github.com/heffcodex/zapex"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/heffcodex/theapp/tcfg"
)

type NewAppFn[C tcfg.IConfig] func() (IApp[C], error)

type Cmd[C tcfg.IConfig] struct {
	newAppFn NewAppFn[C]
	opts     []CmdOption
	commands []*cobra.Command
}

func NewCmd[C tcfg.IConfig](newAppFn NewAppFn[C], opts ...CmdOption) *Cmd[C] {
	return &Cmd[C]{
		newAppFn: newAppFn,
		opts:     opts,
	}
}

func (c *Cmd[C]) Add(commands ...*cobra.Command) {
	c.commands = append(c.commands, commands...)
}

func (c *Cmd[C]) Execute() error {
	defer zapex.OnRecover(func(err error) { zapex.Default().Fatal("panic", zap.Error(err)) })()

	shut := newShutter()
	root := c.makeRoot(shut)

	if err := root.Execute(); err != nil {
		shut.shutdown()
		return err
	}

	return nil
}

func (c *Cmd[C]) makeRoot(shut *shutter) *cobra.Command {
	root := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			app, err := c.newAppFn()
			if err != nil {
				return fmt.Errorf("create app: %w", err)
			}

			cancelFn := cmdInject(cmd, app, shut)
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
