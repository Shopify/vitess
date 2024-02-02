/*
Copyright 2019 The Vitess Authors.

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

package discovery

import "vitess.io/vitess/go/vt/log"

type CircuitBreakerState int

const (
	CircuitBreakerState_CLOSED CircuitBreakerState = iota
	CircuitBreakerState_OPEN
)

// These would come from command line flags
const queryTimeoutWindowSeconds = 3
const queryTimeoutPerSecondThreshold = 1

func FilterStatsByCircuitBreakerState(tabletHealthList []*TabletHealth) []*TabletHealth {
	aliases := make([]string, 0, len(tabletHealthList))
	for _, ts := range tabletHealthList {
		aliases = append(aliases, ts.Tablet.Alias.String())
	}
	log.Infof("Filtering tablet health list by circuit state: %v", aliases)

	list := make([]*TabletHealth, 0, len(tabletHealthList))
	for _, ts := range tabletHealthList {
		if ts.Stats == nil || len(ts.Stats.QueryTimeoutRates) == 0 {
			continue
		}

		timeoutRates := ts.Stats.QueryTimeoutRates
		windowSize := min(queryTimeoutWindowSeconds, len(timeoutRates))
		ratesInWindow := timeoutRates[len(timeoutRates)-windowSize:]
		log.Infof("query timeout rates for tablet %v: %v", ts.Tablet.Alias.String(), ratesInWindow)

		sum := 0.0
		for _, rate := range ratesInWindow {
			sum += rate
		}
		avgRate := sum / float64(windowSize)
		log.Infof("average query timeout rate for tablet %v: %v", ts.Tablet.Alias.String(), avgRate)

		if avgRate > queryTimeoutPerSecondThreshold {
			log.Infof("removing tablet %v from healthcheck list due to high query timeout rate: %v", ts.Tablet.Alias.String(), avgRate)
			ts.CircuitState = CircuitBreakerState_OPEN
			continue
		}

		ts.CircuitState = CircuitBreakerState_CLOSED
		list = append(list, ts)
	}

	return list
}
