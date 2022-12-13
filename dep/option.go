package dep

import "io"

type Option[T any] func(*D[T])

func Singleton[T any]() Option[T] {
	return func(d *D[T]) {
		d.singleton = true
	}
}

func Debug[T any](debug bool) Option[T] {
	return func(d *D[T]) {
		d.debug = debug
	}
}

func DebugWriter[T any](w io.Writer) Option[T] {
	return func(d *D[T]) {
		d.debugWriter = w
	}
}
