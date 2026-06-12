package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

type CallStaticProperty struct {
	*Node    `pp:"-"`
	Stmt     data.GetValue
	Property string // 属性名
}

func NewCallStaticProperty(token *TokenFrom, stmt data.GetValue, property string) *CallStaticProperty {
	return &CallStaticProperty{
		Node:     NewNode(token),
		Stmt:     stmt,
		Property: property,
	}
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallStaticProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	switch expr := pe.Stmt.(type) {
	case data.GetStaticProperty:
		property, ok := expr.GetStaticProperty(pe.Property)
		if ok {
			return property, nil
		}
		// 在父类中查找
		if cs, ok := expr.(data.ClassStmt); ok {
			if prop, found := pe.findStaticPropertyInParents(ctx, cs); found {
				return prop, nil
			}
		}
		return nil, data.NewPHPUncaughtError(pe.GetFrom(), fmt.Sprintf("Uncaught Error: Access to undeclared static property %s::$%s", TryGetCallClassName(pe.Stmt), pe.Property))
	default:
		next, acl := pe.Stmt.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		switch expr := next.(type) {
		case *data.ClassValue:
			if c, ok := expr.Class.(data.GetStaticProperty); ok {
				property, ok := c.GetStaticProperty(pe.Property)
				if ok {
					return property, nil
				}
			}
			if prop, found := pe.findStaticPropertyInParents(ctx, expr.Class); found {
				return prop, nil
			}

		case data.GetStaticProperty:
			property, ok := expr.GetStaticProperty(pe.Property)
			if ok {
				return property, nil
			}
			if cs, ok := expr.(data.ClassStmt); ok {
				if prop, found := pe.findStaticPropertyInParents(ctx, cs); found {
					return prop, nil
				}
			}
		}
	}
	return nil, data.NewPHPUncaughtError(pe.GetFrom(), fmt.Sprintf("Uncaught Error: Access to undeclared static property %s::$%s", pe.getClassName(), pe.Property))
}

// findStaticPropertyInParents 在父类继承链中查找静态属性/常量
func (pe *CallStaticProperty) findStaticPropertyInParents(ctx data.Context, class data.ClassStmt) (data.Value, bool) {
	vm := ctx.GetVM()
	extend := class.GetExtend()
	for extend != nil {
		parent, acl := vm.GetOrLoadClass(*extend)
		if acl != nil || parent == nil {
			break
		}
		if gsp, ok := parent.(data.GetStaticProperty); ok {
			if prop, found := gsp.GetStaticProperty(pe.Property); found {
				return prop, true
			}
		}
		extend = parent.GetExtend()
	}
	return nil, false
}

func (pe *CallStaticProperty) hasDeclaredStaticProperty(vm data.VM) bool {
	class, ok := pe.Stmt.(data.ClassStmt)
	if !ok {
		return false
	}
	// 检查当前类
	if prop, ok := class.GetProperty(pe.Property); ok && prop.GetIsStatic() {
		return true
	}
	// 检查父类
	extend := class.GetExtend()
	seen := map[string]bool{}
	for extend != nil {
		if seen[*extend] {
			break
		}
		seen[*extend] = true
		parent, acl := vm.GetOrLoadClass(*extend)
		if acl != nil || parent == nil {
			break
		}
		if prop, ok := parent.GetProperty(pe.Property); ok && prop.GetIsStatic() {
			return true
		}
		extend = parent.GetExtend()
	}
	return false
}

func (pe *CallStaticProperty) getClassName() string {
	if getName, ok := pe.Stmt.(data.ClassStmt); ok {
		return getName.GetName()
	}
	return ""
}

func (pe *CallStaticProperty) SetProperty(ctx data.Context, name string, value data.Value) data.Control {
	vm := ctx.GetVM()
	if pe.hasDeclaredStaticProperty(vm) {
		switch c := pe.Stmt.(type) {
		case *ClassStatement:
			c.StaticProperty.Store(name, value)
			return nil
		case *ClassGeneric:
			c.ClassStatement.StaticProperty.Store(name, value)
			return nil
		case data.SetProperty:
			return c.SetProperty(name, value)
		}
		return nil
	}
	return data.NewPHPUncaughtError(pe.GetFrom(), fmt.Sprintf("Uncaught Error: Access to undeclared static property %s::$%s", pe.getClassName(), pe.Property))
}

// CallStaticPropertyLater 延迟的静态属性访问（类未加载时）
type CallStaticPropertyLater struct {
	*Node
	className string              // 类名（字符串形式）
	property  string              // 属性名
	namespace string              // 命名空间
	access    *CallStaticProperty `pp:"-"` // 解析后缓存
}

// NewCallStaticPropertyLater 创建延迟的静态属性访问
func NewCallStaticPropertyLater(from *TokenFrom, className, property, namespace string) *CallStaticPropertyLater {
	return &CallStaticPropertyLater{
		Node:      NewNode(from),
		className: className,
		property:  property,
		namespace: namespace,
	}
}

func (pe *CallStaticPropertyLater) resolveAccess(ctx data.Context) (*CallStaticProperty, data.Control) {
	if pe.access != nil {
		return pe.access, nil
	}
	vm := ctx.GetVM()
	target, acl := vm.LoadPkg(pe.className)
	if acl != nil {
		return nil, acl
	}
	if target == nil {
		fullName := pe.className
		if pe.namespace != "" {
			fullName = pe.namespace + "\\" + pe.className
		}
		target, acl = vm.LoadPkg(fullName)
		if acl != nil {
			return nil, acl
		}
	}
	if target == nil {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法访问静态属性(%s::%s), 未找到类或接口", pe.className, pe.property))
	}
	tokenFrom, ok := pe.GetFrom().(*TokenFrom)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法获取TokenFrom信息"))
	}
	pe.access = NewCallStaticProperty(tokenFrom, target, pe.property)
	return pe.access, nil
}

// GetValue 获取延迟静态属性访问的值
func (pe *CallStaticPropertyLater) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	access, acl := pe.resolveAccess(ctx)
	if acl != nil {
		return nil, acl
	}
	return access.GetValue(ctx)
}

// SetProperty 设置延迟静态属性的值
func (pe *CallStaticPropertyLater) SetProperty(ctx data.Context, name string, value data.Value) data.Control {
	access, acl := pe.resolveAccess(ctx)
	if acl != nil {
		return acl
	}
	return access.SetProperty(ctx, name, value)
}
