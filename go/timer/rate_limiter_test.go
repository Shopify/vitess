/*
Copyright 2020 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package timer

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimiterLong(t *testing.T) {
	r := NewRateLimiter(time.Hour)
	require.NotNil(t, r)
	defer r.Stop()
	val := 0
	incr := func() error { val++; return nil }
	for i := 0; i < 10; i++ {
		err := r.Do(incr)
		assert.NoError(t, err)
	}
	assert.Equal(t, 1, val)
}

func TestRateLimiterShort(t *testing.T) {
	r := NewRateLimiter(time.Millisecond * 250)
	require.NotNil(t, r)
	defer r.Stop()
	val := 0
	incr := func() error { val++; return nil }
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 100)
		err := r.Do(incr)
		assert.NoError(t, err)
	}
	// we expect some 3-5 entries; this depends on the CI server performance.
	assert.Greater(t, val, 2)
	assert.Less(t, val, 10)
}

func TestRateLimiterAllowOne(t *testing.T) {
	r := NewRateLimiter(time.Millisecond * 250)
	require.NotNil(t, r)
	defer r.Stop()
	val := 0
	incr := func() error { val++; return nil }
	times := 10
	for range times {
		time.Sleep(time.Millisecond * 100)
		r.AllowOne()
		err := r.Do(incr)
		assert.NoError(t, err)
	}
	// we expect exactly 10 successful entries.
	assert.Equal(t, times, val)
}

func TestRateLimiterStop(t *testing.T) {
	r := NewRateLimiter(time.Millisecond * 10)
	require.NotNil(t, r)
	defer r.Stop()
	val := 0
	incr := func() error { val++; return nil }
	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 10)
		err := r.Do(incr)
		assert.NoError(t, err)
	}
	// we expect some 3-5 entries; this depends on the CI server performance.
	assert.Greater(t, val, 2)
	valSnapshot := val
	r.Stop()
	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 10)
		err := r.Do(incr)
		assert.NoError(t, err)
	}
	assert.Equal(t, valSnapshot, val)
}

func TestRateLimiterDiff(t *testing.T) {
	d := 2 * time.Second
	r := NewRateLimiter(d)
	require.NotNil(t, r)
	defer r.Stop()

	// This assumes the last couple lines of code run faster than 2 seconds, which should be the case.
	// But if you see flakiness due to slow runners, we can revisit the logic.
	assert.Greater(t, r.Diff(), int64(math.MaxInt32))
	time.Sleep(d + time.Second)
	assert.Greater(t, r.Diff(), int64(math.MaxInt32))
	r.DoEmpty()
	assert.LessOrEqual(t, r.Diff(), int64(1))
}

func TestRateLimiterUninitialized(t *testing.T) {
	r := &RateLimiter{}
	r.Stop()
}
