package php

import "github.com/php-any/origami/data"

// PHP error_get_last / error_clear_last 用的最近错误状态。
func lastErrorToArray(info *data.PHPErrorInfo) data.Value {
	return data.NewArrayValueFromSlots([]*data.ZVal{
		data.NewNamedZVal("type", data.NewIntValue(info.Type)),
		data.NewNamedZVal("message", data.NewStringValue(info.Message)),
		data.NewNamedZVal("file", data.NewStringValue(info.File)),
		data.NewNamedZVal("line", data.NewIntValue(info.Line)),
	})

}
