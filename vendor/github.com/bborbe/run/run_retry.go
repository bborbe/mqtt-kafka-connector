// Copyright (c) 2020 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import (
	"context"
	"time"

	"github.com/bborbe/errors"
)

// DefaultWaiter is the default waiter implementation used by the Retry function.
var DefaultWaiter = NewWaiter()

// Backoff configures retry behavior including delays, retry counts, and retry conditions.
type Backoff struct {
	// Delay is the initial delay to wait before the first retry.
	Delay time.Duration `json:"delay"`
	// Factor is the multiplier applied to the delay for each subsequent retry.
	Factor float64 `json:"factor"`
	// Retries is the maximum number of retry attempts.
	Retries int `json:"retries"`
	// IsRetryAble is an optional function that determines if an error is retryable.
	// If nil, all errors are considered retryable.
	IsRetryAble func(error) bool `json:"-"`
}

// Retry wraps a function with retry logic using the specified backoff configuration.
// It uses the DefaultWaiter for delays between retry attempts.
func Retry(backoff Backoff, fn Func) Func {
	return RetryWaiter(backoff, DefaultWaiter, fn)
}

// RetryWaiter wraps a function with retry logic using the specified backoff configuration and custom waiter.
// The waiter controls how delays are implemented, allowing for custom timing behavior.
func RetryWaiter(backoff Backoff, waiter Waiter, fn Func) Func {
	return func(ctx context.Context) error {
		var counter int
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if err := fn(ctx); err != nil {
					if counter == backoff.Retries {
						return errors.Wrapf(ctx, err, "reached try counter(%d)", backoff.Retries)
					}
					if backoff.IsRetryAble != nil && !backoff.IsRetryAble(err) {
						return errors.Wrap(ctx, err, "error is not retryable")
					}
					counter++
					if backoff.Delay > 0 {
						delay := backoff.Delay + backoff.Delay*time.Duration(
							backoff.Factor*float64(counter-1),
						)
						if err := waiter.Wait(ctx, delay); err != nil {
							return errors.Wrapf(ctx, err, "wait %v failed", backoff.Delay)
						}
					}
					continue
				}
				return nil
			}
		}
	}
}
