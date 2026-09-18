package httpfoundation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	fqnRequestStack    = "Symfony\\Component\\HttpFoundation\\RequestStack"
	fqnSessionNotFound = "Symfony\\Component\\HttpFoundation\\Exception\\SessionNotFoundException"
)

// RequestStackClass 实现 Symfony\Component\HttpFoundation\RequestStack。
type RequestStackClass struct {
	node.Node
	methods    map[string]data.Method
	methodList []data.Method
}

func NewRequestStackClass() data.ClassStmt {
	c := &RequestStackClass{}
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{param("requests", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("requests", 0, nil)},
			nil, requestStackConstruct),
		pubMethod("push",
			[]data.GetValue{param("request", 0, nil, nil)},
			[]data.Variable{variable("request", 0, nil)},
			nil, requestStackPush),
		pubMethod("pop", nil, nil, nil, requestStackPop),
		pubMethod("getCurrentRequest", nil, nil, nil, requestStackGetCurrent),
		pubMethod("getMainRequest", nil, nil, nil, requestStackGetMain),
		pubMethod("getParentRequest", nil, nil, nil, requestStackGetParent),
		pubMethod("getSession", nil, nil, nil, requestStackGetSession),
		pubMethod("resetRequestFormats", nil, nil, nil, requestStackResetFormats),
	}
	c.methods = indexMethods(list)
	c.methodList = list
	return c
}

func (c *RequestStackClass) GetName() string         { return fqnRequestStack }
func (c *RequestStackClass) GetExtend() *string      { return nil }
func (c *RequestStackClass) GetImplements() []string { return nil }
func (c *RequestStackClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.GetPropertyList() {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *RequestStackClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "requests", "private", false, data.NewArrayValue(nil)),
	}
}
func (c *RequestStackClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *RequestStackClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *RequestStackClass) GetMethod(name string) (data.Method, bool) {
	return getIndexedMethod(c.methods, name)
}
func (c *RequestStackClass) GetMethods() []data.Method { return c.methodList }

func requestStackRequests(cv *data.ClassValue) *data.ArrayValue {
	if cv == nil {
		return data.NewArrayValue(nil).(*data.ArrayValue)
	}
	v, _ := cv.GetProperty("requests")
	if arr, ok := v.(*data.ArrayValue); ok && arr != nil {
		return arr
	}
	arr := data.NewArrayValue(nil).(*data.ArrayValue)
	_ = cv.SetProperty("requests", arr)
	return arr
}

func requestStackConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := bagClassValue(ctx)
	arr := data.NewArrayValue(nil).(*data.ArrayValue)
	_ = cv.SetProperty("requests", arr)
	raw, _ := ctx.GetIndexValue(0)
	if raw == nil || isNull(raw) {
		return data.NewNullValue(), nil
	}
	if ctl := requestStackPushAll(cv, raw); ctl != nil {
		return nil, ctl
	}
	return data.NewNullValue(), nil
}

func requestStackPushAll(cv *data.ClassValue, raw data.Value) data.Control {
	switch t := raw.(type) {
	case *data.ArrayValue:
		for _, z := range t.List {
			if z == nil || z.Value == nil || isNull(z.Value) {
				continue
			}
			requestStackRequests(cv).List = append(requestStackRequests(cv).List, data.NewZVal(z.Value))
		}
	default:
		m, err := valueToAssocMap(raw)
		if err != nil {
			return nil
		}
		for _, v := range m {
			if v == nil || isNull(v) {
				continue
			}
			requestStackRequests(cv).List = append(requestStackRequests(cv).List, data.NewZVal(v))
		}
	}
	return nil
}

func requestStackPush(ctx data.Context) (data.GetValue, data.Control) {
	cv := bagClassValue(ctx)
	v, _ := ctx.GetIndexValue(0)
	if v != nil && !isNull(v) {
		arr := requestStackRequests(cv)
		arr.List = append(arr.List, data.NewZVal(v))
	}
	return data.NewNullValue(), nil
}

func requestStackPop(ctx data.Context) (data.GetValue, data.Control) {
	arr := requestStackRequests(bagClassValue(ctx))
	if len(arr.List) == 0 {
		return data.NewNullValue(), nil
	}
	last := arr.List[len(arr.List)-1]
	arr.List = arr.List[:len(arr.List)-1]
	if last == nil || last.Value == nil {
		return data.NewNullValue(), nil
	}
	return last.Value, nil
}

func requestStackGetCurrent(ctx data.Context) (data.GetValue, data.Control) {
	arr := requestStackRequests(bagClassValue(ctx))
	if len(arr.List) == 0 {
		return data.NewNullValue(), nil
	}
	last := arr.List[len(arr.List)-1]
	if last == nil || last.Value == nil {
		return data.NewNullValue(), nil
	}
	return last.Value, nil
}

func requestStackGetMain(ctx data.Context) (data.GetValue, data.Control) {
	arr := requestStackRequests(bagClassValue(ctx))
	if len(arr.List) == 0 {
		return data.NewNullValue(), nil
	}
	first := arr.List[0]
	if first == nil || first.Value == nil {
		return data.NewNullValue(), nil
	}
	return first.Value, nil
}

func requestStackGetParent(ctx data.Context) (data.GetValue, data.Control) {
	arr := requestStackRequests(bagClassValue(ctx))
	pos := len(arr.List) - 2
	if pos < 0 {
		return data.NewNullValue(), nil
	}
	item := arr.List[pos]
	if item == nil || item.Value == nil {
		return data.NewNullValue(), nil
	}
	return item.Value, nil
}

func requestStackGetSession(ctx data.Context) (data.GetValue, data.Control) {
	cur, ctl := requestStackGetCurrent(ctx)
	if ctl != nil {
		return nil, ctl
	}
	req, ok := cur.(*data.ClassValue)
	if !ok || req == nil {
		return nil, throwNamed(fqnSessionNotFound, "There is currently no session available.")
	}
	has, ctl := callObjMethod(req, "hasSession")
	if ctl != nil {
		return nil, ctl
	}
	if !valueIsTruthy(has) {
		return nil, throwNamed(fqnSessionNotFound, "There is currently no session available.")
	}
	return callObjMethod(req, "getSession")
}

func requestStackResetFormats(ctx data.Context) (data.GetValue, data.Control) {
	resetRequestFormatsState()
	return data.NewNullValue(), nil
}
