package tdep

import (
	"github.com/heffcodex/theapp/tcfg"
	"go.uber.org/zap"
)

type OptSet struct {
	env       tcfg.Env
	singleton bool
	debug     bool
	debugLog  *zap.Logger
}

func newOptSet(options ...Option) OptSet {
	opts := OptSet{
		debugLog: zap.NewNop(),
	}

	for _, opt := range options {
		opt(&opts)
	}

	return opts
}

func (o *OptSet) Env() tcfg.Env            { return o.env }
func (o *OptSet) IsSingleton() bool        { return o.singleton }
func (o *OptSet) IsDebug() bool            { return o.debug }
func (o *OptSet) DebugLogger() *zap.Logger { return o.debugLog }

type Option func(*OptSet)

func Env(env tcfg.Env) Option {
	return func(o *OptSet) {
		o.env = env
	}
}

func Singleton() Option {
	return func(o *OptSet) {
		o.singleton = true
	}
}

func Debug(enable bool, l *zap.Logger) Option {
	return func(o *OptSet) {
		o.debug = enable
		o.debugLog = l
	}
}
