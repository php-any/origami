package core

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// ResourceClass 表示资源类定义
type ResourceClass struct {
	ResourceType string      // 资源类型，如 "process", "file", "stream" 等
	Resource     interface{} // 实际的资源对象
	id           int         // 资源ID（对于进程资源，这是真实的系统进程ID）
}

// NewResourceClass 创建资源类
func NewResourceClass(resourceType string, resource interface{}, resourceID int) *ResourceClass {
	return &ResourceClass{
		ResourceType: resourceType,
		Resource:     resource,
		id:           resourceID,
	}
}

// 实现 ClassStmt 接口
func (r *ResourceClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return NewResourceValue(r, ctx), nil
}

func (r *ResourceClass) GetFrom() data.From {
	return nil
}

func (r *ResourceClass) GetName() string {
	// 对于进程资源，显示为 "Resource#PID"
	// 对于其他资源类型，也使用资源ID
	return fmt.Sprintf("Resource#%d", r.id)
}

func (r *ResourceClass) GetExtend() *string {
	return nil
}

func (r *ResourceClass) GetImplements() []string {
	// 返回固定的资源接口名称，所有资源都实现 "Resource" 接口
	return []string{"Resource"}
}

func (r *ResourceClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (r *ResourceClass) GetPropertyList() []data.Property {
	return nil
}

func (r *ResourceClass) GetMethod(name string) (data.Method, bool) {
	return nil, false
}

func (r *ResourceClass) GetMethods() []data.Method {
	return nil
}

func (r *ResourceClass) GetConstruct() data.Method {
	return nil
}

// GetResourceID 获取资源ID
func (r *ResourceClass) GetResourceID() int {
	return r.id
}

// GetResourceType 获取资源类型
func (r *ResourceClass) GetResourceType() string {
	return r.ResourceType
}

// GetResource 获取实际的资源对象
func (r *ResourceClass) GetResource() interface{} {
	return r.Resource
}
