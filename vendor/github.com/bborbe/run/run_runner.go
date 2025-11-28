// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import (
	"context"
	"runtime"
	"sync"
)

// CancelOnFirstFinish executes all given functions in parallel and cancels the remaining functions when the first one completes.
// It returns the error from the first function that finishes, or nil if that function succeeds.
func CancelOnFirstFinish(ctx context.Context, funcs ...Func) error {
	if len(funcs) == 0 {
		return nil
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return <-Run(ctx, funcs...)
}

// CancelOnFirstFinishWait executes all given functions in parallel and cancels the remaining functions when the first one completes.
// Unlike CancelOnFirstFinish, it waits for all functions to complete or be canceled and returns an aggregate error of all failures.
func CancelOnFirstFinishWait(ctx context.Context, funcs ...Func) error {
	if len(funcs) == 0 {
		return nil
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var errs []error
	for err := range Run(ctx, funcs...) {
		cancel()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return NewErrorList(errs...)
}

// CancelOnFirstError executes all given functions in parallel and cancels the remaining functions when the first error occurs.
// It returns the first error encountered, providing fail-fast behavior.
func CancelOnFirstError(ctx context.Context, funcs ...Func) error {
	if len(funcs) == 0 {
		return nil
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for err := range Run(ctx, funcs...) {
		if err != nil {
			return err
		}
	}
	return nil
}

// CancelOnFirstErrorWait executes all given functions in parallel and cancels the remaining functions when the first error occurs.
// Unlike CancelOnFirstError, it waits for all functions to complete or be canceled and returns an aggregate error of all failures.
func CancelOnFirstErrorWait(ctx context.Context, funcs ...Func) error {
	if len(funcs) == 0 {
		return nil
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var errs []error
	for err := range Run(ctx, funcs...) {
		if err != nil {
			cancel()
			errs = append(errs, err)
		}
	}
	return NewErrorList(errs...)
}

// All executes all given functions in parallel and waits for all to complete.
// It returns an aggregate error containing all errors that occurred during execution, or nil if all functions succeed.
func All(ctx context.Context, funcs ...Func) error {
	if len(funcs) == 0 {
		return nil
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return NewErrorListByChan(onlyNotNil(Run(ctx, funcs...)))
}

// Sequential executes all given functions one after another in order.
// It stops and returns the first error encountered, or nil if all functions succeed.
func Sequential(ctx context.Context, funcs ...Func) (err error) {
	if len(funcs) == 0 {
		return nil
	}
	for _, fn := range funcs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err = fn(ctx); err != nil {
				return
			}
		}
	}
	return
}

// Run executes all given functions in parallel and returns a channel that receives the result of each function.
// The channel is closed when all functions have completed. This provides the lowest-level access to execution results.
func Run(ctx context.Context, funcs ...Func) <-chan error {
	if len(funcs) == 0 {
		return nil
	}
	errors := make(chan error, runtime.NumCPU())
	var wg sync.WaitGroup
	for _, run := range funcs {
		wg.Add(1)
		go func(run Func) {
			defer wg.Done()
			errors <- run(ctx)
		}(run)
	}
	go func() {
		wg.Wait()
		close(errors)
	}()
	return errors
}

func onlyNotNil(ch <-chan error) <-chan error {
	errors := make(chan error, runtime.NumCPU())
	go func() {
		defer close(errors)
		for err := range ch {
			if err != nil {
				errors <- err
			}
		}
	}()
	return errors
}
