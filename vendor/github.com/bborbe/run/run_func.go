// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import "context"

// Func represents a function that can be executed with a context.
// It forms the foundation of this package, representing any operation
// that may fail and supports cancellation through context.
type Func func(context.Context) error

// Run executes the function with the provided context.
func (r Func) Run(ctx context.Context) error {
	return r(ctx)
}

//counterfeiter:generate -o mocks/runnable.go --fake-name Runnable . Runnable

// Runnable represents any object that can execute an operation with a context.
// This interface allows for more complex execution patterns where state
// or configuration needs to be maintained between calls.
type Runnable interface {
	Run(ctx context.Context) error
}
