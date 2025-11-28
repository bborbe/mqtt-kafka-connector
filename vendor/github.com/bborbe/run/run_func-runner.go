// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import "errors"

//counterfeiter:generate -o mocks/func-runner.go --fake-name FuncRunner . FuncRunner

// FuncRunner is an interface for executing Func with custom behavior.
// It allows wrapping, decorating, or transforming function execution.
type FuncRunner interface {
	Run(runFunc Func) error
}

// FuncRunnerFunc is a function type adapter that implements the FuncRunner interface.
// This allows plain functions to be used wherever FuncRunner is expected, following
// the standard Go pattern of function types implementing single-method interfaces.
type FuncRunnerFunc func(runFunc Func) error

// Run executes the function, implementing the FuncRunner interface.
func (f FuncRunnerFunc) Run(runFunc Func) error {
	if runFunc == nil {
		return errors.New("nil function")
	}
	return f(runFunc)
}
