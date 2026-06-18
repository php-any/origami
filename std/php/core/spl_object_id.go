package core

import (
	"fmt"
	"reflect"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SplObjectIdFunction 实现 spl_object_id 函数。
// 返回对象实例的唯一整数标识，生命周期内保持不变（基于 Go 指针地址）。
type SplObjectIdFunction struct{}

func NewSplObjectIdFunction() data.FuncStmt {
	return &SplObjectIdFunction{}
}

func (f *SplObjectIdFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	objVal, _ := ctx.GetIndexValue(0)
	if objVal == nil {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("spl_object_id(): Argument #1 ($object) must be of type object"), "TypeError")
	}
	val, ok := objVal.(data.Value)
	objType := data.Object{}
	if !ok || !objType.Is(val) {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("spl_object_id(): Argument #1 ($object) must be of type object"), "TypeError")
	}

	rv := reflect.ValueOf(objVal)
	for rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		return data.NewIntValue(int(rv.Pointer())), nil
	}

	return nil, data.NewErrorThrowByName(nil, fmt.Errorf("spl_object_id(): Argument #1 ($object) must be of type object"), "TypeError")
}

func (f *SplObjectIdFunction) GetName() string {
	return "spl_object_id"
}

func (f *SplObjectIdFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
	}
}

func (f *SplObjectIdFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, data.NewBaseType("object")),
	}
}

func (f *SplObjectIdFunction) GetReturnType() data.Types {
	return data.NewBaseType("int")
}
