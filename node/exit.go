package node

import "github.com/php-any/origami/data"

// ExitStatement 表示 exit / die 语句。
type ExitStatement struct {
	*Node `pp:"-"`
	Value data.GetValue // 可选状态码或输出字符串
}

func NewExitStatement(from data.From, value data.GetValue) *ExitStatement {
	return &ExitStatement{
		Node:  NewNode(from),
		Value: value,
	}
}

func (s *ExitStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	code := 0
	if s.Value != nil {
		v, ctl := s.Value.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if v != nil {
			if asInt, ok := v.(data.AsInt); ok {
				if n, err := asInt.AsInt(); err == nil {
					code = n
				}
			} else {
				str := ""
				if asStr, ok := v.(data.AsString); ok {
					str = asStr.AsString()
				} else if val, ok := v.(data.Value); ok {
					str = val.AsString()
				}
				if str != "" {
					data.EmitOutput(ctx, str)
				}
			}
		}
	}
	return data.NewNullValue(), data.NewExitControl(code)
}
