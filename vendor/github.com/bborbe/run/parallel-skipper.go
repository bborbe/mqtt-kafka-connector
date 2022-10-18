// Copyright (c) 2020 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import (
	"context"
	"sync"

	"github.com/golang/glog"
)

// ParallelSkipper prevent execution of the given function at the same time.
type ParallelSkipper interface {
	SkipParallel(action Func) Func
}

// NewParallelSkipper return a new instance of ParallelSkipper.
func NewParallelSkipper() ParallelSkipper {
	return &parallelSkipper{}
}

type parallelSkipper struct {
	running bool
	mux     sync.Mutex
}

func (d *parallelSkipper) SkipParallel(action Func) Func {
	return func(ctx context.Context) error {
		d.mux.Lock()
		if d.running {
			glog.V(2).Infof("skip => already running")
			d.mux.Unlock()
			return nil
		}
		glog.V(2).Infof("run started => locked")
		d.running = true
		d.mux.Unlock()
		err := action(ctx)
		d.mux.Lock()
		glog.V(2).Infof("run finished => unlocked")
		d.running = false
		d.mux.Unlock()
		return err
	}
}
