package data

import "sync"

// chainStore 是「请求级对象」的属性表：本地表为空，读未命中时回落到 parent（全局单例）。
//
// 存在的理由：Laravel 的 Application/Container 上挂着二十来张按请求可能改写的表
// （instances / bindings / resolved / aliases / abstractAliases / deferredServices …）。
// 原先的做法是每请求把它们全量深拷贝一份（CloneSandboxKeys + appSandboxDeepKeys），
// 暖路径上是纯浪费——一次请求真正**写**到的通常只有 instances 等一两个键，
// bindings 这种几百项的表根本不会被写，却要连内层数组一起递归拷。
//
// 语义与「每请求先全量深拷贝」一致，只是把拷贝时机从「请求开始」挪到「第一次访问」：
//   - Get      本地没有就用 parent 的；容器类值升级时拷成请求私有副本（见下）
//   - GetZVal  写目标。首次写之前把 parent 的值升级到本地，之后都写本地
//   - Set      只写本地，不回写 parent
//   - Delete   本地删除 + 记影子，避免 parent 的值又回落上来
//
// 于是「拷多少」从「表有多大」变成「请求期真正碰过多少」。
//
// 为什么 Get 对容器也要升级（拷贝）、不能直接透传 parent 的数组/关联数组：
// PHP 里 `$this->instances[$k] = $v` 是「读 $this->instances 再原地改」，
// node 侧的 IndexExpression.SetValue 拿 GetValue 返回的容器直接 SetStringKey，
// 不会在属性这一层做写前分离；嵌套写（`$this->tags[$t][] = $v`）更是只有最内层
// 那一级会分离。透传或只分离一层，都会让请求改到全局容器的内层。
// 所以这里给调用方的永远是本地副本，深度与旧的 appSandboxDeepKeys 全量深拷贝一致。
// 代价是「读一次容器属性 = 一次深拷贝」，换来的是不必再为没人碰的键付钱。
//
// 对象（*ClassValue，即服务实例）与标量按引用透传：与 PHP 数组按值、对象按引用一致，
// 也与旧实现 deepCloneValue 只认 *ArrayValue / *ObjectValue 的行为一致。
//
// 本结构不做跨 goroutine 共享：请求级对象本就单 goroutine 独占。
type chainStore struct {
	parent PropertyStore
	local  PropertyStore

	// mu 只护 shadow 与「升级」这段临界区。local 自己的 OrderedMap 有锁，
	// 两把锁无嵌套获取顺序，不会死锁。
	mu     sync.Mutex
	shadow map[string]struct{}
	arrays *ArrayOverlayScope
	scope  *RequestObjectScope
}

// NewChainPropertyStore 用 parent 做回落源，建一张空的本地属性表。
func NewChainPropertyStore(parent PropertyStore) PropertyStore {
	return &chainStore{parent: parent, local: NewOrderedMap()}
}

func (s *chainStore) BindZVal(key string, slot *ZVal) {
	s.local.(interface{ BindZVal(string, *ZVal) }).BindZVal(key, slot)
}

func (s *chainStore) isShadowed(key string) bool {
	if s.shadow == nil {
		return false
	}
	_, ok := s.shadow[key]
	return ok
}

// promoteContainer 把 parent 的容器值拷成请求私有的副本后存进本地。
// 非容器值原样返回（第二返回值为 false 表示不需要升级）。
//
// 只分离**顶层一层**（CloneArrayValue / CloneObjectValue），嵌套数组的分离交给写路径：
// node.cowSeparateNestedArray 在 $this->map[$k][$j] = $v 这类嵌套写入前会克隆内层数组，
// 再由 writeBackArrayProperty 把克隆写回父容器。容器里被嵌套写的只有
//
//	Container::alias()          -> $this->abstractAliases[$abstract][] = $alias
//	Container::tag()            -> $this->tags[$tag][] = $abstract
//	Container::resolving()      -> $this->resolvingCallbacks[$abstract][] = $cb
//	Container::rebinding()      -> $this->reboundCallbacks[$abstract][] = $cb
//
// 四处，形态都是 $this->map[$k][] = $v，走的正是那条逐层分离路径。
//
// 旧实现每次都 deepCloneValue（递归整棵子树）：分配画像里
// promoteContainer -> deepCloneValue 占每请求分配量的 28%（≈214KB/请求），
// 而其中绝大多数内层数组这一层根本不会被写。分离一层 +
// 写路径按需逐层分离，得到的隔离语义相同，代价从「表有多大」变成「真正写了哪条路径」。
//
// 对象（服务实例）按引用共享，与 PHP 对象语义、与旧 deepCloneValue 的行为一致。
func (s *chainStore) promoteContainer(v Value) (Value, bool) {
	if s.scope != nil {
		bound := s.scope.Bind(v)
		return bound, bound != v
	}
	switch t := v.(type) {
	case *ArrayValue:
		if s.arrays == nil {
			s.arrays = NewArrayOverlayScope()
		}
		return s.arrays.Array(t), true
	default:
		return v, false
	}
}

func (s *chainStore) Get(key string) (Value, bool) {
	if v, ok := s.local.Get(key); ok {
		return v, true
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isShadowed(key) {
		return nil, false
	}
	v, ok := s.parent.Get(key)
	if !ok {
		return nil, false
	}
	if s.scope != nil {
		if source, ok := s.parent.GetZVal(key); ok && (source.RefCount() > 0 || source.Guard() != nil) {
			target := s.scope.BindSlot(source)
			s.local.(interface{ BindZVal(string, *ZVal) }).BindZVal(key, target)
			return target.ReadValue(), true
		}
	}
	if promoted, isContainer := s.promoteContainer(v); isContainer {
		s.local.Set(key, promoted)
		return promoted, true
	}
	return v, true
}

// GetZVal 取写目标的 zval：本地没有就把 parent 的值升级到本地再返回本地的 zval。
func (s *chainStore) GetZVal(key string) (*ZVal, bool) {
	if z, ok := s.local.GetZVal(key); ok {
		return z, true
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isShadowed(key) {
		return nil, false
	}
	parentSlot, ok := s.parent.GetZVal(key)
	if !ok {
		return nil, false
	}
	if s.scope != nil {
		target := s.scope.BindSlot(parentSlot)
		s.local.(interface{ BindZVal(string, *ZVal) }).BindZVal(key, target)
		return target, true
	}
	pv := parentSlot.ReadValue()
	if promoted, isContainer := s.promoteContainer(pv); isContainer {
		pv = promoted
	}
	s.local.Set(key, pv)
	return s.local.GetZVal(key)
}

func (s *chainStore) Set(key string, value Value) {
	s.local.Set(key, value)
	s.mu.Lock()
	delete(s.shadow, key)
	s.mu.Unlock()
}

func (s *chainStore) Delete(key string) {
	s.local.Delete(key)
	s.mu.Lock()
	if s.shadow == nil {
		s.shadow = make(map[string]struct{})
	}
	s.shadow[key] = struct{}{}
	s.mu.Unlock()
}

// Range 按「parent 的插入顺序优先、本地独有的键随后」遍历，与全量拷贝时的顺序一致。
func (s *chainStore) Range(fn func(key string, value Value) bool) {
	proceed := true
	s.parent.Range(func(key string, value Value) bool {
		v, ok := s.Get(key)
		if !ok {
			return true
		}
		value = v
		if !fn(key, value) {
			proceed = false
			return false
		}
		return true
	})
	if !proceed {
		return
	}
	s.local.Range(func(key string, value Value) bool {
		if _, ok := s.parent.Get(key); ok {
			return true // 前面按 parent 顺序已经报过
		}
		return fn(key, value)
	})
}

func (s *chainStore) Len() int {
	n := 0
	s.Range(func(string, Value) bool {
		n++
		return true
	})
	return n
}

func (s *chainStore) GetByIndex(index int) (string, Value, bool) {
	if index < 0 {
		return "", nil, false
	}
	i := 0
	key, val := "", Value(nil)
	found := false
	s.Range(func(k string, v Value) bool {
		if i == index {
			key, val, found = k, v, true
			return false
		}
		i++
		return true
	})
	return key, val, found
}
