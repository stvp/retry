// Package retry implements a basic retry loop with
// optional exponential back-off. It gracefully handles
// permanent net.Error errors.
//
//     err := retry.Exponential(10, time.Second, func() error {
//       err := someOperation()
//       return err
//     })
//
package retry
