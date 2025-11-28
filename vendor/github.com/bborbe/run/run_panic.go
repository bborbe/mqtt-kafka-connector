// Copyright (c) 2021 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import (
	"context"
	"fmt"
)

// CatchPanic wraps the given function to recover from panics and convert them to errors.
// If the wrapped function panics, the panic is recovered and returned as an error.
func CatchPanic(fn Func) Func {
	return func(ctx context.Context) (err error) {
		defer func() {
			if panic := recover(); panic != nil {
				err = fmt.Errorf("catch panic: %v", panic)
			}
		}()
		return fn(ctx)
	}
}
