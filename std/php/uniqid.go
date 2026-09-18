package php

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// UniqidFunction 实现 uniqid(string $prefix = "", bool $more_entropy = false): string
// 对齐 PHP：基于当前时间的微秒生成唯一 ID；more_entropy 时追加额外熵。
type UniqidFunction struct{}

func NewUniqidFunction() data.FuncStmt { return &UniqidFunction{} }

func (f *UniqidFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	prefix := ""
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		prefix = v.AsString()
	}

	moreEntropy := false
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		if b, ok := v.(data.AsBool); ok {
			if bv, err := b.AsBool(); err == nil {
				moreEntropy = bv
			}
		} else if b, ok := v.(*data.BoolValue); ok {
			moreEntropy = b.Value
		}
	}

	now := time.Now()
	// PHP：无 more_entropy 时为 13 位十六进制（秒+微秒）；有 more_entropy 时再追加「.」+ 8 位十进制熵。
	sec := now.Unix()
	usec := now.Nanosecond() / 1000
	id := fmt.Sprintf("%08x%05x", sec, usec)
	if moreEntropy {
		mtRandMutex.Lock()
		extra := rand.Float64()
		mtRandMutex.Unlock()
		id = fmt.Sprintf("%s.%08.0f", id, extra*100000000)
	}

	return data.NewStringValue(prefix + id), nil
}

func (f *UniqidFunction) GetName() string { return "uniqid" }

func (f *UniqidFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "prefix", 0, data.NewStringValue(""), nil),
		node.NewParameter(nil, "more_entropy", 1, data.NewBoolValue(false), nil),
	}
}

func (f *UniqidFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "prefix", 0, nil),
		node.NewVariable(nil, "more_entropy", 1, nil),
	}
}
