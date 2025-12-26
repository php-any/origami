package core

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// ResourceValue 表示资源类型（如文件句柄、进程句柄等）
// 嵌入 ClassValue，使资源可以作为类实例处理
type ResourceValue struct {
	*data.ClassValue
}

// NewResourceValue 创建资源值
func NewResourceValue(resourceClass *ResourceClass, ctx data.Context) *ResourceValue {
	classValue := data.NewClassValue(resourceClass, ctx)
	return &ResourceValue{
		ClassValue: classValue,
	}
}

// AsString 重写 AsString 方法，显示资源ID
func (r *ResourceValue) AsString() string {
	if resourceClass, ok := r.Class.(*ResourceClass); ok {
		return fmt.Sprintf("Resource id #%d", resourceClass.GetResourceID())
	}
	return "Resource"
}

// GetResourceID 获取资源ID（用于显示）
func (r *ResourceValue) GetResourceID() int {
	if resourceClass, ok := r.Class.(*ResourceClass); ok {
		return resourceClass.GetResourceID()
	}
	return 0
}

// GetResourceType 获取资源类型
func (r *ResourceValue) GetResourceType() string {
	if resourceClass, ok := r.Class.(*ResourceClass); ok {
		return resourceClass.GetResourceType()
	}
	return ""
}

// GetResource 获取实际的资源对象
func (r *ResourceValue) GetResource() interface{} {
	if resourceClass, ok := r.Class.(*ResourceClass); ok {
		return resourceClass.GetResource()
	}
	return nil
}
