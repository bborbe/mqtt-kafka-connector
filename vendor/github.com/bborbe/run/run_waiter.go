// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import (
	"context"
	"time"

	"github.com/golang/glog"
)

//counterfeiter:generate -o mocks/waiter.go --fake-name Waiter . Waiter

// Waiter provides an interface for implementing delays with context cancellation support.
// This abstraction allows for testing and custom wait implementations.
type Waiter interface {
	Wait(ctx context.Context, wait time.Duration) error
}

// WaiterFunc is a function type that implements the Waiter interface.
// It allows converting simple functions into Waiter implementations.
type WaiterFunc func(ctx context.Context, wait time.Duration) error

func (w WaiterFunc) Wait(ctx context.Context, wait time.Duration) error {
	return w(ctx, wait)
}

// NewWaiter creates a default Waiter implementation that uses time.Timer for delays.
// The waiter respects context cancellation and logs the wait duration.
func NewWaiter() Waiter {
	return WaiterFunc(func(ctx context.Context, wait time.Duration) error {
		glog.V(3).Infof("sleep for %v", wait)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.NewTimer(wait).C:
			return nil
		}
	})
}
