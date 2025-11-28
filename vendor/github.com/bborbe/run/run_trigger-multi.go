// Copyright (c) 2021 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import (
	"sync"
)

// AddFire represents the ability to create new triggers dynamically.
type AddFire interface {
	// Add creates and returns a new trigger that becomes part of the multi-trigger group.
	Add() Trigger
}

// MultiTrigger manages multiple triggers and fires when all of them have been triggered.
// It combines the ability to add new triggers dynamically and wait for all triggers to fire.
type MultiTrigger interface {
	Done
	AddFire
}

// NewMultiTrigger creates a new MultiTrigger that waits for all added triggers to fire.
// The Done channel signals when all individual triggers have been fired.
func NewMultiTrigger() MultiTrigger {
	return &multiTrigger{}
}

type multiTrigger struct {
	mux      sync.Mutex
	triggers []Trigger
}

func (m *multiTrigger) Add() Trigger {
	m.mux.Lock()
	defer m.mux.Unlock()

	result := NewTrigger()
	m.triggers = append(m.triggers, result)
	return result
}

func (m *multiTrigger) Done() <-chan struct{} {
	m.mux.Lock()
	defer m.mux.Unlock()

	result := NewTrigger()
	var wg sync.WaitGroup
	counter := 0
	for _, trigger := range m.triggers {
		v := trigger
		select {
		case <-v.Done():
			counter++
		default:
			wg.Add(1)
			go func(done Done) {
				<-done.Done()
				wg.Done()
			}(v)
		}
	}
	if counter == len(m.triggers) {
		result.Fire()
	} else {
		go func() {
			wg.Wait()
			result.Fire()
		}()
	}
	return result.Done()
}
