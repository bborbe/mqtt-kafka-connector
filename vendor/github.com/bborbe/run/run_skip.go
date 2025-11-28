// Copyright (c) 2020 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/getsentry/sentry-go"
	"github.com/golang/glog"
)

// SkipErrors wraps the given function to suppress all errors and always return nil.
// Errors are logged as warnings but do not propagate to the caller.
func SkipErrors(fn Func) Func {
	return func(ctx context.Context) error {
		if err := fn(ctx); err != nil {
			glog.Warningf("run failed: %v", err)
		}
		return nil
	}
}

//counterfeiter:generate -o mocks/has-capture-exception.go --fake-name HasCaptureException . HasCaptureException

// HasCaptureException defines the interface for error reporting services.
// It is compatible with sentry.Client for reporting captured errors.
type HasCaptureException interface {
	CaptureException(
		exception error,
		hint *sentry.EventHint,
		scope sentry.EventModifier,
	) *sentry.EventID
}

// SkipErrorsAndReport wraps the given function to suppress all errors, report them to an error tracking service, and always return nil.
// Context cancellation errors are not reported. The function logs errors as warnings and sends them to the specified error tracker.
func SkipErrorsAndReport(
	fn Func,
	hasCaptureException HasCaptureException,
	tags map[string]string,
) Func {
	return func(ctx context.Context) error {
		if err := fn(ctx); err != nil && !errors.Is(err, context.Canceled) {
			glog.Warningf("run failed: %v", err)
			hasCaptureException.CaptureException(
				errors.Wrapf(ctx, err, "run failed"),
				&sentry.EventHint{
					Context:           ctx,
					Data:              tags,
					OriginalException: err,
				},
				sentry.NewScope(),
			)
		}
		return nil
	}
}
