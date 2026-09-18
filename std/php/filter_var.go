package php

import (
	"net/mail"
	"net/url"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewFilterVarFunction() data.FuncStmt {
	return &FilterVarFunction{}
}

type FilterVarFunction struct{}

func (fn *FilterVarFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	filter, _ := ctx.GetIndexValue(1)

	filterInt := 516 // FILTER_DEFAULT
	if fv, ok := filter.(data.AsInt); ok {
		if n, err := fv.AsInt(); err == nil {
			filterInt = n
		}
	}

	switch filterInt {
	case 257: // FILTER_VALIDATE_INT
		if v, ok := value.(*data.StringValue); ok {
			if _, err := strconv.Atoi(v.Value); err == nil {
				return value, nil
			}
		}
		if _, ok := value.(*data.IntValue); ok {
			return value, nil
		}
		return data.NewBoolValue(false), nil
	case 258: // FILTER_VALIDATE_BOOLEAN
		if _, ok := value.(*data.BoolValue); ok {
			return value, nil
		}
		return data.NewBoolValue(false), nil
	case 273: // FILTER_VALIDATE_URL
		s := phpFilterString(value)
		if phpFilterValidateURL(s) {
			return data.NewStringValue(s), nil
		}
		return data.NewBoolValue(false), nil
	case 274: // FILTER_VALIDATE_EMAIL
		s := phpFilterString(value)
		if phpFilterValidateEmail(s) {
			return data.NewStringValue(s), nil
		}
		return data.NewBoolValue(false), nil
	default:
		return value, nil
	}
}

func phpFilterString(value data.Value) string {
	if value == nil {
		return ""
	}
	if s, ok := value.(interface{ AsString() string }); ok {
		return s.AsString()
	}
	return ""
}

// phpFilterValidateURL 对齐 PHP FILTER_VALIDATE_URL：必须有 scheme 与 host。
// 相对路径（css/app.css、/js/x）必须失败，否则 Laravel UrlGenerator::asset
// 会把资源路径当成「已是 URL」原样输出，后台页变成 /admin/css/... 404。
func phpFilterValidateURL(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	if u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func phpFilterValidateEmail(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || !strings.Contains(s, "@") {
		return false
	}
	_, err := mail.ParseAddress(s)
	return err == nil
}

func (fn *FilterVarFunction) GetName() string {
	return "filter_var"
}

func (fn *FilterVarFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
		node.NewParameter(nil, "filter", 1, data.NewIntValue(516), nil),
	}
}

func (fn *FilterVarFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.Mixed{}),
		node.NewVariable(nil, "filter", 1, data.NewBaseType("int")),
	}
}
