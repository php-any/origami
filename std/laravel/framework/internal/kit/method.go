package kit

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// InstanceMethod 注册期构造实例方法（调用期仍直接 fn(ctx)，无热路径税）。
func InstanceMethod(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) data.Method {
	return buildMethod(name, params, -1, fn, false, false)
}

// InstanceMethodOpt 带可选参数默认值的实例方法。
func InstanceMethodOpt(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) data.Method {
	return buildMethod(name, params, optionalFrom, fn, false, false)
}

// StaticMethod 注册期构造静态方法。
func StaticMethod(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) data.Method {
	return buildMethod(name, params, optionalFrom, fn, true, false)
}

// StaticMethodRef 第一个参数按引用传递（如 Arr::set）。
func StaticMethodRef(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) data.Method {
	return buildMethod(name, params, optionalFrom, fn, true, true)
}

// InstanceMethodVariadic 注册变参实例方法：最后一个形参吸收全部实参（对应 PHP 的 `...$values`）。
// 调用方传 0 个实参时该形参得到空数组，而不是「缺少值和默认值」错误。
// 变参形参是 *node.Parameters，canFastPositionalBind 会拒绝它 → 走慢路径绑定，
// 因此不给无变参的快路径加固定税。
func InstanceMethodVariadic(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) data.Method {
	return buildMethodVariadic(name, params, fn)
}

func buildMethodVariadic(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) data.Method {
	ps := make([]data.GetValue, len(params))
	vs := make([]data.Variable, len(params))
	for i, p := range params {
		if i == len(params)-1 {
			ps[i] = node.NewParameters(nil, p, i, nil, nil)
		} else {
			ps[i] = node.NewParameter(nil, p, i, nil, nil)
		}
		vs[i] = node.NewVariable(nil, p, i, nil)
	}
	return &method{name: name, params: ps, vars: vs, fn: fn}
}

func buildMethod(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control), static, firstByRef bool) data.Method {
	ps := make([]data.GetValue, len(params))
	vs := make([]data.Variable, len(params))
	for i, p := range params {
		var def data.GetValue
		if optionalFrom >= 0 && i >= optionalFrom {
			def = data.NewNullValue()
		}
		if i == 0 && firstByRef {
			ps[i] = node.NewParameterReference(nil, p, i, def, nil)
		} else {
			ps[i] = node.NewParameter(nil, p, i, def, nil)
		}
		vs[i] = node.NewVariable(nil, p, i, nil)
	}
	return &method{name: name, params: ps, vars: vs, fn: fn, static: static}
}

type method struct {
	name   string
	params []data.GetValue
	vars   []data.Variable
	fn     func(data.Context) (data.GetValue, data.Control)
	static bool
}

func (m *method) Call(ctx data.Context) (data.GetValue, data.Control) { return m.fn(ctx) }
func (m *method) GetName() string                                     { return m.name }
func (m *method) GetModifier() data.Modifier                          { return data.ModifierPublic }
func (m *method) GetIsStatic() bool                                   { return m.static }
func (m *method) GetParams() []data.GetValue                          { return m.params }
func (m *method) GetVariables() []data.Variable                       { return m.vars }
func (m *method) GetReturnType() data.Types                           { return nil }
