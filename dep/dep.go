package dep

import (
	"context"
	"fmt"
	"io"
	"log"
	"reflect"
	"sync"
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

type ResolveFn[T any] func() (T, error)

type D[T any] struct {
	l       sync.Mutex
	name    string
	resolve ResolveFn[T]

	// updated in behaviour of Get(), MustGet() or Close()
	instance T
	resolved bool

	// controlled by options:
	singleton   bool
	debug       bool
	debugWriter io.Writer
}

func New[T any](resolve ResolveFn[T], options ...Option[T]) *D[T] {
	tof := reflect.TypeOf(new(T)).Elem()
	if tof.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("type `%s` is not a pointer", tof.String()))
	}

	d := &D[T]{
		name:        fmt.Sprintf("dep(%s)", tof.String()),
		resolve:     resolve,
		debugWriter: log.Writer(),
	}

	for _, opt := range options {
		opt(d)
	}

	return d
}

func (d *D[T]) Get() (T, error) {
	if d == nil {
		panic("nil dep")
	}

	d.l.Lock()
	defer d.l.Unlock()

	if !d.singleton || !d.resolved {
		instance, err := d.resolve()
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
	_, _ = d.debugWriter.Write([]byte(d.name + ": " + msg + "\n"))
}
