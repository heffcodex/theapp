package theapp

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/heffcodex/zapex"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"

	"github.com/heffcodex/theapp/tcfg"
)

type CloseFn func(context.Context) error

type App[C tcfg.Config] interface {
	Config() C
	L() *zap.Logger
	AddCloser(fns ...CloseFn)
	Close(ctx context.Context) error
}

var _ App[tcfg.Config] = (*BaseApp[tcfg.Config])(nil)

type BaseApp[C tcfg.Config] struct {
	cfg C
	log *zap.Logger

	closed     bool
	closers    []CloseFn
	closerLock sync.Mutex
}

func New[C tcfg.Config]() (*BaseApp[C], error) {
	cfg, err := tcfg.LoadConfig[C]()
	if err != nil {
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

	return &BaseApp[C]{
		cfg: cfg,
		log: appLog,
	}, nil
}

func (a *BaseApp[C]) Config() C      { return a.cfg }
func (a *BaseApp[C]) L() *zap.Logger { return a.log }

func (a *BaseApp[C]) AddCloser(fns ...CloseFn) {
	_ = a.safeClose(func() error {
		a.closers = append(a.closers, fns...)
		return nil
	})
}

func (a *BaseApp[C]) Close(ctx context.Context) error {
	return a.safeClose(func() error {
		errs := make([]error, 0, len(a.closers))

		for i := len(a.closers) - 1; i >= 0; i-- {
			closer := a.closers[i]

			if err := closer(ctx); err != nil {
				errs = append(errs, err)
			}
		}

		a.closed = true

		return errors.Join(errs...)
	})
}

func (a *BaseApp[C]) safeClose(f func() error) error {
	a.closerLock.Lock()
	defer a.closerLock.Unlock()

	if a.closed {
		panic("app is closed")
	}

	return f()
}
