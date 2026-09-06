package php

import (
	"math/rand"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

var (
	mtRandMutex sync.Mutex
)

// MtSrandFunction 实现 mt_srand 函数
type MtSrandFunction struct{}

func NewMtSrandFunction() data.FuncStmt { return &MtSrandFunction{} }

func (f *MtSrandFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	mtRandMutex.Lock()
	defer mtRandMutex.Unlock()
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		if iv, ok2 := v.(data.AsInt); ok2 {
			if seed, err := iv.AsInt(); err == nil {
				rand.Seed(int64(seed))
			}
		}
	}
	return data.NewNullValue(), nil
}

func (f *MtSrandFunction) GetName() string { return "mt_srand" }
func (f *MtSrandFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "seed", 0, nil, nil)}
}
func (f *MtSrandFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "seed", 0, data.NewBaseType("int"))}
}

// MtRandFunction 实现 mt_rand 函数
type MtRandFunction struct{}

func NewMtRandFunction() data.FuncStmt { return &MtRandFunction{} }

func (f *MtRandFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	mtRandMutex.Lock()
	defer mtRandMutex.Unlock()

	min, max := 0, int(^uint(0)>>1) // 0..MaxInt

	// 如果有 min 参数
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		if iv, ok2 := v.(data.AsInt); ok2 {
			if n, err := iv.AsInt(); err == nil {
				min = n
			}
		}
	}
	// 如果有 max 参数
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		if iv, ok2 := v.(data.AsInt); ok2 {
			if n, err := iv.AsInt(); err == nil {
				max = n
			}
		}
	}

	if min > max {
		min, max = max, min
	}

	diff := max - min
	if diff < 0 {
		// overflow
		return data.NewIntValue(min), nil
	}
	if diff == 0 {
		return data.NewIntValue(min), nil
	}

	result := min + rand.Intn(diff+1)
	return data.NewIntValue(result), nil
}

func (f *MtRandFunction) GetName() string { return "mt_rand" }
func (f *MtRandFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "min", 0, nil, nil),
		node.NewParameter(nil, "max", 1, nil, nil),
	}
}
func (f *MtRandFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "min", 0, data.NewNullableType(data.NewBaseType("int"))),
		node.NewVariable(nil, "max", 1, data.NewNullableType(data.NewBaseType("int"))),
	}
}

// RandFunction 实现 rand 函数（PHP 的 rand 与 mt_rand 行为兼容）
type RandFunction struct{}

func NewRandFunction() data.FuncStmt { return &RandFunction{} }

func (f *RandFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	mtRandMutex.Lock()
	defer mtRandMutex.Unlock()

	min, max := 0, int(^uint(0)>>1)

	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		if iv, ok2 := v.(data.AsInt); ok2 {
			if n, err := iv.AsInt(); err == nil {
				min = n
			}
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		if iv, ok2 := v.(data.AsInt); ok2 {
			if n, err := iv.AsInt(); err == nil {
				max = n
			}
		}
	}

	if min > max {
		min, max = max, min
	}
	if min == max {
		return data.NewIntValue(min), nil
	}

	diff := max - min
	if diff < 0 {
		return data.NewIntValue(min), nil
	}

	result := min + rand.Intn(diff+1)
	return data.NewIntValue(result), nil
}

func (f *RandFunction) GetName() string { return "rand" }
func (f *RandFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "min", 0, nil, nil),
		node.NewParameter(nil, "max", 1, nil, nil),
	}
}
func (f *RandFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "min", 0, data.NewNullableType(data.NewBaseType("int"))),
		node.NewVariable(nil, "max", 1, data.NewNullableType(data.NewBaseType("int"))),
	}
}

// SrandFunction 实现 srand 函数
type SrandFunction struct{}

func NewSrandFunction() data.FuncStmt { return &SrandFunction{} }

func (f *SrandFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	mtRandMutex.Lock()
	defer mtRandMutex.Unlock()
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		if iv, ok2 := v.(data.AsInt); ok2 {
			if seed, err := iv.AsInt(); err == nil {
				rand.Seed(int64(seed))
			}
		}
	}
	return data.NewNullValue(), nil
}

func (f *SrandFunction) GetName() string { return "srand" }
func (f *SrandFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "seed", 0, nil, nil)}
}
func (f *SrandFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "seed", 0, data.NewBaseType("int"))}
}

// MtGetrandmaxFunction 实现 mt_getrandmax 函数
type MtGetrandmaxFunction struct{}

func NewMtGetrandmaxFunction() data.FuncStmt { return &MtGetrandmaxFunction{} }

func (f *MtGetrandmaxFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(int(^uint(0) >> 1)), nil
}

func (f *MtGetrandmaxFunction) GetName() string { return "mt_getrandmax" }
func (f *MtGetrandmaxFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}
func (f *MtGetrandmaxFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}
