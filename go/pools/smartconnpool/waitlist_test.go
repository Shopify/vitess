/*
Copyright 2023 The Vitess Authors.

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

package smartconnpool

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitlistExpireWithMultipleWaiters(t *testing.T) {
	wait := waitlist[*TestConn]{}
	wait.init()

	ctx := context.Background()

	waiterCount := 2
	expireCount := atomic.Int32{}

	for i := 0; i < waiterCount; i++ {
		go func() {
			fmt.Printf("Running goroutine %d\n", i)
			newCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
			defer cancel()
			_, err := wait.waitForConn(newCtx, nil)
			if err != nil {
				expireCount.Add(1)
			}
		}()
	}

	// Wait for the contexts to expire
	time.Sleep(1 * time.Second)

	// Expire all waiters (we hope)
	wait.expire(false)

	// Wait for the expired goroutines to finish
	time.Sleep(1 * time.Second)

	assert.Equal(t, int32(waiterCount), expireCount.Load())
}
