package retry

import (
	"net"
	"time"
)

// Exponential calls fn up to maxTries number of times with an exponential
// back-off until fn returns nil. If fn returns a permanent net.Error,
// Exponential returns early. Exponential returns the last error (or nil, if fn
// didn't return an error).
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
