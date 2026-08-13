package httpkernel

import (
	"fmt"
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// kernelMethod 是共享的方法适配器。
type kernelMethod struct {
	name     string
	modifier data.Modifier
	static   bool
	params   []data.GetValue
	vars     []data.Variable
	ret      data.Types
	fn       func(ctx data.Context) (data.GetValue, data.Control)
}

func (m *kernelMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return m.fn(ctx)
}
func (m *kernelMethod) GetName() string               { return m.name }
func (m *kernelMethod) GetModifier() data.Modifier    { return m.modifier }
func (m *kernelMethod) GetIsStatic() bool             { return m.static }
func (m *kernelMethod) GetParams() []data.GetValue    { return m.params }
func (m *kernelMethod) GetVariables() []data.Variable { return m.vars }
func (m *kernelMethod) GetReturnType() data.Types     { return m.ret }

func pubMethod(name string, params []data.GetValue, vars []data.Variable, ret data.Types, fn func(ctx data.Context) (data.GetValue, data.Control)) data.Method {
	return &kernelMethod{
		name:     name,
		modifier: data.ModifierPublic,
		params:   params,
		vars:     vars,
		ret:      ret,
		fn:       fn,
	}
}

func param(name string, index int, def data.GetValue, ty data.Types) data.GetValue {
	return node.NewParameter(nil, name, index, def, ty)
}

func variable(name string, index int, ty data.Types) data.Variable {
	return node.NewVariable(nil, name, index, ty)
}

func thisClassValue(ctx data.Context) *data.ClassValue {
	switch c := ctx.(type) {
	case *data.ClassMethodContext:
		return c.ClassValue
	case *data.ClassValue:
		return c
	default:
		return nil
	}
}

func getState(ctx data.Context) (*kernelState, data.Control) {
	cv := thisClassValue(ctx)
	if cv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: 无法获取 Kernel 实例"))
	}
	if s, ok := cv.GetSource().(*kernelState); ok && s != nil {
		return s, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: Kernel 状态未初始化"))
}

func indexValue(ctx data.Context, index int) data.Value {
	v, ok := ctx.GetIndexValue(index)
	if !ok {
		return nil
	}
	return v
}

func valueAsString(v data.Value) string {
	if v == nil {
		return ""
	}
	return v.AsString()
}

func isTruthy(v data.Value) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(data.AsBool); ok {
		bv, err := b.AsBool()
		return err == nil && bv
	}
	s := v.AsString()
	return s != "" && s != "0" && s != "false"
}

func stringsToArrayValue(ss []string) *data.ArrayValue {
	vals := make([]data.Value, len(ss))
	for i, s := range ss {
		vals[i] = data.NewStringValue(s)
	}
	return data.NewArrayValue(vals).(*data.ArrayValue)
}

func stringListFromValue(v data.Value) []string {
	if v == nil {
		return nil
	}
	switch arr := v.(type) {
	case *data.ArrayValue:
		out := make([]string, 0, len(arr.List))
		for _, z := range arr.List {
			if z == nil || z.Value == nil {
				continue
			}
			out = append(out, z.Value.AsString())
		}
		return out
	case *data.ObjectValue:
		out := make([]string, 0)
		arr.RangeProperties(func(_ string, value data.Value) bool {
			if value != nil {
				out = append(out, value.AsString())
			}
			return true
		})
		return out
	default:
		s := v.AsString()
		if s == "" {
			return nil
		}
		return []string{s}
	}
}

func stringMapFromValue(v data.Value) map[string]string {
	out := make(map[string]string)
	if v == nil {
		return out
	}
	switch arr := v.(type) {
	case *data.ArrayValue:
		for i, z := range arr.List {
			if z == nil || z.Value == nil {
				continue
			}
			key := z.Name
			if key == "" {
				key = strconv.Itoa(i)
			}
			out[key] = z.Value.AsString()
		}
	case *data.ObjectValue:
		arr.RangeProperties(func(key string, value data.Value) bool {
			if value != nil {
				out[key] = value.AsString()
			}
			return true
		})
	}
	return out
}

func groupsFromValue(v data.Value) map[string][]string {
	out := make(map[string][]string)
	if v == nil {
		return out
	}
	switch arr := v.(type) {
	case *data.ArrayValue:
		for i, z := range arr.List {
			if z == nil || z.Value == nil {
				continue
			}
			key := z.Name
			if key == "" {
				key = strconv.Itoa(i)
			}
			out[key] = stringListFromValue(z.Value)
		}
	case *data.ObjectValue:
		arr.RangeProperties(func(key string, value data.Value) bool {
			out[key] = stringListFromValue(value)
			return true
		})
	}
	return out
}

func groupsToArrayValue(groups map[string][]string) *data.ArrayValue {
	list := make([]*data.ZVal, 0, len(groups))
	for k, items := range groups {
		list = append(list, &data.ZVal{Name: k, Value: stringsToArrayValue(items)})
	}
	return &data.ArrayValue{List: list}
}

func aliasesToArrayValue(aliases map[string]string) *data.ArrayValue {
	list := make([]*data.ZVal, 0, len(aliases))
	for k, v := range aliases {
		list = append(list, &data.ZVal{Name: k, Value: data.NewStringValue(v)})
	}
	return &data.ArrayValue{List: list}
}

func syncProperties(cv *data.ClassValue, s *kernelState) {
	if cv == nil || s == nil {
		return
	}
	if s.app != nil {
		_ = cv.SetProperty("app", s.app)
	}
	if s.router != nil {
		_ = cv.SetProperty("router", s.router)
	}
	_ = cv.SetProperty("bootstrappers", stringsToArrayValue(s.bootstrappers))
	_ = cv.SetProperty("middleware", stringsToArrayValue(s.middleware))
	_ = cv.SetProperty("middlewareGroups", groupsToArrayValue(s.middlewareGroups))
	_ = cv.SetProperty("middlewareAliases", aliasesToArrayValue(s.middlewareAliases))
	_ = cv.SetProperty("middlewarePriority", stringsToArrayValue(s.middlewarePriority))
}

// callObjectMethod 在 ClassValue 上调用实例方法；失败时原样返回 control。
func callObjectMethod(obj data.Value, name string, args ...data.Value) (data.GetValue, data.Control) {
	return callObjectMethodInContext(nil, obj, name, args...)
}

// callObjectMethodInContext 让常驻 Application/Kernel 的方法调用使用当前请求 VM，
// 但继续共享对象属性和类定义。ctx 为 nil 时保持对象原有上下文。
func callObjectMethodInContext(ctx data.Context, obj data.Value, name string, args ...data.Value) (data.GetValue, data.Control) {
	var cv *data.ClassValue
	switch value := obj.(type) {
	case *data.ClassValue:
		cv = value
	case *data.ThisValue:
		cv = value.ClassValue
	}
	ok := cv != nil
	if !ok || cv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: 期望对象以调用 %s", name))
	}
	if ctx != nil {
		cv = &data.ClassValue{
			Context:     ctx.CreateBaseContext(),
			ObjectValue: cv.ObjectValue,
			Class:       cv.Class,
		}
	}
	method, exists := cv.GetMethod(name)
	if !exists || method == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpkernel: 方法 %s 不存在", name))
	}
	fnCtx := cv.CreateContext(method.GetVariables())
	vars := method.GetVariables()
	for i, arg := range args {
		if i >= len(vars) {
			break
		}
		if arg == nil {
			arg = data.NewNullValue()
		}
		if acl := fnCtx.SetVariableValue(vars[i], arg); acl != nil {
			return nil, acl
		}
	}
	return method.Call(fnCtx)
}

func containsString(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func indexOfString(list []string, item string) int {
	for i, v := range list {
		if v == item {
			return i
		}
	}
	return -1
}

func existingNames(v data.Value) []string {
	if v == nil {
		return nil
	}
	if arr, ok := v.(*data.ArrayValue); ok {
		out := make([]string, 0, len(arr.List))
		for _, z := range arr.List {
			if z != nil && z.Value != nil {
				out = append(out, z.Value.AsString())
			}
		}
		return out
	}
	s := v.AsString()
	if s == "" {
		return nil
	}
	return []string{s}
}
