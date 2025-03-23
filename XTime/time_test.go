// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XTime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test GetMicrosecond
func TestGetMicrosecond(t *testing.T) {
	t.Run("Test Get Microsecond", func(t *testing.T) {
		microsecond := GetMicrosecond()
		assert.True(t, microsecond > 0, "Microsecond should be greater than 0")
	})
}

// Test GetMillisecond
func TestGetMillisecond(t *testing.T) {
	t.Run("Test Get Millisecond", func(t *testing.T) {
		millisecond := GetMillisecond()
		assert.True(t, millisecond > 0, "Millisecond should be greater than 0")
	})
}

// Test GetTimestamp
func TestGetTimestamp(t *testing.T) {
	t.Run("Test Get Timestamp", func(t *testing.T) {
		timestamp := GetTimestamp()
		assert.True(t, timestamp > 0, "Timestamp should be greater than 0")
	})
}

// Test NowTime
func TestNowTime(t *testing.T) {
	t.Run("Test NowTime", func(t *testing.T) {
		now := NowTime()
		assert.WithinDuration(t, now, time.Now(), 1*time.Second, "NowTime should be close to current time")
	})
}

// Test ToTime
func TestToTime(t *testing.T) {
	t.Run("Test ToTime", func(t *testing.T) {
		timestamp := GetTimestamp()
		result := ToTime(timestamp)
		expected := time.Unix(int64(timestamp), 0)

		assert.Equal(t, expected, result, "ToTime should return the correct time")
	})
}

// Test TimeToZero
func TestTimeToZero(t *testing.T) {
	t.Run("Test TimeToZero with current timestamp", func(t *testing.T) {
		timeToZero := TimeToZero()
		assert.True(t, timeToZero >= 0, "TimeToZero should return a non-negative value")
	})

	t.Run("Test TimeToZero with custom timestamp", func(t *testing.T) {
		timestamp := GetTimestamp() - Minute1 // 1 minute before now
		timeToZero := TimeToZero(timestamp)
		assert.True(t, timeToZero >= 0, "TimeToZero should return a non-negative value")
	})
}

// Test ZeroTime
func TestZeroTime(t *testing.T) {
	t.Run("Test ZeroTime with current timestamp", func(t *testing.T) {
		zeroTime := ZeroTime()
		assert.True(t, zeroTime >= 0, "ZeroTime should return a non-negative value")
	})

	t.Run("Test ZeroTime with custom timestamp", func(t *testing.T) {
		timestamp := GetTimestamp() - Minute1 // 1 minute before now
		zeroTime := ZeroTime(timestamp)
		assert.True(t, zeroTime >= 0, "ZeroTime should return a non-negative value")
	})
}

// Test Format function
func TestFormat(t *testing.T) {
	t.Run("Test Format with standard timestamp", func(t *testing.T) {
		timestamp := GetTimestamp()
		formattedTime := Format(timestamp, FormatFull)

		expectedFormat := time.Unix(int64(timestamp), 0).Format(FormatFull)
		assert.Equal(t, expectedFormat, formattedTime, "Format should return the correct formatted time")
	})

	t.Run("Test Format with lite timestamp", func(t *testing.T) {
		timestamp := GetTimestamp()
		formattedTime := Format(timestamp, FormatLite)

		expectedFormat := time.Unix(int64(timestamp), 0).Format(FormatLite)
		assert.Equal(t, expectedFormat, formattedTime, "Format should return the correct formatted time")
	})
}
