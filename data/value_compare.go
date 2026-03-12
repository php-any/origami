package data

// Compare 比较两个值，返回 -1、0 或 1
// 返回值：-1（v1 < v2）、0（v1 == v2）、1（v1 > v2）
func Compare(v1, v2 Value) int {
	// 处理 null
	if _, ok := v1.(AsNull); ok {
		if _, ok := v2.(AsNull); ok {
			return 0
		}
		return -1
	}
	if _, ok := v2.(AsNull); ok {
		return 1
	}

	// 尝试转换为整数
	if int1, ok := v1.(*IntValue); ok {
		if int2, ok := v2.(*IntValue); ok {
			n1 := int1.Value
			n2 := int2.Value
			if n1 < n2 {
				return -1
			} else if n1 > n2 {
				return 1
			}
			return 0
		}
	}

	// 尝试转换为浮点数
	if float1, ok := v1.(*FloatValue); ok {
		if float2, ok := v2.(*FloatValue); ok {
			f1 := float1.Value
			f2 := float2.Value
			if f1 < f2 {
				return -1
			} else if f1 > f2 {
				return 1
			}
			return 0
		}
	}

	// 尝试转换为字符串
	if str1, ok := v1.(*StringValue); ok {
		if str2, ok := v2.(*StringValue); ok {
			s1 := str1.Value
			s2 := str2.Value
			if s1 < s2 {
				return -1
			} else if s1 > s2 {
				return 1
			}
			return 0
		}
	}

	// 尝试转换为布尔值
	if bool1, ok := v1.(*BoolValue); ok {
		if bool2, ok := v2.(*BoolValue); ok {
			b1 := bool1.Value
			b2 := bool2.Value
			if !b1 && b2 {
				return -1
			} else if b1 && !b2 {
				return 1
			}
			return 0
		}
	}

	// 默认情况：返回 0
	return 0
}
