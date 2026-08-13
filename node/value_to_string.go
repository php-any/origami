package node

import (
	"github.com/php-any/origami/data"
)

// ValueToDisplayString 在“字符串上下文”中将值转为字符串（如拼接、echo）。
// 若为类实例且定义了 __toString()，则调用 __toString() 并返回其返回值转成的字符串；
// 否则使用 Value.AsString()。
func ValueToDisplayString(ctx data.Context, v data.GetValue) (string, data.Control) {
	if v == nil {
		return "", nil
	}
	// 类实例或 $this：尝试调用 __toString()
	if obj, ok := v.(data.GetMethod); ok {
		if toString, has := obj.GetMethod("__toString"); has {
			if objCtx, ok := v.(data.Context); ok {
				result, acl := toString.Call(objCtx.CreateContext(toString.GetVariables()))
				if acl != nil {
					return "", acl
				}
				if result != nil {
					if val, ok := result.(data.Value); ok {
						return val.AsString(), nil
					}
				}
				return "", nil
			}
		}
	}
	// 其他类型直接 AsString
	if val, ok := v.(data.Value); ok {
		return val.AsString(), nil
	}
	return "", nil
}
