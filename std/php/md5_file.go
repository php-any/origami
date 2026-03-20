package php

import (
	"crypto/md5"
	"fmt"
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// Md5FileFunction 实现 PHP 内置函数 md5_file
//
//	md5_file(string $filename, bool $raw_output = false): string|false
//
// 计算指定文件的 MD5 散列值。文件不存在或读取失败时返回 false。
func NewMd5FileFunction() data.FuncStmt {
	return &Md5FileFunction{}
}

type Md5FileFunction struct{}

func (f *Md5FileFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	filenameValue, _ := ctx.GetIndexValue(0)
	rawOutputValue, _ := ctx.GetIndexValue(1)

	if filenameValue == nil {
		return data.NewBoolValue(false), nil
	}

	filename := filenameValue.AsString()
	if filename == "" {
		return data.NewBoolValue(false), nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	hash := md5.Sum(content)

	rawOutput := false
	if rawOutputValue != nil {
		if _, ok := rawOutputValue.(*data.NullValue); !ok {
			if rawBool, ok := rawOutputValue.(data.AsBool); ok {
				if r, e := rawBool.AsBool(); e == nil {
					rawOutput = r
				}
			}
		}
	}

	if rawOutput {
		return data.NewStringValue(string(hash[:])), nil
	}
	return data.NewStringValue(fmt.Sprintf("%x", hash)), nil
}

func (f *Md5FileFunction) GetName() string {
	return "md5_file"
}

func (f *Md5FileFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
		node.NewParameter(nil, "raw_output", 1, node.NewNullLiteral(nil), nil),
	}
}

func (f *Md5FileFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "raw_output", 1, data.NewBaseType("bool")),
	}
}
