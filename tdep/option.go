package tdep

import (
	"go.uber.org/zap"

	"github.com/heffcodex/theapp/tcfg"
)

type OptSet struct {
	name      string
	env       tcfg.Env
	singleton bool
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

func (o *OptSet) Name() string      { return o.name }
func (o *OptSet) Env() tcfg.Env     { return o.env }
func (o *OptSet) IsSingleton() bool { return o.singleton }
func (o *OptSet) Log() *zap.Logger  { return o.log }
func (o *OptSet) IsDebug() bool     { return o.log == nil || o.log.Core().Enabled(zap.DebugLevel) }

type Option func(*OptSet)

func Name(name string) Option {
	return func(o *OptSet) {
		o.name = name
	}
}

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

func Log(log *zap.Logger) Option {
	return func(o *OptSet) {
		o.log = log
	}
}
