// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XEvent

import (
	"reflect"
	"sync"

	"github.com/eframework-org/GO.UTIL/XCollect"
	"github.com/eframework-org/GO.UTIL/XLog"
)

// Callback 定义了事件处理函数的类型。
type Callback func(args ...any)

// Manager 管理事件的注册和通知。
type Manager struct {
	sync.Mutex
	Multiple bool
	Events   map[int]*EvtWrap
}

// EvtWrap 包装了特定事件 ID 的事件处理程序。
type EvtWrap struct {
	ID   int
	Hnds []*HndWrap
}

var evtWrapPool sync.Pool = sync.Pool{New: func() any { return new(EvtWrap) }}

// HndWrap 包装了一个事件处理程序。
type HndWrap struct {
	Func Callback
	Ptr  uintptr
	Once bool
}

var hndWrapPool sync.Pool = sync.Pool{New: func() any { return new(HndWrap) }}

// NewManager 创建一个新的 Manager 实例。
func NewManager(multiple bool) *Manager {
	return &Manager{sync.Mutex{}, multiple, map[int]*EvtWrap{}}
}

// Clear 清除所有注册的事件和处理程序。
func (mgr *Manager) Clear() {
	defer mgr.Unlock()
	mgr.Lock()

	for _, m := range mgr.Events {
		for _, n := range m.Hnds {
			hndWrapPool.Put(n)
		}
		evtWrapPool.Put(m)
	}
	mgr.Events = make(map[int]*EvtWrap)
}

// Get 获取给定事件 ID 的事件包装器。
func (mgr *Manager) Get(eid int) *EvtWrap {
	defer mgr.Unlock()
	mgr.Lock()

	return mgr.Events[eid]
}

// Reg 为给定的事件 ID 注册一个新的事件处理程序。
func (mgr *Manager) Reg(eid int, handler Callback, once ...bool) bool {
	if nil == handler {
		XLog.Error("XEvent.Manager.Reg: nil handler, id=%v", eid)
		return false
	}
	defer mgr.Unlock()
	mgr.Lock()

	evtWrap, ok := mgr.Events[eid]
	if !ok {
		evtWrap = evtWrapPool.Get().(*EvtWrap)
		evtWrap.ID = eid
		evtWrap.Hnds = make([]*HndWrap, 0)
		mgr.Events[eid] = evtWrap
	}

	if !mgr.Multiple && len(evtWrap.Hnds) > 0 {
		XLog.Error("XEvent.Manager.Reg: doesn't support multiple register, id=%v", eid)
		return false
	}

	hndWrap := hndWrapPool.Get().(*HndWrap)
	hndWrap.Func = handler
	hndWrap.Ptr = uintptr(reflect.ValueOf(handler).Pointer())
	if len(once) > 0 {
		hndWrap.Once = once[0]
	}

	evtWrap.Hnds = append(evtWrap.Hnds, hndWrap)

	return true
}

// Unreg 为给定的事件 ID 注销一个事件处理程序。
func (mgr *Manager) Unreg(eid int, handler ...Callback) bool {
	defer mgr.Unlock()
	mgr.Lock()

	if evtWrap, ok := mgr.Events[eid]; ok {
		sig := false
		var hptr uintptr = 0
		if len(handler) > 0 && handler[0] != nil {
			hptr = uintptr(reflect.ValueOf(handler[0]).Pointer())
		}
		if hptr != 0 {
			nhnds := XCollect.Remove(evtWrap.Hnds, func(ele *HndWrap) bool {
				ok := ele.Ptr == hptr
				if ok {
					hndWrapPool.Put(ele)
				}
				return ok
			})
			sig = len(nhnds) != len(evtWrap.Hnds)
			evtWrap.Hnds = nhnds
			if len(evtWrap.Hnds) == 0 {
				evtWrapPool.Put(evtWrap)
				delete(mgr.Events, eid)
			}
		} else {
			sig = len(evtWrap.Hnds) > 0
			if sig {
				for _, h := range evtWrap.Hnds {
					hndWrapPool.Put(h)
				}
				evtWrapPool.Put(evtWrap)
				delete(mgr.Events, eid)
			}
		}

		return sig
	} else {
		return false
	}
}

// Notify 通知所有为给定事件 ID 注册的处理程序。
func (mgr *Manager) Notify(eid int, args ...any) bool {
	defer XLog.Caught(false)

	evtWrap := mgr.Get(eid)
	if evtWrap == nil {
		return false
	}

	hnds := evtWrap.Hnds
	var onces []*HndWrap
	for _, h := range hnds {
		if h != nil && h.Func != nil {
			if h.Once {
				onces = append(onces, h)
			}
			h.Func(args...)
		}
	}

	if onces != nil && len(onces) > 0 {
		for _, h := range onces {
			mgr.Unreg(eid, h.Func)
		}
	}
	return true
}
