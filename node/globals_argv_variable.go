package node

import (
	"errors"
	"os"

	"github.com/php-any/origami/data"
)

// $argv

type ArgvVariable struct {
	*Node `pp:"-"`
}

// CLI 进程级缓存：$argv 对应的数组值（os.Args[1:]）
var argvValue *data.ArrayValue

func NewArgvVariable(from data.From) data.Variable {
	return &ArgvVariable{
		Node: NewNode(from),
	}
}

func (v *ArgvVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// CLI 语义：$argv = os.Args[1:]（脚本路径 + 额外参数）
	if argvValue == nil {
		if len(os.Args) <= 1 {
			if av, ok := data.NewArrayValue([]data.Value{}).(*data.ArrayValue); ok {
				argvValue = av
			}
		} else {
			arr := make([]data.Value, 0, len(os.Args)-1)
			for _, s := range os.Args[1:] {
				arr = append(arr, data.NewStringValue(s))
			}
			if av, ok := data.NewArrayValue(arr).(*data.ArrayValue); ok {
				argvValue = av
			} else {
				argvValue = data.NewArrayValue([]data.Value{}).(*data.ArrayValue)
			}
		}
	}
	return argvValue, nil
}

func (v *ArgvVariable) GetIndex() int       { return 0 }
func (v *ArgvVariable) GetName() string     { return "$argv" }
func (v *ArgvVariable) GetType() data.Types { return nil }
func (v *ArgvVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	// 保守实现：$argv 只读，不允许通过赋值修改（避免与 $_SERVER['argv'] 语义打架）
	return data.NewErrorThrow(v.from, errors.New("$argv is read-only in Origami runtime"))
}

// $argc = len($argv)

type ArgcVariable struct {
	*Node `pp:"-"`
}

var argcValue *data.IntValue

func NewArgcVariable(from data.From) data.Variable {
	return &ArgcVariable{
		Node: NewNode(from),
	}
}

func (v *ArgcVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if argcValue == nil {
		// 通过 argvValue 计算，保证与 $argv 一致
		argv, _ := (&ArgvVariable{Node: v.Node}).GetValue(ctx)
		if arr, ok := argv.(*data.ArrayValue); ok {
			argcValue = data.NewIntValue(len(arr.List)).(*data.IntValue)
		} else {
			argcValue = data.NewIntValue(0).(*data.IntValue)
		}
	}
	return argcValue, nil
}

func (v *ArgcVariable) GetIndex() int       { return 0 }
func (v *ArgcVariable) GetName() string     { return "$argc" }
func (v *ArgcVariable) GetType() data.Types { return nil }
func (v *ArgcVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return data.NewErrorThrow(v.from, errors.New("$argc is read-only in Origami runtime"))
}
