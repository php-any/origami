package node

import (
	"github.com/php-any/origami/data"
)

// CheckExecutionTimeLimit 由 std/php 在 Load 时注入（避免 node↔core 循环依赖）。
var CheckExecutionTimeLimit func(file string, line int)

// MarkHeaderOutputStarted 由 std/php 在 Load 时注入。
var MarkHeaderOutputStarted func()

func checkTimeLimit(from data.From) {
	if CheckExecutionTimeLimit == nil {
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
