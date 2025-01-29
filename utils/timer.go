package utils

import (
	"errors"
	"time"

	"0xlayer/go-layer-mesh/utils/gopool"
)

func Timeout(task func() error, timeout time.Duration) error {
	ch := make(chan error, 1)
	gopool.Submit(func() {
		ch <- task()
		close(ch)
	})

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case err := <-ch:
		return err
	case <-timer.C:
		return errors.New("timeout")
	}
}
