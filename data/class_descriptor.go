package data

import (
	"sync"
	"sync/atomic"
)

type ClassID uint32
type ClassFlags uint8

type MethodFlags uint8

const (
	MethodFinal MethodFlags = 1 << iota
	MethodAbstract
)

func MethodDeclarationFlags(method Method) MethodFlags {
	var flags MethodFlags
	if metadata, ok := method.(interface{ DeclarationMethodFlags() MethodFlags }); ok {
		flags = metadata.DeclarationMethodFlags()
	}
	if metadata, ok := method.(AbstractMethodMarker); ok && metadata.IsAbstractMethod() {
		flags |= MethodAbstract
	}
	return flags
}

const (
	ClassInterface ClassFlags = 1 << iota
	ClassTrait
	ClassEnum
	ClassAbstract
	ClassFinal
	ClassReadonly
)

type ParameterDescriptor struct {
	name                            SymbolID
	declared                        TypeRef
	byReference, variadic, optional bool
}
type MethodDescriptor struct {
	name        SymbolID
	modifier    Modifier
	static      bool
	byReference bool
	flags       MethodFlags
	result      TypeRef
	parameters  []ParameterDescriptor
}
type PropertyDescriptor struct {
	name               SymbolID
	modifier           Modifier
	static             bool
	declared           TypeRef
	readonly, constant bool
	defaultValue       GetValue
}
type ClassDescriptor struct {
	id          ClassID
	name        SymbolID
	parent      ClassID
	interfaces  []ClassID
	methods     []MethodDescriptor
	properties  []PropertyDescriptor
	ancestors   Set[ClassID]
	flags       ClassFlags
	declaration *ClassDescriptor // immutable unlinked metadata used for this publication
}

func (d *ClassDescriptor) ID() ClassID                           { return d.id }
func (d *ClassDescriptor) Name() string                          { return Symbols.Name(d.name) }
func (d *ClassDescriptor) Parent() ClassID                       { return d.parent }
func (d *ClassDescriptor) Interfaces() Span[ClassID]             { return NewSpan(d.interfaces) }
func (d *ClassDescriptor) Methods() Span[MethodDescriptor]       { return NewSpan(d.methods) }
func (d *ClassDescriptor) Properties() Span[PropertyDescriptor]  { return NewSpan(d.properties) }
func (d *ClassDescriptor) Flags() ClassFlags                     { return d.flags }
func (d *ClassDescriptor) IsA(id ClassID) bool                   { return d.ancestors.Contains(id) }
func (m MethodDescriptor) Name() string                          { return Symbols.Name(m.name) }
func (m MethodDescriptor) ReturnType() TypeRef                   { return m.result }
func (m MethodDescriptor) Modifier() Modifier                    { return m.modifier }
func (m MethodDescriptor) IsStatic() bool                        { return m.static }
func (m MethodDescriptor) ReturnsByReference() bool              { return m.byReference }
func (m MethodDescriptor) Flags() MethodFlags                    { return m.flags }
func (m MethodDescriptor) Parameters() Span[ParameterDescriptor] { return NewSpan(m.parameters) }
func (p ParameterDescriptor) Name() string                       { return Symbols.Name(p.name) }
func (p ParameterDescriptor) Type() TypeRef                      { return p.declared }
func (p PropertyDescriptor) Name() string                        { return Symbols.Name(p.name) }
func (p PropertyDescriptor) Type() TypeRef                       { return p.declared }
func (p PropertyDescriptor) Modifier() Modifier                  { return p.modifier }
func (p PropertyDescriptor) IsStatic() bool                      { return p.static }
func (p PropertyDescriptor) IsReadonly() bool                    { return p.readonly }
func (p PropertyDescriptor) IsConstant() bool                    { return p.constant }
func (p PropertyDescriptor) DefaultValue() GetValue              { return p.defaultValue }
func (p ParameterDescriptor) IsByReference() bool                { return p.byReference }
func (p ParameterDescriptor) IsVariadic() bool                   { return p.variadic }
func (p ParameterDescriptor) IsOptional() bool                   { return p.optional || p.variadic }

type DescriptorSnapshot struct {
	descriptors IDMap[*ClassDescriptor]
	classes     IDMap[ClassStmt]
	interfaces  IDMap[InterfaceStmt]
	aliases     IDMap[ClassID]
	dependents  IDMap[Set[ClassID]]
	revision    uint64
}
type ClassRegistry struct {
	mu       sync.Mutex
	snapshot atomic.Pointer[DescriptorSnapshot]
}

func NewClassRegistry(parent *DescriptorSnapshot) *ClassRegistry {
	r := &ClassRegistry{}
	if parent == nil {
		parent = &DescriptorSnapshot{}
	}
	r.snapshot.Store(parent)
	return r
}
func (r *ClassRegistry) Snapshot() *DescriptorSnapshot { return r.snapshot.Load() }
func (s *DescriptorSnapshot) FindClass(name string) (ClassStmt, bool) {
	id, ok := Symbols.Lookup(name)
	if !ok {
		return nil, false
	}
	return s.classes.Get(uint32(id))
}
func (s *DescriptorSnapshot) FindInterface(name string) (InterfaceStmt, bool) {
	id, ok := Symbols.Lookup(name)
	if !ok {
		return nil, false
	}
	return s.interfaces.Get(uint32(id))
}
func descriptorMethods(methods []Method) []MethodDescriptor {
	out := make([]MethodDescriptor, len(methods))
	for i, method := range methods {
		if method == nil {
			continue
		}
		m := MethodDescriptor{name: Symbols.Intern(method.GetName()), modifier: method.GetModifier(), static: method.GetIsStatic(), result: DeclaredTypeRef(method.GetReturnType()), flags: MethodDeclarationFlags(method)}
		if ref, ok := method.(interface{ ReturnsByReference() bool }); ok {
			m.byReference = ref.ReturnsByReference()
		}
		for _, parameter := range method.GetParams() {
			if p, ok := parameter.(Parameter); ok {
				entry := ParameterDescriptor{name: Symbols.InternProperty(p.GetName()), declared: DeclaredTypeRef(p.GetType()), optional: p.GetDefaultValue() != nil}
				if ref, ok := parameter.(interface{ IsReferenceParameter() bool }); ok {
					entry.byReference = ref.IsReferenceParameter()
				}
				if rest, ok := parameter.(interface{ IsVariadicParameter() bool }); ok {
					entry.variadic = rest.IsVariadicParameter()
				}
				m.parameters = append(m.parameters, entry)
			}
		}
		out[i] = m
	}
	return out
}

// DescribeClassDeclaration constructs immutable signature/default metadata.
// Request-independent declarations may retain it while each registry links
// its own ancestors, aliases and declaration identities.
func DescribeClassDeclaration(class ClassStmt) *ClassDescriptor {
	canonical := originalClass(class)
	id := ClassID(Symbols.Intern(canonical.GetName()))
	d := &ClassDescriptor{id: id, name: SymbolID(id)}
	if metadata, ok := canonical.(interface{ DeclarationFlags() ClassFlags }); ok {
		d.flags = metadata.DeclarationFlags()
	}
	if parent := canonical.GetExtend(); parent != nil {
		d.parent = ClassID(Symbols.Intern(*parent))
	}
	for _, iface := range canonical.GetImplements() {
		d.interfaces = append(d.interfaces, ClassID(Symbols.Intern(iface)))
	}
	d.methods = descriptorMethods(canonical.GetMethods())
	properties := canonical.GetPropertyList()
	if metadata, ok := canonical.(interface{ StaticPropertyDeclarations() []Property }); ok {
		properties = append(properties, metadata.StaticPropertyDeclarations()...)
	}
	for _, property := range properties {
		entry := PropertyDescriptor{name: Symbols.InternProperty(property.GetName()), modifier: property.GetModifier(), static: property.GetIsStatic(), declared: DeclaredTypeRef(property.GetType()), defaultValue: property.GetDefaultValue()}
		if metadata, ok := property.(interface{ IsReadonlyProperty() bool }); ok {
			entry.readonly = metadata.IsReadonlyProperty()
		}
		if metadata, ok := property.(interface{ IsClassConstant() bool }); ok {
			entry.constant = metadata.IsClassConstant()
		}
		d.properties = append(d.properties, entry)
	}
	return d
}

func (r *ClassRegistry) PublishClass(class ClassStmt) {
	canonical := originalClass(class)
	var d *ClassDescriptor
	if source, ok := canonical.(interface{ DeclarationDescriptor() *ClassDescriptor }); ok {
		d = source.DeclarationDescriptor()
	} else {
		d = DescribeClassDeclaration(canonical)
	}
	name := d.name
	if class != canonical {
		name = Symbols.Intern(class.GetName())
	}
	id := d.id
	r.mu.Lock()
	defer r.mu.Unlock()
	s := *r.snapshot.Load()
	if existing, found := s.classes.Get(uint32(name)); found && existing == class {
		if published, found := s.descriptors.Get(uint32(id)); found && (published == d || published.declaration == d) {
			return
		}
	}
	s.classes = s.classes.With(uint32(name), class)
	if previous, exists := s.aliases.Get(uint32(name)); ClassID(name) != id || exists && previous != id {
		s.aliases = s.aliases.With(uint32(name), id)
	}
	s.revision++
	r.linkPublished(&s, id, d)
	if ClassID(name) != id {
		r.linkPublished(&s, ClassID(name), nil)
	}
	r.snapshot.Store(&s)
}
func (r *ClassRegistry) PublishInterface(iface InterfaceStmt) {
	name := Symbols.Intern(iface.GetName())
	id := ClassID(name)
	d := &ClassDescriptor{id: id, name: name, flags: ClassInterface, methods: descriptorMethods(iface.GetMethods())}
	for _, parent := range iface.GetExtends() {
		d.interfaces = append(d.interfaces, ClassID(Symbols.Intern(parent)))
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	s := *r.snapshot.Load()
	s.interfaces = s.interfaces.With(uint32(name), iface)
	if previous, exists := s.aliases.Get(uint32(name)); exists && previous != id {
		s.aliases = s.aliases.With(uint32(name), id)
	}
	s.revision++
	r.linkPublished(&s, id, d)
	r.snapshot.Store(&s)
}
func (r *ClassRegistry) Descriptor(name SymbolID) (*ClassDescriptor, bool) {
	s := r.snapshot.Load()
	id, ok := s.aliases.Get(uint32(name))
	if !ok {
		id = ClassID(name)
	}
	return s.descriptors.Get(uint32(id))
}
func linkDescriptor(s *DescriptorSnapshot, id ClassID, raw *ClassDescriptor) *ClassDescriptor {
	if raw == nil {
		raw, _ = s.descriptors.Get(uint32(id))
	}
	if raw == nil {
		return nil
	}
	seen := make(map[ClassID]bool)
	var ids []ClassID
	var visit func(ClassID)
	visit = func(current ClassID) {
		if alias, ok := s.aliases.Get(uint32(current)); ok {
			current = alias
		}
		if current == 0 || seen[current] {
			return
		}
		seen[current] = true
		ids = append(ids, current)
		d, ok := s.descriptors.Get(uint32(current))
		if current == id {
			d, ok = raw, true
		}
		if ok {
			visit(d.parent)
			for _, iface := range d.interfaces {
				visit(iface)
			}
		}
		// Internal throws and iterator markers can have parents without a
		// registered declaration. Include those PHP relationships as metadata.
		if parent := internalTypeParents[Symbols.Name(SymbolID(current))]; parent != "" {
			visit(ClassID(Symbols.Intern(parent)))
		}
	}
	visit(id)
	clone := *raw
	if clone.declaration == nil {
		clone.declaration = raw
	}
	clone.ancestors = NewSet(ids...)
	return &clone
}

// All linking and invalidation occur at declaration publication. A type
// predicate only reads a published descriptor and never interns or autoloads.
func (r *ClassRegistry) linkPublished(s *DescriptorSnapshot, id ClassID, declaration *ClassDescriptor) {
	d := declaration
	if d == nil {
		d, _ = s.descriptors.Get(uint32(id))
	}
	if d != nil {
		parents := append([]ClassID{d.parent}, d.interfaces...)
		for _, parent := range parents {
			if parent == 0 {
				continue
			}
			previous, _ := s.dependents.Get(uint32(parent))
			if !previous.Contains(id) {
				values := append(previous.View().AppendTo(nil), id)
				s.dependents = s.dependents.With(uint32(parent), NewSet(values...))
			}
		}
	}
	seen := map[ClassID]bool{}
	var refresh func(ClassID)
	refresh = func(current ClassID) {
		if seen[current] {
			return
		}
		seen[current] = true
		var raw *ClassDescriptor
		if current == id {
			raw = declaration
		}
		if linked := linkDescriptor(s, current, raw); linked != nil {
			s.descriptors = s.descriptors.With(uint32(current), linked)
		}
		children, _ := s.dependents.Get(uint32(current))
		for i := 0; i < children.View().Len(); i++ {
			refresh(children.View().At(i))
		}
	}
	refresh(id)
}
func (r *ClassRegistry) IsA(source ClassStmt, target SymbolID) (bool, bool) {
	id, ok := Symbols.Lookup(source.GetName())
	if !ok {
		return false, false
	}
	d, ok := r.Descriptor(id)
	if !ok {
		return false, false
	}
	if alias, ok := r.snapshot.Load().aliases.Get(uint32(target)); ok {
		target = SymbolID(alias)
	}
	return d.IsA(ClassID(target)), true
}

type DescriptorProvider interface{ ClassRegistry() *ClassRegistry }
