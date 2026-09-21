package node

import (
	"sync/atomic"

	"github.com/php-any/origami/data"
)

// CheckExecutionTimeLimit 由 std/php 在 Load 时注入（避免 node↔core 循环依赖）。
var CheckExecutionTimeLimit func(file string, line int)

// MarkHeaderOutputStarted 由 std/php 在 Load 时注入。
var MarkHeaderOutputStarted func()

// 语句/循环边界不必每次 time.Now；64 次抽样一次仍能在秒级 max_execution_time 内中止。
var timeLimitCheckTick uint32

func checkTimeLimit(from data.From) {
	if CheckExecutionTimeLimit == nil {
		return
	}
	if atomic.AddUint32(&timeLimitCheckTick, 1)&63 != 0 {
		return
	}
	file, line := "Unknown", 0
	if from != nil {
		if src := from.GetSource(); src != "" {
			file = src
		}
		if sl, _ := from.GetStartPosition(); sl >= 0 {
			line = sl + 1
		}
	}
	CheckExecutionTimeLimit(file, line)
}
