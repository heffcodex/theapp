package dep

import (
	"github.com/heffcodex/theapp/cfg"
	"go.uber.org/zap"
)

type OptSet struct {
	env       cfg.Env
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

func (o *OptSet) Env() cfg.Env             { return o.env }
func (o *OptSet) IsSingleton() bool        { return o.singleton }
func (o *OptSet) IsDebug() bool            { return o.debug }
func (o *OptSet) DebugLogger() *zap.Logger { return o.debugLog }

type Option func(*OptSet)

func Env(env cfg.Env) Option {
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
