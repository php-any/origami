package data

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func NewReferenceValue(v Variable, ctx Context) Value {
	return &ReferenceValue{
		Val: v,
		Ctx: ctx,
	}
}

type ReferenceValue struct {
	Val Variable
	Ctx Context
}

func (s *ReferenceValue) GetValue(ctx Context) (GetValue, Control) {
	return s, nil
}

func (s *ReferenceValue) AsString() string {
	v, _ := s.Val.GetValue(s.Ctx)
	return v.(Value).AsString()
}

func (s *ReferenceValue) Marshal(serializer Serializer) ([]byte, error) {
	v, _ := s.Val.GetValue(s.Ctx)
	return v.(ValueSerializer).Marshal(serializer)
}

func (s *ReferenceValue) Unmarshal(data []byte, serializer Serializer) error {
	v, _ := s.Val.GetValue(s.Ctx)
	return v.(ValueSerializer).Unmarshal(data, serializer)
}

func (s *ReferenceValue) ToGoValue(serializer Serializer) (any, error) {
	return s, nil
}

// Scan 接收 database.sql 包的传值
func (s *ReferenceValue) Scan(value any) error {
	if value == nil {
		s.Ctx.SetVariableValue(s.Val, NewNullValue())
		return nil
	}

	// 检查是否为 sql.Null* 类型
	if isSQLNullType(value) {
		if isNullValue(value) {
			s.Ctx.SetVariableValue(s.Val, NewNullValue())
			return nil
		}
		if actualValue := extractSQLNullValue(value); actualValue != nil {
			return s.Scan(actualValue)
		}
	}

	// 先获取变量的类型
	varType := s.Val.GetType()

	// 如果变量有明确的类型，根据类型进行转换和赋值
	if varType != nil {
		return s.assignByType(varType, value)
	}

	// 如果变量没有设置类型，先尝试获取现有值
	temp, _ := s.Val.GetValue(s.Ctx)
	if temp == nil {
		// 变量不存在且没有类型，根据输入值类型创建新值
		got, err := parseToValue(value)
		if err != nil {
			return fmt.Errorf("failed to parse value: %w", err)
		}
		s.Ctx.SetVariableValue(s.Val, got)
		return nil
	}

	// 变量存在但没有明确类型，根据现有值类型进行转换
	// 需要将 GetValue 转换为 Value
	if val, ok := temp.(Value); ok {
		return s.assignByValueType(val, value)
	}

	// 如果无法转换，尝试直接设置值
	got, err := parseToValue(value)
	if err != nil {
		return fmt.Errorf("failed to parse value: %w", err)
	}
	s.Ctx.SetVariableValue(s.Val, got)
	return nil
}

// assignByType 根据变量类型进行赋值
func (s *ReferenceValue) assignByType(varType Types, value any) error {
	switch varType.(type) {
	case Int:
		return s.assignToIntType(value)
	case Float:
		return s.assignToFloatType(value)
	case String:
		return s.assignToStringType(value)
	case Bool:
		return s.assignToBoolType(value)
	case Arrays:
		return s.assignToArrayType(value)
	case Object:
		return s.assignToObjectType(value)
	default:
		// 未知类型，尝试根据值类型创建
		got, err := parseToValue(value)
		if err != nil {
			return fmt.Errorf("failed to parse value for type %T: %w", varType, err)
		}
		s.Ctx.SetVariableValue(s.Val, got)
		return nil
	}
}

// assignByValueType 根据现有值类型进行转换和赋值
func (s *ReferenceValue) assignByValueType(existingValue Value, value any) error {
	switch v := existingValue.(type) {
	case *IntValue:
		return s.updateIntValue(v, value)
	case *FloatValue:
		return s.updateFloatValue(v, value)
	case *StringValue:
		return s.updateStringValue(v, value)
	case *BoolValue:
		return s.updateBoolValue(v, value)
	default:
		// 尝试直接设置值
		got, err := parseToValue(value)
		if err != nil {
			return fmt.Errorf("failed to parse value: %w", err)
		}
		s.Ctx.SetVariableValue(s.Val, got)
		return nil
	}
}

// assignToIntType 赋值给整数类型变量
func (s *ReferenceValue) assignToIntType(value any) error {
	intVal, err := parseToIntValue(value)
	if err != nil {
		return fmt.Errorf("failed to parse to int: %w", err)
	}
	s.Ctx.SetVariableValue(s.Val, intVal)
	return nil
}

// assignToFloatType 赋值给浮点数类型变量
func (s *ReferenceValue) assignToFloatType(value any) error {
	floatVal, err := parseToFloatValue(value)
	if err != nil {
		return fmt.Errorf("failed to parse to float: %w", err)
	}
	s.Ctx.SetVariableValue(s.Val, floatVal)
	return nil
}

// assignToStringType 赋值给字符串类型变量
func (s *ReferenceValue) assignToStringType(value any) error {
	stringVal, err := parseToStringValue(value)
	if err != nil {
		return fmt.Errorf("failed to parse to string: %w", err)
	}
	s.Ctx.SetVariableValue(s.Val, stringVal)
	return nil
}

// assignToBoolType 赋值给布尔类型变量
func (s *ReferenceValue) assignToBoolType(value any) error {
	boolVal, err := parseToBoolValue(value)
	if err != nil {
		return fmt.Errorf("failed to parse to bool: %w", err)
	}
	s.Ctx.SetVariableValue(s.Val, boolVal)
	return nil
}

// assignToArrayType 赋值给数组类型变量
func (s *ReferenceValue) assignToArrayType(value any) error {
	// 这里可以根据需要实现数组类型的赋值逻辑
	// 暂时使用通用解析
	got, err := parseToValue(value)
	if err != nil {
		return fmt.Errorf("failed to parse value for array type: %w", err)
	}
	s.Ctx.SetVariableValue(s.Val, got)
	return nil
}

// assignToObjectType 赋值给对象类型变量
func (s *ReferenceValue) assignToObjectType(value any) error {
	// 这里可以根据需要实现对象类型的赋值逻辑
	// 暂时使用通用解析
	got, err := parseToValue(value)
	if err != nil {
		return fmt.Errorf("failed to parse value for object type: %w", err)
	}
	s.Ctx.SetVariableValue(s.Val, got)
	return nil
}

// updateIntValue 更新 IntValue
func (s *ReferenceValue) updateIntValue(v *IntValue, value any) error {
	intVal, err := parseToIntValue(value)
	if err != nil {
		return fmt.Errorf("failed to parse to int: %w", err)
	}
	if intVal, ok := intVal.(*IntValue); ok {
		v.Value = intVal.Value
	}
	return nil
}

// updateFloatValue 更新 FloatValue
func (s *ReferenceValue) updateFloatValue(v *FloatValue, value any) error {
	floatVal, err := parseToFloatValue(value)
	if err != nil {
		return fmt.Errorf("failed to parse to float: %w", err)
	}
	if floatVal, ok := floatVal.(*FloatValue); ok {
		v.Value = floatVal.Value
	}
	return nil
}

// updateStringValue 更新 StringValue
func (s *ReferenceValue) updateStringValue(v *StringValue, value any) error {
	stringVal, err := parseToStringValue(value)
	if err != nil {
		return fmt.Errorf("failed to parse to string: %w", err)
	}
	if stringVal, ok := stringVal.(*StringValue); ok {
		v.Value = stringVal.Value
	}
	return nil
}

// updateBoolValue 更新 BoolValue
func (s *ReferenceValue) updateBoolValue(v *BoolValue, value any) error {
	boolVal, err := parseToBoolValue(value)
	if err != nil {
		return fmt.Errorf("failed to parse to bool: %w", err)
	}
	if boolVal, ok := boolVal.(*BoolValue); ok {
		v.Value = boolVal.Value
	}
	return nil
}

// parseToIntValue 将任意值解析为 IntValue
func parseToIntValue(value any) (Value, error) {
	if value == nil {
		return NewIntValue(0), nil
	}

	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		// 统一处理所有整数类型
		ival := reflect.ValueOf(v).Int()
		if ival > int64(math.MaxInt) || ival < int64(math.MinInt) {
			return nil, fmt.Errorf("value %d out of int range", ival)
		}
		return NewIntValue(int(ival)), nil
	case uint, uint8, uint16, uint32, uint64:
		// 统一处理所有无符号整数类型
		uval := reflect.ValueOf(v).Uint()
		if uval > uint64(math.MaxInt) {
			return nil, fmt.Errorf("value %d out of int range", uval)
		}
		return NewIntValue(int(uval)), nil
	case float32, float64:
		// 统一处理浮点数类型
		fval := reflect.ValueOf(v).Float()
		if fval > float64(math.MaxInt) || fval < float64(math.MinInt) {
			return nil, fmt.Errorf("value %f out of int range", fval)
		}
		return NewIntValue(int(fval)), nil
	case string:
		return parseStringToInt(v)
	case []byte:
		return parseStringToInt(string(v))
	case bool:
		if v {
			return NewIntValue(1), nil
		}
		return NewIntValue(0), nil
	case Value:
		return parseValueToInt(v)
	default:
		return parseReflectedValueToInt(value)
	}
}

// parseStringToInt 解析字符串为整数
func parseStringToInt(s string) (Value, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse string '%s' to int: %w", s, err)
	}
	if i > int64(math.MaxInt) || i < int64(math.MinInt) {
		return nil, fmt.Errorf("parsed value %d out of int range", i)
	}
	return NewIntValue(int(i)), nil
}

// parseValueToInt 解析 Value 类型为整数
func parseValueToInt(v Value) (Value, error) {
	if intVal, ok := v.(*IntValue); ok {
		return intVal, nil
	}
	if asInt, ok := v.(AsInt); ok {
		if intVal, err := asInt.AsInt(); err == nil {
			return NewIntValue(intVal), nil
		}
	}
	str := v.AsString()
	return parseStringToInt(str)
}

// parseReflectedValueToInt 使用反射解析值为整数
func parseReflectedValueToInt(value any) (Value, error) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewIntValue(int(rv.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uval := rv.Uint()
		if uval > uint64(math.MaxInt) {
			return nil, fmt.Errorf("reflected uint value %d out of int range", uval)
		}
		return NewIntValue(int(uval)), nil
	case reflect.Float32, reflect.Float64:
		fval := rv.Float()
		if fval > float64(math.MaxInt) || fval < float64(math.MinInt) {
			return nil, fmt.Errorf("reflected float value %f out of int range", fval)
		}
		return NewIntValue(int(fval)), nil
	case reflect.Bool:
		if rv.Bool() {
			return NewIntValue(1), nil
		}
		return NewIntValue(0), nil
	case reflect.String:
		return parseStringToInt(rv.String())
	default:
		return nil, fmt.Errorf("cannot convert %T to int", value)
	}
}

// parseToFloatValue 将任意值解析为 FloatValue
func parseToFloatValue(value any) (Value, error) {
	if value == nil {
		return NewFloatValue(0.0), nil
	}

	switch v := value.(type) {
	case float64:
		return NewFloatValue(v), nil
	case float32:
		return NewFloatValue(float64(v)), nil
	case int, int8, int16, int32, int64:
		return NewFloatValue(float64(reflect.ValueOf(v).Int())), nil
	case uint, uint8, uint16, uint32, uint64:
		return NewFloatValue(float64(reflect.ValueOf(v).Uint())), nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse string '%s' to float: %w", v, err)
		}
		return NewFloatValue(f), nil
	case []byte:
		f, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse []byte '%s' to float: %w", string(v), err)
		}
		return NewFloatValue(f), nil
	case bool:
		if v {
			return NewFloatValue(1.0), nil
		}
		return NewFloatValue(0.0), nil
	case Value:
		return parseValueToFloat(v)
	default:
		return parseReflectedValueToFloat(value)
	}
}

// parseValueToFloat 解析 Value 类型为浮点数
func parseValueToFloat(v Value) (Value, error) {
	if floatVal, ok := v.(*FloatValue); ok {
		return floatVal, nil
	}
	if asFloat, ok := v.(AsFloat); ok {
		if floatVal, err := asFloat.AsFloat(); err == nil {
			return NewFloatValue(floatVal), nil
		}
	}
	str := v.AsString()
	return parseToFloatValue(str)
}

// parseReflectedValueToFloat 使用反射解析值为浮点数
func parseReflectedValueToFloat(value any) (Value, error) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Float32, reflect.Float64:
		return NewFloatValue(rv.Float()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewFloatValue(float64(rv.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return NewFloatValue(float64(rv.Uint())), nil
	case reflect.Bool:
		if rv.Bool() {
			return NewFloatValue(1.0), nil
		}
		return NewFloatValue(0.0), nil
	case reflect.String:
		return parseToFloatValue(rv.String())
	default:
		return nil, fmt.Errorf("cannot convert %T to float", value)
	}
}

// parseToStringValue 将任意值解析为 StringValue
func parseToStringValue(value any) (Value, error) {
	if value == nil {
		return NewStringValue(""), nil
	}

	switch v := value.(type) {
	case string:
		return NewStringValue(v), nil
	case []byte:
		return NewStringValue(string(v)), nil
	case int, int8, int16, int32, int64:
		return NewStringValue(strconv.FormatInt(reflect.ValueOf(v).Int(), 10)), nil
	case uint, uint8, uint16, uint32, uint64:
		return NewStringValue(strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)), nil
	case float32, float64:
		return NewStringValue(strconv.FormatFloat(reflect.ValueOf(v).Float(), 'g', -1, 64)), nil
	case bool:
		return NewStringValue(strconv.FormatBool(v)), nil
	case Value:
		return NewStringValue(v.AsString()), nil
	default:
		return parseReflectedValueToString(value)
	}
}

// parseReflectedValueToString 使用反射解析值为字符串
func parseReflectedValueToString(value any) (Value, error) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.String:
		return NewStringValue(rv.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewStringValue(strconv.FormatInt(rv.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return NewStringValue(strconv.FormatUint(rv.Uint(), 10)), nil
	case reflect.Float32, reflect.Float64:
		return NewStringValue(strconv.FormatFloat(rv.Float(), 'g', -1, 64)), nil
	case reflect.Bool:
		return NewStringValue(strconv.FormatBool(rv.Bool())), nil
	case reflect.Slice, reflect.Array:
		return NewStringValue(fmt.Sprintf("%v", value)), nil
	case reflect.Map:
		return NewStringValue(fmt.Sprintf("%v", value)), nil
	case reflect.Ptr:
		if rv.IsNil() {
			return NewStringValue(""), nil
		}
		return parseToStringValue(rv.Elem().Interface())
	default:
		return NewStringValue(fmt.Sprintf("%v", value)), nil
	}
}

// parseToBoolValue 将任意值解析为 BoolValue
func parseToBoolValue(value any) (Value, error) {
	if value == nil {
		return NewBoolValue(false), nil
	}

	switch v := value.(type) {
	case bool:
		return NewBoolValue(v), nil
	case int, int8, int16, int32, int64:
		return NewBoolValue(reflect.ValueOf(v).Int() != 0), nil
	case uint, uint8, uint16, uint32, uint64:
		return NewBoolValue(reflect.ValueOf(v).Uint() != 0), nil
	case float32, float64:
		return NewBoolValue(reflect.ValueOf(v).Float() != 0.0), nil
	case string:
		return parseStringToBool(v)
	case []byte:
		return parseStringToBool(string(v))
	case Value:
		return parseValueToBool(v)
	default:
		return parseReflectedValueToBool(value)
	}
}

// parseStringToBool 解析字符串为布尔值
func parseStringToBool(s string) (Value, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return NewBoolValue(s != ""), nil
	}
	return NewBoolValue(b), nil
}

// parseValueToBool 解析 Value 类型为布尔值
func parseValueToBool(v Value) (Value, error) {
	if boolVal, ok := v.(*BoolValue); ok {
		return boolVal, nil
	}
	if asBool, ok := v.(AsBool); ok {
		if boolVal, err := asBool.AsBool(); err == nil {
			return NewBoolValue(boolVal), nil
		}
	}
	str := v.AsString()
	return parseStringToBool(str)
}

// parseReflectedValueToBool 使用反射解析值为布尔值
func parseReflectedValueToBool(value any) (Value, error) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Bool:
		return NewBoolValue(rv.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewBoolValue(rv.Int() != 0), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return NewBoolValue(rv.Uint() != 0), nil
	case reflect.Float32, reflect.Float64:
		return NewBoolValue(rv.Float() != 0.0), nil
	case reflect.String:
		return NewBoolValue(rv.String() != ""), nil
	case reflect.Slice, reflect.Array:
		return NewBoolValue(rv.Len() > 0), nil
	case reflect.Map:
		return NewBoolValue(rv.Len() > 0), nil
	case reflect.Ptr:
		return NewBoolValue(!rv.IsNil()), nil
	default:
		return NewBoolValue(true), nil
	}
}

// parseToValue 根据值的类型自动选择合适的 Value 类型
func parseToValue(value any) (Value, error) {
	if value == nil {
		return NewNullValue(), nil
	}

	switch value.(type) {
	case bool:
		return parseToBoolValue(value)
	case int, int8, int16, int32, int64:
		return parseToIntValue(value)
	case uint, uint8, uint16, uint32, uint64:
		return parseToIntValue(value)
	case float32, float64:
		return parseToFloatValue(value)
	case string:
		return parseToStringValue(value)
	case []byte:
		return parseToStringValue(value)
	case Value:
		return value.(Value), nil
	default:
		return parseReflectedValueToValue(value)
	}
}

// parseReflectedValueToValue 使用反射解析值为 Value
func parseReflectedValueToValue(value any) (Value, error) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Bool:
		return parseToBoolValue(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return parseToIntValue(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return parseToIntValue(value)
	case reflect.Float32, reflect.Float64:
		return parseToFloatValue(value)
	case reflect.String:
		return parseToStringValue(value)
	case reflect.Slice, reflect.Array:
		return NewStringValue(fmt.Sprintf("%v", value)), nil
	case reflect.Map:
		return NewStringValue(fmt.Sprintf("%v", value)), nil
	case reflect.Ptr:
		if rv.IsNil() {
			return NewNullValue(), nil
		}
		return parseToValue(rv.Elem().Interface())
	default:
		return NewStringValue(fmt.Sprintf("%v", value)), nil
	}
}

// isNumeric 检查字符串是否为数字
func isNumeric(s string) bool {
	if s == "" {
		return false
	}

	// 跳过前导空格
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}

	// 检查是否有符号
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
		if s == "" {
			return false
		}
	}

	// 检查是否为纯数字
	hasDigit := false
	hasDot := false

	for _, c := range s {
		if c >= '0' && c <= '9' {
			hasDigit = true
		} else if c == '.' && !hasDot {
			hasDot = true
		} else {
			return false
		}
	}

	return hasDigit
}

// isSQLNullType 检查是否为 sql.Null* 类型
func isSQLNullType(value any) bool {
	if value == nil {
		return false
	}

	// 使用反射检查类型名称
	rt := reflect.TypeOf(value)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	typeName := rt.String()
	return strings.Contains(typeName, "sql.Null")
}

// isNullValue 检查 sql.Null* 类型是否为空值
func isNullValue(value any) bool {
	if value == nil {
		return true
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// 检查 Valid 字段
	if validField := rv.FieldByName("Valid"); validField.IsValid() {
		return !validField.Bool()
	}

	return false
}

// extractSQLNullValue 从 sql.Null* 类型中提取实际值
func extractSQLNullValue(value any) any {
	if value == nil {
		return nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// 尝试获取各种可能的字段
	fields := []string{"String", "Int64", "Float64", "Bool", "Time"}
	for _, fieldName := range fields {
		if field := rv.FieldByName(fieldName); field.IsValid() {
			return field.Interface()
		}
	}

	return nil
}
