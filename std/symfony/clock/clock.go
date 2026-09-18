package clock

import (
	"strings"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const nativeClockName = "Symfony\\Component\\Clock\\NativeClock"

// NativeClockClass 对齐 Symfony\Component\Clock\NativeClock（v8.1）。
// 不注册 Clock 静态门面，避免挡住 PHP 的 wither/sleep/Modifier。
// now() 返回 DateTimeImmutable（PSR ClockInterface）；不注册 DatePoint。
type NativeClockClass struct {
	node.Node
	methods    map[string]data.Method
	methodList []data.Method
	props      []data.Property
	propIndex  map[string]data.Property
}

var nativeClockOnce = &NativeClockClass{}

func NewNativeClockClass() data.ClassStmt {
	c := nativeClockOnce
	if c.methods != nil {
		return c
	}

	nullableTz := data.NewNullableType(data.NewUnionType([]data.Types{
		data.NewBaseType("string"),
		data.NewBaseType("DateTimeZone"),
	}))
	tzUnion := data.NewUnionType([]data.Types{
		data.NewBaseType("string"),
		data.NewBaseType("DateTimeZone"),
	})
	secondsUnion := data.NewUnionType([]data.Types{
		data.NewBaseType("float"),
		data.NewBaseType("int"),
	})

	tzProp := node.NewProperty(nil, "timezone", "private", false, nil, data.NewBaseType("DateTimeZone"))
	c.props = []data.Property{tzProp}
	c.propIndex = map[string]data.Property{"timezone": tzProp}

	add := func(name string, params []clockParam, ret data.Types, fn func(data.Context) (data.GetValue, data.Control)) {
		m := &clockMethod{name: name, params: params, ret: ret, fn: fn}
		c.methods[strings.ToLower(name)] = m
		c.methodList = append(c.methodList, m)
	}
	c.methods = map[string]data.Method{}
	add("__construct", []clockParam{{
		name: "timezone",
		def:  data.NewNullValue(),
		typ:  nullableTz,
	}}, nil, clockConstruct)
	add("now", nil, data.NewBaseType("DateTimeImmutable"), clockNow)
	add("sleep", []clockParam{{name: "seconds", typ: secondsUnion}}, nil, clockSleep)
	add("withTimeZone", []clockParam{{name: "timezone", typ: tzUnion}}, data.NewBaseType(nativeClockName), clockWithTimeZone)
	return c
}

func (c *NativeClockClass) GetName() string    { return nativeClockName }
func (c *NativeClockClass) GetExtend() *string { return nil }
func (c *NativeClockClass) GetImplements() []string {
	return []string{
		"Symfony\\Component\\Clock\\ClockInterface",
		"Psr\\Clock\\ClockInterface",
	}
}
func (c *NativeClockClass) GetProperty(name string) (data.Property, bool) {
	p, ok := c.propIndex[name]
	return p, ok
}
func (c *NativeClockClass) GetPropertyList() []data.Property { return c.props }
func (c *NativeClockClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *NativeClockClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

// GetMethod PHP 方法名不区分大小写。
func (c *NativeClockClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *NativeClockClass) GetMethods() []data.Method { return c.methodList }
func (c *NativeClockClass) GetStaticMethod(name string) (data.Method, bool) {
	m, ok := c.GetMethod(name)
	if !ok || m == nil || !m.GetIsStatic() {
		return nil, false
	}
	return m, true
}

type clockParam struct {
	name string
	def  data.GetValue
	typ  data.Types
}

type clockMethod struct {
	name   string
	params []clockParam
	ret    data.Types
	fn     func(data.Context) (data.GetValue, data.Control)
}

func (m *clockMethod) Call(ctx data.Context) (data.GetValue, data.Control) { return m.fn(ctx) }
func (m *clockMethod) GetName() string                                     { return m.name }
func (m *clockMethod) GetModifier() data.Modifier                          { return data.ModifierPublic }
func (m *clockMethod) GetIsStatic() bool                                   { return false }
func (m *clockMethod) GetReturnType() data.Types                           { return m.ret }
func (m *clockMethod) GetParams() []data.GetValue {
	out := make([]data.GetValue, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewParameter(nil, p.name, i, p.def, p.typ)
	}
	return out
}
func (m *clockMethod) GetVariables() []data.Variable {
	out := make([]data.Variable, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewVariable(nil, p.name, i, p.typ)
	}
	return out
}

func clockSelf(ctx data.Context) *data.ClassValue {
	if c, ok := ctx.(*data.ClassMethodContext); ok {
		return c.ClassValue
	}
	return nil
}

func isNullValue(v data.Value) bool {
	if v == nil {
		return true
	}
	_, ok := v.(*data.NullValue)
	return ok
}

func asClassValue(v data.Value) *data.ClassValue {
	if cv, ok := v.(*data.ClassValue); ok {
		return cv
	}
	if tv, ok := v.(*data.ThisValue); ok && tv.ClassValue != nil {
		return tv.ClassValue
	}
	return nil
}

func isDateTimeZoneValue(v data.Value) bool {
	cv := asClassValue(v)
	if cv == nil || cv.Class == nil {
		return false
	}
	cls := cv.Class
	for cls != nil {
		name := cls.GetName()
		if name == "DateTimeZone" || strings.HasSuffix(name, "\\DateTimeZone") {
			return true
		}
		ext := cls.GetExtend()
		if ext == nil {
			return false
		}
		vm := cv.GetVM()
		if vm == nil {
			return false
		}
		next, ok := vm.GetClass(*ext)
		if !ok || next == nil {
			return false
		}
		cls = next
	}
	return false
}

func clockTimezoneValue(ctx data.Context) (data.Value, data.Control) {
	self := clockSelf(ctx)
	if self == nil {
		return nil, nil
	}
	v, ctl := self.GetProperty("timezone")
	if ctl != nil {
		return nil, ctl
	}
	if isNullValue(v) {
		return nil, nil
	}
	return v, nil
}

func defaultTimezoneString(ctx data.Context) (string, data.Control) {
	if vm := ctx.GetVM(); vm != nil {
		if fn, ok := vm.GetFunc("date_default_timezone_get"); ok && fn != nil {
			ret, ctl := fn.Call(ctx.CreateContext(fn.GetVariables()))
			if ctl != nil {
				return "", ctl
			}
			if v, ok := ret.(data.Value); ok && !isNullValue(v) {
				if s := v.AsString(); s != "" {
					return s, nil
				}
			}
		}
	}
	name := time.Local.String()
	if name == "" || name == "Local" {
		return "UTC", nil
	}
	return name, nil
}

func newDateTimeZone(ctx data.Context, id string) (data.Value, data.Control) {
	vm := ctx.GetVM()
	if vm == nil {
		return data.NewStringValue(id), nil
	}
	cls, ok := vm.GetClass("DateTimeZone")
	if !ok || cls == nil {
		return data.NewStringValue(id), nil
	}
	gv, acl := cls.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}
	obj, ok := gv.(*data.ClassValue)
	if !ok || obj == nil {
		return data.NewStringValue(id), nil
	}
	ctor := cls.GetConstruct()
	if ctor == nil {
		if ctl := obj.SetProperty("name", data.NewStringValue(id)); ctl != nil {
			return nil, ctl
		}
		return obj, nil
	}
	vars := ctor.GetVariables()
	cctx := obj.CreateContext(vars)
	if len(vars) > 0 {
		if ctl := cctx.SetVariableValue(vars[0], data.NewStringValue(id)); ctl != nil {
			return nil, ctl
		}
	}
	if _, ctl := ctor.Call(cctx); ctl != nil {
		return nil, ctl
	}
	return obj, nil
}

// resolveTimezone 对齐 NativeClock：null → date_default_timezone_get()；string → new DateTimeZone；对象原样使用。
func resolveTimezone(ctx data.Context, tz data.Value, allowDefault bool) (data.Value, data.Control) {
	if isNullValue(tz) {
		if !allowDefault {
			return newDateTimeZone(ctx, "")
		}
		id, ctl := defaultTimezoneString(ctx)
		if ctl != nil {
			return nil, ctl
		}
		return newDateTimeZone(ctx, id)
	}
	if isDateTimeZoneValue(tz) {
		return tz, nil
	}
	return newDateTimeZone(ctx, tz.AsString())
}

func cloneClock(ctx data.Context, self *data.ClassValue) (*data.ClassValue, data.Control) {
	cls := data.ClassStmt(nativeClockOnce)
	base := ctx.CreateBaseContext()
	if self != nil {
		if self.Class != nil {
			cls = self.Class
		}
		if self.Context != nil {
			base = self.Context
		}
	}
	clone := data.NewClassValue(cls, base)
	if self == nil {
		return clone, nil
	}
	var copyCtl data.Control
	self.RangeProperties(func(key string, v data.Value) bool {
		var next data.Value
		switch val := v.(type) {
		case *data.ArrayValue:
			next = data.DeepCloneArrayValue(val)
		case *data.ObjectValue:
			next = data.DeepCloneObjectValue(val)
		default:
			next = v
		}
		copyCtl = clone.SetProperty(key, next)
		return copyCtl == nil
	})
	if copyCtl != nil {
		return nil, copyCtl
	}
	return clone, nil
}

func clockConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := clockSelf(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	tz, ok := ctx.GetIndexValue(0)
	if !ok {
		tz = nil
	}
	// NativeClock.php：$timezone ??= date_default_timezone_get()；string 走 withTimeZone。
	if isNullValue(tz) || !isDateTimeZoneValue(tz) {
		cloned, ctl := applyWithTimeZone(ctx, cv, tz, true)
		if ctl != nil {
			return nil, ctl
		}
		resolved, ctl := cloned.GetProperty("timezone")
		if ctl != nil {
			return nil, ctl
		}
		if ctl := cv.SetProperty("timezone", resolved); ctl != nil {
			return nil, ctl
		}
		return data.NewNullValue(), nil
	}
	if ctl := cv.SetProperty("timezone", tz); ctl != nil {
		return nil, ctl
	}
	return data.NewNullValue(), nil
}

func applyWithTimeZone(ctx data.Context, self *data.ClassValue, tz data.Value, allowDefault bool) (*data.ClassValue, data.Control) {
	resolved, ctl := resolveTimezone(ctx, tz, allowDefault)
	if ctl != nil {
		return nil, ctl
	}
	clone, ctl := cloneClock(ctx, self)
	if ctl != nil {
		return nil, ctl
	}
	if ctl := clone.SetProperty("timezone", resolved); ctl != nil {
		return nil, ctl
	}
	return clone, nil
}

func clockNow(ctx data.Context) (data.GetValue, data.Control) {
	fallback := func() (data.GetValue, data.Control) {
		return data.NewStringValue(time.Now().Format(time.RFC3339)), nil
	}
	vm := ctx.GetVM()
	if vm == nil {
		return fallback()
	}
	cls, ok := vm.GetClass("DateTimeImmutable")
	if !ok || cls == nil {
		return fallback()
	}
	gv, acl := cls.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}
	obj, ok := gv.(*data.ClassValue)
	if !ok || obj == nil {
		return fallback()
	}
	ctor := cls.GetConstruct()
	if ctor == nil {
		return obj, nil
	}
	// DateTime / DateTimeImmutable::__construct(string $datetime = 'now', ?DateTimeZone $timezone = null)
	vars := ctor.GetVariables()
	cctx := obj.CreateContext(vars)
	tz, ctl := clockTimezoneValue(ctx)
	if ctl != nil {
		return nil, ctl
	}
	for i, v := range vars {
		name := strings.ToLower(v.GetName())
		switch name {
		case "datetime", "time":
			if ctl := cctx.SetVariableValue(v, data.NewStringValue("now")); ctl != nil {
				return nil, ctl
			}
		case "timezone":
			if tz != nil {
				if ctl := cctx.SetVariableValue(v, tz); ctl != nil {
					return nil, ctl
				}
			}
		default:
			if i == 0 {
				if ctl := cctx.SetVariableValue(v, data.NewStringValue("now")); ctl != nil {
					return nil, ctl
				}
			} else if i == 1 && tz != nil {
				if ctl := cctx.SetVariableValue(v, tz); ctl != nil {
					return nil, ctl
				}
			}
		}
	}
	if _, ctl := ctor.Call(cctx); ctl != nil {
		return nil, ctl
	}
	return obj, nil
}

func clockSleep(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := ctx.GetIndexValue(0)
	if !ok || isNullValue(v) {
		return data.NewNullValue(), nil
	}
	seconds := 0.0
	if af, ok := v.(data.AsFloat); ok {
		seconds, _ = af.AsFloat()
	}
	// NativeClock::sleep：整数秒 sleep，小数部分 usleep；非正值不睡。(int) 向零截断。
	s := int64(seconds)
	if s > 0 {
		time.Sleep(time.Duration(s) * time.Second)
	}
	us := seconds - float64(s)
	if us > 0 {
		time.Sleep(time.Duration(us*1e6) * time.Microsecond)
	}
	return data.NewNullValue(), nil
}

func clockWithTimeZone(ctx data.Context) (data.GetValue, data.Control) {
	self := clockSelf(ctx)
	tz, ok := ctx.GetIndexValue(0)
	if !ok {
		tz = nil
	}
	clone, ctl := applyWithTimeZone(ctx, self, tz, false)
	if ctl != nil {
		return nil, ctl
	}
	return clone, nil
}
