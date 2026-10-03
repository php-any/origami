package data

// StdClass is PHP's built-in dynamic object declaration. Keeping the empty
// declaration here lets native conversion create real objects without a VM.
type StdClass struct{}

func (s *StdClass) GetValue(ctx Context) (GetValue, Control) {
	return NewClassValue(s, ctx), nil
}
func (*StdClass) GetFrom() From                       { return nil }
func (*StdClass) GetName() string                     { return "stdClass" }
func (*StdClass) GetExtend() *string                  { return nil }
func (*StdClass) GetImplements() []string             { return nil }
func (*StdClass) GetProperty(string) (Property, bool) { return nil, false }
func (*StdClass) GetPropertyList() []Property         { return nil }
func (*StdClass) GetMethod(string) (Method, bool)     { return nil, false }
func (*StdClass) GetMethods() []Method                { return nil }
func (*StdClass) GetConstruct() Method                { return nil }

func NewStdClassValue(ctx Context) *ClassValue {
	var class ClassStmt = &StdClass{}
	if ctx != nil && ctx.GetVM() != nil {
		if registered, ok := ctx.GetVM().GetClass("stdClass"); ok {
			class = registered
		}
	}
	return NewClassValue(class, ctx)
}
