package php

import (
	"strconv"
	"strings"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DateFunction 实现 date 函数
// date(string $format, ?int $timestamp = null): string
type DateFunction struct{}

func NewDateFunction() data.FuncStmt { return &DateFunction{} }

func (f *DateFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	formatV, _ := ctx.GetIndexValue(0)
	tsV, _ := ctx.GetIndexValue(1)

	if formatV == nil {
		return data.NewStringValue(""), nil
	}
	format := formatV.AsString()

	var t time.Time
	// 检查 timestamp 参数是否为 null（默认参数值为 null 时表示使用当前时间）
	if tsV != nil {
		if _, isNull := tsV.(*data.NullValue); isNull {
			t = time.Now()
		} else if iv, ok := tsV.(*data.IntValue); ok {
			t = time.Unix(int64(iv.Value), 0)
		} else if fv, ok := tsV.(*data.FloatValue); ok {
			t = time.Unix(int64(fv.Value), 0)
		} else if ai, ok := tsV.(data.AsInt); ok {
			iv, err := ai.AsInt()
			if err == nil {
				t = time.Unix(int64(iv), 0)
			} else {
				t = time.Now()
			}
		} else {
			t = time.Now()
		}
	} else {
		t = time.Now()
	}

	return data.NewStringValue(formatDate(t, format)), nil
}

// formatDate 实现 PHP date() 格式化的核心逻辑
func formatDate(t time.Time, format string) string {
	var result strings.Builder
	for i := 0; i < len(format); i++ {
		c := format[i]
		switch c {
		case 'Y':
			result.WriteString(strconv.Itoa(t.Year()))
		case 'y':
			result.WriteString(strconv.Itoa(t.Year() % 100))
		case 'm':
			result.WriteString(twoDigit(int(t.Month())))
		case 'n':
			result.WriteString(strconv.Itoa(int(t.Month())))
		case 'M':
			result.WriteString(t.Month().String()[:3])
		case 'F':
			result.WriteString(t.Month().String())
		case 'd':
			result.WriteString(twoDigit(t.Day()))
		case 'j':
			result.WriteString(strconv.Itoa(t.Day()))
		case 'D':
			result.WriteString(t.Weekday().String()[:3])
		case 'l':
			result.WriteString(t.Weekday().String())
		case 'N':
			wd := int(t.Weekday())
			if wd == 0 {
				wd = 7
			}
			result.WriteString(strconv.Itoa(wd))
		case 'w':
			result.WriteString(strconv.Itoa(int(t.Weekday())))
		case 'H':
			result.WriteString(twoDigit(t.Hour()))
		case 'G':
			result.WriteString(strconv.Itoa(t.Hour()))
		case 'h':
			h := t.Hour() % 12
			if h == 0 {
				h = 12
			}
			result.WriteString(twoDigit(h))
		case 'g':
			h := t.Hour() % 12
			if h == 0 {
				h = 12
			}
			result.WriteString(strconv.Itoa(h))
		case 'i':
			result.WriteString(twoDigit(t.Minute()))
		case 's':
			result.WriteString(twoDigit(t.Second()))
		case 'A':
			if t.Hour() < 12 {
				result.WriteString("AM")
			} else {
				result.WriteString("PM")
			}
		case 'a':
			if t.Hour() < 12 {
				result.WriteString("am")
			} else {
				result.WriteString("pm")
			}
		case 'U':
			result.WriteString(strconv.FormatInt(t.Unix(), 10))
		case 'z':
			result.WriteString(strconv.Itoa(t.YearDay() - 1))
		case 't':
			firstDayNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
			lastDay := firstDayNextMonth.AddDate(0, 0, -1)
			result.WriteString(strconv.Itoa(lastDay.Day()))
		case 'L':
			year := t.Year()
			isLeap := (year%4 == 0 && year%100 != 0) || year%400 == 0
			if isLeap {
				result.WriteString("1")
			} else {
				result.WriteString("0")
			}
		case 'e':
			result.WriteString(t.Location().String())
		case 'P':
			_, offset := t.Zone()
			sign := "+"
			if offset < 0 {
				sign = "-"
				offset = -offset
			}
			result.WriteString(sign + twoDigit(offset/3600) + ":" + twoDigit((offset%3600)/60))
		case 'O':
			_, offset := t.Zone()
			sign := "+"
			if offset < 0 {
				sign = "-"
				offset = -offset
			}
			result.WriteString(sign + twoDigit(offset/3600) + twoDigit((offset%3600)/60))
		case 'c':
			result.WriteString(t.Format("2006-01-02T15:04:05-07:00"))
		case 'r':
			result.WriteString(t.Format("Mon, 02 Jan 2006 15:04:05 -0700"))
		case '\\':
			if i+1 < len(format) {
				i++
				result.WriteByte(format[i])
			}
		default:
			result.WriteByte(c)
		}
	}
	return result.String()
}

func twoDigit(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

func (f *DateFunction) GetName() string { return "date" }
var dateFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "format", 0, nil, nil),
	node.NewParameter(nil, "timestamp", 1, node.NewNullLiteral(nil), nil),
}

func (f *DateFunction) GetParams() []data.GetValue {
	return dateFunctionGetParams
}
var dateFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "format", 0, nil),
	node.NewVariable(nil, "timestamp", 1, nil),
}

func (f *DateFunction) GetVariables() []data.Variable {
	return dateFunctionGetVariables
}

// MktimeFunction 实现 mktime 函数
// mktime(?int $hour = null, ?int $minute = null, ?int $second = null, ?int $month = null, ?int $day = null, ?int $year = null): int|false
type MktimeFunction struct{}

func NewMktimeFunction() data.FuncStmt { return &MktimeFunction{} }

func (f *MktimeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	hourV, _ := ctx.GetIndexValue(0)
	minuteV, _ := ctx.GetIndexValue(1)
	secondV, _ := ctx.GetIndexValue(2)
	monthV, _ := ctx.GetIndexValue(3)
	dayV, _ := ctx.GetIndexValue(4)
	yearV, _ := ctx.GetIndexValue(5)

	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	month := int(now.Month())
	day := now.Day()
	year := now.Year()

	getInt := func(v data.GetValue, defaultVal int) int {
		if v == nil {
			return defaultVal
		}
		if iv, ok := v.(*data.IntValue); ok {
			return iv.Value
		}
		if ai, ok := v.(data.AsInt); ok {
			if iv, err := ai.AsInt(); err == nil {
				return iv
			}
		}
		return defaultVal
	}

	hour = getInt(hourV, hour)
	minute = getInt(minuteV, minute)
	second = getInt(secondV, second)
	month = getInt(monthV, month)
	day = getInt(dayV, day)
	year = getInt(yearV, year)

	t := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
	return data.NewIntValue(int(t.Unix())), nil
}

func (f *MktimeFunction) GetName() string { return "mktime" }
var mktimeFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "hour", 0, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "minute", 1, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "second", 2, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "month", 3, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "day", 4, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "year", 5, node.NewNullLiteral(nil), nil),
}

func (f *MktimeFunction) GetParams() []data.GetValue {
	return mktimeFunctionGetParams
}
var mktimeFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "hour", 0, nil),
	node.NewVariable(nil, "minute", 1, nil),
	node.NewVariable(nil, "second", 2, nil),
	node.NewVariable(nil, "month", 3, nil),
	node.NewVariable(nil, "day", 4, nil),
	node.NewVariable(nil, "year", 5, nil),
}

func (f *MktimeFunction) GetVariables() []data.Variable {
	return mktimeFunctionGetVariables
}

// GmmktimeFunction 实现 gmmktime 函数
// gmmktime(?int $hour = null, ?int $minute = null, ?int $second = null, ?int $month = null, ?int $day = null, ?int $year = null): int|false
type GmmktimeFunction struct{}

func NewGmmktimeFunction() data.FuncStmt { return &GmmktimeFunction{} }

func (f *GmmktimeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	hourV, _ := ctx.GetIndexValue(0)
	minuteV, _ := ctx.GetIndexValue(1)
	secondV, _ := ctx.GetIndexValue(2)
	monthV, _ := ctx.GetIndexValue(3)
	dayV, _ := ctx.GetIndexValue(4)
	yearV, _ := ctx.GetIndexValue(5)

	now := time.Now().UTC()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	month := int(now.Month())
	day := now.Day()
	year := now.Year()

	getInt := func(v data.GetValue, defaultVal int) int {
		if v == nil {
			return defaultVal
		}
		if iv, ok := v.(*data.IntValue); ok {
			return iv.Value
		}
		if ai, ok := v.(data.AsInt); ok {
			if iv, err := ai.AsInt(); err == nil {
				return iv
			}
		}
		return defaultVal
	}

	hour = getInt(hourV, hour)
	minute = getInt(minuteV, minute)
	second = getInt(secondV, second)
	month = getInt(monthV, month)
	day = getInt(dayV, day)
	year = getInt(yearV, year)

	t := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
	return data.NewIntValue(int(t.Unix())), nil
}

func (f *GmmktimeFunction) GetName() string { return "gmmktime" }
var gmmktimeFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "hour", 0, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "minute", 1, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "second", 2, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "month", 3, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "day", 4, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "year", 5, node.NewNullLiteral(nil), nil),
}

func (f *GmmktimeFunction) GetParams() []data.GetValue {
	return gmmktimeFunctionGetParams
}
var gmmktimeFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "hour", 0, nil),
	node.NewVariable(nil, "minute", 1, nil),
	node.NewVariable(nil, "second", 2, nil),
	node.NewVariable(nil, "month", 3, nil),
	node.NewVariable(nil, "day", 4, nil),
	node.NewVariable(nil, "year", 5, nil),
}

func (f *GmmktimeFunction) GetVariables() []data.Variable {
	return gmmktimeFunctionGetVariables
}

// CheckdateFunction 实现 checkdate 函数
// checkdate(int $month, int $day, int $year): bool
type CheckdateFunction struct{}

func NewCheckdateFunction() data.FuncStmt { return &CheckdateFunction{} }

func (f *CheckdateFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	monthV, _ := ctx.GetIndexValue(0)
	dayV, _ := ctx.GetIndexValue(1)
	yearV, _ := ctx.GetIndexValue(2)

	if monthV == nil || dayV == nil || yearV == nil {
		return data.NewBoolValue(false), nil
	}

	getInt := func(v data.GetValue) int {
		if iv, ok := v.(*data.IntValue); ok {
			return iv.Value
		}
		if ai, ok := v.(data.AsInt); ok {
			if iv, err := ai.AsInt(); err == nil {
				return iv
			}
		}
		return 0
	}

	month := getInt(monthV)
	day := getInt(dayV)
	year := getInt(yearV)

	if month < 1 || month > 12 {
		return data.NewBoolValue(false), nil
	}
	if day < 1 || day > 31 {
		return data.NewBoolValue(false), nil
	}
	if year < 1 || year > 32767 {
		return data.NewBoolValue(false), nil
	}

	// 检查具体月份的天数
	nextMonth := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, time.UTC)
	lastDay := nextMonth.AddDate(0, 0, -1).Day()
	if day > lastDay {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func (f *CheckdateFunction) GetName() string { return "checkdate" }
var checkdateFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "month", 0, nil, data.NewBaseType("int")),
	node.NewParameter(nil, "day", 1, nil, data.NewBaseType("int")),
	node.NewParameter(nil, "year", 2, nil, data.NewBaseType("int")),
}

func (f *CheckdateFunction) GetParams() []data.GetValue {
	return checkdateFunctionGetParams
}
var checkdateFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "month", 0, data.NewBaseType("int")),
	node.NewVariable(nil, "day", 1, data.NewBaseType("int")),
	node.NewVariable(nil, "year", 2, data.NewBaseType("int")),
}

func (f *CheckdateFunction) GetVariables() []data.Variable {
	return checkdateFunctionGetVariables
}

// GetdateFunction 实现 getdate 函数
// getdate(?int $timestamp = null): array
type GetdateFunction struct{}

func NewGetdateFunction() data.FuncStmt { return &GetdateFunction{} }

func (f *GetdateFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	tsV, _ := ctx.GetIndexValue(0)

	var t time.Time
	if tsV != nil {
		if _, isNull := tsV.(*data.NullValue); isNull {
			t = time.Now()
		} else if iv, ok := tsV.(*data.IntValue); ok {
			t = time.Unix(int64(iv.Value), 0)
		} else if ai, ok := tsV.(data.AsInt); ok {
			iv, err := ai.AsInt()
			if err == nil {
				t = time.Unix(int64(iv), 0)
			} else {
				t = time.Now()
			}
		} else {
			t = time.Now()
		}
	} else {
		t = time.Now()
	}

	// PHP 的 getdate 返回带字符串键的关联数组
	arr := &data.ArrayValue{List: []*data.ZVal{
		data.NewNamedZVal("seconds", data.NewIntValue(t.Second())),
		data.NewNamedZVal("minutes", data.NewIntValue(t.Minute())),
		data.NewNamedZVal("hours", data.NewIntValue(t.Hour())),
		data.NewNamedZVal("mday", data.NewIntValue(t.Day())),
		data.NewNamedZVal("wday", data.NewIntValue(int(t.Weekday()))),
		data.NewNamedZVal("mon", data.NewIntValue(int(t.Month()))),
		data.NewNamedZVal("year", data.NewIntValue(t.Year())),
		data.NewNamedZVal("yday", data.NewIntValue(t.YearDay()-1)),
		data.NewNamedZVal("weekday", data.NewStringValue(t.Weekday().String())),
		data.NewNamedZVal("month", data.NewStringValue(t.Month().String())),
		data.NewNamedZVal(data.IntArrayKeyName(0), data.NewIntValue(int(t.Unix()))),
	}}

	return arr, nil
}

func (f *GetdateFunction) GetName() string { return "getdate" }
var getdateFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "timestamp", 0, node.NewNullLiteral(nil), nil),
}

func (f *GetdateFunction) GetParams() []data.GetValue {
	return getdateFunctionGetParams
}
var getdateFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "timestamp", 0, nil),
}

func (f *GetdateFunction) GetVariables() []data.Variable {
	return getdateFunctionGetVariables
}
