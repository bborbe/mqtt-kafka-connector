// Copyright (c) 2023-2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import (
	"context"

	"github.com/golang/glog"
)

// BackgroundRunner executes functions in background goroutines with parallel execution prevention.
// It embeds FuncRunner and inherits the Run method, executing functions asynchronously.
// It ensures that only one instance of a function runs at a time, skipping subsequent calls if already running.
type BackgroundRunner interface {
	FuncRunner
}

// NewBackgroundRunner creates a new BackgroundRunner that uses the provided context for all background operations.
// The returned runner will skip parallel executions and log the results of background operations.
func NewBackgroundRunner(ctx context.Context) BackgroundRunner {
	parallelSkipper := NewParallelSkipper()
	return FuncRunnerFunc(func(runFunc Func) error {
		go func() {
			action := parallelSkipper.SkipParallel(func(ctx context.Context) error {
				if err := runFunc(ctx); err != nil {
					return err
				}
				return nil
			})
			glog.V(3).Infof("run started")
			if err := action(ctx); err != nil {
				glog.Warningf("run failed: %v", err)
			}
			glog.V(3).Infof("run completed")
		}()
		return nil
	})
}
