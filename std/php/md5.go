package php

import (
	"crypto/md5"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewMd5Function() data.FuncStmt {
	return &Md5Function{}
}

type Md5Function struct{}

func (f *Md5Function) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	rawOutputValue, _ := ctx.GetIndexValue(1)

	if stringValue == nil {
		return data.NewStringValue(""), nil
	}

	str := stringValue.AsString()

	// 计算 MD5 哈希
	hash := md5.Sum([]byte(str))

	// 处理 raw_output 参数
	rawOutput := false
	if rawOutputValue != nil {
		if _, ok := rawOutputValue.(*data.NullValue); !ok {
			if rawBool, ok := rawOutputValue.(data.AsBool); ok {
				if r, err := rawBool.AsBool(); err == nil {
					rawOutput = r
				}
			}
		}
	}

	if rawOutput {
		// 返回原始二进制数据
		return data.NewStringValue(string(hash[:])), nil
	}

	// 返回十六进制字符串
	return data.NewStringValue(fmt.Sprintf("%x", hash)), nil
}

func (f *Md5Function) GetName() string {
	return "md5"
}

func (f *Md5Function) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "raw_output", 1, node.NewNullLiteral(nil), nil),
	}
}

func (f *Md5Function) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "raw_output", 1, data.NewBaseType("bool")),
	}
}
