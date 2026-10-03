package data

import "testing"

type strictCompatibilityContext struct{ Context }

func (strictCompatibilityContext) StrictTypes() bool { return true }

func TestLegacyCompoundCoercionUsesCompactMembers(t *testing.T) {
	tests := []struct {
		name   string
		typeOf Types
		input  Value
		want   string
	}{
		{"union", UnionType{Types: []Types{TypeInt, TypeFloat}}, NewStringValue("2.5"), "2.5"},
		{"nullable", NullableType{BaseType: TypeString}, NewIntValue(42), "42"},
		{"nested", UnionType{Types: []Types{NullableType{BaseType: TypeFloat}, TypeBool}}, NewStringValue("2.5"), "2.5"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, accepted, ctl := PrepareTypedValueWithControl(test.typeOf, test.input)
			if ctl != nil || !accepted || value == nil || value.AsString() != test.want {
				t.Fatalf("conversion = %v, %v, %v; want %q", value, accepted, ctl, test.want)
			}
			compact, accepted, ctl := PrepareDeclaredValueInContext(DeclaredTypeRef(test.typeOf), test.input, nil)
			if ctl != nil || !accepted || ValueKindOf(value) != ValueKindOf(compact) {
				t.Fatal("compatibility conversion differs from the declaration checker")
			}
		})
	}
}

func TestLegacyUnionStrictFloatPromotion(t *testing.T) {
	typeOf := UnionType{Types: []Types{TypeFloat, TypeBool}}
	ctx := strictCompatibilityContext{}
	value, accepted, ctl := PrepareTypedValueInContext(typeOf, NewIntValue(42), ctx)
	if ctl != nil || !accepted || ValueKindOf(value) != ValueFloat || value.AsString() != "42" {
		t.Fatalf("strict int-to-float promotion = %v, %v, %v", value, accepted, ctl)
	}
	if _, accepted, ctl := PrepareTypedValueInContext(typeOf, NewStringValue("42"), ctx); ctl != nil || accepted {
		t.Fatal("strict compatibility conversion accepted a numeric string")
	}
}

func TestLegacyDeclarationKeywordIdentity(t *testing.T) {
	for name, expected := range map[string]TypeRef{
		"mixed": TypeMixed, "void": TypeVoid, "never": TypeNever,
		"false": TypeFalse, "true": TypeTrue, "null": TypeNull,
		"self": TypeSelf, "parent": TypeParent, "static": TypeStatic,
		"iterable": TypeIterable,
	} {
		if actual := DeclaredTypeRef(NewBaseType(name)); actual != expected {
			t.Errorf("legacy %s = %s; want %s", name, actual, expected)
		}
		if actual := DeclaredTypeRef(LegacyType(expected)); actual != expected {
			t.Errorf("legacy round trip %s = %s; want %s", name, actual, expected)
		}
	}
	if DeclaredTypeRef(NewGenericType("void", nil)) != TypeVoid || DeclaredTypeRef(NewGenericType("mixed", nil)) != TypeMixed {
		t.Fatal("generic compatibility constructor lost a declared keyword")
	}
}

func TestDeclarationTextPreservesGroupedMembers(t *testing.T) {
	a := DeclaredTypeRef(NewDeclaredType("TextContractA"))
	b := DeclaredTypeRef(NewDeclaredType("TextContractB"))
	c := DeclaredTypeRef(NewDeclaredType("TextContractC"))
	want := DeclaredTypeRef(NewDeclaredUnionType([]Types{NewDeclaredIntersectionType([]Types{a, b}), c, TypeNull}))
	for _, text := range []string{"(TextContractA&TextContractB)|TextContractC|null", " ( TextContractA & TextContractB ) | TextContractC | null ", "((TextContractA&TextContractB)|TextContractC|null)"} {
		for _, constructor := range []func(string) Types{NewDeclaredType, NewBaseType} {
			if got := DeclaredTypeRef(constructor(text)); got != want {
				t.Errorf("%q = %s; want %s", text, got, want)
			}
		}
	}
	if got := DeclaredTypeRef(NewBaseType("?TextContractA")); got != DeclaredTypeRef(NewDeclaredUnionType([]Types{a, TypeNull})) {
		t.Errorf("nullable constructor = %s", got)
	}
}

func TestToolingInferenceCannotBecomeRuntimeDeclaration(t *testing.T) {
	for _, analysis := range []Types{NewLspTypes(TypeInt), NewMultipleReturnType([]Types{TypeInt, TypeString})} {
		ref := DeclaredTypeRef(analysis)
		if ref != TypeNever || ref.Matches(NewIntValue(1), nil) || ref.Matches(NewArrayValue(nil), nil) {
			t.Fatal("tooling result entered runtime declaration matching")
		}
		if _, exists := any(analysis).(interface{ Is(Value) bool }); exists {
			t.Fatal("tooling result still exposes a runtime value predicate")
		}
	}
}
