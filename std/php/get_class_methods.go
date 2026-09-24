package php

import (
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GetClassMethodsFunction 实现 get_class_methods 函数
// get_class_methods(object|string $object_or_class): array
// 返回指定类中已定义的方法名数组（含继承链，不区分大小写去重）。
type GetClassMethodsFunction struct{}

func NewGetClassMethodsFunction() data.FuncStmt {
	return &GetClassMethodsFunction{}
}

func (f *GetClassMethodsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	objectOrClassValue, ok := ctx.GetIndexValue(0)
	if !ok || objectOrClassValue == nil {
		return data.NewArrayValue(nil), nil
	}

	vm := ctx.GetVM()

	var classStmt data.ClassStmt

	// 对象实例（ClassValue）→ 直接取类
	if cv, ok := objectOrClassValue.(*data.ClassValue); ok {
		classStmt = cv.Class
	} else if tv, ok := objectOrClassValue.(*data.ThisValue); ok {
		classStmt = tv.ClassValue.Class
	} else if gv, ok := objectOrClassValue.(data.GetName); ok {
		// 类名
		name := gv.GetName()
		if stmt, acl := vm.GetOrLoadClass(name); acl == nil && stmt != nil {
			classStmt = stmt
		}
	} else if sv, ok := objectOrClassValue.(*data.StringValue); ok {
		if stmt, acl := vm.GetOrLoadClass(sv.Value); acl == nil && stmt != nil {
			classStmt = stmt
		}
	}

	if classStmt == nil {
		return data.NewArrayValue(nil), nil
	}

	// 收集方法名（GetMethods 含实例方法、静态方法与构造函数）
	methods := make(map[string]bool)
	collect := func(stmt data.ClassStmt) {
		for _, m := range stmt.GetMethods() {
			if m != nil {
				methods[m.GetName()] = true
			}
		}
	}

	collect(classStmt)
	last := classStmt
	for last.GetExtend() != nil {
		ext := last.GetExtend()
		parent, acl := vm.GetOrLoadClass(*ext)
		if acl != nil || parent == nil {
			break
		}
		collect(parent)
		last = parent
	}

	names := make([]string, 0, len(methods))
	for name := range methods {
		names = append(names, name)
	}
	sort.Strings(names)

	list := make([]data.Value, len(names))
	for i, n := range names {
		list[i] = data.NewStringValue(n)
	}
	return data.NewArrayValue(list), nil
}

func (f *GetClassMethodsFunction) GetName() string {
	return "get_class_methods"
}

var getClassMethodsFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "object_or_class", 0, nil, nil),
}

func (f *GetClassMethodsFunction) GetParams() []data.GetValue {
	return getClassMethodsFunctionGetParams
}

var getClassMethodsFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "object_or_class", 0, data.NewBaseType("object|string")),
}

func (f *GetClassMethodsFunction) GetVariables() []data.Variable {
	return getClassMethodsFunctionGetVariables
}
