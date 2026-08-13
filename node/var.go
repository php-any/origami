package node

import "github.com/php-any/origami/data"

func (u *VarStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// use语句本身不返回值
	return nil, nil
}

// VarStatement 表示变量声明语句
type VarStatement struct {
	*Node       `pp:"-"`
	Name        string
	Initializer data.GetValue
}

// NewVarStatement 创建一个新的变量声明语句
func NewVarStatement(token *TokenFrom, name string, initializer data.GetValue) *VarStatement {
	return &VarStatement{
		Node:        NewNode(token),
		Name:        name,
		Initializer: initializer,
	}
}

// StaticVarStatement 表示静态局部变量声明语句
type StaticVarStatement struct {
	*Node       `pp:"-"`
	Var         data.Variable
	Initializer data.GetValue
}

// NewStaticVarStatement 创建一个新的静态局部变量声明语句
func NewStaticVarStatement(token *TokenFrom, variable data.Variable, initializer data.GetValue) *StaticVarStatement {
	return &StaticVarStatement{
		Node:        NewNode(token),
		Var:         variable,
		Initializer: initializer,
	}
}

func bindStaticLocals(ctx data.Context, store *data.StaticLocals) {
	if b, ok := ctx.(data.StaticLocalsBinder); ok {
		b.BindStaticLocals(store)
	}
}

func staticLocalsFromCtx(ctx data.Context) *data.StaticLocals {
	if b, ok := ctx.(data.StaticLocalsBinder); ok {
		return b.StaticLocalsStore()
	}
	return nil
}

func (s *StaticVarStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	store := staticLocalsFromCtx(ctx)
	idx := s.Var.GetIndex()
	if store != nil {
		if _, ok := store.Get(idx); !ok {
			val := data.NewNullValue()
			if s.Initializer != nil {
				init, ctl := s.Initializer.GetValue(ctx)
				if ctl != nil {
					return nil, ctl
				}
				if v, ok := init.(data.Value); ok {
					val = v
				}
			}
			store.Init(idx, val)
		}
		if v, ok := store.Get(idx); ok {
			if ctl := s.Var.SetValue(ctx, v); ctl != nil {
				return nil, ctl
			}
		}
		return nil, nil
	}
	if s.Initializer != nil {
		init, ctl := s.Initializer.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if v, ok := init.(data.Value); ok {
			if ctl := s.Var.SetValue(ctx, v); ctl != nil {
				return nil, ctl
			}
		}
	}
	return nil, nil
}
