// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package systimemon

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// StartMonitor calls systimeErrHandler if system time jump backward.
func StartMonitor(ctx context.Context, logger *zap.Logger, now func() time.Time, systimeErrHandler func(), successCallback func()) {
	logger.Info("start system time monitor")
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	tickCount := 0
	for {
		last := now().UnixNano()
		select {
		case <-tick.C:
		case <-ctx.Done():
			return
		}
		if now().UnixNano() < last {
			logger.Error("system time jump backward", zap.Int64("last", last))
			systimeErrHandler()
		}
		// call successCallback per second.
		tickCount++
		if tickCount >= 10 {
			tickCount = 0
			successCallback()
		}
	}
}
