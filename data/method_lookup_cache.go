package data

import (
	"strings"
	"sync"
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

func MethodLookupKey(name string) string {
	return strings.ToLower(name)
}

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
