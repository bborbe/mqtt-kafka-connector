// Copyright (c) 2020 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import (
	"context"
	"time"
)

// Retry on error n times and wait between the given delay.
func Retry(fn Func, limit int, delay time.Duration) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		i := 0
		for {
			i++
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if err := fn(ctx); i > limit || err == nil {
					return err
				}
				if delay == 0 {
					continue
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(delay):
				}
			}
		}
	}
}
