package theapp

import (
	"context"
	"fmt"
	"sync"

	"github.com/heffcodex/zapex"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/heffcodex/theapp/tcfg"
)

type CloseFn func(context.Context) error

type App[C tcfg.IConfig] struct {
	cfg       C
	closeFns  []CloseFn
	closeLock sync.Mutex
	log       *zap.Logger
}

func New[C tcfg.IConfig]() (*App[C], error) {
	var cfg C
	if err := cfg.Load(); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

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

	return &App[C]{
		cfg: cfg,
		log: appLog,
	}, nil
}

func (a *App[C]) Config() C      { return a.cfg }
func (a *App[C]) IsDebug() bool  { return a.cfg.LogLevel() == zap.DebugLevel.String() }
func (a *App[C]) L() *zap.Logger { return a.log }

func (a *App[C]) AddCloser(fns ...CloseFn) {
	a.closeLock.Lock()
	a.closeFns = append(a.closeFns, fns...)
	a.closeLock.Unlock()
}

func (a *App[C]) Close(ctx context.Context) error {
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
