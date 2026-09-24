package data

import (
	"strings"
	"sync"
	"sync/atomic"
)

// MethodLookupCache 缓存实例方法解析结果（含继承链），key 为小写方法名。
// PHP 方法名不区分大小写；trait 合并后必须 Invalidate。
type MethodLookupCache struct {
	mu   sync.RWMutex
	gen  uint64
	hits map[string]Method
	miss map[string]struct{}
}

func NewMethodLookupCache() *MethodLookupCache {
	return &MethodLookupCache{
		hits: make(map[string]Method),
		miss: make(map[string]struct{}),
	}
}

type MethodLookupCacher interface {
	MethodLookupCache() *MethodLookupCache
}

// MethodLookupKey 返回方法名的小写形式（PHP 方法名不区分大小写）。
//
// 绝大多数方法名本来就是全小写（`handle`、`__construct`），此时直接原样返回：
// 省掉 strings.ToLower 的整串扫描。含大写字母或非 ASCII 字节时才真正转换，
// 且转换结果按名字缓存 —— `$obj->setRequest()` 这类调用每次都会走到这里，
// 重复转换会为同一个名字反复分配新字符串，是热路径上可观的一笔垃圾。
func MethodLookupKey(name string) string {
	for i := 0; i < len(name); i++ {
		c := name[i]
		// 非 ASCII 字节可能是大写 Unicode 字母，交给 strings.ToLower 的完整规则
		if c >= 0x80 || (c >= 'A' && c <= 'Z') {
			return lowerMethodKey(name)
		}
	}
	return name
}

func lowerMethodKey(name string) string {
	if v, ok := methodKeyCache.Load(name); ok {
		s, _ := v.(string)
		return s
	}
	lowered := strings.ToLower(name)
	if methodKeyCacheSize.Load() < methodKeyCacheLimit {
		if _, loaded := methodKeyCache.LoadOrStore(name, lowered); !loaded {
			methodKeyCacheSize.Add(1)
		}
	}
	return lowered
}

// methodKeyCacheLimit 上限：方法名来自代码，总量有限；设上限只为防御动态方法名。
const methodKeyCacheLimit = 8192

var (
	methodKeyCache     sync.Map // string -> string
	methodKeyCacheSize atomic.Int64
)

func (c *MethodLookupCache) Lookup(key string) (Method, bool, bool, uint64) {
	if c == nil {
		return nil, false, false, 0
	}
	c.mu.RLock()
	gen := c.gen
	defer c.mu.RUnlock()
	if m, ok := c.hits[key]; ok {
		return m, true, true, gen
	}
	if _, miss := c.miss[key]; miss {
		return nil, false, true, gen
	}
	return nil, false, false, gen
}

func (c *MethodLookupCache) Store(key string, m Method, found bool, gen uint64) {
	if c == nil || key == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.gen != gen {
		return
	}
	if found && m != nil {
		c.hits[key] = m
		delete(c.miss, key)
		return
	}
	c.miss[key] = struct{}{}
	delete(c.hits, key)
}

func (c *MethodLookupCache) Invalidate() {
	if c == nil {
		return
	}
	c.mu.Lock()
	c.gen++
	c.hits = make(map[string]Method)
	c.miss = make(map[string]struct{})
	c.mu.Unlock()
}
