package php

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type UniqidFunction struct{}

func NewUniqidFunction() data.FuncStmt {
	return &UniqidFunction{}
}

var uniqidCounter uint64

func (f *UniqidFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	prefixVal, _ := ctx.GetIndexValue(0)
	moreEntropyVal, _ := ctx.GetIndexValue(1)

	prefix := ""
	if prefixVal != nil {
		prefix = prefixVal.AsString()
	}

	moreEntropy := false
	if moreEntropyVal != nil {
		if b, ok := moreEntropyVal.(*data.BoolValue); ok && b.Value {
			moreEntropy = true
		}
	}

	now := time.Now()
	sec := now.Unix()
	usec := now.UnixNano() / 1000

	if moreEntropy {
		counter := atomic.AddUint64(&uniqidCounter, 1)
		// Match PHP's more_entropy format: prefix + %08x%05x + "0.xxxxxxxx" (10 chars)
		id := fmt.Sprintf("%s%08x%05x0.%08d", prefix, sec, usec%100000, counter%100000000)
		return data.NewStringValue(id), nil
	}

	id := fmt.Sprintf("%s%08x%05x", prefix, sec, usec%100000)
	return data.NewStringValue(id), nil
}

func (f *UniqidFunction) GetName() string {
	return "uniqid"
}

func (f *UniqidFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "prefix", 0, data.NewStringValue(""), data.String{}),
		node.NewParameter(nil, "more_entropy", 1, data.NewBoolValue(false), data.Bool{}),
	}
}

func (f *UniqidFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "prefix", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "more_entropy", 1, data.NewBaseType("bool")),
	}
}
