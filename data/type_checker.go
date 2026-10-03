package data

// ValueKind describes PHP values; references are storage wrappers, not kinds.
// It deliberately does not add a virtual Kind method to every Value operation.
type ValueKind uint8

// Installed by the language layer; cold callable checks use the same resolver
// as callback registration, including visibility and magic dispatch.
var CallableResolver func(Context, Value) (Value, Control)

const (
	ValueUnknown ValueKind = iota
	ValueNull
	ValueBool
	ValueInt
	ValueFloat
	ValueString
	ValueArray
	ValueObject
	ValueResource
)

func ValueKindOf(value Value) ValueKind {
	switch value.(type) {
	case *NullValue:
		return ValueNull
	case *BoolValue:
		return ValueBool
	case *IntValue:
		return ValueInt
	case *FloatValue:
		return ValueFloat
	case *StringValue:
		return ValueString
	case *ArrayValue:
		return ValueArray
	case *ClassValue, *ThisValue, *FuncValue, *BoundFuncValue, *ThrowValue:
		return ValueObject
	default:
		if _, ok := value.(interface{ PHPResourceKind() }); ok {
			return ValueResource
		}
		return ValueUnknown
	}
}

func (a *TypeArena) Matches(ref TypeRef, value Value, ctx Context) bool {
	if ref >= typeBuiltinEnd {
		snapshot := a.snapshot.Load()
		node := snapshot.nodes[ref-typeBuiltinEnd]
		switch node.kind {
		case TypeKindNominal:
			if object, ok := value.(*ClassValue); ok {
				if provider, ok := object.GetVM().(DescriptorProvider); ok {
					if matched, linked := provider.ClassRegistry().IsA(object.Class, node.symbol); linked {
						return matched
					}
				}
			}
			if object, ok := value.(*ThisValue); ok {
				if provider, ok := object.GetVM().(DescriptorProvider); ok {
					if matched, linked := provider.ClassRegistry().IsA(object.Class, node.symbol); linked {
						return matched
					}
				}
			}
			return NominalValueMatches(value, Symbols.Name(node.symbol), ctx)
		case TypeKindUnion:
			for _, member := range snapshot.members[node.first : node.first+node.count] {
				if a.Matches(member, value, ctx) {
					return true
				}
			}
			return false
		case TypeKindIntersection:
			for _, member := range snapshot.members[node.first : node.first+node.count] {
				if !a.Matches(member, value, ctx) {
					return false
				}
			}
			return true
		}
		return false
	}
	switch ref {
	case TypeInvalid, TypeMixed:
		return true
	case TypeVoid, TypeNever:
		return false
	case TypeNull:
		return ValueKindOf(value) == ValueNull
	case TypeBool:
		return ValueKindOf(value) == ValueBool
	case TypeFalse:
		b, ok := value.(*BoolValue)
		return ok && !b.Value
	case TypeTrue:
		b, ok := value.(*BoolValue)
		return ok && b.Value
	case TypeInt:
		return ValueKindOf(value) == ValueInt
	case TypeFloat:
		return ValueKindOf(value) == ValueFloat
	case TypeString:
		return ValueKindOf(value) == ValueString
	case TypeArray:
		return ValueKindOf(value) == ValueArray
	case TypeObject:
		return ValueKindOf(value) == ValueObject
	case TypeCallable:
		if ctx != nil && CallableResolver != nil {
			resolved, ctl := CallableResolver(ctx, value)
			return resolved != nil && ctl == nil
		}
		return (Callable{}).Is(value)
	case TypeIterable:
		return ValueKindOf(value) == ValueArray || NominalValueMatches(value, "Traversable", ctx)
	case TypeAST:
		return (AST{}).Is(value)
	case TypeSelf, TypeParent, TypeStatic:
		owner := typeClassContext(ctx)
		if owner == nil {
			return false
		}
		class := owner.SelfClass
		if ref == TypeStatic {
			class = owner.StaticClass
		}
		if class == nil {
			class = owner.Class
		}
		if class == nil {
			return false
		}
		name := class.GetName()
		if ref == TypeParent {
			parent := class.GetExtend()
			if parent == nil {
				return false
			}
			name = *parent
		}
		return NominalValueMatches(value, name, ctx)
	default:
		return false
	}
}
func typeClassContext(ctx Context) *ClassMethodContext {
	for ctx != nil {
		switch current := ctx.(type) {
		case *ClassMethodContext:
			return current
		case *BoundContext:
			if current.ScopeClass != "" && current.GetVM() != nil {
				if class, ok := current.GetVM().GetClass(current.ScopeClass); ok {
					return &ClassMethodContext{ClassValue: current.BoundThis, Context: current.Context, SelfClass: class, StaticClass: class}
				}
			}
			ctx = current.Context
		default:
			return nil
		}
	}
	return nil
}

// Prepare implements exact matching and coercion for compact declarations.
// Snapshots and compound members are read-only; no interning occurs here.
func (a *TypeArena) Prepare(ref TypeRef, value Value, ctx Context) (Value, bool, Control) {
	if a.Matches(ref, value, ctx) {
		return value, true, nil
	}
	if ref < typeBuiltinEnd {
		if ctx != nil && ctx.StrictTypes() {
			if integer, ok := value.(*IntValue); ok && ref == TypeFloat {
				return NewFloatValue(float64(integer.Value)), true, nil
			}
			return nil, false, nil
		}
		// Single builtin declarations need no union mask or member traversal.
		// Reuse the same conversion routines as the compatibility entry point.
		switch ref {
		case TypeInt:
			prepared, ok := coerceToIntValue(value)
			return prepared, ok, nil
		case TypeFloat:
			prepared, ok := coerceToFloatValue(value)
			return prepared, ok, nil
		case TypeString:
			return coerceToStringValue(value, ctx)
		case TypeBool:
			return coerceToBoolValue(value)
		default:
			return nil, false, nil
		}
	}
	var scalarKinds uint8
	snapshot := a.snapshot.Load()
	node := snapshot.nodes[ref-typeBuiltinEnd]
	if node.kind != TypeKindUnion {
		return nil, false, nil
	}
	for _, member := range snapshot.members[node.first : node.first+node.count] {
		scalarKinds |= typeScalarMask(member)
	}
	if ctx != nil && ctx.StrictTypes() {
		if integer, ok := value.(*IntValue); ok && scalarKinds&2 != 0 {
			return NewFloatValue(float64(integer.Value)), true, nil
		}
		return nil, false, nil
	}
	return coerceScalarKinds(scalarKinds, value, ctx)
}
func typeScalarMask(ref TypeRef) uint8 {
	switch ref {
	case TypeInt:
		return 1
	case TypeFloat:
		return 2
	case TypeString:
		return 4
	case TypeBool:
		return 8
	}
	return 0
}

func AllowsImplicitReturn(ty Types) bool { return ty == nil || ty == TypeVoid }

func TypeAllowsNull(ty Types) bool {
	if ref, ok := ty.(TypeRef); ok {
		return declarationTypes.Matches(ref, NewNullValue(), nil)
	}
	return ty == nil || DeclaredTypeRef(ty).Matches(NewNullValue(), nil)
}
func NullableDeclaredBase(ref TypeRef) (TypeRef, bool) {
	if ref.Kind() != TypeKindUnion {
		return ref, false
	}
	members := ref.Members()
	if members.Len() == 2 {
		if members.At(0) == TypeNull {
			return members.At(1), true
		}
		if members.At(1) == TypeNull {
			return members.At(0), true
		}
	}
	return ref, false
}

func (ref TypeRef) RequiresClassScope() bool {
	if ref == TypeSelf || ref == TypeParent || ref == TypeStatic {
		return true
	}
	if ref.Kind() == TypeKindUnion || ref.Kind() == TypeKindIntersection {
		view := ref.Members()
		for i := 0; i < view.Len(); i++ {
			if view.At(i).RequiresClassScope() {
				return true
			}
		}
	}
	return false
}

// ParameterTypeScope is implemented only by callables which capture a PHP
// declaration scope. Native callback binders retain their symbol-table frame.
type ParameterTypeScope interface{ ParameterTypeContext(Context) Context }
