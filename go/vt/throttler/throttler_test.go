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

package throttler

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"vitess.io/vitess/go/vt/discovery"
	"vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/proto/topodata"
)

// testTabletTypes is the list of tablet types to test.
var testTabletTypes = []topodata.TabletType{
	topodata.TabletType_REPLICA,
	topodata.TabletType_RDONLY,
}

// The main purpose of the benchmarks below is to demonstrate the functionality
// of the throttler in the real-world (using a non-faked time.Now).
// The benchmark values should be as close as possible to the request interval
// (which is the inverse of the QPS).
// For example, 1k QPS should result into 1,000,000 ns/op.
//
// Example for benchmark results on Lenovo Thinkpad X250, i7-5600U, Quadcore.
//
// $ go test -run=XXX -bench=. -cpu=4 --benchtime=30s
// PASS
// BenchmarkThrottler_1kQPS-4          	   50000	   1000040 ns/op
// BenchmarkThrottler_10kQPS-4         	 1000000	     99999 ns/op
// BenchmarkThrottler_100kQPS-4        	 5000000	      9999 ns/op
// BenchmarkThrottlerParallel_1kQPS-4  	   50000	    999903 ns/op
// BenchmarkThrottlerParallel_10kQPS-4 	  500000	    100060 ns/op
// BenchmarkThrottlerParallel_100kQPS-4	 5000000	      9999 ns/op
// BenchmarkThrottlerDisabled-4	500000000	        94.9 ns/op
// ok  	vitess.io/vitess/go/vt/throttler	448.282

func BenchmarkThrottler_1kQPS(b *testing.B) {
	benchmarkThrottler(b, 1*1000)
}

func BenchmarkThrottler_10kQPS(b *testing.B) {
	benchmarkThrottler(b, 10*1000)
}

func BenchmarkThrottler_100kQPS(b *testing.B) {
	benchmarkThrottler(b, 100*1000)
}

// benchmarkThrottler shows that Throttler actually throttles requests.
func benchmarkThrottler(b *testing.B, qps int64) {
	throttler, _ := NewThrottler("test", "queries", 1, qps, ReplicationLagModuleDisabled)
	defer throttler.Close()
	backoffs := 0
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		backedOff := 0
		for {
			backoff := throttler.Throttle(0)
			if backoff == NotThrottled {
				break
			}
			backedOff++
			backoffs++
			if backedOff > 1 {
				b.Logf("did not wait long enough after backoff. n = %v, last backoff = %v", i, backoff)
			}
			time.Sleep(backoff)
		}
	}
}

func BenchmarkThrottlerParallel_1kQPS(b *testing.B) {
	benchmarkThrottlerParallel(b, 1*1000)
}

func BenchmarkThrottlerParallel_10kQPS(b *testing.B) {
	benchmarkThrottlerParallel(b, 10*1000)
}

func BenchmarkThrottlerParallel_100kQPS(b *testing.B) {
	benchmarkThrottlerParallel(b, 100*1000)
}

// benchmarkThrottlerParallel is the parallel version of benchmarkThrottler.
// Set -cpu to change the number of threads. The QPS should be distributed
// across all threads and the reported benchmark value should be similar
// to the value of benchmarkThrottler.
func benchmarkThrottlerParallel(b *testing.B, qps int64) {
	threadCount := runtime.GOMAXPROCS(0)
	throttler, _ := NewThrottler("test", "queries", threadCount, qps, ReplicationLagModuleDisabled)
	defer throttler.Close()
	threadIDs := make(chan int, threadCount)
	for id := 0; id < threadCount; id++ {
		threadIDs <- id
	}
	close(threadIDs)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		threadID := <-threadIDs

		for pb.Next() {
			backedOff := 0
			for {
				backoff := throttler.Throttle(threadID)
				if backoff == NotThrottled {
					break
				}
				backedOff++
				if backedOff > 1 {
					b.Logf("did not wait long enough after backoff. threadID = %v, last backoff = %v", threadID, backoff)
				}
				time.Sleep(backoff)
			}
		}
		throttler.ThreadFinished(threadID)
	})
}

// BenchmarkThrottlerDisabled is the unthrottled version of
// BenchmarkThrottler. It should report a much lower ns/op value.
func BenchmarkThrottlerDisabled(b *testing.B) {
	throttler, _ := NewThrottler("test", "queries", 1, MaxRateModuleDisabled, ReplicationLagModuleDisabled)
	defer throttler.Close()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		backoff := throttler.Throttle(0)
		if backoff != NotThrottled {
			b.Fatalf("unthrottled throttler should never have throttled us: %v", backoff)
		}
	}
}

type fakeClock struct {
	nowValue time.Time
}

func (f *fakeClock) now() time.Time {
	return f.nowValue
}

func (f *fakeClock) setNow(timeOffset time.Duration) {
	f.nowValue = sinceZero(timeOffset)
}

func sinceZero(sinceZero time.Duration) time.Time {
	return time.Time{}.Add(sinceZero)
}

// Due to limitations of golang.org/x/time/rate.Limiter the 'now' parameter of
// threadThrottler.throttle() must be at least 1 second. See the comment in
// threadThrottler.newThreadThrottler() for more details.

// newThrottlerWithClock should only be used for testing.
func newThrottlerWithClock(name, unit string, threadCount int, maxRate int64, maxReplicationLag int64, nowFunc func() time.Time) (Throttler, error) {
	return newThrottler(GlobalManager, name, unit, threadCount, maxRate, maxReplicationLag, nowFunc)
}

func TestThrottle(t *testing.T) {
	fc := &fakeClock{}
	// 1 Thread, 2 QPS.
	throttler, _ := newThrottlerWithClock("test", "queries", 1, 2, ReplicationLagModuleDisabled, fc.now)
	defer throttler.Close()

	fc.setNow(1000 * time.Millisecond)
	// 2 QPS should divide the current second into two chunks of 500 ms:
	// a) [1s, 1.5s), b) [1.5s, 2s)
	// First call goes through since the chunk is not "used" yet.
	gotBackoff := throttler.Throttle(0)
	require.Equal(t, NotThrottled, gotBackoff, "throttler should not have throttled us")

	// Next call should tell us to backoff until we reach the second chunk.
	fc.setNow(1000 * time.Millisecond)
	wantBackoff := 500 * time.Millisecond
	gotBackoff = throttler.Throttle(0)
	require.Equal(t, wantBackoff, gotBackoff, "throttler should have throttled us")

	// Some time elpased, but we are still in the first chunk and must backoff.
	fc.setNow(1111 * time.Millisecond)
	wantBackoff2 := 389 * time.Millisecond
	gotBackoff = throttler.Throttle(0)
	require.Equal(t, wantBackoff2, gotBackoff, "throttler should have still throttled us")

	// Enough time elapsed that we are in the second chunk now.
	fc.setNow(1500 * time.Millisecond)
	gotBackoff = throttler.Throttle(0)
	require.Equal(t, NotThrottled, gotBackoff, "throttler should not have throttled us")

	// We're in the third chunk and are allowed to issue the third request.
	fc.setNow(2001 * time.Millisecond)
	gotBackoff = throttler.Throttle(0)
	require.Equal(t, NotThrottled, gotBackoff, "throttler should not have throttled us")
}

func TestThrottle_RateRemainderIsDistributedAcrossThreads(t *testing.T) {
	fc := &fakeClock{}
	// 3 Threads, 5 QPS.
	throttler, _ := newThrottlerWithClock("test", "queries", 3, 5, ReplicationLagModuleDisabled, fc.now)
	defer throttler.Close()

	fc.setNow(1000 * time.Millisecond)
	// Out of 5 QPS, each thread gets 1 and two threads get 1 query extra.
	for threadID := 0; threadID < 2; threadID++ {
		gotBackoff := throttler.Throttle(threadID)
		require.Equalf(t, NotThrottled, gotBackoff, "throttler should not have throttled thread %d", threadID)
	}

	fc.setNow(1500 * time.Millisecond)
	// Find the thread which got one extra query.
	threadsWithMoreThanOneQPS := 0
	for threadID := 0; threadID < 2; threadID++ {
		if gotBackoff := throttler.Throttle(threadID); gotBackoff == NotThrottled {
			threadsWithMoreThanOneQPS++
		} else {
			wantBackoff := 500 * time.Millisecond
			require.Equal(t, wantBackoff, gotBackoff, "throttler did throttle us with the wrong backoff time")
		}
	}
	if want := 2; threadsWithMoreThanOneQPS != want {
		require.Equal(t, want, threadsWithMoreThanOneQPS, "wrong number of threads were throttled")
	}

	// Now, all threads are throttled.
	for threadID := 0; threadID < 2; threadID++ {
		wantBackoff := 500 * time.Millisecond
		gotBackoff := throttler.Throttle(threadID)
		require.Equalf(t, wantBackoff, gotBackoff, "throttler should have throttled thread %d", threadID)
	}
}

func TestThreadFinished(t *testing.T) {
	fc := &fakeClock{}
	// 2 Threads, 2 QPS.
	throttler, _ := newThrottlerWithClock("test", "queries", 2, 2, ReplicationLagModuleDisabled, fc.now)
	defer throttler.Close()

	// [1000ms, 2000ms):  Each thread consumes their 1 QPS.
	fc.setNow(1000 * time.Millisecond)
	for threadID := 0; threadID < 2; threadID++ {
		gotBackoff := throttler.Throttle(threadID)
		require.Equalf(t, NotThrottled, gotBackoff, "throttler should not have throttled thread %d", threadID)
	}
	// Now they would be throttled.
	wantBackoff := 1000 * time.Millisecond
	for threadID := 0; threadID < 2; threadID++ {
		gotBackoff := throttler.Throttle(threadID)
		require.Equalf(t, wantBackoff, gotBackoff, "throttler should have throttled thread %d", threadID)
	}

	// [2000ms, 3000ms): One thread finishes, other one gets remaining 1 QPS extra.
	fc.setNow(2000 * time.Millisecond)
	throttler.ThreadFinished(1)

	// Max rate update to threadThrottlers happens asynchronously. Wait for it.
	timer := time.NewTimer(2 * time.Second)
	throttlerImpl, ok := throttler.(*ThrottlerImpl)
	require.True(t, ok)
	for {
		if throttlerImpl.threadThrottlers[0].getMaxRate() == 2 {
			timer.Stop()
			break
		}
		select {
		case <-timer.C:
			t.Fatalf("max rate was not propapgated to threadThrottler[0] in time: %v", throttlerImpl.threadThrottlers[0].getMaxRate())
		default:
			// Timer not up yet. Try again.
		}
	}

	// Consume 2 QPS.
	gotBackoff := throttler.Throttle(0)
	require.Equal(t, NotThrottled, gotBackoff, "throttler should not have throttled us")

	fc.setNow(2500 * time.Millisecond)
	gotBackoff = throttler.Throttle(0)
	require.Equal(t, NotThrottled, gotBackoff, "throttler should not have throttled us")

	// 2 QPS are consumed. Thread 0 should be throttled now.
	wantBackoff2 := 500 * time.Millisecond
	gotBackoff = throttler.Throttle(0)
	require.Equal(t, wantBackoff2, gotBackoff, "throttler should have throttled us")

	// Throttle() from a finished thread will panic.
	defer func() {
		msg := recover()
		require.NotNil(t, msg)
		require.Contains(t, msg, "already finished", "Throttle() after ThreadFinished() panic'd for wrong reason")
	}()
	throttler.Throttle(1)
}

// TestThrottle_MaxRateIsZero tests the behavior if max rate is set to zero.
// In this case, the throttler won't let any requests through until the rate
// changes.
func TestThrottle_MaxRateIsZero(t *testing.T) {
	fc := &fakeClock{}
	// 1 Thread, 0 QPS.
	throttler, _ := newThrottlerWithClock("test", "queries", 1, ZeroRateNoProgess, ReplicationLagModuleDisabled, fc.now)
	defer throttler.Close()

	fc.setNow(1000 * time.Millisecond)
	wantBackoff := 1000 * time.Millisecond
	gotBackoff := throttler.Throttle(0)
	require.Equal(t, wantBackoff, gotBackoff, "throttler should have throttled us")

	fc.setNow(1111 * time.Millisecond)
	wantBackoff2 := 1000 * time.Millisecond
	gotBackoff = throttler.Throttle(0)
	require.Equal(t, wantBackoff2, gotBackoff, "throttler should have throttled us")

	fc.setNow(2000 * time.Millisecond)
	wantBackoff3 := 1000 * time.Millisecond
	gotBackoff = throttler.Throttle(0)
	require.Equal(t, wantBackoff3, gotBackoff, "throttler should have throttled us")
}

func TestThrottle_MaxRateDisabled(t *testing.T) {
	fc := &fakeClock{}
	throttler, _ := newThrottlerWithClock("test", "queries", 1, MaxRateModuleDisabled, ReplicationLagModuleDisabled, fc.now)
	defer throttler.Close()

	fc.setNow(1000 * time.Millisecond)
	// No QPS set. 10 requests in a row are fine.
	for i := 0; i < 10; i++ {
		gotBackoff := throttler.Throttle(0)
		require.Equal(t, NotThrottled, gotBackoff, "throttler should not have throttled us")
	}
}

// TestThrottle_MaxRateLowerThanThreadCount tests the behavior that maxRate
// must not be lower than threadCount. If this is the case, maxRate will be
// set to threadCount.
func TestThrottle_MaxRateLowerThanThreadCount(t *testing.T) {
	fc := &fakeClock{}
	// 2 Thread, 1 QPS.
	throttler, _ := newThrottlerWithClock("test", "queries", 2, 1, ReplicationLagModuleDisabled, fc.now)
	defer throttler.Close()

	// 2 QPS instead of configured 1 QPS allowed since there are 2 threads which
	// must not starve.
	fc.setNow(1000 * time.Millisecond)
	for threadID := 0; threadID < 1; threadID++ {
		gotBackoff := throttler.Throttle(threadID)
		require.Equalf(t, NotThrottled, gotBackoff, "throttler should not have throttled thread %d", threadID)
	}
	wantBackoff := 1000 * time.Millisecond
	for threadID := 0; threadID < 1; threadID++ {
		gotBackoff := throttler.Throttle(threadID)
		require.Equalf(t, wantBackoff, gotBackoff, "throttler should have throttled thread %d", threadID)
	}
}

func TestUpdateMaxRate_AllThreadsFinished(t *testing.T) {
	fc := &fakeClock{}
	throttler, _ := newThrottlerWithClock("test", "queries", 2, 1e9, ReplicationLagModuleDisabled, fc.now)
	defer throttler.Close()

	throttler.ThreadFinished(0)
	throttler.ThreadFinished(1)

	// Make sure that there's no division by zero error (threadsRunning == 0).
	throttlerImpl, ok := throttler.(*ThrottlerImpl)
	require.True(t, ok)
	throttlerImpl.updateMaxRate()
	// We don't care about the Throttler state at this point.
}

func TestClose(t *testing.T) {
	fc := &fakeClock{}
	throttler, _ := newThrottlerWithClock("test", "queries", 1, 1, ReplicationLagModuleDisabled, fc.now)
	throttler.Close()

	defer func() {
		msg := recover()
		require.NotNil(t, msg)
		require.Contains(t, msg, "must not access closed Throttler", "Throttle() after ThreadFinished() panic'd for wrong reason")
	}()
	throttler.Throttle(0)
}

func TestThreadFinished_SecondCallPanics(t *testing.T) {
	fc := &fakeClock{}
	throttler, _ := newThrottlerWithClock("test", "queries", 1, 1, ReplicationLagModuleDisabled, fc.now)
	throttler.ThreadFinished(0)

	defer func() {
		msg := recover()
		require.NotNil(t, msg)
		require.Contains(t, msg, "already finished", "Throttle() after ThreadFinished() panic'd for wrong reason")
	}()
	throttler.ThreadFinished(0)
}

func TestThrottlerMaxLag(t *testing.T) {
	fc := &fakeClock{}
	th, err := newThrottlerWithClock(t.Name(), "queries", 1, 1, 10, fc.now)
	require.NoError(t, err)
	throttler := th.(*ThrottlerImpl)
	defer throttler.Close()

	require.NotNil(t, throttler)
	require.NotNil(t, throttler.maxReplicationLagModule)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// run .add() and .MaxLag() concurrently to detect races
	for _, tabletType := range testTabletTypes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					throttler.MaxLag(tabletType)
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					cache := throttler.maxReplicationLagModule.lagCacheByType(tabletType)
					require.NotNil(t, cache)
					cache.add(replicationLagRecord{
						time: time.Now(),
						TabletHealth: discovery.TabletHealth{
							Serving: true,
							Stats: &query.RealtimeStats{
								ReplicationLagSeconds: 5,
							},
							Tablet: &topodata.Tablet{
								Hostname: t.Name(),
								Type:     tabletType,
								PortMap: map[string]int32{
									"test": 15999,
								},
							},
						},
					})
				}
			}
		}()
	}
	time.Sleep(time.Second)
	cancel()
	wg.Wait()

	// check .MaxLag()
	for _, tabletType := range testTabletTypes {
		require.Equal(t, uint32(5), throttler.MaxLag(tabletType))
	}
}
