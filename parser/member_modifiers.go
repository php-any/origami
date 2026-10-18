package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/token"
)

func (p *ClassParser) parseMemberModifiers() (modifier, setVisibility string, isAbstractMethod, isFinalMethod, isStatic, isReadonly bool, ctl data.Control) {
	isAbstractMethod, isFinalMethod, isStatic, isReadonly = false, false, false, false
	modifier = "public"
	setVisibility = ""
	for {
		switch p.current().Type() {
		case token.FINAL:
			isFinalMethod = true
			p.next()
		case token.ABSTRACT:
			isAbstractMethod = true
			p.next()
		case token.STATIC:
			isStatic = true
			p.next()
		case token.READONLY:
			isReadonly = true
			p.next()
		case token.PUBLIC, token.PROTECTED, token.PRIVATE:
			visibility := p.parseModifier()
			if p.checkPositionIs(0, token.LPAREN) && p.peek(1).Literal() == "set" {
				p.next()
				p.next()
				if ctl := p.nextAndCheck(token.RPAREN); ctl != nil {
					return "", "", false, false, false, false, ctl
				}
				setVisibility = visibility
			} else {
				modifier = visibility
			}
		default:
			goto memberModifiersDone
		}
	}
memberModifiersDone:
	if isAbstractMethod && isFinalMethod {
		return "", "", false, false, false, false, data.NewCompileFatal(p.newFrom(), "Cannot use the final modifier on an abstract method")
	}

	return
}
