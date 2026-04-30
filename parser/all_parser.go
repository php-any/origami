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
	token.DO:         NewDoWhileParser,
	token.FOR:        NewForParser,
	token.FOREACH:    NewForeachParser,
	token.RETURN:     NewReturnParser,
	token.BREAK:      NewBreakParser,
	token.CONTINUE:   NewContinueParser,
	token.GOTO:       NewGotoParser,
	token.VAR:        NewVarParser,
	token.CONST:      NewConstParser,
	token.FUNC:       NewFunctionParser,
	token.CLASS:      NewClassParser,
	token.ABSTRACT:   NewAbstractClassParser,
	token.ENUM:       NewEnumParser,
	token.INTERFACE:  NewInterfaceParser,
	token.TRAIT:      NewTraitParser,
	token.NAMESPACE:  NewNamespaceParser,
	token.USE:        NewUseParser,
	token.NEW:        NewNewParser,
	token.CLONE:      NewCloneParser,
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
	token.ISSET:      NewIssetParser,
	token.UNSET:      NewUnsetParser,
	token.COMPACT:    NewCompactParser,
	token.MATCH:      NewMatchParser,
	token.SWITCH:     NewSwitchParser,
	token.ARRAY:      NewArrayParser,

	token.TERNARY:       NewNullableParser,
	token.AT:            NewAnnotationParser,
	token.HASH:          NewAnnotationParser, // #[] 格式的注解也使用 AnnotationParser
	token.JS_SERVER:     NewJsServerParser,
	token.STATIC:        NewStaticParser,
	token.FN:            NewFnParser,
	token.DECLARE:       NewDeclareParser,
	token.INCLUDE:       NewIncludeParser,
	token.INCLUDE_ONCE:  NewIncludeParser,
	token.REQUIRE:       NewIncludeParser,
	token.REQUIRE_ONCE:  NewIncludeParser,
	token.FINAL:         NewFinalParser,
	token.YIELD:         NewYieldParser,
	token.FUNC_GET_ARGS: NewFuncGetArgsParser,
	token.FUNC_NUM_ARGS: NewFuncNumArgsParser,
	token.GLOBAL:        NewGlobalParser,
}

func AddParse(t token.TokenType, parser func(parser *Parser) StatementParser) {
	parserRouter[t] = parser
}
