package node

import "github.com/php-any/origami/data"

// 用于检查和补充 from 跟踪信息
func checkThrowControlFrom(pre data.GetValue, v data.Control) data.Control {
	if e, ok := v.(*data.ThrowValue); ok {
		if e.Error.From == nil {
			if pre, ok := pre.(GetFrom); ok {
				if pre.GetFrom() != nil {
					e.Error.From = pre.GetFrom()
				}
			}
		}
	}
	return v
}
