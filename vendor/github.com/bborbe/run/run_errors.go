// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import (
	"github.com/bborbe/errors"
)

// NewErrorList creates an aggregate error from the given errors.
// It returns nil if no errors are provided or all errors are nil.
func NewErrorList(errs ...error) error {
	return errors.Join(errs...)
}

// NewErrorListByChan creates an aggregate error from all errors received from the given channel.
// It blocks until the channel is closed and returns nil if no errors were received.
func NewErrorListByChan(ch <-chan error) error {
	var errs []error
	for err := range ch {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
