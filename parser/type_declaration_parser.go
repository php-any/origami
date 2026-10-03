package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/token"
)

func declarationStarts(p *Parser) bool {
	return isIdentOrTypeToken(p.current().Type()) || p.checkPositionIs(0, token.STATIC, token.TERNARY, token.LPAREN)
}

func declarationError(p *Parser, message string) data.Control {
	return data.NewCompileFatal(p.FromCurrentToken(), message)
}

// DNF allows intersections as parenthesized union members. An ampersand before
// a variable or ellipsis belongs to the parameter, not to the type expression.
func parseDeclaredType(p *Parser, position data.DeclarationPosition) (data.Types, data.Control) {
	nullable := p.checkPositionIs(0, token.TERNARY)
	if nullable {
		p.next()
	}
	var members []data.Types
	first, grouped, ctl := parseDeclarationTerm(p)
	if ctl != nil {
		return nil, ctl
	}
	members = append(members, first)
	union := p.checkPositionIs(0, token.BIT_OR)
	if nullable && (union || grouped || data.DeclaredTypeRef(first).Kind() == data.TypeKindIntersection) {
		return nil, declarationError(p, "Nullable shorthand must name a single type")
	}
	if grouped && !union {
		return nil, declarationError(p, "Parenthesized intersection must be part of a union")
	}
	if union && !grouped && data.DeclaredTypeRef(first).Kind() == data.TypeKindIntersection {
		return nil, declarationError(p, "Intersection in a union must be parenthesized")
	}
	for p.checkPositionIs(0, token.BIT_OR) {
		p.next()
		member, parenthesized, ctl := parseDeclarationTerm(p)
		if ctl != nil {
			return nil, ctl
		}
		if !parenthesized && data.DeclaredTypeRef(member).Kind() == data.TypeKindIntersection {
			return nil, declarationError(p, "Intersection in a union must be parenthesized")
		}
		members = append(members, member)
	}
	if len(members) > 1 {
		refs := declarationRefs(members)
		if err := data.ValidateDeclarationMembers(data.TypeKindUnion, refs); err != nil {
			return nil, declarationError(p, err.Error())
		}
		first = data.NewDeclaredUnionType(members)
	}
	ref := data.DeclaredTypeRef(first)
	if nullable {
		if ref == data.TypeNull || ref == data.TypeMixed || ref == data.TypeVoid || ref == data.TypeNever {
			return nil, declarationError(p, "Type cannot use nullable shorthand: "+ref.String())
		}
		first = data.NewDeclaredNullableType(first)
		ref = data.DeclaredTypeRef(first)
	}
	if err := data.ValidateDeclaration(ref, position, p.currentClass != ""); err != nil {
		return nil, declarationError(p, err.Error())
	}
	if p.currentClassParent != nil && *p.currentClassParent == "" && declarationHasParent(ref) {
		return nil, declarationError(p, "Cannot use parent when current class scope has no parent")
	}
	return first, nil
}

func declarationHasParent(ref data.TypeRef) bool {
	if ref == data.TypeParent {
		return true
	}
	members := ref.Members()
	for i := 0; i < members.Len(); i++ {
		if declarationHasParent(members.At(i)) {
			return true
		}
	}
	return false
}

func (p *Parser) enterClassDeclaration(name string, parent *string) func() {
	previous, previousParent := p.currentClass, p.currentClassParent
	p.currentClass, p.currentClassParent = name, parent
	return func() { p.currentClass, p.currentClassParent = previous, previousParent }
}

func declarationRefs(types []data.Types) []data.TypeRef {
	refs := make([]data.TypeRef, len(types))
	for i, ty := range types {
		refs[i] = data.DeclaredTypeRef(ty)
	}
	return refs
}

func parseDeclarationTerm(p *Parser) (data.Types, bool, data.Control) {
	grouped := p.checkPositionIs(0, token.LPAREN)
	if grouped {
		p.next()
	}
	first, ctl := parseDeclarationAtom(p)
	if ctl != nil {
		return nil, false, ctl
	}
	members := []data.Types{first}
	for p.checkPositionIs(0, token.BIT_AND) && (isIdentOrTypeToken(p.peek(1).Type()) || p.peek(1).Type() == token.STATIC) {
		p.next()
		next, ctl := parseDeclarationAtom(p)
		if ctl != nil {
			return nil, false, ctl
		}
		members = append(members, next)
	}
	if len(members) > 1 {
		if err := data.ValidateDeclarationMembers(data.TypeKindIntersection, declarationRefs(members)); err != nil {
			return nil, false, declarationError(p, err.Error())
		}
		first = data.NewDeclaredIntersectionType(members)
	}
	if grouped {
		if len(members) < 2 || !p.checkPositionIs(0, token.RPAREN) {
			return nil, false, declarationError(p, "Expected a parenthesized intersection type")
		}
		p.next()
	}
	return first, grouped, nil
}

func parseDeclarationAtom(p *Parser) (data.Types, data.Control) {
	if !isIdentOrTypeToken(p.current().Type()) && !p.checkPositionIs(0, token.STATIC) {
		return nil, declarationError(p, "Expected a type name")
	}
	name := p.current().Literal()
	if p.checkPositionIs(0, token.SELF, token.PARENT, token.STATIC) {
		p.next()
		return data.NewDeclaredType(name), nil
	}
	if p.peek(1).Type() == token.LT {
		return nil, declarationError(p, "Generic types are not PHP declarations")
	}
	p.next()
	if !data.ISBaseType(name) {
		if full, ok := p.findFullClassNameByNamespace(name); ok {
			name = full
		}
	}
	return data.NewDeclaredType(name), nil
}

func parseDeclaredReturn(p *Parser) (data.Types, data.Control) {
	if !p.checkPositionIs(0, token.COLON) {
		return nil, nil
	}
	p.next()
	ty, ctl := parseDeclaredType(p, data.DeclarationReturn)
	if ctl != nil {
		return nil, ctl
	}
	if p.checkPositionIs(0, token.COMMA) {
		return nil, declarationError(p, "Multiple return types are not PHP declarations")
	}
	return ty, nil
}
