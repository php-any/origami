package collections

import "github.com/php-any/origami/data"

// laravelValue 对齐 Laravel value($value, ...$args)：仅 Closure 会被调用。
func laravelValue(ctx data.Context, v data.Value, args ...data.Value) (data.GetValue, data.Control) {
	if v == nil {
		return data.NewNullValue(), nil
	}
	switch v.(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return callValue(ctx, v, args...)
	default:
		return v, nil
	}
}
