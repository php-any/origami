package data

import "fmt"

type DeclarationPosition uint8

const (
	DeclarationParameter DeclarationPosition = iota
	DeclarationProperty
	DeclarationReturn
)

// Declaration validation runs before arena normalization. Otherwise duplicate
// and redundant source members disappear during interning.
func ValidateDeclarationMembers(kind TypeKind, members []TypeRef) error {
	seen := make(map[TypeRef]bool, len(members))
	for _, member := range members {
		if seen[member] {
			return fmt.Errorf("Duplicate type %s is redundant", member.String())
		}
		seen[member] = true
		if kind == TypeKindIntersection && member.Kind() != TypeKindNominal {
			return fmt.Errorf("Type %s cannot be part of an intersection type", member.String())
		}
		if kind == TypeKindUnion && (member == TypeMixed || member == TypeVoid || member == TypeNever) {
			return fmt.Errorf("Type %s can only be used as a standalone type", member.String())
		}
	}
	if kind != TypeKindUnion {
		return nil
	}
	if seen[TypeBool] && (seen[TypeFalse] || seen[TypeTrue]) || seen[TypeFalse] && seen[TypeTrue] {
		return fmt.Errorf("Redundant boolean union type")
	}
	for i, member := range members {
		if seen[TypeObject] && (member.Kind() == TypeKindNominal || member.Kind() == TypeKindIntersection || member.RequiresClassScope()) {
			return fmt.Errorf("Type %s is redundant with object", member.String())
		}
		if seen[TypeIterable] && (member == TypeArray || member.Kind() == TypeKindNominal && nominalKey(member.String()) == "traversable") {
			return fmt.Errorf("Type %s is redundant with iterable", member.String())
		}
		for _, other := range members[i+1:] {
			if declarationContains(member, other) || declarationContains(other, member) {
				return fmt.Errorf("Redundant union types %s and %s", member.String(), other.String())
			}
		}
	}
	return nil
}

// An intersection with more conjuncts is covered by the less restrictive
// member. This uses declaration identities only and never loads classes.
func declarationContains(broad, narrow TypeRef) bool {
	if narrow.Kind() != TypeKindIntersection {
		return false
	}
	parts := narrow.Members()
	if broad.Kind() != TypeKindIntersection {
		for i := 0; i < parts.Len(); i++ {
			if parts.At(i) == broad {
				return true
			}
		}
		return false
	}
	broadParts := broad.Members()
	for i := 0; i < broadParts.Len(); i++ {
		found := false
		for j := 0; j < parts.Len(); j++ {
			if broadParts.At(i) == parts.At(j) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func ValidateDeclaration(ref TypeRef, position DeclarationPosition, classScope bool) error {
	if ref == TypeInvalid {
		return fmt.Errorf("Missing type declaration")
	}
	if ref == TypeVoid || ref == TypeNever {
		if position != DeclarationReturn {
			return fmt.Errorf("%s can only be used as a return type", ref.String())
		}
	}
	if ref == TypeCallable && position == DeclarationProperty {
		return fmt.Errorf("Property cannot have type callable")
	}
	if ref.RequiresClassScope() && !classScope {
		return fmt.Errorf("Cannot use %s when no class scope is active", ref.String())
	}
	if ref == TypeStatic && position != DeclarationReturn {
		return fmt.Errorf("static can only be used as a return type")
	}
	if ref.Kind() == TypeKindUnion || ref.Kind() == TypeKindIntersection {
		members := ref.Members()
		for i := 0; i < members.Len(); i++ {
			if err := ValidateDeclaration(members.At(i), position, classScope); err != nil {
				return err
			}
		}
	}
	return nil
}
