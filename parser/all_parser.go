package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/token"
)

// StatementParser 表示语句解析器接口
type StatementParser interface {
	Parse() (data.GetValue, data.Control)
}

var parserRouter = map[token.TokenType]func(parser *Parser) StatementParser{

	token.IF:         NewIfParser,
	token.WHILE:      NewWhileParser,
	token.FOR:        NewForParser,
	token.FOREACH:    NewForeachParser,
	token.RETURN:     NewReturnParser,
	token.BREAK:      NewBreakParser,
	token.CONTINUE:   NewContinueParser,
	token.VAR:        NewVarParser,
	token.CONST:      NewConstParser,
	token.FUNC:       NewFunctionParser,
	token.CLASS:      NewClassParser,
	token.INTERFACE:  NewInterfaceParser,
	token.NAMESPACE:  NewNamespaceParser,
	token.USE:        NewUseParser,
	token.NEW:        NewNewParser,
	token.ECHO:       NewEchoParser,
	token.THIS:       NewThisParser,
	token.PARENT:     NewParentParser,
	token.SELF:       NewSelfParser,
	token.VARIABLE:   NewVariableParser,
	token.LPAREN:     NewLparenParser,
	token.LBRACKET:   NewLbracketParser,
	token.LBRACE:     NewLbraceParser,
	token.IDENTIFIER: NewIdentParser,
	token.BOOL:       NewBoolParser,
	token.DIR:        NewDirParser,
	token.FILE:       NewFileParser,
	token.LINE:       NewLineParser,
	token.TRY:        NewTryParser,
	token.THROW:      NewThrowParser,
	token.SPAWN:      NewSpawnParser,
	token.UNUSED:     NewUnusedParser,
	token.MATCH:      NewMatchParser,
	token.SWITCH:     NewSwitchParser,
	token.ARRAY:      NewArrayParser,
	token.GLOBALS:    NewGlobalsParser,
	token.ENV:        NewGlobalsParser,
	token.SESSION:    NewGlobalsParser,
	token.COOKIE:     NewGlobalsParser,
	token.POST:       NewGlobalsParser,
	token.GET:        NewGlobalsParser,
	token.REQUEST:    NewGlobalsParser,
	token.SERVER:     NewGlobalsParser,

	token.TERNARY:      NewNullableParser,
	token.AT:           NewAnnotationParser,
	token.JS_SERVER:    NewJsServerParser,
	token.STATIC:       NewStaticParser,
	token.DEFINE:       NewDefineParser,
	token.INCLUDE:      NewIncludeParser,
	token.INCLUDE_ONCE: NewIncludeParser,
	token.REQUIRE:      NewIncludeParser,
	token.REQUIRE_ONCE: NewIncludeParser,
}

func AddParse(t token.TokenType, parser func(parser *Parser) StatementParser) {
	parserRouter[t] = parser
}
