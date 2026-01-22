package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArraySliceFunction 实现 array_slice 函数
// 从数组中取出一段
type ArraySliceFunction struct{}

func NewArraySliceFunction() data.FuncStmt {
	return &ArraySliceFunction{}
}

func (f *ArraySliceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：数组
	arrayValue, _ := ctx.GetIndexValue(0)
	if arrayValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 获取第二个参数：偏移量（必需）
	offsetValue, _ := ctx.GetIndexValue(1)
	if offsetValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	offset := 0
	if intVal, ok := offsetValue.(data.AsInt); ok {
		offset, _ = intVal.AsInt()
	}

	// 获取第三个参数：长度（可选）
	lengthValue, hasLength := ctx.GetIndexValue(2)
	length := -999 // 特殊值表示未提供长度参数（到数组末尾），null 也表示到数组末尾
	if hasLength && lengthValue != nil {
		// 检查是否是 null
		if _, isNull := lengthValue.(*data.NullValue); !isNull {
			if intVal, ok := lengthValue.(data.AsInt); ok {
				length, _ = intVal.AsInt()
			}
		}
		// 如果是 null，length 保持为 -999（表示到数组末尾）
	}

	// 获取第四个参数：preserve_keys（可选，默认 false）
	preserveKeysValue, hasPreserveKeys := ctx.GetIndexValue(3)
	preserveKeys := false
	if hasPreserveKeys && preserveKeysValue != nil {
		// 检查是否是 null
		if _, isNull := preserveKeysValue.(*data.NullValue); !isNull {
			if boolVal, ok := preserveKeysValue.(data.AsBool); ok {
				preserveKeys, _ = boolVal.AsBool()
			} else if intVal, ok := preserveKeysValue.(data.AsInt); ok {
				// 也支持整数 0/1 作为布尔值
				if i, err := intVal.AsInt(); err == nil && i != 0 {
					preserveKeys = true
				}
			}
		}
	}

	// 处理数组
	if arrayVal, ok := arrayValue.(*data.ArrayValue); ok {
		arrLen := len(arrayVal.List)

		// 处理负偏移量
		if offset < 0 {
			offset = arrLen + offset
			if offset < 0 {
				offset = 0
			}
		}

		// 如果偏移量超出数组范围，返回空数组
		if offset >= arrLen {
			return data.NewArrayValue([]data.Value{}), nil
		}

		// 计算结束位置
		end := arrLen
		if length == -999 {
			// 未提供长度参数或为 null，表示到数组末尾
			end = arrLen
		} else if length >= 0 {
			end = offset + length
			if end > arrLen {
				end = arrLen
			}
		} else {
			// 负长度：从 offset 开始，到距离末尾 |length| 的位置
			// 例如 length=-1 表示排除最后1个元素，end = arrLen - 1
			// length=-2 表示排除最后2个元素，end = arrLen - 2
			end = arrLen + length
			if end < offset {
				end = offset
			}
		}

		// 提取切片
		valueList := arrayVal.ToValueList()
		result := valueList[offset:end]
		return data.NewArrayValue(result), nil
	}

	// 处理对象（关联数组）
	if objectVal, ok := arrayValue.(*data.ObjectValue); ok {
		properties := objectVal.GetProperties()
		keys := make([]string, 0, len(properties))
		values := make([]data.Value, 0, len(properties))

		// 收集键和值
		for k, v := range properties {
			keys = append(keys, k)
			values = append(values, v)
		}

		// 处理负偏移量
		arrLen := len(keys)
		if offset < 0 {
			offset = arrLen + offset
			if offset < 0 {
				offset = 0
			}
		}

		// 如果偏移量超出数组范围，返回空对象
		if offset >= arrLen {
			return data.NewObjectValue(), nil
		}

		// 计算结束位置
		end := arrLen
		if length == -999 {
			// 未提供长度参数或为 null，表示到数组末尾
			end = arrLen
		} else if length >= 0 {
			end = offset + length
			if end > arrLen {
				end = arrLen
			}
		} else {
			// 负长度：从 offset 开始，到距离末尾 |length| 的位置
			// 例如 length=-1 表示排除最后1个元素，end = arrLen - 1
			// length=-2 表示排除最后2个元素，end = arrLen - 2
			end = arrLen + length
			if end < offset {
				end = offset
			}
		}

		// 提取切片
		if preserveKeys {
			// 保留键
			resultObj := data.NewObjectValue()
			for i := offset; i < end; i++ {
				resultObj.SetProperty(keys[i], values[i])
			}
			return resultObj, nil
		} else {
			// 重新索引（从 0 开始）
			result := make([]data.Value, 0, end-offset)
			for i := offset; i < end; i++ {
				result = append(result, values[i])
			}
			return data.NewArrayValue(result), nil
		}
	}

	// 不是数组类型，返回空数组
	return data.NewArrayValue([]data.Value{}), nil
}

func (f *ArraySliceFunction) GetName() string {
	return "array_slice"
}

func (f *ArraySliceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
		node.NewParameter(nil, "offset", 1, nil, nil),
		node.NewParameter(nil, "length", 2, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "preserve_keys", 3, node.NewIntLiteral(nil, "0"), nil),
	}
}

func (f *ArraySliceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "offset", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "length", 2, data.NewBaseType("int|null")),
		node.NewVariable(nil, "preserve_keys", 3, data.NewBaseType("bool")),
	}
}
