package node

import "github.com/php-any/origami/data"

// SpawnStatement 表示 spawn 异步执行语句
type SpawnStatement struct {
	*Node `pp:"-"`
	Call  data.GetValue // spawn 后面的方法调用
}

// NewSpawnStatement 创建一个新的 spawn 语句
func NewSpawnStatement(token *TokenFrom, call data.GetValue) *SpawnStatement {
	return &SpawnStatement{
		Node: NewNode(token),
		Call: call,
	}
}

// GetValue 获取 spawn 语句的值
func (s *SpawnStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 异步执行方法调用
	go func() {
		_, acl := s.Call.GetValue(ctx)
		if acl != nil {
			ctx.GetVM().ThrowControl(acl)
		}
	}()

	return nil, nil
}
