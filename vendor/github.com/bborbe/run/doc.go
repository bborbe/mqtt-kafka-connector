// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package run provides utilities for parallel function execution with different error handling strategies.
//
// The package offers various execution patterns for running multiple functions concurrently:
//   - CancelOnFirstError: stops all functions when the first error occurs
//   - CancelOnFirstFinish: stops all functions when the first one completes
//   - All: runs all functions and collects all errors
//   - Sequential: runs functions one after another
//
// Additional utilities include retry mechanisms with configurable backoff,
// delayed execution, panic recovery, parallel execution prevention,
// and background task management.
//
// Basic usage:
//
//	ctx := context.Background()
//	err := run.All(ctx,
//		func(ctx context.Context) error { return doTask1(ctx) },
//		func(ctx context.Context) error { return doTask2(ctx) },
//	)
//
// The package is built around the Func type which represents any operation
// that can be executed with a context and may return an error.
package run
