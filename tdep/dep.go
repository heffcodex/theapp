package tdep

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

type ResolveFn[T any] func(opts OptSet) (T, error)

type D[T any] struct {
	l       sync.Mutex
	typ     string
	resolve ResolveFn[T]
	opts    OptSet
	health  func(ctx context.Context, d *D[T]) error

	// updated in behaviour of Get(), MustGet() or Close()
	instance T
	resolved bool
}

func New[T any](resolve ResolveFn[T], options ...Option) *D[T] {
	tof := reflect.TypeOf(new(T)).Elem()
	tofStr := tof.String()

	if tof.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("type `%s` is not a pointer", tofStr))
	}

	d := &D[T]{
		typ:     tofStr,
		resolve: resolve,
		opts:    newOptSet(options...),
	}

	return d
}

func (d *D[T]) WithHealthCheck(fn func(ctx context.Context, d *D[T]) error) *D[T] {
	d.health = fn
	return d
}

func (d *D[T]) Options() OptSet {
	return d.opts
}

func (d *D[T]) Get() (T, error) {
	if d == nil {
		panic("nil dep")
	}

	d.l.Lock()
	defer d.l.Unlock()

	if !d.opts.singleton || !d.resolved {
		instance, err := d.resolve(d.opts)
		if err != nil {
			return *new(T), err
		}

		d.instance = instance
		d.resolved = true

		d.debugWrite("resolved")
	}

	return d.instance, nil
}

func (d *D[T]) MustGet() T {
	v, err := d.Get()
	if err != nil {
		panic(err)
	}

	return v
}

func (d *D[T]) Health(ctx context.Context) error {
	if d.health == nil {
		return nil
	}

	return d.health(ctx, d)
}

func (d *D[T]) Close(ctx context.Context) error {
	if d == nil {
		return nil
	}

	d.l.Lock()
	defer d.l.Unlock()

	if !d.resolved {
		d.debugWrite("close (nop: unresolved)")
		return nil
	}

	defer func() {
		d.instance = *new(T)
		d.resolved = false
	}()

	vof := reflect.ValueOf(d.instance)

	if vof.CanConvert(closerT) {
		d.debugWrite("close (closerT)")
		return vof.Convert(closerT).Interface().(closer).Close()
	} else if vof.CanConvert(ctxCloserT) {
		d.debugWrite("close (ctxCloserT)")
		return vof.Convert(ctxCloserT).Interface().(ctxCloser).Close(ctx)
	}

	d.debugWrite("close (nop: no closer)")

	return nil
}

func (d *D[T]) debugWrite(msg string) {
	if log := d.opts.Log(); log != nil && d.opts.IsDebug() {
		log.Debug(msg, zap.String("typ", d.typ))
	}
}
