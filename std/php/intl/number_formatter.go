package intl

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// PHP NumberFormatter 样式/属性常量（与 ext/intl 对齐的常用子集）
const (
	nfStylePatternDecimal   = 0
	nfStyleDecimal          = 1
	nfStyleCurrency         = 2
	nfStylePercent          = 3
	nfStyleScientific       = 4
	nfStyleSpellout         = 5
	nfStyleOrdinal          = 6
	nfStyleDuration         = 7
	nfStylePatternRuleBased = 8
	nfStyleCurrencyAccounting = 12
	nfStyleDefault          = 1

	nfAttrParseIntOnly         = 0
	nfAttrGroupingUsed         = 1
	nfAttrDecimalAlwaysShown   = 2
	nfAttrMaxIntegerDigits     = 3
	nfAttrMinIntegerDigits     = 4
	nfAttrIntegerDigits        = 5
	nfAttrMaxFractionDigits    = 6
	nfAttrMinFractionDigits    = 7
	nfAttrFractionDigits       = 8
	nfAttrMultiplier           = 9
	nfAttrGroupingSize         = 10
	nfAttrRoundingMode         = 11
	nfAttrRoundingIncrement    = 12
	nfAttrFormatWidth          = 13
	nfAttrPaddingPosition      = 14
	nfAttrSecondaryGroupingSize = 15
	nfAttrSignificantDigitsUsed = 16
	nfAttrMinSignificantDigits = 17
	nfAttrMaxSignificantDigits = 18
	nfAttrLenientParse         = 19

	nfTextAttrDefaultRuleset = 0

	nfTypeDefault = 0
	nfTypeInt32   = 1
	nfTypeInt64   = 2
	nfTypeDouble  = 3
	nfTypeCurrency = 4
)

const (
	nfPropLocale = "__nf_locale__"
	nfPropStyle  = "__nf_style__"
	nfPropAttrs  = "__nf_attrs__"
	nfPropTexts  = "__nf_texts__"
)

// NumberFormatterClass 最小 intl NumberFormatter，满足 Laravel Number::currency/format。
type NumberFormatterClass struct {
	node.Node
	StaticProperty map[string]data.Value
}

func NewNumberFormatterClass() *NumberFormatterClass {
	return &NumberFormatterClass{
		StaticProperty: map[string]data.Value{
			"PATTERN_DECIMAL":     data.NewIntValue(nfStylePatternDecimal),
			"DECIMAL":             data.NewIntValue(nfStyleDecimal),
			"CURRENCY":            data.NewIntValue(nfStyleCurrency),
			"PERCENT":             data.NewIntValue(nfStylePercent),
			"SCIENTIFIC":          data.NewIntValue(nfStyleScientific),
			"SPELLOUT":            data.NewIntValue(nfStyleSpellout),
			"ORDINAL":             data.NewIntValue(nfStyleOrdinal),
			"DURATION":            data.NewIntValue(nfStyleDuration),
			"PATTERN_RULEBASED":   data.NewIntValue(nfStylePatternRuleBased),
			"CURRENCY_ACCOUNTING": data.NewIntValue(nfStyleCurrencyAccounting),
			"DEFAULT_STYLE":       data.NewIntValue(nfStyleDefault),
			"IGNORE":              data.NewIntValue(0),

			"PARSE_INT_ONLY":           data.NewIntValue(nfAttrParseIntOnly),
			"GROUPING_USED":            data.NewIntValue(nfAttrGroupingUsed),
			"DECIMAL_ALWAYS_SHOWN":     data.NewIntValue(nfAttrDecimalAlwaysShown),
			"MAX_INTEGER_DIGITS":       data.NewIntValue(nfAttrMaxIntegerDigits),
			"MIN_INTEGER_DIGITS":       data.NewIntValue(nfAttrMinIntegerDigits),
			"INTEGER_DIGITS":           data.NewIntValue(nfAttrIntegerDigits),
			"MAX_FRACTION_DIGITS":      data.NewIntValue(nfAttrMaxFractionDigits),
			"MIN_FRACTION_DIGITS":      data.NewIntValue(nfAttrMinFractionDigits),
			"FRACTION_DIGITS":          data.NewIntValue(nfAttrFractionDigits),
			"MULTIPLIER":               data.NewIntValue(nfAttrMultiplier),
			"GROUPING_SIZE":            data.NewIntValue(nfAttrGroupingSize),
			"ROUNDING_MODE":            data.NewIntValue(nfAttrRoundingMode),
			"ROUNDING_INCREMENT":       data.NewIntValue(nfAttrRoundingIncrement),
			"FORMAT_WIDTH":             data.NewIntValue(nfAttrFormatWidth),
			"PADDING_POSITION":         data.NewIntValue(nfAttrPaddingPosition),
			"SECONDARY_GROUPING_SIZE":  data.NewIntValue(nfAttrSecondaryGroupingSize),
			"SIGNIFICANT_DIGITS_USED":  data.NewIntValue(nfAttrSignificantDigitsUsed),
			"MIN_SIGNIFICANT_DIGITS":   data.NewIntValue(nfAttrMinSignificantDigits),
			"MAX_SIGNIFICANT_DIGITS":   data.NewIntValue(nfAttrMaxSignificantDigits),
			"LENIENT_PARSE":            data.NewIntValue(nfAttrLenientParse),

			"DEFAULT_RULESET": data.NewIntValue(nfTextAttrDefaultRuleset),

			"TYPE_DEFAULT":  data.NewIntValue(nfTypeDefault),
			"TYPE_INT32":    data.NewIntValue(nfTypeInt32),
			"TYPE_INT64":    data.NewIntValue(nfTypeInt64),
			"TYPE_DOUBLE":   data.NewIntValue(nfTypeDouble),
			"TYPE_CURRENCY": data.NewIntValue(nfTypeCurrency),
		},
	}
}

func (c *NumberFormatterClass) GetName() string    { return "NumberFormatter" }
func (c *NumberFormatterClass) GetExtend() *string { return nil }
func (c *NumberFormatterClass) GetImplements() []string {
	return nil
}
func (c *NumberFormatterClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *NumberFormatterClass) GetPropertyList() []data.Property              { return nil }
func (c *NumberFormatterClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := c.StaticProperty[name]; ok {
		return v, true
	}
	return nil, false
}
func (c *NumberFormatterClass) GetConstruct() data.Method {
	return &NumberFormatterConstructMethod{}
}
func (c *NumberFormatterClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *NumberFormatterClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &NumberFormatterConstructMethod{}, true
	case "format":
		return &NumberFormatterFormatMethod{}, true
	case "formatCurrency":
		return &NumberFormatterFormatCurrencyMethod{}, true
	case "setAttribute":
		return &NumberFormatterSetAttributeMethod{}, true
	case "setTextAttribute":
		return &NumberFormatterSetTextAttributeMethod{}, true
	case "getAttribute":
		return &NumberFormatterGetAttributeMethod{}, true
	}
	return nil, false
}

func (c *NumberFormatterClass) GetMethods() []data.Method {
	return []data.Method{
		&NumberFormatterConstructMethod{},
		&NumberFormatterFormatMethod{},
		&NumberFormatterFormatCurrencyMethod{},
		&NumberFormatterSetAttributeMethod{},
		&NumberFormatterSetTextAttributeMethod{},
		&NumberFormatterGetAttributeMethod{},
	}
}

type NumberFormatterConstructMethod struct{}

func (m *NumberFormatterConstructMethod) GetName() string            { return "__construct" }
func (m *NumberFormatterConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *NumberFormatterConstructMethod) GetIsStatic() bool          { return false }
var numberFormatterConstructMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "locale", 0, nil, data.String{}),
	node.NewParameter(nil, "style", 1, nil, data.Int{}),
	node.NewParameter(nil, "pattern", 2, data.NewNullValue(), nil),
}

func (m *NumberFormatterConstructMethod) GetParams() []data.GetValue {
	return numberFormatterConstructMethodGetParams
}
var numberFormatterConstructMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "locale", 0, data.String{}),
	node.NewVariable(nil, "style", 1, data.Int{}),
	node.NewVariable(nil, "pattern", 2, nil),
}

func (m *NumberFormatterConstructMethod) GetVariables() []data.Variable {
	return numberFormatterConstructMethodGetVariables
}
func (m *NumberFormatterConstructMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

func nfInstance(ctx data.Context) *data.ClassValue {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		return cmc.ClassValue
	}
	if cv, ok := ctx.(*data.ClassValue); ok {
		return cv
	}
	return nil
}

func (m *NumberFormatterConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	locale, _ := utils.ConvertFromIndex[string](ctx, 0)
	style, _ := utils.ConvertFromIndex[int](ctx, 1)
	if locale == "" {
		locale = "en"
	}
	obj := nfInstance(ctx)
	if obj == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("NumberFormatter::__construct 需要实例上下文"))
	}
	obj.SetProperty(nfPropLocale, data.NewStringValue(locale))
	obj.SetProperty(nfPropStyle, data.NewIntValue(style))
	obj.SetProperty(nfPropAttrs, data.NewObjectValue())
	obj.SetProperty(nfPropTexts, data.NewObjectValue())
	return nil, nil
}

type NumberFormatterSetAttributeMethod struct{}

func (m *NumberFormatterSetAttributeMethod) GetName() string            { return "setAttribute" }
func (m *NumberFormatterSetAttributeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *NumberFormatterSetAttributeMethod) GetIsStatic() bool          { return false }
var numberFormatterSetAttributeMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "attr", 0, nil, data.Int{}),
	node.NewParameter(nil, "value", 1, nil, nil),
}

func (m *NumberFormatterSetAttributeMethod) GetParams() []data.GetValue {
	return numberFormatterSetAttributeMethodGetParams
}
var numberFormatterSetAttributeMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "attr", 0, data.Int{}),
	node.NewVariable(nil, "value", 1, nil),
}

func (m *NumberFormatterSetAttributeMethod) GetVariables() []data.Variable {
	return numberFormatterSetAttributeMethodGetVariables
}
func (m *NumberFormatterSetAttributeMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *NumberFormatterSetAttributeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	obj := nfInstance(ctx)
	if obj == nil {
		return data.NewBoolValue(false), nil
	}
	attr, _ := utils.ConvertFromIndex[int](ctx, 0)
	val, acl := ctx.GetVariableValue(node.NewVariable(nil, "value", 1, nil))
	if acl != nil {
		return nil, acl
	}
	attrs := nfAttrsObject(obj)
	attrs.SetProperty(strconv.Itoa(attr), val)
	return data.NewBoolValue(true), nil
}

type NumberFormatterSetTextAttributeMethod struct{}

func (m *NumberFormatterSetTextAttributeMethod) GetName() string { return "setTextAttribute" }
func (m *NumberFormatterSetTextAttributeMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *NumberFormatterSetTextAttributeMethod) GetIsStatic() bool { return false }
var numberFormatterSetTextAttributeMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "attr", 0, nil, data.Int{}),
	node.NewParameter(nil, "value", 1, nil, data.String{}),
}

func (m *NumberFormatterSetTextAttributeMethod) GetParams() []data.GetValue {
	return numberFormatterSetTextAttributeMethodGetParams
}
var numberFormatterSetTextAttributeMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "attr", 0, data.Int{}),
	node.NewVariable(nil, "value", 1, data.String{}),
}

func (m *NumberFormatterSetTextAttributeMethod) GetVariables() []data.Variable {
	return numberFormatterSetTextAttributeMethodGetVariables
}
func (m *NumberFormatterSetTextAttributeMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *NumberFormatterSetTextAttributeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	obj := nfInstance(ctx)
	if obj == nil {
		return data.NewBoolValue(false), nil
	}
	attr, _ := utils.ConvertFromIndex[int](ctx, 0)
	val, _ := utils.ConvertFromIndex[string](ctx, 1)
	texts := nfTextsObject(obj)
	texts.SetProperty(strconv.Itoa(attr), data.NewStringValue(val))
	return data.NewBoolValue(true), nil
}

type NumberFormatterGetAttributeMethod struct{}

func (m *NumberFormatterGetAttributeMethod) GetName() string            { return "getAttribute" }
func (m *NumberFormatterGetAttributeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *NumberFormatterGetAttributeMethod) GetIsStatic() bool          { return false }
var numberFormatterGetAttributeMethodGetParams = []data.GetValue{node.NewParameter(nil, "attr", 0, nil, data.Int{})}

func (m *NumberFormatterGetAttributeMethod) GetParams() []data.GetValue {
	return numberFormatterGetAttributeMethodGetParams
}
var numberFormatterGetAttributeMethodGetVariables = []data.Variable{node.NewVariable(nil, "attr", 0, data.Int{})}

func (m *NumberFormatterGetAttributeMethod) GetVariables() []data.Variable {
	return numberFormatterGetAttributeMethodGetVariables
}
func (m *NumberFormatterGetAttributeMethod) GetReturnType() data.Types {
	return data.NewBaseType("int|float")
}

func (m *NumberFormatterGetAttributeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	obj := nfInstance(ctx)
	if obj == nil {
		return data.NewIntValue(0), nil
	}
	attr, _ := utils.ConvertFromIndex[int](ctx, 0)
	if v, ok := nfGetAttr(obj, attr); ok {
		return v, nil
	}
	return data.NewIntValue(0), nil
}

type NumberFormatterFormatMethod struct{}

func (m *NumberFormatterFormatMethod) GetName() string            { return "format" }
func (m *NumberFormatterFormatMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *NumberFormatterFormatMethod) GetIsStatic() bool          { return false }
var numberFormatterFormatMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "num", 0, nil, nil),
	node.NewParameter(nil, "type", 1, data.NewIntValue(nfTypeDefault), data.Int{}),
}

func (m *NumberFormatterFormatMethod) GetParams() []data.GetValue {
	return numberFormatterFormatMethodGetParams
}
var numberFormatterFormatMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "num", 0, nil),
	node.NewVariable(nil, "type", 1, data.Int{}),
}

func (m *NumberFormatterFormatMethod) GetVariables() []data.Variable {
	return numberFormatterFormatMethodGetVariables
}
func (m *NumberFormatterFormatMethod) GetReturnType() data.Types {
	return data.NewBaseType("string|false")
}

func (m *NumberFormatterFormatMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	obj := nfInstance(ctx)
	if obj == nil {
		return data.NewBoolValue(false), nil
	}
	num := nfToFloat(ctx, 0)
	style := nfStyle(obj)
	digits := nfFractionDigits(obj)
	switch style {
	case nfStylePercent:
		return data.NewStringValue(nfFormatFixed(num*100, digits) + "%"), nil
	case nfStyleCurrency:
		return data.NewStringValue(nfFormatFixed(num, digits)), nil
	default:
		return data.NewStringValue(nfFormatFixed(num, digits)), nil
	}
}

type NumberFormatterFormatCurrencyMethod struct{}

func (m *NumberFormatterFormatCurrencyMethod) GetName() string { return "formatCurrency" }
func (m *NumberFormatterFormatCurrencyMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *NumberFormatterFormatCurrencyMethod) GetIsStatic() bool { return false }
var numberFormatterFormatCurrencyMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "amount", 0, nil, nil),
	node.NewParameter(nil, "currency", 1, nil, data.String{}),
}

func (m *NumberFormatterFormatCurrencyMethod) GetParams() []data.GetValue {
	return numberFormatterFormatCurrencyMethodGetParams
}
var numberFormatterFormatCurrencyMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "amount", 0, nil),
	node.NewVariable(nil, "currency", 1, data.String{}),
}

func (m *NumberFormatterFormatCurrencyMethod) GetVariables() []data.Variable {
	return numberFormatterFormatCurrencyMethodGetVariables
}
func (m *NumberFormatterFormatCurrencyMethod) GetReturnType() data.Types {
	return data.NewBaseType("string|false")
}

func (m *NumberFormatterFormatCurrencyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	obj := nfInstance(ctx)
	if obj == nil {
		return data.NewBoolValue(false), nil
	}
	amount := nfToFloat(ctx, 0)
	currency, _ := utils.ConvertFromIndex[string](ctx, 1)
	if currency == "" {
		currency = "USD"
	}
	digits := nfFractionDigits(obj)
	if digits < 0 {
		digits = 2
	}
	formatted := nfFormatFixed(amount, digits)
	symbol := nfCurrencySymbol(currency)
	if symbol != "" {
		return data.NewStringValue(symbol + formatted), nil
	}
	return data.NewStringValue(currency + " " + formatted), nil
}

func nfAttrsObject(obj *data.ClassValue) *data.ObjectValue {
	v, _ := obj.GetProperty(nfPropAttrs)
	if ov, ok := v.(*data.ObjectValue); ok {
		return ov
	}
	ov := data.NewObjectValue()
	obj.SetProperty(nfPropAttrs, ov)
	return ov
}

func nfTextsObject(obj *data.ClassValue) *data.ObjectValue {
	v, _ := obj.GetProperty(nfPropTexts)
	if ov, ok := v.(*data.ObjectValue); ok {
		return ov
	}
	ov := data.NewObjectValue()
	obj.SetProperty(nfPropTexts, ov)
	return ov
}

func nfStyle(obj *data.ClassValue) int {
	v, _ := obj.GetProperty(nfPropStyle)
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return n
		}
	}
	return nfStyleDecimal
}

func nfGetAttr(obj *data.ClassValue, attr int) (data.Value, bool) {
	attrs := nfAttrsObject(obj)
	v, ctl := attrs.GetProperty(strconv.Itoa(attr))
	if ctl != nil || v == nil {
		return nil, false
	}
	if _, isNull := v.(*data.NullValue); isNull {
		return nil, false
	}
	return v, true
}

func nfFractionDigits(obj *data.ClassValue) int {
	if v, ok := nfGetAttr(obj, nfAttrFractionDigits); ok {
		if iv, ok := v.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				return n
			}
		}
	}
	if v, ok := nfGetAttr(obj, nfAttrMaxFractionDigits); ok {
		if iv, ok := v.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				return n
			}
		}
	}
	style := nfStyle(obj)
	if style == nfStyleCurrency || style == nfStyleCurrencyAccounting {
		return 2
	}
	if style == nfStylePercent {
		return 0
	}
	return -1 // 默认不强制小数位
}

func nfToFloat(ctx data.Context, index int) float64 {
	if v, err := utils.ConvertFromIndex[float64](ctx, index); err == nil {
		return v
	}
	if iv, err := utils.ConvertFromIndex[int](ctx, index); err == nil {
		return float64(iv)
	}
	if sv, err := utils.ConvertFromIndex[string](ctx, index); err == nil {
		f, err2 := strconv.ParseFloat(sv, 64)
		if err2 == nil {
			return f
		}
	}
	return 0
}

func nfFormatFixed(num float64, digits int) string {
	if math.IsNaN(num) {
		return "NaN"
	}
	if math.IsInf(num, 1) {
		return "∞"
	}
	if math.IsInf(num, -1) {
		return "-∞"
	}
	if digits < 0 {
		// 尽量去掉多余尾零
		s := strconv.FormatFloat(num, 'f', -1, 64)
		return s
	}
	s := strconv.FormatFloat(num, 'f', digits, 64)
	return s
}

func nfCurrencySymbol(code string) string {
	switch strings.ToUpper(code) {
	case "CNY", "RMB":
		return "¥"
	case "USD":
		return "$"
	case "EUR":
		return "€"
	case "GBP":
		return "£"
	case "JPY":
		return "¥"
	default:
		return ""
	}
}
