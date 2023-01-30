package dep

import (
	"github.com/heffcodex/theapp/cfg"
	"go.uber.org/zap"
)

type OptSet struct {
	keyEnv    cfg.Env
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

func (o *OptSet) KeyEnv() cfg.Env          { return o.keyEnv }
func (o *OptSet) IsSingleton() bool        { return o.singleton }
func (o *OptSet) IsDebug() bool            { return o.debug }
func (o *OptSet) DebugLogger() *zap.Logger { return o.debugLog }

type Option func(*OptSet)

func KeyEnv(env cfg.Env) Option {
	return func(o *OptSet) {
		o.keyEnv = env
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
