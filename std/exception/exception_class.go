package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

func NewExceptionClass() *ExceptionClass {
	source := &Exception{}

	return &ExceptionClass{
		exception:        &ExceptionExceptionMethod{source},
		error:            &ExceptionErrorMethod{source},
		getMessage:       &ExceptionGetMessageMethod{source},
		getCode:          &ExceptionGetCodeMethod{source},
		getPrevious:      &ExceptionGetPreviousMethod{source},
		getTrace:         &ExceptionGetTraceMethod{source},
		getFile:          &ExceptionGetFileMethod{source},
		getLine:          &ExceptionGetLineMethod{source},
		getTraceAsString: &ExceptionGetTraceAsStringMethod{source},
		// 对齐 PHP Exception 受保护属性，供子类 $this->message .= ... 等使用
		propMessage:  node.NewProperty(nil, "message", "protected", false, data.NewStringValue("")),
		propCode:     node.NewProperty(nil, "code", "protected", false, data.NewIntValue(0)),
		propFile:     node.NewProperty(nil, "file", "protected", false, data.NewStringValue("")),
		propLine:     node.NewProperty(nil, "line", "protected", false, data.NewIntValue(0)),
		propPrevious: node.NewProperty(nil, "previous", "protected", false, data.NewNullValue()),
	}
}

type ExceptionClass struct {
	node.Node
	exception        data.Method
	error            data.Method
	getMessage       *ExceptionGetMessageMethod
	getCode          *ExceptionGetCodeMethod
	getPrevious      *ExceptionGetPreviousMethod
	getTrace         *ExceptionGetTraceMethod
	getFile          *ExceptionGetFileMethod
	getLine          *ExceptionGetLineMethod
	getTraceAsString *ExceptionGetTraceAsStringMethod
	propMessage      data.Property
	propCode         data.Property
	propFile         data.Property
	propLine         data.Property
	propPrevious     data.Property
}

func (s *ExceptionClass) AsString() string {
	return s.getMessage.source.msg
}

func (s *ExceptionClass) IsThrow() bool {
	//TODO implement me
	panic("implement me")
}

func (s *ExceptionClass) GetError() *data.Error {
	//TODO implement me
	panic("implement me")
}

func (s *ExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	newS := *s
	return data.NewClassValue(&newS, ctx), nil
}

func (s *ExceptionClass) GetName() string {
	return "Exception"
}

func (s *ExceptionClass) GetExtend() *string {
	return nil
}

func (s *ExceptionClass) GetImplements() []string {
	// PHP 中 Exception 实现 Throwable 接口
	return []string{"Throwable"}
}

func (s *ExceptionClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "message":
		return s.propMessage, true
	case "code":
		return s.propCode, true
	case "file":
		return s.propFile, true
	case "line":
		return s.propLine, true
	case "previous":
		return s.propPrevious, true
	}
	return nil, false
}

func (s *ExceptionClass) GetPropertyList() []data.Property {
	return []data.Property{s.propMessage, s.propCode, s.propFile, s.propLine, s.propPrevious}
}

func (s *ExceptionClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case token.ConstructName:
		return s.exception, true
	case "error":
		return s.error, true
	case "getMessage":
		return s.getMessage, true
	case "getCode":
		return s.getCode, true
	case "getPrevious":
		return s.getPrevious, true
	case "getTrace":
		return s.getTrace, true
	case "getFile":
		return s.getFile, true
	case "getLine":
		return s.getLine, true
	case "getTraceAsString":
		return s.getTraceAsString, true
	}
	return nil, false
}

func (s *ExceptionClass) GetMethods() []data.Method {
	return []data.Method{
		s.error,
		s.getMessage,
		s.getCode,
		s.getPrevious,
		s.getTrace,
		s.getFile,
		s.getLine,
		s.getTraceAsString,
	}
}

func (s *ExceptionClass) GetConstruct() data.Method {
	return s.exception
}

// instanceObjectValue 从方法调用上下文取出当前异常实例的属性存储
func instanceObjectValue(ctx data.Context) *data.ObjectValue {
	switch c := ctx.(type) {
	case *data.ClassMethodContext:
		return c.ObjectValue
	case *data.ThisValue:
		return c.ObjectValue
	case *data.ClassValue:
		return c.ObjectValue
	default:
		return nil
	}
}

func setInstanceProperty(ctx data.Context, name string, value data.Value) {
	if sp, ok := ctx.(data.SetProperty); ok {
		_ = sp.SetProperty(name, value)
		return
	}
	if ov := instanceObjectValue(ctx); ov != nil {
		_ = ov.SetProperty(name, value)
	}
}

func instancePropertyString(ctx data.Context, name string) (string, bool) {
	ov := instanceObjectValue(ctx)
	if ov == nil || !ov.HasProperty(name) {
		return "", false
	}
	v, _ := ov.GetProperty(name)
	if v == nil {
		return "", false
	}
	return v.AsString(), true
}

func instancePropertyInt(ctx data.Context, name string) (int, bool) {
	ov := instanceObjectValue(ctx)
	if ov == nil || !ov.HasProperty(name) {
		return 0, false
	}
	v, _ := ov.GetProperty(name)
	if v == nil {
		return 0, false
	}
	if iv, ok := v.(data.AsInt); ok {
		n, err := iv.AsInt()
		if err == nil {
			return n, true
		}
	}
	return 0, false
}
