package node

import (
	"errors"
	"fmt"
	"os"

	"github.com/php-any/origami/data"
)

// emitUndefinedArrayKeyWarning 输出 PHP 8+ 未定义数组键 Warning
func emitUndefinedArrayKeyWarning(from data.From, key string, intKey bool) {
	file := "Unknown"
	line := 0
	if from != nil {
		if src := from.GetSource(); src != "" {
			file = src
		}
		if sl, _ := from.GetStartPosition(); sl >= 0 {
			line = sl + 1
		}
	}
	if intKey {
		fmt.Fprintf(os.Stderr, "\nWarning: Undefined array key %s in %s on line %d\n", key, file, line)
	} else {
		fmt.Fprintf(os.Stderr, "\nWarning: Undefined array key \"%s\" in %s on line %d\n", key, file, line)
	}
}

// emitNullOffsetDeprecation prints a PHP 8.1 deprecation when null is used as an array offset.
func emitNullOffsetDeprecation(from data.From) {
	file := "Unknown"
	line := 0
	if from != nil {
		if src := from.GetSource(); src != "" {
			file = src
		}
		if sl, _ := from.GetStartPosition(); sl >= 0 {
			line = sl + 1
		}
	}
	fmt.Printf("Deprecated: Using null as an array offset is deprecated, use an empty string instead in %s on line %d\n", file, line)
}

// callArrayAccessOffsetExists 调用 ArrayAccess::offsetExists
func callArrayAccessOffsetExists(ctx data.Context, classValue *data.ClassValue, index data.Value) (bool, data.Control) {
	method, exists := classValue.GetMethod("offsetExists")
	if !exists {
		return false, nil
	}
	fnCtx := classValue.CreateContext(method.GetVariables())
	if len(method.GetVariables()) > 0 {
		fnCtx.SetVariableValue(method.GetVariables()[0], index)
	}
	ret, ctl := method.Call(fnCtx)
	if ctl != nil {
		return false, ctl
	}
	if bv, ok := ret.(*data.BoolValue); ok {
		return bv.Value, nil
	}
	return false, nil
}

// EmptyViaArrayAccess 对 empty($obj[$k]) 按 PHP 语义走 offsetExists/offsetGet
func EmptyViaArrayAccess(ctx data.Context, ie *IndexExpression) (bool, bool) {
	array, acl := ie.Array.GetValue(ctx)
	if acl != nil {
		return true, true
	}
	index, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return true, true
	}
	iv, ok := index.(data.Value)
	if !ok {
		return false, false
	}
	var obj *data.ClassValue
	switch v := array.(type) {
	case *data.ClassValue:
		obj = v
	case *data.ThisValue:
		obj = v.ClassValue
	default:
		return false, false
	}
	if !checkArrayAccess(ctx, obj.Class) {
		return false, false
	}
	exists, ctl := callArrayAccessOffsetExists(ctx, obj, iv)
	if ctl != nil {
		return true, true
	}
	if !exists {
		return true, true
	}
	val, ctl := callArrayAccessOffsetGet(ctx, obj, index)
	if ctl != nil {
		return true, true
	}
	return isEmptyPHPValue(val), true
}

// EmptyIndexExpression 计算 empty($a[$k]) / empty($a[$k]['x'])（未定义键视为 empty，不触发 Warning）
func EmptyIndexExpression(ctx data.Context, ie *IndexExpression) (isEmpty bool, handled bool) {
	if inner, ok := ie.Array.(*IndexExpression); ok {
		exists, handled := indexExpressionKeyExists(ctx, inner)
		if !handled {
			return true, false
		}
		if !exists {
			return true, true
		}
		parentVal, ok := indexExpressionReadNoWarn(ctx, inner)
		if !ok {
			return true, true
		}
		index, acl := ie.Index.GetValue(ctx)
		if acl != nil {
			return true, true
		}
		return emptyOnContainer(parentVal, index), true
	}
	exists, handled := indexExpressionKeyExists(ctx, ie)
	if !handled {
		return true, false
	}
	if !exists {
		return true, true
	}
	container, acl := indexExpressionContainer(ctx, ie.Array)
	if acl != nil {
		return true, true
	}
	index, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return true, true
	}
	return emptyOnContainer(container, index), true
}

func emptyOnContainer(container data.GetValue, index data.GetValue) bool {
	val, ok := readIndexNoWarn(container, index)
	if !ok {
		return true
	}
	return isEmptyPHPValue(val)
}

func isEmptyPHPValue(v data.GetValue) bool {
	if v == nil {
		return true
	}
	if _, ok := v.(*data.NullValue); ok {
		return true
	}
	if s, ok := v.(data.AsString); ok {
		str := s.AsString()
		return str == "" || str == "0"
	}
	if i, ok := v.(data.AsInt); ok {
		n, _ := i.AsInt()
		return n == 0
	}
	if f, ok := v.(*data.FloatValue); ok {
		return f.Value == 0
	}
	if b, ok := v.(*data.BoolValue); ok {
		return !b.Value
	}
	if a, ok := v.(*data.ArrayValue); ok {
		return len(a.List) == 0
	}
	return false
}

// callArrayAccessOffsetGet 调用 ArrayAccess 接口的 offsetGet 方法
func callArrayAccessOffsetGet(ctx data.Context, classValue *data.ClassValue, index data.GetValue) (data.GetValue, data.Control) {
	method, exists := classValue.GetMethod("offsetGet")
	if !exists {
		return nil, nil
	}

	// 创建参数上下文
	fnCtx := classValue.CreateContext(method.GetVariables())

	// 设置参数值
	params := method.GetParams()
	if len(params) >= 1 {
		if p0, ok := params[0].(interface {
			SetValue(data.Context, data.Value) data.Control
		}); ok {
			if ctl := p0.SetValue(fnCtx, index.(data.Value)); ctl != nil {
				return nil, ctl
			}
		}
	}

	// 调用方法并获取返回值
	result, ctl := method.Call(fnCtx)
	if ctl != nil {
		return nil, ctl
	}

	if result != nil {
		if v, ok := result.(data.Value); ok {
			return tagArrayAccessOffsetCopy(v, classValue.Class.GetName()), nil
		}
		return result, nil
	}
	return nil, nil
}

func tagArrayAccessOffsetCopy(val data.Value, className string) data.Value {
	if val == nil || className == "" {
		return val
	}
	switch v := val.(type) {
	case *data.ArrayValue:
		if v.IndirectOverloadClass != "" {
			return v
		}
		tagged := data.CloneArrayValue(v)
		tagged.IndirectOverloadClass = className
		return tagged
	case *data.ObjectValue:
		if v.IndirectOverloadClass != "" {
			return v
		}
		tagged := data.CloneObjectValue(v)
		tagged.IndirectOverloadClass = className
		return tagged
	default:
		return val
	}
}

// arrayAccessOverloadedClass 若索引表达式根对象为 ArrayAccess 则返回类名
func arrayAccessOverloadedClass(ctx data.Context, ie *IndexExpression) (string, bool) {
	array, acl := ie.Array.GetValue(ctx)
	if acl != nil {
		return "", false
	}
	switch v := array.(type) {
	case *data.ClassValue:
		if checkArrayAccess(ctx, v.Class) {
			return v.Class.GetName(), true
		}
	case *data.ThisValue:
		if v.ClassValue != nil && checkArrayAccess(ctx, v.Class) {
			return v.Class.GetName(), true
		}
	}
	return "", false
}

func indirectOverloadTag(v data.GetValue) string {
	if arr, ok := v.(*data.ArrayValue); ok {
		return arr.IndirectOverloadClass
	}
	if obj, ok := v.(*data.ObjectValue); ok {
		return obj.IndirectOverloadClass
	}
	return ""
}

// blockIndirectOverloadAssign 对 ArrayAccess offsetGet 返回的副本做嵌套赋值时发出 Notice 并阻止写入
func blockIndirectOverloadAssign(ie *IndexExpression, arrayVal data.GetValue) bool {
	if _, ok := ie.Array.(data.Variable); ok {
		return false
	}
	if _, ok := ie.Array.(*IndexExpression); !ok {
		return false
	}
	if className := indirectOverloadTag(arrayVal); className != "" {
		emitIndirectModificationNoticeAssign(ie.GetFrom(), className)
		return true
	}
	return false
}

func emitIndirectModificationNotice(from data.From, className string) {
	emitIndirectModificationNoticeLine(from, className, 0)
}

// emitIndirectModificationNoticeAssign 用于 .= / = 等赋值左侧（行号比 PostfixIncr 多 1）
func emitIndirectModificationNoticeAssign(from data.From, className string) {
	emitIndirectModificationNoticeLine(from, className, 1)
}

func emitIndirectModificationNoticeLine(from data.From, className string, lineAdjust int) {
	file := "Unknown"
	line := 0
	if from != nil {
		if src := from.GetSource(); src != "" {
			file = src
		}
		if sl, _ := from.GetStartPosition(); sl >= 0 {
			line = sl + lineAdjust
		}
	}
	fmt.Fprintf(os.Stderr, "\nNotice: Indirect modification of overloaded element of %s has no effect in %s on line %d\n", className, file, line)
}

// CallArrayAccessOffsetUnset 调用 ArrayAccess::offsetUnset
func CallArrayAccessOffsetUnset(ctx data.Context, classValue *data.ClassValue, index data.Value) data.Control {
	return callArrayAccessOffsetUnset(ctx, classValue, index)
}

func callArrayAccessOffsetUnset(ctx data.Context, classValue *data.ClassValue, index data.Value) data.Control {
	method, exists := classValue.GetMethod("offsetUnset")
	if !exists {
		return nil
	}
	fnCtx := classValue.CreateContext(method.GetVariables())
	if len(method.GetVariables()) > 0 {
		fnCtx.SetVariableValue(method.GetVariables()[0], index)
	}
	_, ctl := method.Call(fnCtx)
	return ctl
}

// callArrayAccessOffsetSet 调用 ArrayAccess 接口的 offsetSet 方法
func callArrayAccessOffsetSet(ctx data.Context, classValue *data.ClassValue, index data.GetValue, value data.Value) data.Control {
	method, exists := classValue.GetMethod("offsetSet")
	if !exists {
		return nil
	}

	// 创建参数上下文
	fnCtx := classValue.CreateContext(method.GetVariables())

	// 设置参数值
	params := method.GetParams()
	if len(params) >= 2 {
		// 设置第一个参数 (offset)
		if p0, ok := params[0].(interface {
			SetValue(data.Context, data.Value) data.Control
		}); ok {
			if ctl := p0.SetValue(fnCtx, index.(data.Value)); ctl != nil {
				return ctl
			}
		}
		// 设置第二个参数 (value)
		if p1, ok := params[1].(interface {
			SetValue(data.Context, data.Value) data.Control
		}); ok {
			if ctl := p1.SetValue(fnCtx, value); ctl != nil {
				return ctl
			}
		}
	}

	// 调用方法
	_, ctl := method.Call(fnCtx)
	return ctl
}

// CheckArrayAccess 检查类是否实现了 ArrayAccess 接口
func CheckArrayAccess(ctx data.Context, classStmt data.ClassStmt) bool {
	return checkArrayAccess(ctx, classStmt)
}

func checkArrayAccess(ctx data.Context, classStmt data.ClassStmt) bool {
	if asBool, ctl := instanceof(ctx, "ArrayAccess", &data.ClassValue{Class: classStmt}); ctl == nil {
		if boolVal, ok := asBool.(*data.BoolValue); ok {
			return boolVal.Value
		}
	}
	return false
}

// IndexExpression 表示数组访问表达式
type IndexExpression struct {
	*Node `pp:"-"`
	Array data.GetValue // 数组表达式
	Index data.GetValue // 索引表达式
}

// NewIndexExpression 创建一个新的数组访问表达式
func NewIndexExpression(token *TokenFrom, array data.GetValue, index data.GetValue) *IndexExpression {
	return &IndexExpression{
		Node:  NewNode(token),
		Array: array,
		Index: index,
	}
}

func (ie *IndexExpression) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	temp, acl := ie.Array.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	index, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}

	switch v := temp.(type) {
	case *data.ArrayValue:
		i := 0
		switch iv := index.(type) {
		case *data.IntValue:
			var err error
			i, err = iv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
			if i >= len(v.List) {
				return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("数组索引超出范围"), "UndefinedIndexExpression")
			}
		case *data.StringValue:
			if len(v.List) == 0 {
				return data.NewZVal(data.NewNullValue()), nil
			}
			// 通过 ZVal.Name 查找字符串键
			key := iv.AsString()
			for _, zval := range v.List {
				if zval.Name == key {
					return zval, nil
				}
			}
			// 未找到，返回 null（PHP 行为）
			return data.NewZVal(data.NewNullValue()), nil
		case data.AsString:
			if len(v.List) == 0 {
				return data.NewZVal(data.NewNullValue()), nil
			}
			key := iv.AsString()
			for _, zval := range v.List {
				if zval.Name == key {
					return zval, nil
				}
			}
			return data.NewZVal(data.NewNullValue()), nil
		case data.AsInt:
			var err error
			i, err = iv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
			if i >= len(v.List) {
				return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("数组索引超出范围"), "UndefinedIndexExpression")
			}
		case *data.BoolValue:
			if iv.Value {
				i = 1
			}
			if i >= len(v.List) {
				return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("数组索引超出范围"), "UndefinedIndexExpression")
			}
		default:
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("无法处理索引的类型值"))
		}

		return v.List[i], nil
	case *data.ObjectValue:
		switch iv := index.(type) {
		case data.AsString:
			ov, acl := v.GetZVal(iv.AsString())
			if acl != nil {
				return nil, acl
			}
			if ov == nil {
				return data.NewZVal(data.NewNullValue()), nil
			}
			return ov, nil
		}
	case *data.NullValue:
		return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("NULL值不支持数组索引操作"), "UndefinedIndexExpression")
	}
	return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("无法处理索引的类型值"), "UndefinedIndexExpression")
}

func (ie *IndexExpression) SetValue(ctx data.Context, value data.Value) data.Control {
	indexVal, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return acl
	}

	var arrayVal data.GetValue
	var arrayAcl data.Control
	// $parent[$key][$sub] = $val：父键不存在时 PHP 自动 vivify 子数组，且不触发未定义键 Warning
	if inner, ok := ie.Array.(*IndexExpression); ok {
		exists, handled := indexExpressionKeyExists(ctx, inner)
		if handled && !exists {
			emptyArr := data.NewArrayValue(nil).(*data.ArrayValue)
			if ctl := inner.SetValue(ctx, emptyArr); ctl != nil {
				return ctl
			}
		}
		if v, ok := indexExpressionReadNoWarn(ctx, inner); ok {
			arrayVal = v
		} else {
			arrayVal = data.NewNullValue()
		}
	}
	if arrayVal == nil {
		arrayVal, arrayAcl = ie.Array.GetValue(ctx)
	}
	if arrayAcl != nil {
		// 在赋值语境下，未定义索引（UndefinedIndexExpression）应视为 null，
		// 以便支持 PHP 风格的链式自动创建：
		// $namespace['commands'][1] = 'foo';
		if tv, ok := acl.(*data.ThrowValue); ok && tv.Name == "UndefinedIndexExpression" {
			arrayVal = data.NewNullValue()
		} else {
			return acl
		}
	}

	// 如果数组当前为 null，根据左侧表达式类型决定如何自动创建容器
	if _, ok := arrayVal.(*data.NullValue); ok {
		switch base := ie.Array.(type) {
		case data.Variable:
			// 场景：$namespace['commands'] = ... 或 $namespace['commands'][1] = ...
			obj := data.NewObjectValue()
			if ctl := base.SetValue(ctx, obj); ctl != nil {
				return ctl
			}
			// 变量赋值时会对 ObjectValue 做 clone，这里需要重新从上下文读取，
			// 确保后续写入操作作用在真实存储的容器上。
			var errCtl data.Control
			arrayVal, errCtl = ie.Array.GetValue(ctx)
			if errCtl != nil {
				return errCtl
			}
		case *IndexExpression:
			// 多级访问，自动创建空数组并挂到上一层（PHP：$arr[$k][$k2]= 会 vivify 子数组）
			// $this->routes[$method][$uri] = $route 等
			subArr := data.NewArrayValue(nil).(*data.ArrayValue)
			if _, errCtl := NewBinaryAssign(ie.GetFrom(), base, subArr).GetValue(ctx); errCtl != nil {
				return errCtl
			}
			// 同样需要重新读取当前数组值，拿到实际挂在上一层上的容器（可能经历了 clone）
			var errCtl2 data.Control
			arrayVal, errCtl2 = ie.Array.GetValue(ctx)
			if errCtl2 != nil {
				return errCtl2
			}
		default:
			// 其它情况保持 null，后续类型分支会给出统一错误信息
		}
	}

	if blockIndirectOverloadAssign(ie, arrayVal) {
		return nil
	}

	switch arr := arrayVal.(type) {
	case *data.NullValue:
		// $null[] = value → 初始化为数组后追加（PHP 语义）
		newArr := data.NewArrayValue([]data.Value{value}).(*data.ArrayValue)
		_, acl = NewBinaryAssign(ie.GetFrom(), ie.Array, newArr).GetValue(ctx)
		return acl
	case *data.ArrayValue:
		// 数组索引赋值
		// Handle null index first (before interface type assertions)
		if _, isNull := indexVal.(*data.NullValue); isNull {
			// $arr[] = $val（PHP 追加，不触发 null 下标弃用提示）
			arr.List = append(arr.List, data.NewZVal(value))
			writeBackArrayProperty(ctx, ie.Array, arr)
			return nil
		}
		i := 0
		if iv, ok := indexVal.(data.AsInt); ok {
			var err error
			i, err = iv.AsInt()
			if err != nil {
				return data.NewErrorThrow(ie.GetFrom(), err)
			}
		} else if iv, ok := indexVal.(data.AsString); ok {
			// 字符串键：查找匹配 Name 的项并更新，找不到则追加
			key := iv.AsString()
			for _, zval := range arr.List {
				if zval != nil && zval.Name == key {
					zval.Value = value
					writeBackArrayProperty(ctx, ie.Array, arr)
					return nil
				}
			}
			// 未找到，追加新项
			arr.List = append(arr.List, &data.ZVal{Name: key, Value: value})
			writeBackArrayProperty(ctx, ie.Array, arr)
			return nil
		} else {
			return data.NewErrorThrow(ie.GetFrom(), errors.New("数组索引不是整数类型"))
		}

		if i < 0 {
			return data.NewErrorThrow(ie.GetFrom(), errors.New("数组索引不能为负数"))
		}

		// PHP 数组是稀疏的：arr[22]=val 不填充 0-21。超出长度时仅当推入末尾才扩容，否则转为 ObjectValue
		arr.SetIntKey(i, value)
		writeBackArrayProperty(ctx, ie.Array, arr)
		return nil

	case *data.ClassValue:
		// 检查是否实现了 ArrayAccess 接口
		if checkArrayAccess(ctx, arr.Class) {
			// 实现了 ArrayAccess，调用 offsetSet 方法
			return callArrayAccessOffsetSet(ctx, arr, indexVal, value)
		}
		// 类实例属性赋值
		if iv, ok := indexVal.(data.AsString); ok {
			name := iv.AsString()
			if prop, ok := arr.GetPropertyStmt(name); ok {
				if prop.GetModifier() != data.ModifierPublic {
					return data.NewErrorThrow(ie.GetFrom(), errors.New("对象属性不是公开的"))
				}
				// 使用 SetProperty 方法设置属性值
				return arr.SetProperty(name, value)
			}
			return data.NewErrorThrow(ie.GetFrom(), errors.New("对象不存在指定属性"))
		}
		return data.NewErrorThrow(ie.GetFrom(), errors.New("ClassValue 索引必须是字符串"))

	case *data.ThisValue:
		// 检查是否实现了 ArrayAccess 接口
		if checkArrayAccess(ctx, arr.Class) {
			// 实现了 ArrayAccess，调用 offsetSet 方法
			return callArrayAccessOffsetSet(ctx, arr.ClassValue, indexVal, value)
		}
		// $this[$name] 动态设置当前对象属性
		if iv, ok := indexVal.(data.AsString); ok {
			name := iv.AsString()
			if prop, ok := arr.Class.GetProperty(name); ok {
				if prop.GetModifier() != data.ModifierPublic {
					return data.NewErrorThrow(ie.GetFrom(), errors.New("对象属性不是公开的"))
				}
				// ThisValue 包含 ClassValue，直接使用 SetProperty
				return arr.ClassValue.SetProperty(name, value)
			}
			return data.NewErrorThrow(ie.GetFrom(), errors.New("对象不存在指定属性"))
		}
		return data.NewErrorThrow(ie.GetFrom(), errors.New("ThisValue 索引必须是字符串"))

	case data.SetProperty:
		// 对象属性赋值
		if _, isNull := indexVal.(*data.NullValue); isNull {
			emitNullOffsetDeprecation(ie.GetFrom())
			arr.SetProperty("", value)
			return nil
		}
		if iv, ok := indexVal.(data.AsString); ok {
			arr.SetProperty(iv.AsString(), value)
			return nil
		} else if iv, ok := indexVal.(data.AsInt); ok {
			// 整数索引转换为字符串
			if i, err := iv.AsInt(); err == nil {
				arr.SetProperty(fmt.Sprintf("%d", i), value)
				return nil
			}
		}
		return data.NewErrorThrow(ie.GetFrom(), errors.New("对象索引必须是字符串或整数"))

	case *data.IndexReferenceValue:
		sub := &IndexExpression{
			Node:  ie.Node,
			Array: arr.Expr,
			Index: ie.Index,
		}
		return sub.SetValue(arr.Ctx, value)

	case *data.ReferenceValue:
		sub := &IndexExpression{
			Node:  ie.Node,
			Array: arr.Val,
			Index: ie.Index,
		}
		return sub.SetValue(arr.Ctx, value)

	default:
		return data.NewErrorThrow(ie.GetFrom(), errors.New("无法设置索引表达式的值"))
	}
}

// GetValue 获取数组访问表达式的值
func (ie *IndexExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	temp, acl := ie.Array.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	index, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}

	switch v := temp.(type) {
	case *data.ArrayValue:
		i := 0
		switch iv := index.(type) {
		case *data.NullValue:
			emitNullOffsetDeprecation(ie.GetFrom())
			// null is treated as empty string
			for _, zval := range v.List {
				if zval != nil && zval.Name == "" {
					return zval.Value, nil
				}
			}
			return data.NewNullValue(), nil
		case *data.IntValue:
			var err error
			i, err = iv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
			if val, ok := arrayIntKeyValue(v, i); ok {
				return val, nil
			}
			emitUndefinedArrayKeyWarning(ie.GetFrom(), fmt.Sprintf("%d", i), true)
			return data.NewNullValue(), nil
		case *data.StringValue:
			// 字符串键：在数组中搜索匹配的 Name
			key := iv.Value
			for _, zval := range v.List {
				if zval != nil && zval.Name == key {
					return zval.Value, nil
				}
			}
			emitUndefinedArrayKeyWarning(ie.GetFrom(), key, false)
			return data.NewNullValue(), nil
		case data.AsInt:
			var err error
			i, err = iv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
			if val, ok := arrayIntKeyValue(v, i); ok {
				return val, nil
			}
			emitUndefinedArrayKeyWarning(ie.GetFrom(), fmt.Sprintf("%d", i), true)
			return data.NewNullValue(), nil
		case *data.BoolValue:
			if iv.Value {
				i = 1
			}
			if val, ok := arrayIntKeyValue(v, i); ok {
				return val, nil
			}
			emitUndefinedArrayKeyWarning(ie.GetFrom(), fmt.Sprintf("%d", i), true)
			return data.NewNullValue(), nil
		default:
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("无法处理索引的类型值"))
		}
	case *data.ObjectValue:
		// 支持整数索引（转换为字符串）和字符串索引
		var key string
		if _, isNull := index.(*data.NullValue); isNull {
			emitNullOffsetDeprecation(ie.GetFrom())
			key = ""
		} else if iv, ok := index.(data.AsString); ok {
			key = iv.AsString()
		} else if iv, ok := index.(data.AsInt); ok {
			// 将整数索引转换为字符串
			if i, err := iv.AsInt(); err == nil {
				key = fmt.Sprintf("%d", i)
			} else {
				return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("无法处理索引的类型值"))
			}
		} else {
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("ObjectValue无法处理索引的类型值"))
		}
		ov, acl := v.GetProperty(key)
		if acl != nil {
			return nil, acl
		}
		return ov, nil
	case *data.ClassValue:
		// 检查是否实现了 ArrayAccess 接口
		if checkArrayAccess(ctx, v.Class) {
			// 实现了 ArrayAccess，调用 offsetGet 方法
			return callArrayAccessOffsetGet(ctx, v, index)
		}
		// 支持对类实例通过字符串索引访问公开属性：
		// $obj[$name]，在动态属性语法 $obj->$name 降级为索引访问后会走到这里
		if iv, ok := index.(data.AsString); ok {
			name := iv.AsString()
			if prop, ok := v.GetPropertyStmt(name); ok {
				if prop.GetModifier() != data.ModifierPublic {
					return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("对象属性不是公开的"))
				}
				return prop.GetValue(v)
			} else if prop, acl := v.ObjectValue.GetProperty(name); acl == nil {
				switch prop := prop.(type) {
				case *data.NullValue:
					return nil, data.NewErrorThrow(ie.GetFrom(), fmt.Errorf("%s对象不存在指定属性：%s", v.Class.GetName(), name))
				default:
					return prop.GetValue(ctx)
				}
			}
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("对象不存在指定属性"))
		}
		return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("ClassValue 无法处理索引的类型值"))
	case *data.ThisValue:
		// 检查是否实现了 ArrayAccess 接口
		if checkArrayAccess(ctx, v.Class) {
			// 实现了 ArrayAccess，调用 offsetGet 方法
			return callArrayAccessOffsetGet(ctx, v.ClassValue, index)
		}
		// $this[$name] 动态访问当前对象属性
		if iv, ok := index.(data.AsString); ok {
			name := iv.AsString()
			if prop, ok := v.Class.GetProperty(name); ok {
				if prop.GetModifier() != data.ModifierPublic {
					return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("对象属性不是公开的"))
				}
				return prop.GetValue(ctx)
			} else if prop, acl := v.ObjectValue.GetProperty(name); acl == nil {
				switch prop := prop.(type) {
				case *data.NullValue:
					return nil, data.NewErrorThrow(ie.GetFrom(), fmt.Errorf("this(%s) 对象不存在指定属性：%s", v.Class.GetName(), name))
				default:
					return prop.GetValue(ctx)
				}
			}
		}
		return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("ThisValue 无法处理索引的类型值"))
	case *data.StringValue:
		// 获取字符串指定位置符号
		if iv, ok := index.(data.AsInt); ok {
			var err error
			i, err := iv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(ie.GetFrom(), err)
			}
			if i < 0 {
				i = len(v.Value) + i
			}
			if i < 0 || i >= len(v.Value) {
				return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("字符串索引超出范围"))
			}
			return data.NewStringValue(string(v.Value[i])), nil
		} else {
			return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("字符串无法处理非int值"))
		}
	case *data.NullValue:
		return data.NewNullValue(), nil
	case *data.BoolValue:
		return data.NewNullValue(), nil
	}
	return nil, data.NewErrorThrowByName(ie.GetFrom(), errors.New("无法处理索引的类型值"), "UndefinedIndexExpression")
}

// assignIndexConcat 处理 $obj[$i][$k] .= $rhs，只获取一次容器，避免 ArrayAccess 重复 offsetGet
func assignIndexConcat(ctx data.Context, ie *IndexExpression, dot *BinaryDot) (data.GetValue, data.Control) {
	rv, rCtl := dot.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}
	if rv == nil {
		rv = data.NewNullValue()
	}

	arrayVal, acl := ie.Array.GetValue(ctx)
	if acl != nil {
		if tv, ok := acl.(*data.ThrowValue); ok && tv.Name == "UndefinedIndexExpression" {
			arrayVal = data.NewNullValue()
		} else {
			return nil, acl
		}
	}

	if blockIndirectOverloadAssign(ie, arrayVal) {
		return indexValueOnContainer(ctx, arrayVal, ie.Index, ie.GetFrom())
	}

	indexVal, acl := ie.Index.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}

	lv, acl := indexValueOnContainer(ctx, arrayVal, indexVal, ie.GetFrom())
	if acl != nil {
		return nil, acl
	}

	newVal := concatPHPValues(lv, rv)
	v, ok := newVal.(data.Value)
	if !ok {
		return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("concat assign failed"))
	}
	if ctl := indexSetValueOnContainer(ctx, ie, arrayVal, indexVal, v); ctl != nil {
		return nil, ctl
	}
	return v, nil
}

func indexValueOnContainer(ctx data.Context, container data.GetValue, indexExpr data.GetValue, from data.From) (data.GetValue, data.Control) {
	index, acl := indexExpr.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	switch v := container.(type) {
	case *data.ArrayValue:
		tmp := &IndexExpression{Node: NewNode(from), Index: indexExpr}
		return tmp.readArrayIndex(ctx, v, index)
	case *data.ObjectValue:
		key, ok := indexKeyString(index)
		if !ok {
			return nil, data.NewErrorThrow(from, errors.New("ObjectValue无法处理索引的类型值"))
		}
		val, acl := v.GetProperty(key)
		if acl != nil {
			return nil, acl
		}
		return val, nil
	case *data.ClassValue:
		if checkArrayAccess(ctx, v.Class) {
			return callArrayAccessOffsetGet(ctx, v, index)
		}
	case *data.ThisValue:
		if checkArrayAccess(ctx, v.Class) {
			return callArrayAccessOffsetGet(ctx, v.ClassValue, index)
		}
	case data.GetProperty:
		key, ok := indexKeyString(index)
		if !ok {
			return nil, data.NewErrorThrow(from, errors.New("对象索引必须是字符串或整数"))
		}
		val, acl := v.GetProperty(key)
		if acl != nil {
			return nil, acl
		}
		return val, nil
	}
	return nil, data.NewErrorThrow(from, errors.New("无法读取容器索引值"))
}

func (ie *IndexExpression) readArrayIndex(ctx data.Context, arr *data.ArrayValue, index data.GetValue) (data.GetValue, data.Control) {
	switch iv := index.(type) {
	case *data.StringValue:
		key := iv.Value
		for _, zval := range arr.List {
			if zval != nil && zval.Name == key {
				return zval.Value, nil
			}
		}
		emitUndefinedArrayKeyWarning(ie.GetFrom(), key, false)
		return data.NewNullValue(), nil
	case data.AsInt:
		i, err := iv.AsInt()
		if err != nil {
			return nil, data.NewErrorThrow(ie.GetFrom(), err)
		}
		if val, ok := arrayIntKeyValue(arr, i); ok {
			return val, nil
		}
		emitUndefinedArrayKeyWarning(ie.GetFrom(), fmt.Sprintf("%d", i), true)
		return data.NewNullValue(), nil
	default:
		return nil, data.NewErrorThrow(ie.GetFrom(), errors.New("无法处理索引的类型值"))
	}
}

// arrayIntKeyValue 按键读取数组元素；键存在时返回值（可为 null），不存在时 ok=false
func arrayIntKeyValue(arr *data.ArrayValue, i int) (data.GetValue, bool) {
	z, _ := arr.FindSlotByIntKey(i)
	if z == nil {
		return nil, false
	}
	if z.Value == nil {
		return data.NewNullValue(), true
	}
	return z.Value, true
}

func indexSetValueOnContainer(ctx data.Context, ie *IndexExpression, container data.GetValue, indexVal data.GetValue, value data.Value) data.Control {
	switch arr := container.(type) {
	case *data.ArrayValue:
		if _, isNull := indexVal.(*data.NullValue); isNull {
			arr.List = append(arr.List, data.NewZVal(value))
			writeBackArrayProperty(ctx, ie.Array, arr)
			return nil
		}
		if iv, ok := indexVal.(data.AsString); ok {
			key := iv.AsString()
			for _, zval := range arr.List {
				if zval != nil && zval.Name == key {
					zval.Value = value
					writeBackArrayProperty(ctx, ie.Array, arr)
					return nil
				}
			}
			arr.List = append(arr.List, &data.ZVal{Name: key, Value: value})
			writeBackArrayProperty(ctx, ie.Array, arr)
			return nil
		}
		if iv, ok := indexVal.(data.AsInt); ok {
			i, err := iv.AsInt()
			if err != nil {
				return data.NewErrorThrow(ie.GetFrom(), err)
			}
			arr.SetIntKey(i, value)
			writeBackArrayProperty(ctx, ie.Array, arr)
			return nil
		}
	case *data.ClassValue:
		if checkArrayAccess(ctx, arr.Class) {
			return callArrayAccessOffsetSet(ctx, arr, indexVal, value)
		}
	case *data.ThisValue:
		if checkArrayAccess(ctx, arr.Class) {
			return callArrayAccessOffsetSet(ctx, arr.ClassValue, indexVal, value)
		}
	case data.SetProperty:
		if iv, ok := indexVal.(data.AsString); ok {
			arr.SetProperty(iv.AsString(), value)
			return nil
		}
		if iv, ok := indexVal.(data.AsInt); ok {
			if i, err := iv.AsInt(); err == nil {
				arr.SetProperty(fmt.Sprintf("%d", i), value)
				return nil
			}
		}
	}
	return data.NewErrorThrow(ie.GetFrom(), errors.New("无法在容器上设置索引值"))
}

func indexKeyString(index data.GetValue) (string, bool) {
	if _, isNull := index.(*data.NullValue); isNull {
		return "", true
	}
	if iv, ok := index.(data.AsString); ok {
		return iv.AsString(), true
	}
	if iv, ok := index.(data.AsInt); ok {
		if i, err := iv.AsInt(); err == nil {
			return fmt.Sprintf("%d", i), true
		}
	}
	return "", false
}

func concatPHPValues(left, right data.GetValue) data.Value {
	dot := &BinaryDot{Left: &literalGetValue{left}, Right: &literalGetValue{right}}
	v, _ := dot.GetValue(nil)
	if val, ok := v.(data.Value); ok {
		return val
	}
	return data.NewNullValue()
}

type literalGetValue struct{ v data.GetValue }

func (l *literalGetValue) GetValue(data.Context) (data.GetValue, data.Control) { return l.v, nil }
