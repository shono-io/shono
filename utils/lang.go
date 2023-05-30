package utils

import (
	"github.com/sirupsen/logrus"
	"time"
)

func MustReturn[T any](res T, err error) T {
	if err != nil {
		panic(err)
	}

	return res
}

func Must(e ...any) {
	if len(e) == 0 {
		return
	}

	if err, ok := e[len(e)-1].(error); ok {
		if err != nil {
			panic(err)
		}
	}
}

func WaitFor[T any](fn func() *T, timeout time.Duration) (*T, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-timeoutTimer.C:
			return nil, ErrTimeout
		case <-ticker.C:
			if res := fn(); res != nil {
				return res, nil
			}
		}
	}
}

func LogIfError(level logrus.Level, err error) {
	if err != nil {
		switch level {
		case logrus.PanicLevel:
			logrus.Panic(err)
		case logrus.ErrorLevel:
			logrus.Error(err)
		case logrus.WarnLevel:
			logrus.Warn(err)
		case logrus.InfoLevel:
			logrus.Info(err)
		case logrus.DebugLevel:
			logrus.Debug(err)
		default:
			logrus.Trace(err)
		}
	}
}
