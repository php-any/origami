package node

import (
	"github.com/php-any/origami/data"
)

// getDateTimeTimestamp 尝试将 value 视为实现了 DateTimeInterface 的对象，
// 返回其 Unix 时间戳（秒）。若不是日期时间对象，返回 ok=false。
//
// PHP 语义：两个 DateTime 对象比较时，比较其底层时间戳。
// Carbon 继承自 DateTime，因此 Carbon 对象之间的 <、>、==、<=> 等也按时间戳比较。
func getDateTimeTimestamp(ctx data.Context, value data.GetValue) (int64, bool) {
	cv, ok := value.(*data.ClassValue)
	if !ok {
		return 0, false
	}
	// 检查是否实现了 DateTimeInterface
	isDT, ctl := checkClassIs(ctx, cv.Class, "DateTimeInterface")
	if ctl != nil || !isDT {
		return 0, false
	}
	// 读取 timestamp 属性（DateTime 实现把 Unix 秒存于该属性）
	ts, ctl := cv.GetProperty("timestamp")
	if ctl != nil {
		return 0, false
	}
	if ts == nil {
		return 0, false
	}
	if iv, ok := ts.(*data.IntValue); ok {
		return int64(iv.Value), true
	}
	// 兼容浮点/字符串形式
	if ai, ok := ts.(data.AsInt); ok {
		v, err := ai.AsInt()
		if err != nil {
			return 0, false
		}
		return int64(v), true
	}
	return 0, false
}
