package data

// RequestObjectScope preserves object identities and references across retained
// static values, closures and nested property/array graphs. It is request-owned.
type RequestObjectScope struct {
	Context  Context
	Objects  map[*ObjectValue]*ClassValue
	arrays   *ArrayOverlayScope
	captures *RequestCaptureScope
	native   map[any]any
	failure  any
}

func NewRequestObjectScope(ctx Context) *RequestObjectScope {
	s := &RequestObjectScope{Context: ctx, Objects: make(map[*ObjectValue]*ClassValue), arrays: NewArrayOverlayScope(), native: make(map[any]any)}
	s.captures = NewRequestCaptureScope(ctx, s.Objects)
	s.captures.ScopeObject = s.Object
	s.captures.ScopeArray = s.arrays.Array
	s.captures.ScopeSlot = s.BindSlot
	s.captures.ScopeValue = s.Bind
	s.arrays.bindSlot = s.BindSlot
	s.captures.slots = s.arrays.refs
	s.arrays.bindValue = s.Bind
	return s
}

func (s *RequestObjectScope) BindSlot(source *ZVal) *ZVal {
	if s.failure != nil {
		panic(s.failure)
	}
	return s.arrays.copySlot(source, func(source *ReferenceGuard) *ReferenceGuard {
		if source == nil {
			return nil
		}
		guard := &ReferenceGuard{properties: make([]propertyConstraint, len(source.properties))}
		for i, constraint := range source.properties {
			if constraint.object != nil {
				object := s.Object(constraint.object)
				constraint.owner, constraint.object = object.ObjectValue, object
			}
			guard.properties[i] = constraint
		}
		return guard
	})
}

func (s *RequestObjectScope) Object(source *ClassValue) *ClassValue {
	if s.failure != nil {
		panic(s.failure)
	}
	if source == nil || source.ObjectValue == nil {
		return source
	}
	if target := s.Objects[source.ObjectValue]; target != nil {
		return target
	}
	target := source.CloneRequestScoped(s.Context)
	s.Objects[source.ObjectValue] = target
	s.Objects[target.ObjectValue] = target
	if store, ok := target.property.(*chainStore); ok {
		store.scope = s
	}
	if native := source.GetSource(); native != nil {
		target.InstanceSource = s.BindNativeState(native)
	}
	return target
}

// Remember lets a host publish the final identity after applying a native
// instance policy; subsequent graph reads reuse that exact request object.
func (s *RequestObjectScope) Remember(source, target *ClassValue) {
	if source == nil || target == nil {
		return
	}
	s.Objects[source.ObjectValue] = target
	s.Objects[target.ObjectValue] = target
	if store, ok := target.property.(*chainStore); ok {
		store.scope = s
	}
}

func (s *RequestObjectScope) Bind(value Value) Value {
	if s.failure != nil {
		panic(s.failure)
	}
	switch v := value.(type) {
	case *ClassValue:
		return s.Object(v)
	case *ThisValue:
		return NewThisValue(s.Object(v.ClassValue))
	case *ArrayValue:
		return s.arrays.Array(v)
	case *AnyValue:
		return NewAnyValue(s.BindNativeState(v.Value))
	case *FuncValue, *BoundFuncValue, *ArraySlotRef:
		return s.captures.Bind(v)
	case interface {
		BindRequestValue(*RequestObjectScope) Value
	}:
		return v.BindRequestValue(s)
	default:
		return value
	}
}

// RequestScopeProvider is used at setup and retained-value boundaries only.
type RequestScopeProvider interface{ RequestObjectScope() *RequestObjectScope }

func BindRequestScopeContext(ctx Context) {
	if ctx == nil || RequestGoid == nil || requestStaticActive.Load() == 0 {
		return
	}
	raw, ok := requestStatics.Load(RequestGoid())
	if !ok {
		return
	}
	state := raw.(*requestStaticMap)
	if state.scope != nil {
		return
	}
	if host, ok := ctx.GetVM().(RequestScopeProvider); ok {
		state.scope = host.RequestObjectScope()
	}
}
