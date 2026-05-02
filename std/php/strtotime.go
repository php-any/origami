package php

import (
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrtotimeFunction 实现 strtotime 函数（简化版）
type StrtotimeFunction struct{}

func NewStrtotimeFunction() data.FuncStmt { return &StrtotimeFunction{} }

func (f *StrtotimeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	datetimeValue, _ := ctx.GetIndexValue(0)
	baseTimestampValue, _ := ctx.GetIndexValue(1)

	if datetimeValue == nil {
		return data.NewBoolValue(false), nil
	}

	datetimeStr := datetimeValue.AsString()

	baseTime := time.Now()
	if baseTimestampValue != nil {
		if iv, ok := baseTimestampValue.(*data.IntValue); ok {
			baseTime = time.Unix(int64(iv.Value), 0)
		}
	}

	t, err := parseRelativeTime(datetimeStr, baseTime)
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	return data.NewIntValue(int(t.Unix())), nil
}

func parseRelativeTime(s string, base time.Time) (time.Time, error) {
	s = trimLower(s)

	switch {
	case s == "now":
		return base, nil

	case s == "yesterday":
		return base.AddDate(0, 0, -1), nil
	case s == "tomorrow":
		return base.AddDate(0, 0, 1), nil

	case s == "midnight":
		y, m, d := base.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, base.Location()), nil
	case s == "today":
		y, m, d := base.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, base.Location()), nil

	// Try "now - X seconds/minutes/hours/days"
	case hasPrefix(s, "now "):
		rest := s[4:]
		t, err := parseModifier(rest, base)
		if err == nil {
			return t, nil
		}
	}

	// Try standard Go formats
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006-01-02 15:04",
		"15:04:05",
		"2006/01/02",
		"01/02/2006",
		"02 Jan 2006",
		"02 Jan 06",
	}
	for _, fmt := range formats {
		if t, err := time.Parse(fmt, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, &parseError{s}
}

type parseError struct{ s string }

func (e *parseError) Error() string { return "cannot parse: " + e.s }

func parseModifier(s string, base time.Time) (time.Time, error) {
	s = trimLower(s)
	if s == "" {
		return base, nil
	}

	// Handle "now + 1 hour", "now - 7200 seconds"
	op := '+'
	if s[0] == '-' || s[0] == '+' {
		op = rune(s[0])
		s = trimLower(s[1:])
	}

	var num int
	var unit string
	for i, c := range s {
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0')
		} else if c == ' ' {
			unit = trimLower(s[i:])
			break
		} else {
			unit = trimLower(s[i:])
			break
		}
	}

	if op == '-' {
		num = -num
	}

	switch {
	case unit == "second" || unit == "seconds":
		return base.Add(time.Duration(num) * time.Second), nil
	case unit == "minute" || unit == "minutes":
		return base.Add(time.Duration(num) * time.Minute), nil
	case unit == "hour" || unit == "hours":
		return base.Add(time.Duration(num) * time.Hour), nil
	case unit == "day" || unit == "days":
		return base.AddDate(0, 0, num), nil
	case unit == "week" || unit == "weeks":
		return base.AddDate(0, 0, 7*num), nil
	case unit == "month" || unit == "months":
		return base.AddDate(0, num, 0), nil
	case unit == "year" || unit == "years":
		return base.AddDate(num, 0, 0), nil
	default:
		return base, nil
	}
}

func trimLower(s string) string {
	result := make([]byte, 0, len(s))
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	for i := start; i < end; i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		result = append(result, c)
	}
	return string(result)
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func (f *StrtotimeFunction) GetName() string { return "strtotime" }
func (f *StrtotimeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "datetime", 0, nil, nil),
		node.NewParameter(nil, "baseTimestamp", 1, node.NewNullLiteral(nil), nil),
	}
}
func (f *StrtotimeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "datetime", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "baseTimestamp", 1, data.NewNullableType(data.NewBaseType("int"))),
	}
}
