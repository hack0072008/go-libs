package common

import (
	"time"

	"github.com/pkg/errors"
)

/*
为执行函数设置超时时间
 */
func DoWithTimeout(handler func() error, timeout int) error {

	ch := make(chan error, 1)
	go func() {
		ch <- handler()
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		return errors.WithMessage(errors.New("connection timed out"), "connect fail")
	case err := <-ch:
		return err
	}
}
