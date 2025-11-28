// Copyright (c) 2021 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import "sync"

// Fire represents the ability to trigger an event.
type Fire interface {
	// Fire signals the trigger, notifying any listeners waiting on the Done channel.
	Fire()
}

// Done represents the ability to wait for a trigger event.
type Done interface {
	// Done returns a channel that receives a signal when the trigger is fired.
	Done() <-chan struct{}
}

// Trigger combines the ability to fire events and wait for them.
// It provides a simple synchronization mechanism for coordinating between goroutines.
type Trigger interface {
	Fire
	Done
}

// NewTrigger creates a new Trigger that can be used to coordinate between goroutines.
// The trigger starts in an unfired state and can be fired multiple times safely.
func NewTrigger() Trigger {
	return &trigger{}
}

type trigger struct {
	mux sync.Mutex
	ch  chan struct{}
}

func (t *trigger) Fire() {
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.ch == nil {
		t.ch = make(chan struct{})
		close(t.ch)
		return
	}
	select {
	case <-t.ch:
	default:
		close(t.ch)
	}
}

func (t *trigger) Done() <-chan struct{} {
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.ch == nil {
		t.ch = make(chan struct{})
	}
	return t.ch
}
