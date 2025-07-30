package reflect

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectClass 脚本反射类
type ReflectClass struct {
	node.Node
}

// GetValue 获取反射类的值
func (r *ReflectClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(&ReflectClass{}, ctx), nil
}

// GetName 返回类名
func (r *ReflectClass) GetName() string {
	return "Reflect"
}

// GetExtend 返回父类名
func (r *ReflectClass) GetExtend() *string {
	return nil
}

// GetImplements 返回实现的接口列表
func (r *ReflectClass) GetImplements() []string {
	return nil
}

// GetProperty 获取属性
func (r *ReflectClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetProperties 获取所有属性
func (r *ReflectClass) GetProperties() map[string]data.Property {
	return nil
}

// GetMethod 获取方法
func (r *ReflectClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "getClassInfo":
		return &GetClassInfoMethod{}, true
	case "getMethodInfo":
		return &GetMethodInfoMethod{}, true
	case "getPropertyInfo":
		return &GetPropertyInfoMethod{}, true
	case "listClasses":
		return &ListClassesMethod{}, true
	case "listMethods":
		return &ListMethodsMethod{}, true
	case "listProperties":
		return &ListPropertiesMethod{}, true
	case "getClassAnnotations":
		return &GetClassAnnotationsMethod{}, true
	case "getMethodAnnotations":
		return &GetMethodAnnotationsMethod{}, true
	case "getPropertyAnnotations":
		return &GetPropertyAnnotationsMethod{}, true
	case "getAllAnnotations":
		return &GetAllAnnotationsMethod{}, true
	case "getAnnotationDetails":
		return &GetAnnotationDetailsMethod{}, true
	}
	return nil, false
}

// GetMethods 获取所有方法
func (r *ReflectClass) GetMethods() []data.Method {
	return []data.Method{
		&GetClassInfoMethod{},
		&GetMethodInfoMethod{},
		&GetPropertyInfoMethod{},
		&ListClassesMethod{},
		&ListMethodsMethod{},
		&ListPropertiesMethod{},
		&GetClassAnnotationsMethod{},
		&GetMethodAnnotationsMethod{},
		&GetPropertyAnnotationsMethod{},
		&GetAllAnnotationsMethod{},
		&GetAnnotationDetailsMethod{},
	}
}

// GetConstruct 获取构造函数
func (r *ReflectClass) GetConstruct() data.Method {
	return nil
}
