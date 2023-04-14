package tdep

import (
	"go.uber.org/zap"

	"github.com/heffcodex/theapp/tcfg"
)

type OptSet struct {
	env       tcfg.Env
	singleton bool
	debug     bool
	log       *zap.Logger
}

func newOptSet(options ...Option) OptSet {
	opts := OptSet{
		log: zap.NewNop(),
	}

	for _, opt := range options {
		opt(&opts)
	}

	return opts
}

func (o *OptSet) Env() tcfg.Env     { return o.env }
func (o *OptSet) IsSingleton() bool { return o.singleton }
func (o *OptSet) IsDebug() bool     { return o.debug }
func (o *OptSet) Log() *zap.Logger  { return o.log }

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

func Debug(enable ...bool) Option {
	enabled := len(enable) == 0 || enable[0]

	return func(o *OptSet) {
		o.debug = enabled
	}
}

func Log(log *zap.Logger) Option {
	return func(o *OptSet) {
		o.log = log
	}
}
