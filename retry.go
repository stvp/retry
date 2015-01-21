package retry

import (
	"net"
	"time"
)

// Do runs a function up to maxTries number of times with an exponential
// back-off. It returns when: the function doesn't return an error, the
// function returns an error after maxTries number of attempts, or immediately
// if the function returns a non-temporary net.Error.
func Exponential(maxTries int, sleep time.Duration, fn func() error) (err error) {
	for try := 0; try < maxTries; try++ {
		err = fn()
		if err == nil {
			break
		}

		// Don't retry permanent network errors
		nerr, ok := err.(net.Error)
		if ok && !nerr.Temporary() {
			break
		}

		// Back off
		time.Sleep(sleep)
		sleep = 2 * sleep
	}

	return err
}
