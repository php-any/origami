package math

import (
	"math"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MinFunction 实现 min 函数
type MinFunction struct{}

func NewMinFunction() data.FuncStmt {
	return &MinFunction{}
}

func (f *MinFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	var minVal float64 = math.MaxFloat64
	var hasValue bool

	for i := 0; ; i++ {
		v, ok := ctx.GetIndexValue(i)
		if !ok || v == nil {
			break
		}
		// 处理数组参数
		if arr, isArr := v.(*data.ArrayValue); isArr {
			for _, z := range arr.List {
				if intVal, ok := z.Value.(data.AsInt); ok {
					if iv, err := intVal.AsInt(); err == nil {
						fv := float64(iv)
						if fv < minVal || !hasValue {
							minVal = fv
							hasValue = true
						}
					}
				} else if floatVal, ok := z.Value.(data.AsFloat); ok {
					fv, _ := floatVal.AsFloat()
					if fv < minVal || !hasValue {
						minVal = fv
						hasValue = true
					}
				}
			}
			continue
		}
		if intVal, ok := v.(data.AsInt); ok {
			if iv, err := intVal.AsInt(); err == nil {
				fv := float64(iv)
				if fv < minVal || !hasValue {
					minVal = fv
					hasValue = true
				}
			}
		} else if floatVal, ok := v.(data.AsFloat); ok {
			fv, _ := floatVal.AsFloat()
			if fv < minVal || !hasValue {
				minVal = fv
				hasValue = true
			}
		}
	}

	if !hasValue {
		return data.NewNullValue(), nil
	}

	if minVal == float64(int64(minVal)) {
		return data.NewIntValue(int(minVal)), nil
	}
	return data.NewFloatValue(minVal), nil
}

func (f *MinFunction) GetName() string {
	return "min"
}

func (f *MinFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "values", 0, nil, data.Mixed{}),
	}
}

func (f *MinFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "values", 0, data.Mixed{}),
	}
}
