// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XEvent

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShared(t *testing.T) {
	// 确保 Shared() 返回的实例是非空的
	manager := Shared()
	assert.NotNil(t, manager, "Shared() should return a non-nil Manager")

	// 确保多次调用 Shared() 返回同一个实例
	manager2 := Shared()
	assert.Equal(t, manager, manager2, "Shared() should return the same instance on multiple calls")

	// 并发测试，确保在并发环境下 Shared() 仍然返回同一个实例
	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	results := make(chan *Manager, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			results <- Shared()
		}()
	}

	wg.Wait()
	close(results)

	for result := range results {
		assert.Equal(t, manager, result, "All goroutines should receive the same Manager instance")
	}
}
