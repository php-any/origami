package container

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

func classValueFromCtx(ctx data.Context) (*data.ClassValue, bool) {
	switch c := ctx.(type) {
	case *data.ClassValue:
		return c, true
	case *data.ClassMethodContext:
		return c.ClassValue, true
	default:
		return nil, false
	}
}

func engineFromCtx(ctx data.Context) (*Engine, data.Control) {
	cv, ok := classValueFromCtx(ctx)
	if !ok {
		return nil, utils.NewThrow(errors.New("Container 方法必须在 Container 实例上调用"))
	}
	c, ok := cv.Class.(*ContainerClass)
	if !ok {
		return nil, utils.NewThrow(errors.New("Container 方法必须在 Container 实例上调用"))
	}
	if c.engine == nil {
		return nil, utils.NewThrow(errors.New("Container 引擎未初始化"))
	}
	return c.engine, nil
}

func stringArg(ctx data.Context, index int) (string, data.Control) {
	v, ok := ctx.GetIndexValue(index)
	if !ok {
		return "", utils.NewThrow(errors.New("缺少参数"))
	}
	s, ok := v.(data.AsString)
	if !ok {
		return "", utils.NewThrow(errors.New("参数必须是字符串"))
	}
	return s.AsString(), nil
}

func optionalStringArg(ctx data.Context, index int) (string, bool, data.Control) {
	v, ok := ctx.GetIndexValue(index)
	if !ok {
		return "", false, nil
	}
	if _, isNull := v.(*data.NullValue); isNull {
		return "", false, nil
	}
	s, ok := v.(data.AsString)
	if !ok {
		return "", false, utils.NewThrow(errors.New("参数必须是字符串"))
	}
	return s.AsString(), true, nil
}

func resolveConcreteArg(ctx data.Context) (abstract, concrete string, factory data.GetValue, acl data.Control) {
	abstract, acl = stringArg(ctx, 0)
	if acl != nil {
		return "", "", nil, acl
	}
	concrete = abstract

	v, ok := ctx.GetIndexValue(1)
	if !ok {
		return abstract, concrete, nil, nil
	}
	if _, isNull := v.(*data.NullValue); isNull {
		return abstract, concrete, nil, nil
	}
	if stmt, ok := callableValue(v); ok && stmt != nil {
		return abstract, "", v, nil
	}
	if s, ok := v.(data.AsString); ok {
		return abstract, s.AsString(), nil, nil
	}
	return "", "", nil, utils.NewThrow(errors.New("concrete 必须是类名或闭包"))
}
