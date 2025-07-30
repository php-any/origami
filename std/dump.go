package std

import (
	"fmt"
	"github.com/php-any/origami/node"
)
import "github.com/php-any/origami/data"

func NewDumpFunction() data.FuncStmt {
	return &DumpFunction{}
}

type DumpFunction struct{}

func (f *DumpFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	for _, argument := range f.GetParams() {
		// argument data.Parameters
		argv, _ := argument.GetValue(ctx)
		switch temp := argv.(type) {
		case data.Variable:
			v, acl := ctx.GetVariableValue(temp)
			if acl != nil {
				return nil, acl
			}
			switch arg := v.(type) {
			case data.AsString:
				fmt.Println(arg.AsString())
			default:
				fmt.Println(arg)
			}
		default:
			switch arg := temp.(type) {
			case data.AsString:
				fmt.Println(arg.AsString())
			default:
				fmt.Println(arg)
			}
		}

	}
	return nil, nil
}
func (f *DumpFunction) GetName() string {
	return "dump"
}
func (f *DumpFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "args", 0, nil, nil),
	}
}
func (f *DumpFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "args", 0, nil),
	}
}
