package theapp

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type shutter struct {
	log          *zap.Logger
	cancelFn     context.CancelFunc
	timeout      time.Duration
	signals      []os.Signal
	notifyChan   chan os.Signal
	completeChan chan struct{}
}

func newShutter(signals ...os.Signal) *shutter {
	if len(signals) == 0 {
		signals = []os.Signal{os.Interrupt, syscall.SIGTERM}
	}

	return &shutter{
		signals:      signals,
		notifyChan:   make(chan os.Signal, len(signals)),
		completeChan: make(chan struct{}),
	}
}

func (s *shutter) setup(log *zap.Logger, cancelFn context.CancelFunc, timeout time.Duration) *shutter {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	s.log = log
	s.cancelFn = cancelFn
	s.timeout = timeout

	return s
}

func (s *shutter) waitShutdown(execCloser CloserFn) {
	signal.Notify(s.notifyChan, s.signals...)
	<-s.notifyChan

	s.log.Info("shutdown started", zap.Duration("timeout", s.timeout))
	s.cancelFn()

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer func() { cancel(); s.completeChan <- struct{}{} }()

	err := (error)(nil)
	shutdownOK := make(chan any)

	go func() {
		defer close(shutdownOK)
		if execCloser != nil {
			err = execCloser(ctx)
		}
	}()

	select {
	case <-ctx.Done():
		s.log.Fatal("shutdown timeout")
	case <-shutdownOK:
		//
	}

	if err != nil {
		s.log.Fatal("shutdown error", zap.Error(err))
	}

	s.log.Info("shutdown complete")
}

func (s *shutter) shutdown(waitCtx context.Context) {
	s.log.Debug("shutdown triggered")

	signal.Stop(s.notifyChan)
	close(s.notifyChan)

	<-waitCtx.Done()
	<-s.completeChan

	s.log.Debug("app exit")
	_ = s.log.Sync()

	os.Exit(0)
}
