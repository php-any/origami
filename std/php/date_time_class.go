package php

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// DateTimeClass 实现 PHP DateTime 类（简化版）
type DateTimeClass struct {
	node.Node
}

func NewDateTimeClass() *DateTimeClass {
	return &DateTimeClass{}
}

func (c *DateTimeClass) GetName() string                               { return "DateTime" }
func (c *DateTimeClass) GetExtend() *string                            { return nil }
func (c *DateTimeClass) GetImplements() []string                       { return []string{"DateTimeInterface"} }
func (c *DateTimeClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *DateTimeClass) GetPropertyList() []data.Property              { return nil }
func (c *DateTimeClass) GetConstruct() data.Method                     { return &DateTimeConstructMethod{} }
func (c *DateTimeClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "getLastErrors":
		return &DateTimeGetLastErrorsMethod{}, true
	case "createFromFormat":
		return &DateTimeCreateFromFormatMethod{}, true
	}
	return nil, false
}
func (c *DateTimeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *DateTimeClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &DateTimeConstructMethod{}, true
	case "getTimestamp":
		return &DateTimeGetTimestampMethod{}, true
	case "setTimestamp":
		return &DateTimeSetTimestampMethod{}, true
	case "setDate":
		return &DateTimeSetDateMethod{}, true
	case "setTime":
		return &DateTimeSetTimeMethod{}, true
	case "setISODate":
		return &DateTimeSetISODateMethod{}, true
	case "setTimezone":
		return &DateTimeSetTimezoneMethod{}, true
	case "getTimezone":
		return &DateTimeGetTimezoneMethod{}, true
	case "format":
		return &DateTimeFormatMethod{}, true
	case "add":
		return &DateTimeAddMethod{}, true
	case "sub":
		return &DateTimeSubMethod{}, true
	case "diff":
		return &DateTimeDiffMethod{}, true
	case "modify":
		return &DateTimeModifyMethod{}, true
	case "createFromFormat":
		return &DateTimeCreateFromFormatMethod{}, true
	case "__toString":
		return &DateTimeToStringMethod{}, true
	}
	return nil, false
}

func (c *DateTimeClass) GetMethods() []data.Method {
	return []data.Method{
		&DateTimeConstructMethod{}, &DateTimeGetTimestampMethod{}, &DateTimeSetTimestampMethod{},
		&DateTimeSetDateMethod{}, &DateTimeSetTimeMethod{}, &DateTimeSetISODateMethod{},
		&DateTimeSetTimezoneMethod{}, &DateTimeGetTimezoneMethod{}, &DateTimeFormatMethod{},
		&DateTimeAddMethod{}, &DateTimeSubMethod{}, &DateTimeDiffMethod{},
		&DateTimeModifyMethod{}, &DateTimeCreateFromFormatMethod{}, &DateTimeToStringMethod{},
	}
}

// timeFromTimestampValue 把 DateTime 内部 timestamp 属性转成 UTC 时间。
func timeFromTimestampValue(t data.Value) (time.Time, bool) {
	if t == nil {
		return time.Time{}, false
	}
	switch v := t.(type) {
	case *data.IntValue:
		return time.Unix(int64(v.Value), 0).UTC(), true
	case *data.FloatValue:
		sec, frac := math.Modf(v.Value)
		return time.Unix(int64(sec), int64(frac*1e9)).UTC(), true
	}
	if ai, ok := t.(data.AsInt); ok {
		n, err := ai.AsInt()
		if err == nil {
			return time.Unix(int64(n), 0).UTC(), true
		}
	}
	return time.Time{}, false
}

// dateTimeFromValue 从 DateTime / DateTimeImmutable / Carbon 对象取出时间。
func dateTimeFromValue(v data.Value) (time.Time, bool) {
	cv, ok := toClassValue(v)
	if !ok || cv == nil {
		return time.Time{}, false
	}
	if ts, ctl := cv.GetProperty("timestamp"); ctl == nil && ts != nil {
		if t, ok := timeFromTimestampValue(ts); ok {
			return t, true
		}
	}
	if m, ok := cv.GetMethod("getTimestamp"); ok {
		callCtx := cv.CreateContext(m.GetVariables())
		result, ctl := m.Call(callCtx)
		if ctl == nil && result != nil {
			if val, ok := result.(data.Value); ok {
				if t, ok := timeFromTimestampValue(val); ok {
					return t, true
				}
			}
		}
	}
	return time.Time{}, false
}

// 从 context 获取 ClassValue 并提取时间戳
func getDateTime(ctx data.Context) (time.Time, data.Control) {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		t, ctl := cmc.ObjectValue.GetProperty("timestamp")
		if ctl == nil && t != nil {
			if ts, ok := timeFromTimestampValue(t); ok {
				return ts, nil
			}
		}
		if dt, ok := dateTimeFromValue(cmc.ClassValue); ok {
			return dt, nil
		}
	}
	return time.Unix(0, 0).UTC(), nil
}

func setDateTime(ctx data.Context, t time.Time) data.Control {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		cmc.ObjectValue.SetProperty("timestamp", data.NewIntValue(int(t.Unix())))
		return nil
	}
	return data.NewErrorThrow(nil, errors.New("无法在非对象上下文中设置时间戳"))
}

// toClassValue 把 ClassValue 或其包装（ThisValue）解包为 *data.ClassValue。
func toClassValue(v data.Value) (*data.ClassValue, bool) {
	if cv, ok := v.(*data.ClassValue); ok {
		return cv, true
	}
	if tv, ok := v.(*data.ThisValue); ok && tv.ClassValue != nil {
		return tv.ClassValue, true
	}
	return nil, false
}

// dateIntervalIntProp 从 DateInterval 对象读取整型分量（y/m/d/h/i/s/invert 等）。
func dateIntervalIntProp(interval *data.ClassValue, name string) int {
	v, ctl := interval.GetProperty(name)
	if ctl != nil || v == nil {
		return 0
	}
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	if ai, ok := v.(data.AsInt); ok {
		n, err := ai.AsInt()
		if err != nil {
			return 0
		}
		return n
	}
	return 0
}

// dateIntervalFloatProp 读取 DateInterval::$f（秒的小数部分）。
func dateIntervalFloatProp(interval *data.ClassValue, name string) float64 {
	v, ctl := interval.GetProperty(name)
	if ctl != nil || v == nil {
		return 0
	}
	if fv, ok := v.(*data.FloatValue); ok {
		return fv.Value
	}
	if ff, ok := v.(data.AsFloat); ok {
		n, err := ff.AsFloat()
		if err != nil {
			return 0
		}
		return n
	}
	return float64(dateIntervalIntProp(interval, name))
}

// applyDateInterval 把 DateInterval 分量应用到时间 t 上。negate=true 表示减法。
func applyDateInterval(t time.Time, interval *data.ClassValue, negate bool) time.Time {
	y := dateIntervalIntProp(interval, "y")
	m := dateIntervalIntProp(interval, "m")
	d := dateIntervalIntProp(interval, "d")
	h := dateIntervalIntProp(interval, "h")
	i := dateIntervalIntProp(interval, "i")
	s := dateIntervalIntProp(interval, "s")
	invert := dateIntervalIntProp(interval, "invert")
	sign := 1
	if invert != 0 {
		sign = -1
	}
	if negate {
		sign = -sign
	}
	t = t.AddDate(sign*y, sign*m, sign*d)
	t = t.Add(time.Duration(sign) * (time.Duration(h)*time.Hour + time.Duration(i)*time.Minute + time.Duration(s)*time.Second))
	if f := dateIntervalFloatProp(interval, "f"); f != 0 {
		t = t.Add(time.Duration(float64(sign) * f * float64(time.Second)))
	}
	return t
}

// calendarDateDiff 对齐 PHP DateTime::diff 的日历进位（y/m/d/h/i/s/f/invert/days）。
func calendarDateDiff(from, to time.Time, absolute bool) (y, m, d, h, i, s int, f float64, invert int, days int) {
	if to.Before(from) {
		if !absolute {
			invert = 1
		}
		from, to = to, from
	}
	if dur := to.Sub(from); dur > 0 {
		days = int(dur / (24 * time.Hour))
	}
	y = to.Year() - from.Year()
	m = int(to.Month()) - int(from.Month())
	d = to.Day() - from.Day()
	h = to.Hour() - from.Hour()
	i = to.Minute() - from.Minute()
	s = to.Second() - from.Second()
	nsec := to.Nanosecond() - from.Nanosecond()
	if nsec < 0 {
		s--
		nsec += 1e9
	}
	f = float64(nsec) / 1e9
	if s < 0 {
		i--
		s += 60
	}
	if i < 0 {
		h--
		i += 60
	}
	if h < 0 {
		d--
		h += 24
	}
	if d < 0 {
		prev := time.Date(to.Year(), to.Month(), 0, 0, 0, 0, 0, to.Location())
		d += prev.Day()
		m--
	}
	if m < 0 {
		y--
		m += 12
	}
	return
}

// ---- __construct ----
type DateTimeConstructMethod struct{}

func (m *DateTimeConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	datetime, _ := ctx.GetIndexValue(0)
	if datetime == nil {
		setDateTime(ctx, time.Now().UTC())
		return data.NewNullValue(), nil
	}
	str := datetime.AsString()
	if str == "" || str == "now" || strings.HasPrefix(str, "now ") {
		// 支持 "now", "now - 3600 seconds" 等
		t, err := parseRelativeTime(str, time.Now().UTC())
		if err != nil {
			setDateTime(ctx, time.Now().UTC())
		} else {
			setDateTime(ctx, t)
		}
		return data.NewNullValue(), nil
	}
	// 尝试用 strtotime 解析
	t, err := parseRelativeTime(str, time.Now().UTC())
	if err != nil {
		// 如果相对时间解析失败，尝试常见的日期格式
		formats := []string{
			"2006-01-02 15:04:05",
			"2006-01-02",
			"2006-01-02 15:04",
			"2006/01/02 15:04:05",
			"2006/01/02",
		}
		for _, f := range formats {
			if t, err := time.Parse(f, str); err == nil {
				setDateTime(ctx, t.UTC())
				return data.NewNullValue(), nil
			}
		}
		// 整数：当作 Unix 时间戳
		if ts, err := strconv.ParseInt(str, 10, 64); err == nil {
			setDateTime(ctx, time.Unix(ts, 0).UTC())
			return data.NewNullValue(), nil
		}
		setDateTime(ctx, time.Now().UTC())
	} else {
		setDateTime(ctx, t)
	}
	return data.NewNullValue(), nil
}

func (m *DateTimeConstructMethod) GetName() string            { return "__construct" }
func (m *DateTimeConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeConstructMethod) GetIsStatic() bool          { return false }
func (m *DateTimeConstructMethod) GetReturnType() data.Types  { return nil }
func (m *DateTimeConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "datetime", 0, node.NewStringLiteralByAst(nil, "now"), nil),
		node.NewParameter(nil, "timezone", 1, node.NewNullLiteral(nil), nil),
	}
}
func (m *DateTimeConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "datetime", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "timezone", 1, data.NewNullableType(data.NewBaseType("DateTimeZone"))),
	}
}

// ---- getTimestamp ----
type DateTimeGetTimestampMethod struct{}

func (m *DateTimeGetTimestampMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	t, ctl := getDateTime(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewIntValue(int(t.Unix())), nil
}
func (m *DateTimeGetTimestampMethod) GetName() string               { return "getTimestamp" }
func (m *DateTimeGetTimestampMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DateTimeGetTimestampMethod) GetIsStatic() bool             { return false }
func (m *DateTimeGetTimestampMethod) GetReturnType() data.Types     { return data.NewBaseType("int") }
func (m *DateTimeGetTimestampMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DateTimeGetTimestampMethod) GetVariables() []data.Variable { return []data.Variable{} }

// ---- setTimestamp ----
type DateTimeSetTimestampMethod struct{}

func (m *DateTimeSetTimestampMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	tsV, _ := ctx.GetIndexValue(0)
	if tsV == nil {
		return nil, utils.NewThrow(errors.New("缺少时间戳参数"))
	}
	ts := int64(0)
	if iv, ok := tsV.(*data.IntValue); ok {
		ts = int64(iv.Value)
	} else {
		s := tsV.AsString()
		parsed, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, utils.NewThrow(errors.New("无效的时间戳"))
		}
		ts = parsed
	}
	setDateTime(ctx, time.Unix(ts, 0).UTC())
	cmc := ctx.(*data.ClassMethodContext)
	return cmc.ClassValue, nil
}
func (m *DateTimeSetTimestampMethod) GetName() string            { return "setTimestamp" }
func (m *DateTimeSetTimestampMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeSetTimestampMethod) GetIsStatic() bool          { return false }
func (m *DateTimeSetTimestampMethod) GetReturnType() data.Types  { return nil }
func (m *DateTimeSetTimestampMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "timestamp", 0, nil, nil)}
}
func (m *DateTimeSetTimestampMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "timestamp", 0, data.NewBaseType("int"))}
}

// ---- setDate ----
type DateTimeSetDateMethod struct{}

func (m *DateTimeSetDateMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	yv, _ := ctx.GetIndexValue(0)
	mv, _ := ctx.GetIndexValue(1)
	dv, _ := ctx.GetIndexValue(2)
	if yv == nil || mv == nil || dv == nil {
		return nil, utils.NewThrow(errors.New("缺少参数"))
	}
	y, _ := strconv.Atoi(yv.AsString())
	mo, _ := strconv.Atoi(mv.AsString())
	d, _ := strconv.Atoi(dv.AsString())
	t, _ := getDateTime(ctx)
	t = time.Date(y, time.Month(mo), d, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	setDateTime(ctx, t)
	cmc := ctx.(*data.ClassMethodContext)
	return cmc.ClassValue, nil
}
func (m *DateTimeSetDateMethod) GetName() string            { return "setDate" }
func (m *DateTimeSetDateMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeSetDateMethod) GetIsStatic() bool          { return false }
func (m *DateTimeSetDateMethod) GetReturnType() data.Types  { return nil }
func (m *DateTimeSetDateMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "year", 0, nil, nil),
		node.NewParameter(nil, "month", 1, nil, nil),
		node.NewParameter(nil, "day", 2, nil, nil),
	}
}
func (m *DateTimeSetDateMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "year", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "month", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "day", 2, data.NewBaseType("int")),
	}
}

// ---- setTime ----
type DateTimeSetTimeMethod struct{}

func (m *DateTimeSetTimeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	hv, _ := ctx.GetIndexValue(0)
	miv, _ := ctx.GetIndexValue(1)
	sv, _ := ctx.GetIndexValue(2)
	usv, _ := ctx.GetIndexValue(3)
	if hv == nil || miv == nil || sv == nil {
		return nil, utils.NewThrow(errors.New("缺少参数"))
	}
	h, _ := strconv.Atoi(hv.AsString())
	mi, _ := strconv.Atoi(miv.AsString())
	s, _ := strconv.Atoi(sv.AsString())
	us := 0
	if usv != nil {
		us, _ = strconv.Atoi(usv.AsString())
	}
	t, _ := getDateTime(ctx)
	t = time.Date(t.Year(), t.Month(), t.Day(), h, mi, s, us*1000, t.Location())
	setDateTime(ctx, t)
	cmc := ctx.(*data.ClassMethodContext)
	return cmc.ClassValue, nil
}
func (m *DateTimeSetTimeMethod) GetName() string            { return "setTime" }
func (m *DateTimeSetTimeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeSetTimeMethod) GetIsStatic() bool          { return false }
func (m *DateTimeSetTimeMethod) GetReturnType() data.Types  { return nil }
func (m *DateTimeSetTimeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "hour", 0, nil, nil),
		node.NewParameter(nil, "min", 1, nil, nil),
		node.NewParameter(nil, "second", 2, node.NewIntLiteral(nil, "0"), nil),
		node.NewParameter(nil, "microsecond", 3, node.NewIntLiteral(nil, "0"), nil),
	}
}
func (m *DateTimeSetTimeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "hour", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "min", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "second", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "microsecond", 3, data.NewBaseType("int")),
	}
}

// ---- setISODate ----
type DateTimeSetISODateMethod struct{}

func (m *DateTimeSetISODateMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	yv, _ := ctx.GetIndexValue(0)
	wv, _ := ctx.GetIndexValue(1)
	dv, _ := ctx.GetIndexValue(2)
	if yv == nil || wv == nil {
		return nil, utils.NewThrow(errors.New("缺少参数"))
	}
	y, _ := strconv.Atoi(yv.AsString())
	w, _ := strconv.Atoi(wv.AsString())
	d := 1
	if dv != nil {
		d, _ = strconv.Atoi(dv.AsString())
	}
	// 使用 ISO week 计算: Jan 4th is always in week 1
	jan4 := time.Date(y, 1, 4, 0, 0, 0, 0, time.UTC)
	_, isoWeek := jan4.ISOWeek()
	// 计算第一天
	startOfYear := jan4.AddDate(0, 0, -(int(jan4.Weekday()+6)%7)-(isoWeek-1)*7)
	target := startOfYear.AddDate(0, 0, (w-1)*7+(d-1))
	setDateTime(ctx, target)
	cmc := ctx.(*data.ClassMethodContext)
	return cmc.ClassValue, nil
}
func (m *DateTimeSetISODateMethod) GetName() string            { return "setISODate" }
func (m *DateTimeSetISODateMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeSetISODateMethod) GetIsStatic() bool          { return false }
func (m *DateTimeSetISODateMethod) GetReturnType() data.Types  { return nil }
func (m *DateTimeSetISODateMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "year", 0, nil, nil),
		node.NewParameter(nil, "week", 1, nil, nil),
		node.NewParameter(nil, "dayOfWeek", 2, node.NewIntLiteral(nil, "1"), nil),
	}
}
func (m *DateTimeSetISODateMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "year", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "week", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "dayOfWeek", 2, data.NewBaseType("int")),
	}
}

// ---- format ----
type DateTimeFormatMethod struct{}

func (m *DateTimeFormatMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	fv, _ := ctx.GetIndexValue(0)
	format := "Y-m-d H:i:s"
	if fv != nil {
		format = fv.AsString()
	}
	t, ctl := getDateTime(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewStringValue(phpDateFormat(format, t)), nil
}
func (m *DateTimeFormatMethod) GetName() string            { return "format" }
func (m *DateTimeFormatMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeFormatMethod) GetIsStatic() bool          { return false }
func (m *DateTimeFormatMethod) GetReturnType() data.Types  { return data.NewBaseType("string") }
func (m *DateTimeFormatMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "format", 0, nil, nil)}
}
func (m *DateTimeFormatMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "format", 0, data.NewBaseType("string"))}
}

// ---- getTimezone / setTimezone / add / sub / diff / modify / createFromFormat / __toString ----
type DateTimeGetTimezoneMethod struct{}

func (m *DateTimeGetTimezoneMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 简化：始终返回 UTC
	return data.NewStringValue("UTC"), nil
}
func (m *DateTimeGetTimezoneMethod) GetName() string               { return "getTimezone" }
func (m *DateTimeGetTimezoneMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DateTimeGetTimezoneMethod) GetIsStatic() bool             { return false }
func (m *DateTimeGetTimezoneMethod) GetReturnType() data.Types     { return nil }
func (m *DateTimeGetTimezoneMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DateTimeGetTimezoneMethod) GetVariables() []data.Variable { return []data.Variable{} }

type DateTimeSetTimezoneMethod struct{}

func (m *DateTimeSetTimezoneMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 简化：忽略时区
	cmc := ctx.(*data.ClassMethodContext)
	return cmc.ClassValue, nil
}
func (m *DateTimeSetTimezoneMethod) GetName() string            { return "setTimezone" }
func (m *DateTimeSetTimezoneMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeSetTimezoneMethod) GetIsStatic() bool          { return false }
func (m *DateTimeSetTimezoneMethod) GetReturnType() data.Types  { return nil }
func (m *DateTimeSetTimezoneMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "timezone", 0, nil, nil)}
}
func (m *DateTimeSetTimezoneMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "timezone", 0, nil)}
}

type DateTimeAddMethod struct{}

func (m *DateTimeAddMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cmc := ctx.(*data.ClassMethodContext)
	if interval, ok := ctx.GetIndexValue(0); ok {
		if iv, isCV := toClassValue(interval); isCV {
			t, _ := getDateTime(ctx)
			setDateTime(ctx, applyDateInterval(t, iv, false))
		}
	}
	return cmc.ClassValue, nil
}
func (m *DateTimeAddMethod) GetName() string            { return "add" }
func (m *DateTimeAddMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeAddMethod) GetIsStatic() bool          { return false }
func (m *DateTimeAddMethod) GetReturnType() data.Types  { return nil }
func (m *DateTimeAddMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "interval", 0, nil, nil)}
}
func (m *DateTimeAddMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "interval", 0, nil)}
}

type DateTimeSubMethod struct{}

func (m *DateTimeSubMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cmc := ctx.(*data.ClassMethodContext)
	if interval, ok := ctx.GetIndexValue(0); ok {
		if iv, isCV := toClassValue(interval); isCV {
			t, _ := getDateTime(ctx)
			setDateTime(ctx, applyDateInterval(t, iv, true))
		}
	}
	return cmc.ClassValue, nil
}
func (m *DateTimeSubMethod) GetName() string            { return "sub" }
func (m *DateTimeSubMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeSubMethod) GetIsStatic() bool          { return false }
func (m *DateTimeSubMethod) GetReturnType() data.Types  { return nil }
func (m *DateTimeSubMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "interval", 0, nil, nil)}
}
func (m *DateTimeSubMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "interval", 0, nil)}
}

type DateTimeDiffMethod struct{}

func (m *DateTimeDiffMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	from, ctl := getDateTime(ctx)
	if ctl != nil {
		return nil, ctl
	}
	targetV, _ := ctx.GetIndexValue(0)
	to, ok := dateTimeFromValue(targetV)
	if !ok {
		return nil, utils.NewThrow(fmt.Errorf("DateTime::diff(): Argument #1 ($targetObject) must be of type DateTimeInterface"))
	}
	absolute := false
	if absV, ok := ctx.GetIndexValue(1); ok && absV != nil {
		if b, ok := absV.(data.AsBool); ok {
			if bv, err := b.AsBool(); err == nil {
				absolute = bv
			}
		}
	}
	y, mo, d, h, i, s, f, invert, _ := calendarDateDiff(from, to, absolute)
	// days 保持 false：CarbonInterval::instance 在 PHP 8.2+ 会对 days!==false 走 serialize 复制，
	// Origami 尚不能正确 unserialize DateInterval，会把登录页打成 500。
	return newDateIntervalValue(ctx, y, mo, d, h, i, s, f, invert, false)
}

func (m *DateTimeDiffMethod) GetName() string            { return "diff" }
func (m *DateTimeDiffMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeDiffMethod) GetIsStatic() bool          { return false }
func (m *DateTimeDiffMethod) GetReturnType() data.Types  { return nil }
func (m *DateTimeDiffMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "target", 0, nil, nil),
		node.NewParameter(nil, "absolute", 1, node.NewBooleanLiteral(nil, false), nil),
	}
}
func (m *DateTimeDiffMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "target", 0, nil),
		node.NewVariable(nil, "absolute", 1, data.NewBaseType("bool")),
	}
}

type DateTimeModifyMethod struct{}

func (m *DateTimeModifyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	modV, _ := ctx.GetIndexValue(0)
	if modV == nil {
		return nil, utils.NewThrow(errors.New("缺少 modifier 参数"))
	}
	modStr := modV.AsString()
	t, _ := getDateTime(ctx)
	newT, err := parseRelativeTime(modStr, t)
	if err != nil {
		return nil, utils.NewThrow(fmt.Errorf("无效的修改器: %s", modStr))
	}
	setDateTime(ctx, newT)
	cmc := ctx.(*data.ClassMethodContext)
	return cmc.ClassValue, nil
}
func (m *DateTimeModifyMethod) GetName() string            { return "modify" }
func (m *DateTimeModifyMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeModifyMethod) GetIsStatic() bool          { return false }
func (m *DateTimeModifyMethod) GetReturnType() data.Types  { return nil }
func (m *DateTimeModifyMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "modifier", 0, nil, nil)}
}
func (m *DateTimeModifyMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "modifier", 0, data.NewBaseType("string"))}
}

type DateTimeCreateFromFormatMethod struct{}

func (m *DateTimeCreateFromFormatMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	formatV, _ := ctx.GetIndexValue(0)
	datetimeV, _ := ctx.GetIndexValue(1)
	if formatV == nil || datetimeV == nil {
		return data.NewBoolValue(false), nil
	}
	format := formatV.AsString()
	datetimeStr := datetimeV.AsString()
	t, err := parsePHPCreateFromFormat(format, datetimeStr)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	stmt := data.ClassStmt(NewDateTimeClass())
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		if cmc.StaticClass != nil {
			stmt = cmc.StaticClass
		} else if cmc.Class != nil {
			stmt = cmc.Class
		}
	}
	cv := data.NewClassValue(stmt, ctx.CreateBaseContext())
	cv.ObjectValue.SetProperty("timestamp", data.NewIntValue(int(t.Unix())))
	return cv, nil
}
func (m *DateTimeCreateFromFormatMethod) GetName() string            { return "createFromFormat" }
func (m *DateTimeCreateFromFormatMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeCreateFromFormatMethod) GetIsStatic() bool          { return true }
func (m *DateTimeCreateFromFormatMethod) GetReturnType() data.Types  { return nil }
func (m *DateTimeCreateFromFormatMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "format", 0, nil, nil),
		node.NewParameter(nil, "datetime", 1, nil, nil),
		node.NewParameter(nil, "timezone", 2, node.NewNullLiteral(nil), nil),
	}
}
func (m *DateTimeCreateFromFormatMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "datetime", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "timezone", 2, data.NewNullableType(data.NewBaseType("DateTimeZone"))),
	}
}

type DateTimeToStringMethod struct{}

func (m *DateTimeToStringMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	t, ctl := getDateTime(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewStringValue(t.Format("2006-01-02 15:04:05")), nil
}
func (m *DateTimeToStringMethod) GetName() string               { return "__toString" }
func (m *DateTimeToStringMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DateTimeToStringMethod) GetIsStatic() bool             { return false }
func (m *DateTimeToStringMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
func (m *DateTimeToStringMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DateTimeToStringMethod) GetVariables() []data.Variable { return []data.Variable{} }

// ---- getLastErrors ----
type DateTimeGetLastErrorsMethod struct{}

func (m *DateTimeGetLastErrorsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// PHP：array{warning_count, warnings, error_count, errors}；Carbon 会读 $lastErrors['errors']。
	return &data.ArrayValue{List: []*data.ZVal{
		data.NewNamedZVal("warning_count", data.NewIntValue(0)),
		data.NewNamedZVal("warnings", data.NewArrayValue(nil)),
		data.NewNamedZVal("error_count", data.NewIntValue(0)),
		data.NewNamedZVal("errors", data.NewArrayValue(nil)),
	}}, nil
}
func (m *DateTimeGetLastErrorsMethod) GetName() string            { return "getLastErrors" }
func (m *DateTimeGetLastErrorsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateTimeGetLastErrorsMethod) GetIsStatic() bool          { return true }
func (m *DateTimeGetLastErrorsMethod) GetReturnType() data.Types {
	return data.NewNullableType(data.NewBaseType("array"))
}
func (m *DateTimeGetLastErrorsMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *DateTimeGetLastErrorsMethod) GetVariables() []data.Variable { return []data.Variable{} }

// convertDateFormat 简单的时间格式转换
func parsePHPCreateFromFormat(format, str string) (time.Time, error) {
	if t, err := time.Parse(format, str); err == nil {
		return t, nil
	}
	if format == "U" || format == "U.u" {
		parts := strings.SplitN(str, ".", 2)
		sec, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		nsec := int64(0)
		if len(parts) == 2 && parts[1] != "" {
			frac := parts[1]
			for len(frac) < 6 {
				frac += "0"
			}
			if len(frac) > 6 {
				frac = frac[:6]
			}
			us, err := strconv.ParseInt(frac, 10, 64)
			if err == nil {
				nsec = us * 1000
			}
		}
		return time.Unix(sec, nsec).UTC(), nil
	}
	return convertDateFormat(format, str)
}

func convertDateFormat(format, str string) (time.Time, error) {
	goFormat := ""
	for i := 0; i < len(format); i++ {
		c := format[i]
		switch c {
		case 'Y':
			goFormat += "2006"
		case 'm':
			goFormat += "01"
		case 'd':
			goFormat += "02"
		case 'H':
			goFormat += "15"
		case 'i':
			goFormat += "04"
		case 's':
			goFormat += "05"
		case 'u':
			goFormat += "000000"
		case 'v':
			goFormat += "000"
		case 'P':
			goFormat += "Z07:00"
		case '!':
			// 重置字段
		case '\\':
			if i+1 < len(format) {
				i++
				goFormat += string(format[i])
			}
		default:
			goFormat += string(c)
		}
	}
	return time.Parse(goFormat, str)
}
