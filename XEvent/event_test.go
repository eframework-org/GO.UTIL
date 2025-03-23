// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XEvent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	mgr := NewManager(true)
	assert.NotNil(t, mgr)
	assert.Equal(t, true, mgr.Multiple)
	assert.NotNil(t, mgr.Events)
}

func TestManagerClear(t *testing.T) {
	mgr := NewManager(true)
	mgr.Reg(1, func(args ...any) {})
	mgr.Clear()
	assert.Equal(t, 0, len(mgr.Events))
}

func TestManagerGet(t *testing.T) {
	mgr := NewManager(true)
	mgr.Reg(1, func(args ...any) {})
	evtWrap := mgr.Get(1)
	assert.NotNil(t, evtWrap)
	assert.Equal(t, 1, evtWrap.ID)
}

func TestManagerReg(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		mgr := NewManager(true)
		assert.False(t, mgr.Reg(1, nil))
	})

	t.Run("Singleton", func(t *testing.T) {
		mgr := NewManager(false)
		handler1 := func(args ...any) {}
		handler2 := func(args ...any) {}

		ok := mgr.Reg(1, handler1)
		assert.True(t, ok)
		assert.Equal(t, 1, len(mgr.Events[1].Hnds))

		// 在Multiple=false时不能注册多个处理器
		ok = mgr.Reg(1, handler2)
		assert.False(t, ok)
		assert.Equal(t, 1, len(mgr.Events[1].Hnds))
	})

	t.Run("Multiple", func(t *testing.T) {
		mgr := NewManager(true)
		handler1 := func(args ...any) {}
		handler2 := func(args ...any) {}

		ok := mgr.Reg(1, handler1)
		assert.True(t, ok)
		assert.Equal(t, 1, len(mgr.Events[1].Hnds))

		// 在Multiple=true时可以注册多个处理器
		ok = mgr.Reg(1, handler2)
		assert.True(t, ok)
		assert.Equal(t, 2, len(mgr.Events[1].Hnds))
	})
}

func TestManagerUnreg(t *testing.T) {
	mgr := NewManager(true)
	handler := func(args ...any) {}
	mgr.Reg(1, handler)
	ok := mgr.Unreg(1, handler)
	assert.True(t, ok)
	assert.Equal(t, 0, len(mgr.Events))

	mgr.Reg(1, handler)
	ok = mgr.Unreg(1)
	assert.True(t, ok)
	assert.Equal(t, 0, len(mgr.Events))
}

func TestManagerNotify(t *testing.T) {
	t.Run("Multiple", func(t *testing.T) {
		mgr := NewManager(true)
		called := false
		handler := func(args ...any) {
			called = true
		}
		mgr.Reg(1, handler)
		ok := mgr.Notify(1)
		assert.True(t, ok)
		assert.True(t, called)
	})

	t.Run("Once", func(t *testing.T) {
		mgr := NewManager(true)
		called := 0
		handler := func(args ...any) {
			called++
		}
		mgr.Reg(1, handler, true)
		mgr.Notify(1)
		mgr.Notify(1)
		assert.Equal(t, 1, called)
	})
}
