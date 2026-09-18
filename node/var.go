package node

import (
	"sync"

	"github.com/php-any/origami/data"
)

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

// StaticLocalsHolder 函数/方法级 static 局部变量存储。解析期每个含 static 的函数共用一份，
// 运行时只在真正执行到 static 语句时才分配内部 map。
type StaticLocalsHolder struct {
	once  sync.Once
	store *data.StaticLocals
}

// Store 返回（必要时创建）该函数的 static 存储。
func (h *StaticLocalsHolder) Store() *data.StaticLocals {
	if h == nil {
		return nil
	}
	h.once.Do(func() {
		h.store = data.NewStaticLocals()
	})
	return h.store
}

// StaticVarStatement 表示静态局部变量声明语句
type StaticVarStatement struct {
	*Node       `pp:"-"`
	Var         data.Variable
	Initializer data.GetValue
	owner       *StaticLocalsHolder
}

// NewStaticVarStatement 创建一个新的静态局部变量声明语句
func NewStaticVarStatement(token *TokenFrom, variable data.Variable, initializer data.GetValue, owner *StaticLocalsHolder) *StaticVarStatement {
	return &StaticVarStatement{
		Node:        NewNode(token),
		Var:         variable,
		Initializer: initializer,
		owner:       owner,
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
	if b, ok := ctx.(data.StaticLocalsBinder); ok {
		if b.StaticLocalsStore() == nil {
			b.BindStaticLocals(s.owner.Store())
		}
	}
	store := staticLocalsFromCtx(ctx)
	idx := s.Var.GetIndex()
	if store != nil {
		if slot := store.Slot(idx); slot != nil {
			ctx.SetIndexZVal(idx, slot)
			return nil, nil
		}
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
		ctx.SetIndexZVal(idx, store.EnsureSlot(idx, val))
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
