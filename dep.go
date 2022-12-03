package theapp

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"go.uber.org/zap"
)

var (
	closerT    = reflect.TypeOf((*closer)(nil)).Elem()
	ctxCloserT = reflect.TypeOf((*ctxCloser)(nil)).Elem()
)

type closer interface {
	Close() error
}
type ctxCloser interface {
	Close(ctx context.Context) error
}

type ResolveFn[T any] func() T

type Dep[T any] struct {
	l         sync.Mutex
	log       *zap.Logger
	instance  T
	resolve   ResolveFn[T]
	resolved  bool
	singleton bool
}

func NewDep[T any](resolve ResolveFn[T], singleton bool, log *zap.Logger) *Dep[T] {
	tof := reflect.TypeOf(new(T)).Elem()
	if tof.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("type `%s` is not a pointer", tof.String()))
	}

	return &Dep[T]{
		resolve:   resolve,
		singleton: singleton,
		log:       log.Named(fmt.Sprintf("dep(%s)", tof.String())),
	}
}

func (d *Dep[T]) Get() T {
	if d == nil {
		panic("nil dep")
	}

	d.l.Lock()
	defer d.l.Unlock()

	if !d.singleton || !d.resolved {
		d.instance = d.resolve()
		d.resolved = true
		d.log.Debug("resolve")
	}

	return d.instance
}

func (d *Dep[T]) Close(ctx context.Context) error {
	if d == nil {
		return nil
	}

	d.l.Lock()
	defer d.l.Unlock()

	if !d.resolved {
		d.log.Debug("close (nop: unresolved)")
		return nil
	}

	defer func() {
		d.instance = *new(T)
		d.resolved = false
	}()

	vof := reflect.ValueOf(d.instance)

	if vof.CanConvert(closerT) {
		d.log.Debug("close (closerT)")
		return vof.Convert(closerT).Interface().(closer).Close()
	} else if vof.CanConvert(ctxCloserT) {
		d.log.Debug("close (ctxCloserT)")
		return vof.Convert(ctxCloserT).Interface().(ctxCloser).Close(ctx)
	}

	d.log.Debug("close (nop: no closer)")

	return nil
}
