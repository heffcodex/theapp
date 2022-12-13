package dep

import (
	"go.uber.org/zap"
)

type Option[T any] func(*D[T])

func Singleton[T any]() Option[T] {
	return func(d *D[T]) {
		d.singleton = true
	}
}

func Debug[T any](enable bool, l *zap.Logger) Option[T] {
	return func(d *D[T]) {
		d.debug = enable
		d.debugLog = l
	}
}
