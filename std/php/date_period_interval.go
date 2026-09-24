package php

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// DatePeriodClass 最小存根，满足 CarbonPeriod 对原生 DatePeriod 的继承需求。
type DatePeriodClass struct {
	node.Node
}

func NewDatePeriodClass() *DatePeriodClass { return &DatePeriodClass{} }

func (c *DatePeriodClass) GetName() string                               { return "DatePeriod" }
func (c *DatePeriodClass) GetExtend() *string                            { return nil }
func (c *DatePeriodClass) GetImplements() []string                       { return []string{"IteratorAggregate"} }
func (c *DatePeriodClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *DatePeriodClass) GetPropertyList() []data.Property              { return nil }
func (c *DatePeriodClass) GetConstruct() data.Method                     { return &DatePeriodConstructMethod{} }
func (c *DatePeriodClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *DatePeriodClass) GetMethods() []data.Method {
	return []data.Method{&DatePeriodConstructMethod{}, &DatePeriodGetIteratorMethod{}}
}
func (c *DatePeriodClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return &DatePeriodConstructMethod{}, true
	}
	if name == "getIterator" {
		return &DatePeriodGetIteratorMethod{}, true
	}
	return nil, false
}
func (c *DatePeriodClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

type DatePeriodConstructMethod struct{}

func (m *DatePeriodConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}
func (m *DatePeriodConstructMethod) GetName() string            { return "__construct" }
func (m *DatePeriodConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DatePeriodConstructMethod) GetIsStatic() bool          { return false }
func (m *DatePeriodConstructMethod) GetReturnType() data.Types  { return nil }
var datePeriodConstructMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "start", 0, nil, nil),
	node.NewParameter(nil, "interval", 1, nil, nil),
	node.NewParameter(nil, "end", 2, nil, nil),
	node.NewParameter(nil, "options", 3, data.NewIntValue(0), nil),
}

func (m *DatePeriodConstructMethod) GetParams() []data.GetValue {
	return datePeriodConstructMethodGetParams
}
var datePeriodConstructMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "start", 0, nil),
	node.NewVariable(nil, "interval", 1, nil),
	node.NewVariable(nil, "end", 2, nil),
	node.NewVariable(nil, "options", 3, nil),
}

func (m *DatePeriodConstructMethod) GetVariables() []data.Variable {
	return datePeriodConstructMethodGetVariables
}

// DatePeriodGetIteratorMethod 提供 DatePeriod 的 getIterator()，
// 满足 IteratorAggregate 接口要求。原生 DatePeriod::getIterator() 返回
// 迭代日期区间的迭代器，此处返回 $this（调用方按 Iterator 或对象属性遍历）。
type DatePeriodGetIteratorMethod struct{}

func (m *DatePeriodGetIteratorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		return data.NewThisValue(cmc.ClassValue), nil
	}
	return nil, nil
}
func (m *DatePeriodGetIteratorMethod) GetName() string            { return "getIterator" }
func (m *DatePeriodGetIteratorMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DatePeriodGetIteratorMethod) GetIsStatic() bool          { return false }
func (m *DatePeriodGetIteratorMethod) GetReturnType() data.Types  { return nil }
func (m *DatePeriodGetIteratorMethod) GetParams() []data.GetValue { return nil }
func (m *DatePeriodGetIteratorMethod) GetVariables() []data.Variable {
	return nil
}

// DateIntervalClass 最小存根。
type DateIntervalClass struct {
	node.Node
}

func NewDateIntervalClass() *DateIntervalClass { return &DateIntervalClass{} }

func (c *DateIntervalClass) GetName() string         { return "DateInterval" }
func (c *DateIntervalClass) GetExtend() *string      { return nil }
func (c *DateIntervalClass) GetImplements() []string { return nil }

// dateIntervalPropertyNames 原生 DateInterval 的公开属性。
// 这些属性在 DateInterval 上声明，因此子类（如 CarbonInterval）里
// 直接 $this->d = $value 应写入属性，而不是触发 __set 魔术方法。
// 注意：days 默认值为 false（未参与日差计算时），与原生 PHP 一致。
var dateIntervalPropertyNames = []string{"y", "m", "d", "h", "i", "s", "f", "invert", "days"}

func dateIntervalPropertyDefault(name string) data.GetValue {
	if name == "days" {
		return data.NewBoolValue(false)
	}
	return nil
}

func (c *DateIntervalClass) GetProperty(name string) (data.Property, bool) {
	for _, n := range dateIntervalPropertyNames {
		if n == name {
			return node.NewProperty(nil, n, "public", false, dateIntervalPropertyDefault(n)), true
		}
	}
	return nil, false
}
func (c *DateIntervalClass) GetPropertyList() []data.Property {
	props := make([]data.Property, 0, len(dateIntervalPropertyNames))
	for _, n := range dateIntervalPropertyNames {
		props = append(props, node.NewProperty(nil, n, "public", false, dateIntervalPropertyDefault(n)))
	}
	return props
}
func (c *DateIntervalClass) GetConstruct() data.Method { return &DateIntervalConstructMethod{} }
func (c *DateIntervalClass) GetStaticMethod(name string) (data.Method, bool) {
	if name == "createFromDateString" {
		return &DateIntervalCreateFromDateStringMethod{}, true
	}
	return nil, false
}
func (c *DateIntervalClass) GetMethods() []data.Method {
	return []data.Method{
		&DateIntervalConstructMethod{},
		&DateIntervalFormatMethod{},
		&DateIntervalCreateFromDateStringMethod{},
	}
}
func (c *DateIntervalClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &DateIntervalConstructMethod{}, true
	case "format":
		return &DateIntervalFormatMethod{}, true
	case "createFromDateString":
		return &DateIntervalCreateFromDateStringMethod{}, true
	}
	return nil, false
}
func (c *DateIntervalClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

type DateIntervalConstructMethod struct{}

// iso8601DurationRegexp 匹配 PHP DateInterval 的 ISO 8601 时长规范：
// P[nY][nM][nW][nD][T[nH][nM][nS]]
var iso8601DurationRegexp = regexp.MustCompile(`^P(?:(?P<year>\d+)Y)?(?:(?P<month>\d+)M)?(?:(?P<week>\d+)W)?(?:(?P<day>\d+)D)?(?:T(?:(?P<hour>\d+)H)?(?:(?P<minute>\d+)M)?(?:(?P<second>\d+(?:\.\d+)?)S)?)?$`)

// parseISO8601Duration 解析 ISO 8601 时长并返回各分量（单位数）。
func parseISO8601Duration(spec string) (year, month, week, day, hour, minute int, second float64, ok bool) {
	m := iso8601DurationRegexp.FindStringSubmatch(spec)
	if m == nil {
		return 0, 0, 0, 0, 0, 0, 0, false
	}
	atoi := func(name string) int {
		i := iso8601DurationRegexp.SubexpIndex(name)
		if i < 0 || i >= len(m) || m[i] == "" {
			return 0
		}
		v, _ := strconv.Atoi(m[i])
		return v
	}
	f := func(name string) float64 {
		i := iso8601DurationRegexp.SubexpIndex(name)
		if i < 0 || i >= len(m) || m[i] == "" {
			return 0
		}
		v, _ := strconv.ParseFloat(m[i], 64)
		return v
	}
	return atoi("year"), atoi("month"), atoi("week"), atoi("day"), atoi("hour"), atoi("minute"), f("second"), true
}

func (m *DateIntervalConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cmc, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, nil
	}
	duration, _ := ctx.GetIndexValue(0)
	if duration == nil {
		return nil, nil
	}
	spec := duration.AsString()
	if spec == "" {
		return nil, nil
	}
	year, month, week, day, hour, minute, second, parsed := parseISO8601Duration(spec)
	if !parsed {
		// 无效规格：PHP 抛 DateMalformedIntervalStringException。
		// 这里保持空间隔，避免 CarbonInterval 构造回退路径在 Origami 上把整页打成 500。
		return nil, nil
	}
	// 周折算为天（PHP：W → 天）。
	day += week * 7
	whole := int(second)
	frac := second - float64(whole)
	applyDateIntervalProps(cmc.ClassValue, year, month, day, hour, minute, whole, frac, 0, false)
	return nil, nil
}
func (m *DateIntervalConstructMethod) GetName() string            { return "__construct" }
func (m *DateIntervalConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateIntervalConstructMethod) GetIsStatic() bool          { return false }
func (m *DateIntervalConstructMethod) GetReturnType() data.Types  { return nil }
var dateIntervalConstructMethodGetParams = []data.GetValue{node.NewParameter(nil, "duration", 0, nil, data.String{})}

func (m *DateIntervalConstructMethod) GetParams() []data.GetValue {
	return dateIntervalConstructMethodGetParams
}
var dateIntervalConstructMethodGetVariables = []data.Variable{node.NewVariable(nil, "duration", 0, data.String{})}

func (m *DateIntervalConstructMethod) GetVariables() []data.Variable {
	return dateIntervalConstructMethodGetVariables
}

func applyDateIntervalProps(cv *data.ClassValue, y, m, d, h, i, s int, f float64, invert int, days any) {
	if cv == nil {
		return
	}
	cv.SetProperty("y", data.NewIntValue(y))
	cv.SetProperty("m", data.NewIntValue(m))
	cv.SetProperty("d", data.NewIntValue(d))
	cv.SetProperty("h", data.NewIntValue(h))
	cv.SetProperty("i", data.NewIntValue(i))
	cv.SetProperty("s", data.NewIntValue(s))
	cv.SetProperty("f", data.NewFloatValue(f))
	cv.SetProperty("invert", data.NewIntValue(invert))
	switch t := days.(type) {
	case int:
		cv.SetProperty("days", data.NewIntValue(t))
	case bool:
		cv.SetProperty("days", data.NewBoolValue(t))
	default:
		cv.SetProperty("days", data.NewBoolValue(false))
	}
}

// newDateIntervalValue 构造原生 DateInterval 实例（DateTime::diff / createFromDateString）。
func newDateIntervalValue(ctx data.Context, y, m, d, h, i, s int, f float64, invert int, days any) (data.GetValue, data.Control) {
	var stmt data.ClassStmt = NewDateIntervalClass()
	if ctx != nil && ctx.GetVM() != nil {
		loaded, acl := ctx.GetVM().GetOrLoadClass("DateInterval")
		if acl != nil {
			return nil, acl
		}
		if loaded != nil {
			stmt = loaded
		}
	}
	base := ctx
	if ctx != nil {
		base = ctx.CreateBaseContext()
	}
	cv := data.NewClassValue(stmt, base)
	applyDateIntervalProps(cv, y, m, d, h, i, s, f, invert, days)
	return cv, nil
}

type DateIntervalFormatMethod struct{}

func (m *DateIntervalFormatMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cmc, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return data.NewStringValue(""), nil
	}
	fmtV, _ := ctx.GetIndexValue(0)
	spec := ""
	if fmtV != nil {
		spec = fmtV.AsString()
	}
	return data.NewStringValue(formatDateInterval(cmc.ClassValue, spec)), nil
}
func (m *DateIntervalFormatMethod) GetName() string            { return "format" }
func (m *DateIntervalFormatMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateIntervalFormatMethod) GetIsStatic() bool          { return false }
func (m *DateIntervalFormatMethod) GetReturnType() data.Types  { return data.NewBaseType("string") }
var dateIntervalFormatMethodGetParams = []data.GetValue{node.NewParameter(nil, "format", 0, nil, data.String{})}

func (m *DateIntervalFormatMethod) GetParams() []data.GetValue {
	return dateIntervalFormatMethodGetParams
}
var dateIntervalFormatMethodGetVariables = []data.Variable{node.NewVariable(nil, "format", 0, data.String{})}

func (m *DateIntervalFormatMethod) GetVariables() []data.Variable {
	return dateIntervalFormatMethodGetVariables
}

func formatDateInterval(cv *data.ClassValue, spec string) string {
	if cv == nil {
		return ""
	}
	y := dateIntervalIntProp(cv, "y")
	mo := dateIntervalIntProp(cv, "m")
	d := dateIntervalIntProp(cv, "d")
	h := dateIntervalIntProp(cv, "h")
	i := dateIntervalIntProp(cv, "i")
	s := dateIntervalIntProp(cv, "s")
	f := dateIntervalFloatProp(cv, "f")
	invert := dateIntervalIntProp(cv, "invert")
	us := int(math.Round(f * 1e6))
	if us < 0 {
		us = -us
	}
	sign, signR := "+", ""
	if invert != 0 {
		sign, signR = "-", "-"
	}
	var b strings.Builder
	for idx := 0; idx < len(spec); idx++ {
		if spec[idx] == '%' && idx+1 < len(spec) {
			idx++
			switch spec[idx] {
			case '%':
				b.WriteByte('%')
			case 'Y':
				b.WriteString(fmt.Sprintf("%02d", y))
			case 'y':
				b.WriteString(strconv.Itoa(y))
			case 'M':
				b.WriteString(fmt.Sprintf("%02d", mo))
			case 'm':
				b.WriteString(strconv.Itoa(mo))
			case 'D':
				b.WriteString(fmt.Sprintf("%02d", d))
			case 'd':
				b.WriteString(strconv.Itoa(d))
			case 'H':
				b.WriteString(fmt.Sprintf("%02d", h))
			case 'h':
				b.WriteString(strconv.Itoa(h))
			case 'I':
				b.WriteString(fmt.Sprintf("%02d", i))
			case 'i':
				b.WriteString(strconv.Itoa(i))
			case 'S':
				b.WriteString(fmt.Sprintf("%02d", s))
			case 's':
				b.WriteString(strconv.Itoa(s))
			case 'F':
				b.WriteString(fmt.Sprintf("%06d", us))
			case 'f':
				b.WriteString(strconv.Itoa(us))
			case 'R':
				b.WriteString(sign)
			case 'r':
				b.WriteString(signR)
			case 'a':
				daysV, _ := cv.GetProperty("days")
				if daysV == nil {
					b.WriteString("(unknown)")
					break
				}
				if bv, ok := daysV.(*data.BoolValue); ok && !bv.Value {
					b.WriteString("(unknown)")
					break
				}
				b.WriteString(strconv.Itoa(dateIntervalIntProp(cv, "days")))
			default:
				b.WriteByte(spec[idx])
			}
			continue
		}
		b.WriteByte(spec[idx])
	}
	return b.String()
}

type DateIntervalCreateFromDateStringMethod struct{}

func (m *DateIntervalCreateFromDateStringMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	datetimeV, _ := ctx.GetIndexValue(0)
	if datetimeV == nil {
		return nil, utils.NewThrow(fmt.Errorf("DateInterval::createFromDateString(): Argument #1 ($datetime) must be of type string"))
	}
	raw := datetimeV.AsString()
	y, mo, d, h, i, s, f, invert, ok := parseDateStringInterval(raw)
	if !ok {
		return nil, utils.NewThrow(fmt.Errorf("DateInterval::createFromDateString(): Unknown or bad format (%s)", raw))
	}
	return newDateIntervalValue(ctx, y, mo, d, h, i, s, f, invert, false)
}
func (m *DateIntervalCreateFromDateStringMethod) GetName() string { return "createFromDateString" }
func (m *DateIntervalCreateFromDateStringMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *DateIntervalCreateFromDateStringMethod) GetIsStatic() bool         { return true }
func (m *DateIntervalCreateFromDateStringMethod) GetReturnType() data.Types { return nil }
var dateIntervalCreateFromDateStringMethodGetParams = []data.GetValue{node.NewParameter(nil, "datetime", 0, nil, data.String{})}

func (m *DateIntervalCreateFromDateStringMethod) GetParams() []data.GetValue {
	return dateIntervalCreateFromDateStringMethodGetParams
}
var dateIntervalCreateFromDateStringMethodGetVariables = []data.Variable{node.NewVariable(nil, "datetime", 0, data.String{})}

func (m *DateIntervalCreateFromDateStringMethod) GetVariables() []data.Variable {
	return dateIntervalCreateFromDateStringMethodGetVariables
}

// dateStringIntervalRE 解析 PHP DateInterval::createFromDateString 的相对时长片段。
var dateStringIntervalRE = regexp.MustCompile(`(?i)([+-]?\d+(?:\.\d+)?)\s*(microseconds?|milliseconds?|seconds?|minutes?|hours?|days?|weeks?|months?|years?|secs?|mins?|hrs?|ms|us|µs)`)

func parseDateStringInterval(s string) (y, m, d, h, i, sec int, f float64, invert int, ok bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, 0, 0, 0, 0, 0, 0, 0, true
	}
	s = strings.ReplaceAll(s, ",", " ")
	lower := strings.ToLower(s)
	lower = strings.ReplaceAll(lower, " and ", " ")
	ago := false
	if strings.HasSuffix(lower, " ago") {
		ago = true
		lower = strings.TrimSpace(lower[:len(lower)-4])
	}
	matches := dateStringIntervalRE.FindAllStringSubmatch(lower, -1)
	if len(matches) == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, false
	}
	ok = true
	for _, match := range matches {
		n, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			continue
		}
		unit := strings.ToLower(match[2])
		switch {
		case strings.HasPrefix(unit, "year"):
			y += int(n)
		case strings.HasPrefix(unit, "month"):
			m += int(n)
		case strings.HasPrefix(unit, "week"):
			d += int(n) * 7
		case strings.HasPrefix(unit, "day"):
			d += int(n)
		case strings.HasPrefix(unit, "hour") || unit == "hrs" || unit == "hr":
			h += int(n)
		case strings.HasPrefix(unit, "min"):
			i += int(n)
		case unit == "ms" || strings.HasPrefix(unit, "millisecond"):
			f += n / 1e3
		case unit == "us" || unit == "µs" || strings.HasPrefix(unit, "micro"):
			f += n / 1e6
		case strings.HasPrefix(unit, "sec"):
			sec += int(n)
			frac := n - float64(int(n))
			if frac != 0 {
				f += frac
			}
		}
	}
	if extra := int(f); extra != 0 {
		sec += extra
		f -= float64(extra)
	}
	if ago {
		invert = 1
	}
	return
}
