package retry

import (
	"testing"
)

type regularError struct{}

func (e regularError) Error() string { return "" }

type tempNetError struct{ error }

func (e tempNetError) Temporary() bool { return true }
func (e tempNetError) Timeout() bool   { return true }

type permNetError struct{ error }

func (e permNetError) Temporary() bool { return false }
func (e permNetError) Timeout() bool   { return false }

type testcase struct {
	max         int
	err         func(try int) error
	expectTries int
	expectError bool
}

func TestExponential(t *testing.T) {
	testcases := []testcase{
		// When an error isn't returned
		testcase{
			5,
			func(try int) error { return nil },
			1,
			false,
		},
		// When an error is always returned
		testcase{
			5,
			func(try int) error { return regularError{} },
			5,
			true,
		},
		// When a temporary net error is always returned
		testcase{
			5,
			func(try int) error { return tempNetError{} },
			5,
			true,
		},
		// When a permanent net error is always returned
		testcase{
			5,
			func(try int) error { return permNetError{} },
			1,
			true,
		},
		// When an error is returned only a couple times
		testcase{
			5,
			func(try int) error {
				if try == 3 {
					return nil
				} else {
					return regularError{}
				}
			},
			3,
			false,
		},
	}

	for i, test := range testcases {
		actual := 0
		err := Exponential(test.max, 0, func() error {
			actual += 1
			return test.err(actual)
		})
		if actual != test.expectTries {
			t.Errorf("testcases[%d]: tried %d times", i, actual)
		}
		if err == nil && test.expectError {
			t.Errorf("testcases[%d]: didn't get error", i)
		} else if err != nil && !test.expectError {
			t.Errorf("testcases[%d]: got error", i)
		}
	}
}
