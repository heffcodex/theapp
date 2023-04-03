package theapp

import (
	"context"
	"fmt"
	"github.com/heffcodex/theapp/tcfg"
	"github.com/heffcodex/zapex"
	"sync"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type CloseFn func(context.Context) error

type IApp interface {
	IConfig() tcfg.IConfig
	IsDebug() bool
	L() *zap.Logger
	AddCloser(fns ...CloseFn)
	Close(ctx context.Context) error
}

var _ IApp = (*App)(nil)

type App struct {
	cfg       tcfg.IConfig
	closeFns  []CloseFn
	closeLock sync.Mutex
	log       *zap.Logger
}

func NewApp(cfg tcfg.IConfig) (*App, error) {
	log, err := zapex.New(cfg.LogLevel())
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	zapex.SetDefault(log)

	appLog := log.Named(cfg.AppName())
	maxprocsLog := appLog.Named("maxprocs")

	_, err = maxprocs.Set(
		maxprocs.Logger(
			func(format string, args ...any) { maxprocsLog.Debug(fmt.Sprintf(format, args...)) },
		),
	)
	if err != nil {
		return nil, fmt.Errorf("set maxprocs: %w", err)
	}

	return &App{
		cfg: cfg,
		log: appLog,
	}, nil
}

func (a *App) IConfig() tcfg.IConfig { return a.cfg }
func (a *App) IsDebug() bool         { return a.cfg.LogLevel() == zap.DebugLevel.String() }

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
