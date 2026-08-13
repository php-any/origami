package php

import (
	"runtime"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewMemoryGetPeakUsageFunction() data.FuncStmt {
	return &MemoryGetPeakUsageFunction{}
}

type MemoryGetPeakUsageFunction struct {
	data.Function
}

func (f *MemoryGetPeakUsageFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	realUsage := false
	if temp, ok := ctx.GetIndexValue(0); ok {
		realUsage, _ = temp.(data.AsBool).AsBool()
	}

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	var bytes uint64
	if realUsage {
		bytes = stats.Sys
	} else {
		bytes = stats.Alloc
	}

	return data.NewIntValue(int(bytes)), nil
}

func (f *MemoryGetPeakUsageFunction) GetName() string {
	return "memory_get_peak_usage"
}

func (f *MemoryGetPeakUsageFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "real_usage", 0, data.NewBoolValue(false), nil),
	}
}

func (f *MemoryGetPeakUsageFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "real_usage", 0, nil),
	}
}
