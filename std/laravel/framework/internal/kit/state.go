package kit

import "github.com/php-any/origami/data"

// CachedState 读取挂在对象隐蔽属性 prop 上的 Go 侧状态，没有就用 newState 建一份并挂上。
//
// 为什么不用「全局表 + 自增 id 做键」：那种写法遇到「每请求 new 一个对象」只会增不会删，
// 而 state 里通常存着 container（请求级 Application 克隆）/ passable / pipes，
// 于是每请求泄漏一条 state 加整棵请求对象图。examples/laravel13 实测：
// 每请求 ~400KB 常驻不回收，6000 请求后 RSS 4.5GB 且 GC 不降；
// 常驻堆爬到 26GB 触发换页，吞吐从 ~600rps 塌到 56rps、单请求卡到几十秒。
//
// 挂在对象上由 GC 随对象回收，不需要任何清理记账；克隆对象（属性表按值拷贝）
// 依旧共享同一份状态，与原 id 方案语义一致。
func CachedState[T any](cv *data.ClassValue, prop string, newState func() T) T {
	if cv == nil {
		return newState()
	}
	if v, ctl := cv.GetProperty(prop); ctl == nil && v != nil {
		if av, ok := Unwrap(v).(*data.AnyValue); ok && av != nil {
			if s, ok := av.Value.(T); ok {
				return s
			}
		}
	}
	s := newState()
	_ = cv.SetProperty(prop, data.NewAnyValue(s))
	return s
}
