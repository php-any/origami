package reflection

import (
	"github.com/php-any/origami/data"
)

// newPhpReflectionType 对齐 PHP ReflectionParameter/Property/Method::get*Type()：
// 联合类型返回 ReflectionUnionType，交集返回 ReflectionIntersectionType，其余为 ReflectionNamedType。
// Laravel Reflector::getParameterClassName 只对 ReflectionNamedType 取类名；
// 若把 int|string 做成 NamedType，会 new ReflectionClass("int|string") 抛 Class does not exist。
func newPhpReflectionType(ctx data.Context, typeInfo data.Types) data.Value {
	if typeInfo == nil {
		return data.NewNullValue()
	}
	if u, ok := asUnionType(typeInfo); ok && len(u.Types) > 0 {
		return newReflectionUnionType(ctx, u)
	}
	if i, ok := asIntersectionType(typeInfo); ok && len(i.Types) > 0 {
		return newReflectionIntersectionType(ctx, i)
	}
	return newReflectionNamedType(ctx, typeInfo)
}

func asUnionType(t data.Types) (data.UnionType, bool) {
	switch u := t.(type) {
	case data.UnionType:
		return u, true
	case *data.UnionType:
		if u != nil {
			return *u, true
		}
	}
	return data.UnionType{}, false
}

func asIntersectionType(t data.Types) (data.IntersectionType, bool) {
	switch i := t.(type) {
	case data.IntersectionType:
		return i, true
	case *data.IntersectionType:
		if i != nil {
			return *i, true
		}
	}
	return data.IntersectionType{}, false
}

func compoundAllowsNull(members []data.Types) bool {
	for _, t := range members {
		if typeAllowsNullMember(t) {
			return true
		}
	}
	return false
}

func typeAllowsNullMember(t data.Types) bool {
	if t == nil {
		return false
	}
	switch x := t.(type) {
	case data.NullableType:
		return true
	case *data.NullableType:
		return true
	case data.NullType:
		return true
	case *data.NullType:
		return true
	default:
		s := x.String()
		return s == "null" || s == "?null"
	}
}
