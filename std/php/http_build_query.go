package php

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

const (
	phpQueryRFC1738 = 1
	phpQueryRFC3986 = 2
)

// HttpBuildQueryFunction 实现 PHP 内置函数 http_build_query
//
//	http_build_query(array|object $data, string $numeric_prefix = "", ?string $arg_separator = null, int $encoding_type = PHP_QUERY_RFC1738): string
func NewHttpBuildQueryFunction() data.FuncStmt {
	return &HttpBuildQueryFunction{}
}

type HttpBuildQueryFunction struct{}

type queryEncodeFunc func(string) string

type httpBuildQueryConfig struct {
	numericPrefix string
	separator     string
	encode        queryEncodeFunc
}

func (f *HttpBuildQueryFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	dataVal, _ := ctx.GetIndexValue(0)
	if dataVal == nil {
		return data.NewStringValue(""), nil
	}
	val, ok := dataVal.(data.Value)
	if !ok {
		return data.NewStringValue(""), nil
	}

	numericPrefix := ""
	if v, _ := ctx.GetIndexValue(1); v != nil {
		if _, ok := v.(*data.NullValue); !ok {
			numericPrefix = v.AsString()
		}
	}

	separator := "&"
	if v, _ := ctx.GetIndexValue(2); v != nil {
		if _, ok := v.(*data.NullValue); !ok {
			separator = v.AsString()
		}
	}

	encodingType := phpQueryRFC1738
	if v, _ := ctx.GetIndexValue(3); v != nil {
		if _, ok := v.(*data.NullValue); !ok {
			if n, err := utils.ConvertFromIndex[int](ctx, 3); err == nil {
				encodingType = n
			}
		}
	}

	var encode queryEncodeFunc
	switch encodingType {
	case phpQueryRFC3986:
		encode = url.PathEscape
	default:
		encode = url.QueryEscape
	}

	cfg := httpBuildQueryConfig{
		numericPrefix: numericPrefix,
		separator:     separator,
		encode:        encode,
	}

	pairs := make([]string, 0)
	appendQueryPairs("", val, true, cfg, &pairs)
	return data.NewStringValue(strings.Join(pairs, cfg.separator)), nil
}

func appendQueryPairs(prefix string, val data.Value, topLevel bool, cfg httpBuildQueryConfig, pairs *[]string) {
	switch v := val.(type) {
	case *data.ArrayValue:
		for i, zval := range v.List {
			if zval == nil {
				continue
			}
			keyStr, numKey, isNumeric := arrayEntryKey(zval, i)
			name := formatQueryName(prefix, keyStr, numKey, isNumeric, topLevel, cfg.numericPrefix)
			appendQueryValue(name, zval.Value, false, cfg, pairs)
		}
	case *data.ObjectValue:
		v.RangeProperties(func(key string, value data.Value) bool {
			keyStr, numKey, isNumeric := objectEntryKey(key)
			name := formatQueryName(prefix, keyStr, numKey, isNumeric, topLevel, cfg.numericPrefix)
			appendQueryValue(name, value, false, cfg, pairs)
			return true
		})
	case *data.ClassValue:
		v.RangeProperties(func(key string, value data.Value) bool {
			keyStr, numKey, isNumeric := objectEntryKey(key)
			name := formatQueryName(prefix, keyStr, numKey, isNumeric, topLevel, cfg.numericPrefix)
			appendQueryValue(name, value, false, cfg, pairs)
			return true
		})
	default:
		if prefix != "" {
			appendScalarPair(prefix, val, cfg, pairs)
		}
	}
}

func appendQueryValue(name string, val data.Value, topLevel bool, cfg httpBuildQueryConfig, pairs *[]string) {
	switch val.(type) {
	case *data.ArrayValue, *data.ObjectValue, *data.ClassValue:
		appendQueryPairs(name, val, false, cfg, pairs)
	default:
		appendScalarPair(name, val, cfg, pairs)
	}
}

func appendScalarPair(name string, val data.Value, cfg httpBuildQueryConfig, pairs *[]string) {
	value := ""
	if val != nil {
		if _, ok := val.(*data.NullValue); !ok {
			value = val.AsString()
		}
	}
	*pairs = append(*pairs, cfg.encode(name)+"="+cfg.encode(value))
}

func objectEntryKey(key string) (keyStr string, numKey int, isNumeric bool) {
	if n, ok := data.ParseIntArrayKeyName(key); ok {
		return key, n, true
	}
	return key, 0, false
}

func arrayEntryKey(zval *data.ZVal, index int) (keyStr string, numKey int, isNumeric bool) {
	if zval.Name != "" {
		keyStr = zval.Name
		if n, ok := data.ParseIntArrayKeyName(zval.Name); ok {
			return keyStr, n, true
		}
		return keyStr, 0, false
	}
	return strconv.Itoa(index), index, true
}

func formatQueryName(prefix, keyStr string, numKey int, isNumeric, topLevel bool, numericPrefix string) string {
	if prefix == "" {
		if topLevel && isNumeric && numericPrefix != "" {
			return numericPrefix + strconv.Itoa(numKey)
		}
		if isNumeric {
			return strconv.Itoa(numKey)
		}
		return keyStr
	}
	if isNumeric {
		return prefix + "[" + strconv.Itoa(numKey) + "]"
	}
	return prefix + "[" + keyStr + "]"
}

func (f *HttpBuildQueryFunction) GetName() string {
	return "http_build_query"
}

func (f *HttpBuildQueryFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, nil),
		node.NewParameter(nil, "numeric_prefix", 1, node.NewStringLiteral(nil, ""), nil),
		node.NewParameter(nil, "arg_separator", 2, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "encoding_type", 3, node.NewIntLiteral(nil, "1"), nil),
	}
}

func (f *HttpBuildQueryFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "numeric_prefix", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "arg_separator", 2, data.NewNullableType(data.NewBaseType("string"))),
		node.NewVariable(nil, "encoding_type", 3, data.NewBaseType("int")),
	}
}
