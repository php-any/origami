package node

import (
	"reflect"

	"github.com/php-any/origami/data"
)

// CheckExecutionTimeLimit 由 std/php 在 Load 时注入（避免 node↔core 循环依赖）。
var CheckExecutionTimeLimit func(file string, line int)

// MarkHeaderOutputStarted 由 std/php 在 Load 时注入。
var MarkHeaderOutputStarted func()

// fromForStatement 安全获取语句的 From（避免 typed-nil 或 nil 嵌入字段调用 GetFrom 时 panic）
func fromForStatement(program *Program, statement data.GetValue) (from data.From) {
	if statement == nil {
		if program != nil {
			return program.GetFrom()
		}
		return nil
	}
	gf, ok := statement.(GetFrom)
	if !ok {
		if program != nil {
			return program.GetFrom()
		}
		return nil
	}
	rv := reflect.ValueOf(gf)
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		if program != nil {
			return program.GetFrom()
		}
		return nil
	}
	func() {
		defer func() {
			if recover() != nil && program != nil {
				from = program.GetFrom()
			}
		}()
		from = gf.GetFrom()
	}()
	return from
}

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
