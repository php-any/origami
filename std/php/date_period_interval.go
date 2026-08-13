package php

import (
	"regexp"
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
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
func (m *DatePeriodConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "start", 0, nil, nil),
		node.NewParameter(nil, "interval", 1, nil, nil),
		node.NewParameter(nil, "end", 2, nil, nil),
		node.NewParameter(nil, "options", 3, data.NewIntValue(0), nil),
	}
}
func (m *DatePeriodConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "start", 0, nil),
		node.NewVariable(nil, "interval", 1, nil),
		node.NewVariable(nil, "end", 2, nil),
		node.NewVariable(nil, "options", 3, nil),
	}
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

func (c *DateIntervalClass) GetName() string    { return "DateInterval" }
func (c *DateIntervalClass) GetExtend() *string { return nil }
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
func (c *DateIntervalClass) GetConstruct() data.Method                     { return &DateIntervalConstructMethod{} }
func (c *DateIntervalClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *DateIntervalClass) GetMethods() []data.Method {
	return []data.Method{&DateIntervalConstructMethod{}}
}
func (c *DateIntervalClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return &DateIntervalConstructMethod{}, true
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
		// 无效规格：PHP 抛异常，这里按空间隔处理（保持 DateInterval 最小存根语义）
		return nil, nil
	}
	// 周折算为天（PHP：W → 天）。
	day += week * 7
	cmc.ObjectValue.SetProperty("y", data.NewIntValue(year))
	cmc.ObjectValue.SetProperty("m", data.NewIntValue(month))
	cmc.ObjectValue.SetProperty("d", data.NewIntValue(day))
	cmc.ObjectValue.SetProperty("h", data.NewIntValue(hour))
	cmc.ObjectValue.SetProperty("i", data.NewIntValue(minute))
	if second != 0 {
		// 秒含小数时写入 f
		whole := int(second)
		frac := second - float64(whole)
		cmc.ObjectValue.SetProperty("s", data.NewIntValue(whole))
		if frac != 0 {
			cmc.ObjectValue.SetProperty("f", data.NewFloatValue(frac))
		}
	} else {
		cmc.ObjectValue.SetProperty("s", data.NewIntValue(0))
	}
	return nil, nil
}
func (m *DateIntervalConstructMethod) GetName() string            { return "__construct" }
func (m *DateIntervalConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateIntervalConstructMethod) GetIsStatic() bool          { return false }
func (m *DateIntervalConstructMethod) GetReturnType() data.Types  { return nil }
func (m *DateIntervalConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "duration", 0, nil, data.String{})}
}
func (m *DateIntervalConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "duration", 0, data.String{})}
}
