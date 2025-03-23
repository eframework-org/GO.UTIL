// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XCollect

import (
	"testing"
)

func TestIndex(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		ele      any
		expected int
	}{
		{"Element found", []int{1, 2, 3, 4}, 3, 2},
		{"Element not found", []int{1, 2, 3, 4}, 5, -1},
		{"Function match", []int{1, 2, 3, 4}, func(x int) bool { return x%2 == 0 }, 1},
		{"Nil array", nil, 3, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Index(tt.arr, tt.ele)
			if got != tt.expected {
				t.Errorf("Index() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		ele      any
		expected bool
	}{
		{"Element found", []int{1, 2, 3, 4}, 3, true},
		{"Element not found", []int{1, 2, 3, 4}, 5, false},
		{"Function match", []int{1, 2, 3, 4}, func(x int) bool { return x%2 == 0 }, true},
		{"Nil array", nil, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Contains(tt.arr, tt.ele)
			if got != tt.expected {
				t.Errorf("Contains() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		ele      any
		expected []int
	}{
		{"Remove element", []int{1, 2, 3, 4}, 3, []int{1, 2, 4}},
		{"Remove non-existing element", []int{1, 2, 3, 4}, 5, []int{1, 2, 3, 4}},
		{"Nil array", nil, 3, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Remove(tt.arr, tt.ele)
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("Remove() = %v, want %v", got, tt.expected)
					break
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		idx      int
		expected []int
	}{
		{"Delete valid index", []int{1, 2, 3, 4}, 2, []int{1, 2, 4}},
		{"Delete invalid index", []int{1, 2, 3, 4}, 10, []int{1, 2, 3, 4}},
		{"Nil array", nil, 2, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Delete(tt.arr, tt.idx)
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("Delete() = %v, want %v", got, tt.expected)
					break
				}
			}
		})
	}
}

func TestAppend(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		ele      int
		expected []int
	}{
		{"Append to array", []int{1, 2, 3}, 4, []int{1, 2, 3, 4}},
		{"Append to empty array", []int{}, 1, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Append(tt.arr, tt.ele)
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("Append() = %v, want %v", got, tt.expected)
					break
				}
			}
		})
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		idx      int
		ele      int
		expected []int
	}{
		{"Insert at valid index", []int{1, 2, 3, 4}, 2, 5, []int{1, 2, 5, 3, 4}},
		{"Insert at start", []int{1, 2, 3}, 0, 0, []int{0, 1, 2, 3}},
		{"Insert at end", []int{1, 2, 3}, 3, 4, []int{1, 2, 3, 4}},
		{"Invalid index", []int{1, 2, 3}, -1, 5, []int{1, 2, 3}},
		{"Nil array", nil, 2, 5, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Insert(tt.arr, tt.idx, tt.ele)
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("Insert() = %v, want %v", got, tt.expected)
					break
				}
			}
		})
	}
}
