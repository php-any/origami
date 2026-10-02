package data

import (
	"slices"
	"strings"
	"sync"
	"sync/atomic"
)

// TypeRef is an immutable, process-local declaration handle. The zero handle
// means no declaration, distinct from mixed, void and never.
type TypeRef uint32
type SymbolID uint32
type TypeKind uint8

const (
	TypeInvalid TypeRef = iota
	TypeMixed
	TypeVoid
	TypeNever
	TypeNull
	TypeFalse
	TypeTrue
	TypeBool
	TypeInt
	TypeFloat
	TypeString
	TypeArray
	TypeObject
	TypeCallable
	TypeIterable
	TypeSelf
	TypeParent
	TypeStatic
	typeBuiltinEnd
)

const (
	TypeKindBuiltin TypeKind = iota
	TypeKindNominal
	TypeKindUnion
	TypeKindIntersection
)

type typeNode struct {
	symbol SymbolID
	first  uint32
	count  uint32
	kind   TypeKind
}
type typeSymbol struct{ display, lookup string }
type typeSnapshot struct {
	nodes   []typeNode
	members []TypeRef
	symbols []typeSymbol
}

// Construction is serialized; checks read only append-only published prefixes.
// Appending never changes an existing node, member or symbol, so snapshots may
// safely share their underlying arrays without a lock in the runtime checker.
type TypeArena struct {
	mu        sync.Mutex
	nodes     []typeNode
	members   []TypeRef
	symbols   []typeSymbol
	names     map[string]TypeRef
	compounds map[string]TypeRef
	snapshot  atomic.Pointer[typeSnapshot]
}

var declarationTypes = NewTypeArena()

func nominalKey(name string) string {
	for i := 0; i < len(name); i++ {
		if name[i] >= 'A' && name[i] <= 'Z' {
			bytes := []byte(name)
			for j := i; j < len(bytes); j++ {
				if bytes[j] >= 'A' && bytes[j] <= 'Z' {
					bytes[j] += 'a' - 'A'
				}
			}
			return string(bytes)
		}
	}
	return name
}

func NewTypeArena() *TypeArena {
	a := &TypeArena{names: make(map[string]TypeRef), compounds: make(map[string]TypeRef)}
	a.publish()
	return a
}
func (a *TypeArena) publish() {
	a.snapshot.Store(&typeSnapshot{nodes: a.nodes, members: a.members, symbols: a.symbols})
}
func (a *TypeArena) Nominal(name string) TypeRef {
	name = strings.TrimPrefix(name, "\\")
	key := nominalKey(name)
	a.mu.Lock()
	defer a.mu.Unlock()
	if ref, ok := a.names[key]; ok {
		return ref
	}
	symbol := SymbolID(len(a.symbols))
	a.symbols = append(a.symbols, typeSymbol{display: name, lookup: key})
	ref := typeBuiltinEnd + TypeRef(len(a.nodes))
	a.nodes = append(a.nodes, typeNode{symbol: symbol, kind: TypeKindNominal})
	a.names[key] = ref
	a.publish()
	return ref
}
func (a *TypeArena) Union(members ...TypeRef) TypeRef { return a.compound(TypeKindUnion, members) }
func (a *TypeArena) Intersection(members ...TypeRef) TypeRef {
	return a.compound(TypeKindIntersection, members)
}
func (a *TypeArena) compound(kind TypeKind, members []TypeRef) TypeRef {
	a.mu.Lock()
	defer a.mu.Unlock()
	flat := make([]TypeRef, 0, len(members))
	var appendMember func(TypeRef)
	appendMember = func(ref TypeRef) {
		if ref >= typeBuiltinEnd {
			node := a.nodes[ref-typeBuiltinEnd]
			if node.kind == kind {
				for _, member := range a.members[node.first : node.first+node.count] {
					appendMember(member)
				}
				return
			}
		}
		flat = append(flat, ref)
	}
	for _, member := range members {
		appendMember(member)
	}
	slices.Sort(flat)
	flat = slices.Compact(flat)
	if kind == TypeKindUnion {
		if slices.Contains(flat, TypeMixed) {
			return TypeMixed
		}
		flat = slices.DeleteFunc(flat, func(ref TypeRef) bool { return ref == TypeNever })
		if slices.Contains(flat, TypeBool) || (slices.Contains(flat, TypeTrue) && slices.Contains(flat, TypeFalse)) {
			flat = slices.DeleteFunc(flat, func(ref TypeRef) bool { return ref == TypeTrue || ref == TypeFalse })
			if !slices.Contains(flat, TypeBool) {
				flat = append(flat, TypeBool)
				slices.Sort(flat)
			}
		}
	}
	if len(flat) == 0 {
		return TypeNever
	}
	if len(flat) == 1 {
		return flat[0]
	}
	// Binary keys retain all members; no hash collision can alias declarations.
	key := make([]byte, 1+4*len(flat))
	key[0] = byte(kind)
	for i, ref := range flat {
		n := 1 + 4*i
		key[n] = byte(ref)
		key[n+1] = byte(ref >> 8)
		key[n+2] = byte(ref >> 16)
		key[n+3] = byte(ref >> 24)
	}
	if ref, ok := a.compounds[string(key)]; ok {
		return ref
	}
	ref := typeBuiltinEnd + TypeRef(len(a.nodes))
	a.nodes = append(a.nodes, typeNode{kind: kind, first: uint32(len(a.members)), count: uint32(len(flat))})
	a.members = append(a.members, flat...)
	a.compounds[string(key)] = ref
	a.publish()
	return ref
}
func (a *TypeArena) Kind(ref TypeRef) TypeKind {
	if ref < typeBuiltinEnd {
		return TypeKindBuiltin
	}
	return a.snapshot.Load().nodes[ref-typeBuiltinEnd].kind
}

// Members returns a structural, read-only view, without a temporary slice.
func (a *TypeArena) Members(ref TypeRef) Span[TypeRef] {
	if ref < typeBuiltinEnd {
		return NewSpan[TypeRef](nil)
	}
	snapshot := a.snapshot.Load()
	node := snapshot.nodes[ref-typeBuiltinEnd]
	return NewSpan(snapshot.members[node.first : node.first+node.count])
}
func (a *TypeArena) Name(ref TypeRef) string {
	if ref < typeBuiltinEnd {
		return builtinTypeNames[ref]
	}
	snapshot := a.snapshot.Load()
	node := snapshot.nodes[ref-typeBuiltinEnd]
	if node.kind == TypeKindNominal {
		return snapshot.symbols[node.symbol].display
	}
	var result strings.Builder
	separator := "|"
	if node.kind == TypeKindIntersection {
		separator = "&"
	}
	for i, member := range snapshot.members[node.first : node.first+node.count] {
		if i > 0 {
			result.WriteString(separator)
		}
		nested := a.Kind(member) == TypeKindIntersection && node.kind == TypeKindUnion
		if nested {
			result.WriteByte('(')
		}
		result.WriteString(a.Name(member))
		if nested {
			result.WriteByte(')')
		}
	}
	return result.String()
}

var builtinTypeNames = [...]string{"", "mixed", "void", "never", "null", "false", "true", "bool", "int", "float", "string", "array", "object", "callable", "iterable", "self", "parent", "static"}

func (ref TypeRef) String() string         { return declarationTypes.Name(ref) }
func (ref TypeRef) Is(value Value) bool    { return declarationTypes.Matches(ref, value, nil) }
func (ref TypeRef) Kind() TypeKind         { return declarationTypes.Kind(ref) }
func (ref TypeRef) Members() Span[TypeRef] { return declarationTypes.Members(ref) }

// NewDeclaredType is the parser/native declaration entry point. Existing Types
// constructors remain an adapter for extensions during migration.
func NewDeclaredType(name string) Types {
	if name == "" {
		return nil
	}
	if strings.HasPrefix(name, "?") {
		return NewDeclaredNullableType(NewDeclaredType(name[1:]))
	}
	if strings.Contains(name, "|") {
		parts := strings.Split(name, "|")
		members := make([]Types, len(parts))
		for i, part := range parts {
			members[i] = NewDeclaredType(part)
		}
		return NewDeclaredUnionType(members)
	}
	if strings.Contains(name, "&") {
		parts := strings.Split(name, "&")
		members := make([]Types, len(parts))
		for i, part := range parts {
			members[i] = NewDeclaredType(part)
		}
		return NewDeclaredIntersectionType(members)
	}
	if ref, ok := declaredBuiltinRef(name); ok {
		return ref
	}
	return declarationTypes.Nominal(name)
}

func declaredBuiltinRef(name string) (TypeRef, bool) {
	key := nominalKey(name)
	for ref, builtin := range builtinTypeNames {
		if key == builtin {
			return TypeRef(ref), true
		}
	}
	return TypeInvalid, false
}
func DeclaredTypeRef(ty Types) TypeRef {
	if ty == nil {
		return TypeInvalid
	}
	if ref, ok := ty.(TypeRef); ok {
		return ref
	}
	switch ty := ty.(type) {
	case NullableType:
		return declarationTypes.Union(DeclaredTypeRef(ty.BaseType), TypeNull)
	case UnionType:
		members := make([]TypeRef, len(ty.Types))
		for i, member := range ty.Types {
			members[i] = DeclaredTypeRef(member)
		}
		return declarationTypes.Union(members...)
	case IntersectionType:
		members := make([]TypeRef, len(ty.Types))
		for i, member := range ty.Types {
			members[i] = DeclaredTypeRef(member)
		}
		return declarationTypes.Intersection(members...)
	default:
		return NewDeclaredType(ty.String()).(TypeRef)
	}
}
func NewDeclaredUnionType(types []Types) Types {
	members := make([]TypeRef, len(types))
	for i, ty := range types {
		members[i] = DeclaredTypeRef(ty)
	}
	return declarationTypes.Union(members...)
}
func NewDeclaredIntersectionType(types []Types) Types {
	members := make([]TypeRef, len(types))
	for i, ty := range types {
		members[i] = DeclaredTypeRef(ty)
	}
	return declarationTypes.Intersection(members...)
}
func NewDeclaredNullableType(ty Types) Types {
	return declarationTypes.Union(DeclaredTypeRef(ty), TypeNull)
}

// LegacyType is for tooling and extension migration, never for runtime checks.
func LegacyType(ty Types) Types {
	ref, ok := ty.(TypeRef)
	if !ok {
		return ty
	}
	if base, nullable := NullableDeclaredBase(ref); nullable {
		return NullableType{BaseType: LegacyType(base)}
	}
	if ref == TypeMixed {
		return Mixed{}
	}
	switch ref.Kind() {
	case TypeKindUnion, TypeKindIntersection:
		view := ref.Members()
		members := make([]Types, view.Len())
		for i := 0; i < view.Len(); i++ {
			members[i] = LegacyType(view.At(i))
		}
		if ref.Kind() == TypeKindUnion {
			return UnionType{Types: members}
		}
		return IntersectionType{Types: members}
	default:
		return NewBaseType(ref.String())
	}
}
