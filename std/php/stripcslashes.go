package php

import (
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StripslashesFunction 实现 stripslashes 函数
type StripslashesFunction struct{}

func NewStripslashesFunction() data.FuncStmt { return &StripslashesFunction{} }

func (f *StripslashesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewStringValue(""), nil
	}
	s := v.AsString()
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			result.WriteByte(s[i+1])
			i++
		} else {
			result.WriteByte(s[i])
		}
	}
	return data.NewStringValue(result.String()), nil
}

func (f *StripslashesFunction) GetName() string { return "stripslashes" }
func (f *StripslashesFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}
func (f *StripslashesFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}

// StripsCslashesFunction 实现 stripcslashes 函数
type StripsCslashesFunction struct{}

func NewStripsCslashesFunction() data.FuncStmt { return &StripsCslashesFunction{} }

func (f *StripsCslashesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewStringValue(""), nil
	}
	s := v.AsString()
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			i++
			ch := s[i]
			switch ch {
			case 'n':
				result.WriteByte('\n')
			case 'r':
				result.WriteByte('\r')
			case 't':
				result.WriteByte('\t')
			case 'v':
				result.WriteByte('\v')
			case 'f':
				result.WriteByte('\f')
			case 'e':
				result.WriteByte('\x1b')
			case '\\':
				result.WriteByte('\\')
			case '"':
				result.WriteByte('"')
			case '\'':
				result.WriteByte('\'')
			case 'x':
				// hex escape: \xHH
				if i+2 < len(s) {
					hex := s[i+1 : i+3]
					if val, err := strconv.ParseUint(hex, 16, 8); err == nil {
						result.WriteByte(byte(val))
						i += 2
					} else {
						result.WriteByte('\\')
						result.WriteByte(ch)
					}
				} else {
					result.WriteByte('\\')
					result.WriteByte(ch)
				}
			case '0', '1', '2', '3', '4', '5', '6', '7':
				// octal escape: \OOO (up to 3 digits)
				octal := string(ch)
				j := i + 1
				for j < len(s) && j < i+3 && s[j] >= '0' && s[j] <= '7' {
					octal += string(s[j])
					j++
				}
				if val, err := strconv.ParseUint(octal, 8, 8); err == nil {
					result.WriteByte(byte(val))
					i = j - 1
				} else {
					result.WriteByte('\\')
					result.WriteByte(ch)
				}
			default:
				result.WriteByte('\\')
				result.WriteByte(ch)
			}
		} else {
			result.WriteByte(s[i])
		}
	}
	return data.NewStringValue(result.String()), nil
}

func (f *StripsCslashesFunction) GetName() string { return "stripcslashes" }
func (f *StripsCslashesFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}
func (f *StripsCslashesFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
