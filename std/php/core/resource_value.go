package core

import (
	"context"
	"fmt"
	"io"

	"github.com/php-any/origami/data"
)

// ResourceValue 表示资源类型（如文件句柄、进程句柄等）
// 嵌入 ClassValue，使资源可以作为类实例处理
type ResourceValue struct {
	*data.ClassValue
}

func (*ResourceValue) PHPResourceKind() {}
func (r *ResourceValue) IsPHPResourceOpen() bool {
	if resource := r.GetResource(); resource != nil {
		if state, ok := resource.(interface{ IsClosed() bool }); ok {
			return !state.IsClosed()
		}
		return true
	}
	return false
}

// The wrapper is request-owned; retained handles require an explicit policy.
func (r *ResourceValue) BindRequestValue(scope *data.RequestObjectScope) data.Value {
	class := r.Class.(*ResourceClass)
	resource := scope.BindNativeState(class.Resource)
	return NewResourceValue(NewResourceClass(class.ResourceType, resource, class.id), scope.Context)
}

// NewResourceValue 创建资源值
func NewResourceValue(resourceClass *ResourceClass, ctx data.Context) *ResourceValue {
	classValue := data.NewClassValue(resourceClass, ctx)
	if closer, ok := resourceClass.Resource.(io.Closer); ok {
		BindOwnedResource(ctx, closer)
	}
	return &ResourceValue{ClassValue: classValue}
}

// BindOwnedResource also covers handles held by native PHP objects, such as
// SplFileObject, which do not expose a PHP resource wrapper.
func BindOwnedResource(ctx data.Context, closer io.Closer) {
	if ctx != nil {
		request := ctx.GoContext()
		if request.Done() != nil {
			stop := context.AfterFunc(request, func() { _ = closer.Close() })
			if resource, ok := closer.(interface{ BindRequestCancel(func() bool) }); ok {
				resource.BindRequestCancel(stop)
			}
		}
	}
}

// AsString 重写 AsString 方法，显示资源ID
func (r *ResourceValue) AsString() string {
	if r.ClassValue != nil && r.ClassValue.Class != nil {
		if resourceClass, ok := r.ClassValue.Class.(*ResourceClass); ok {
			return fmt.Sprintf("Resource id #%d", resourceClass.GetResourceID())
		}
	}
	return "Resource"
}

// GetResourceID 获取资源ID（用于显示）
func (r *ResourceValue) GetResourceID() int {
	if r.ClassValue != nil && r.ClassValue.Class != nil {
		if resourceClass, ok := r.ClassValue.Class.(*ResourceClass); ok {
			return resourceClass.GetResourceID()
		}
	}
	return 0
}

// GetResourceType 获取资源类型
func (r *ResourceValue) GetResourceType() string {
	if !r.IsPHPResourceOpen() {
		return "Unknown"
	}
	if r.ClassValue != nil && r.ClassValue.Class != nil {
		if resourceClass, ok := r.ClassValue.Class.(*ResourceClass); ok {
			return resourceClass.GetResourceType()
		}
	}
	return ""
}

// GetResource 获取实际的资源对象
func (r *ResourceValue) GetResource() interface{} {
	// ResourceValue 嵌入了 ClassValue，通过 ClassValue.Class 获取 ResourceClass
	// 由于嵌入，可以直接访问 ClassValue 的字段
	if r.ClassValue != nil && r.ClassValue.Class != nil {
		if resourceClass, ok := r.ClassValue.Class.(*ResourceClass); ok {
			return resourceClass.GetResource()
		}
	}
	return nil
}
