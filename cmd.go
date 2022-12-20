package theapp

import (
	"github.com/heffcodex/zapex"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type NewAppFn func() (IApp, error)

type Cmd struct {
	newAppFn NewAppFn
	opts     []CmdOption
	commands []*cobra.Command
}

func NewCmd(newAppFn NewAppFn, opts ...CmdOption) *Cmd {
	return &Cmd{
		newAppFn: newAppFn,
		opts:     opts,
	}
}

func (c *Cmd) Add(commands ...*cobra.Command) {
	c.commands = append(c.commands, commands...)
}

func (c *Cmd) Execute() error {
	defer zapex.OnRecover(func(err error) { zapex.Default().Fatal("panic", zap.Error(err)) })()
	return c.makeRoot().Execute()
}

func (c *Cmd) makeRoot() *cobra.Command {
	shutter := newShutter()

	root := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			app, err := c.newAppFn()
			if err != nil {
				return errors.Wrap(err, "can't create app")
			}

			cancelFn := injectAppAndCancelIntoCmd(cmd, app)
			timeout := app.IConfig().ShutdownTimeout()

			shutter.setup(app.L(), cancelFn, timeout)
			go shutter.waitShutdown(app.Close)

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, _ []string) {
			shutter.shutdown(cmd.Context())
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
