// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XEvent

import "sync"

var (
	shared     *Manager
	sharedOnce sync.Once
)

func Shared() *Manager {
	sharedOnce.Do(func() { shared = NewManager(true) })
	return shared
}
