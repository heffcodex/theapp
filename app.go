package theapp

import (
	"context"
	"fmt"
	"github.com/heffcodex/zapex"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type CloseFn func(context.Context) error

type IApp interface {
	IConfig() IConfig
	IsDebug() bool
	L() *zap.Logger
	AddCloser(fns ...CloseFn)
	Close(ctx context.Context) error
}

var _ IApp = (*App)(nil)

type App struct {
	cfg       IConfig
	closeFns  []CloseFn
	closeLock sync.Mutex
	log       *zap.Logger
}

func NewApp(cfg IConfig) (*App, error) {
	log, err := zapex.New(cfg.LogLevel())
	if err != nil {
		return nil, errors.Wrap(err, "can't create logger")
	}

	appLog := log.Named(cfg.AppName())
	maxprocsLog := appLog.Named("maxprocs")

	_, err = maxprocs.Set(
		maxprocs.Logger(
			func(format string, args ...any) { maxprocsLog.Debug(fmt.Sprintf(format, args...)) },
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "can't set maxprocs")
	}

	return &App{
		cfg: cfg,
		log: appLog,
	}, nil
}

func (a *App) IConfig() IConfig { return a.cfg }
func (a *App) IsDebug() bool    { return a.cfg.LogLevel() == zap.DebugLevel.String() }

func (a *App) L() *zap.Logger { return a.log }

func (a *App) AddCloser(fns ...CloseFn) {
	a.closeLock.Lock()
	defer a.closeLock.Unlock()

	a.closeFns = append(a.closeFns, fns...)
}

func (a *App) Close(ctx context.Context) error {
	a.closeLock.Lock()
	defer a.closeLock.Unlock()

	errs := make([]error, 0, len(a.closeFns))

	for i := len(a.closeFns) - 1; i >= 0; i-- {
		if err := a.closeFns[i](ctx); err != nil {
			errs = append(errs, err)
		}
	}

	return multierr.Combine(errs...)
}
