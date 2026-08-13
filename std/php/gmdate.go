package php

import (
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GmdateFunction 实现 gmdate 函数
type GmdateFunction struct{}

func NewGmdateFunction() data.FuncStmt { return &GmdateFunction{} }

func (f *GmdateFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	formatValue, _ := ctx.GetIndexValue(0)
	timestampValue, _ := ctx.GetIndexValue(1)

	format := "Y-m-d H:i:s"
	if formatValue != nil {
		format = formatValue.AsString()
	}

	var t time.Time
	if timestampValue != nil {
		if iv, ok := timestampValue.(*data.IntValue); ok {
			t = time.Unix(int64(iv.Value), 0).UTC()
		} else {
			t = time.Now().UTC()
		}
	} else {
		t = time.Now().UTC()
	}

	result := phpDateFormat(format, t)
	return data.NewStringValue(result), nil
}

func phpDateFormat(format string, t time.Time) string {
	result := make([]byte, 0, len(format)*2)
	for i := 0; i < len(format); i++ {
		c := format[i]
		if c == '\\' && i+1 < len(format) {
			i++
			result = append(result, format[i])
			continue
		}
		switch c {
		case 'd':
			result = append(result, pad2(t.Day())...)
		case 'D':
			result = append(result, t.Format("Mon")...)
		case 'j':
			result = append(result, itoa(t.Day())...)
		case 'l':
			result = append(result, t.Format("Monday")...)
		case 'N':
			w := t.Weekday()
			if w == 0 {
				w = 7
			}
			result = append(result, byte('0'+w))
		case 'w':
			result = append(result, byte('0'+t.Weekday()))
		case 'z':
			result = append(result, itoa(t.YearDay()-1)...)
		case 'W':
			_, w := t.ISOWeek()
			result = append(result, pad2(w)...)
		case 'F':
			result = append(result, t.Format("January")...)
		case 'm':
			result = append(result, pad2(int(t.Month()))...)
		case 'M':
			result = append(result, t.Format("Jan")...)
		case 'n':
			result = append(result, itoa(int(t.Month()))...)
		case 't':
			lastDay := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location())
			result = append(result, itoa(lastDay.Day())...)
		case 'Y':
			result = append(result, itoa(t.Year())...)
		case 'y':
			y := t.Year() % 100
			result = append(result, pad2(y)...)
		case 'a':
			result = append(result, t.Format("pm")...)
		case 'A':
			result = append(result, t.Format("PM")...)
		case 'g':
			h := t.Hour() % 12
			if h == 0 {
				h = 12
			}
			result = append(result, itoa(h)...)
		case 'G':
			result = append(result, itoa(t.Hour())...)
		case 'h':
			h := t.Hour() % 12
			if h == 0 {
				h = 12
			}
			result = append(result, pad2(h)...)
		case 'H':
			result = append(result, pad2(t.Hour())...)
		case 'i':
			result = append(result, pad2(t.Minute())...)
		case 's':
			result = append(result, pad2(t.Second())...)
		case 'v':
			result = append(result, pad3(t.Nanosecond()/1000000)...)
		case 'u':
			result = append(result, pad6(t.Nanosecond()/1000)...)
		case 'e':
			result = append(result, t.Location().String()...)
		case 'O':
			result = append(result, t.Format("-0700")...)
		case 'P':
			result = append(result, t.Format("-07:00")...)
		case 'T':
			result = append(result, t.Format("MST")...)
		case 'Z':
			_, offset := t.Zone()
			result = append(result, itoa(offset)...)
		case 'c':
			result = append(result, t.Format("2006-01-02T15:04:05-07:00")...)
		case 'r':
			result = append(result, t.Format("Mon, 02 Jan 2006 15:04:05 -0700")...)
		case 'U':
			result = append(result, itoa(int(t.Unix()))...)
		default:
			result = append(result, c)
		}
	}
	return string(result)
}

func pad2(n int) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}

func pad3(n int) string {
	switch {
	case n < 10:
		return "00" + itoa(n)
	case n < 100:
		return "0" + itoa(n)
	default:
		return itoa(n)
	}
}

func pad6(n int) string {
	s := itoa(n)
	for len(s) < 6 {
		s = "0" + s
	}
	return s
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		s = string(byte('0'+n%10)) + s
		n /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}

func (f *GmdateFunction) GetName() string { return "gmdate" }
func (f *GmdateFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "format", 0, nil, nil),
		node.NewParameter(nil, "timestamp", 1, node.NewNullLiteral(nil), nil),
	}
}
func (f *GmdateFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "timestamp", 1, data.NewNullableType(data.NewBaseType("int"))),
	}
}
