package php

import (
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PHP token constants — unique arbitrary values used internally.
// These are NOT PHP 8 canonical values; only internal consistency matters
// since PHP code references them by constant name, not by numeric value.
const (
	T_REQUIRE_ONCE         = 258
	T_LOGICAL_AND          = 259
	T_LOGICAL_OR           = 260
	T_LOGICAL_XOR          = 261
	T_CLONE                = 262
	T_ELSEIF               = 263
	T_ELSE                 = 264
	T_ENDIF                = 265
	T_ECHO                 = 266
	T_DO                   = 267
	T_WHILE                = 268
	T_ENDWHILE             = 269
	T_FOR                  = 270
	T_ENDFOR               = 271
	T_FOREACH              = 272
	T_ENDFOREACH           = 273
	T_DECLARE              = 274
	T_ENDDECLARE           = 275
	T_AS                   = 276
	T_SWITCH               = 277
	T_ENDSWITCH            = 278
	T_CASE                 = 279
	T_DEFAULT              = 280
	T_BREAK                = 281
	T_CONTINUE             = 282
	T_GOTO                 = 283
	T_FUNCTION             = 284
	T_CONST                = 285
	T_RETURN               = 286
	T_YIELD                = 287
	T_TRY                  = 288
	T_CATCH                = 289
	T_FINALLY              = 290
	T_THROW                = 291
	T_IF                   = 292
	T_INSTANCEOF           = 293
	T_NEW                  = 294
	T_EXIT                 = 295
	T_EMPTY                = 296
	T_EVAL                 = 297
	T_INCLUDE              = 298
	T_INCLUDE_ONCE         = 299
	T_REQUIRE              = 300
	T_USE                  = 301
	T_GLOBAL               = 302
	T_ISSET                = 303
	T_UNSET                = 304
	T_LIST                 = 305
	T_ARRAY                = 306
	T_PRINT                = 307
	T_NAMESPACE            = 308
	T_OBJECT_OPERATOR      = 309
	T_PAAMAYIM_NEKUDOTAYIM = 310
	T_DOUBLE_ARROW         = 311
	T_DOUBLE_COLON         = 312
	T_NS_SEPARATOR         = 313
	T_ELLIPSIS             = 314
	T_COALESCE             = 315
	T_SPACESHIP            = 316
	T_POW                  = 317
	T_POW_EQUAL            = 318

	T_INLINE_HTML              = 319
	T_OPEN_TAG                 = 320
	T_OPEN_TAG_WITH_ECHO       = 321
	T_CLOSE_TAG                = 322
	T_WHITESPACE               = 323
	T_COMMENT                  = 324
	T_DOC_COMMENT              = 325
	T_STRING                   = 326
	T_VARIABLE                 = 327
	T_LNUMBER                  = 328
	T_DNUMBER                  = 329
	T_NUM_STRING               = 330
	T_CONSTANT_ENCAPSED_STRING = 331
	T_ENCAPSED_AND_WHITESPACE  = 332
	T_CHARACTER                = 333
	T_BAD_CHARACTER            = 334
	T_ABSTRACT                 = 335
	T_FINAL                    = 336
	T_PRIVATE                  = 337
	T_PROTECTED                = 338
	T_PUBLIC                   = 339
	T_STATIC                   = 340
	T_TRAIT                    = 341
	T_INTERFACE                = 342
	T_CLASS                    = 343
	T_CALLABLE                 = 344
	T_EXTENDS                  = 345
	T_IMPLEMENTS               = 346
	T_VAR                      = 347
	T_READONLY                 = 348
	T_ENUM                     = 349
	T_NAME_FULLY_QUALIFIED     = 350
	T_NAME_RELATIVE            = 351
	T_NAME_QUALIFIED           = 352
	T_MATCH                    = 353
	T_ATTRIBUTE                = 354
	T_NULLSAFE_OBJECT_OPERATOR = 355
	T_FN                       = 356

	T_TRUE  = 357
	T_FALSE = 358
	T_NULL  = 359
)

// keywordMap maps PHP keywords to their token constants
var keywordMap = map[string]int{
	"abstract":     T_ABSTRACT,
	"and":          T_LOGICAL_AND,
	"array":        T_ARRAY,
	"as":           T_AS,
	"break":        T_BREAK,
	"callable":     T_CALLABLE,
	"case":         T_CASE,
	"catch":        T_CATCH,
	"class":        T_CLASS,
	"clone":        T_CLONE,
	"const":        T_CONST,
	"continue":     T_CONTINUE,
	"declare":      T_DECLARE,
	"default":      T_DEFAULT,
	"die":          T_EXIT,
	"do":           T_DO,
	"echo":         T_ECHO,
	"else":         T_ELSE,
	"elseif":       T_ELSEIF,
	"empty":        T_EMPTY,
	"enddeclare":   T_ENDDECLARE,
	"endfor":       T_ENDFOR,
	"endforeach":   T_ENDFOREACH,
	"endif":        T_ENDIF,
	"endswitch":    T_ENDSWITCH,
	"endwhile":     T_ENDWHILE,
	"enum":         T_ENUM,
	"eval":         T_EVAL,
	"exit":         T_EXIT,
	"extends":      T_EXTENDS,
	"final":        T_FINAL,
	"finally":      T_FINALLY,
	"fn":           T_FN,
	"for":          T_FOR,
	"foreach":      T_FOREACH,
	"function":     T_FUNCTION,
	"global":       T_GLOBAL,
	"goto":         T_GOTO,
	"if":           T_IF,
	"implements":   T_IMPLEMENTS,
	"include":      T_INCLUDE,
	"include_once": T_INCLUDE_ONCE,
	"instanceof":   T_INSTANCEOF,
	"interface":    T_INTERFACE,
	"isset":        T_ISSET,
	"list":         T_LIST,
	"match":        T_MATCH,
	"namespace":    T_NAMESPACE,
	"new":          T_NEW,
	"or":           T_LOGICAL_OR,
	"print":        T_PRINT,
	"private":      T_PRIVATE,
	"protected":    T_PROTECTED,
	"public":       T_PUBLIC,
	"readonly":     T_READONLY,
	"require":      T_REQUIRE,
	"require_once": T_REQUIRE_ONCE,
	"return":       T_RETURN,
	"static":       T_STATIC,
	"switch":       T_SWITCH,
	"throw":        T_THROW,
	"trait":        T_TRAIT,
	"try":          T_TRY,
	"unset":        T_UNSET,
	"use":          T_USE,
	"var":          T_VAR,
	"while":        T_WHILE,
	"xor":          T_LOGICAL_XOR,
	"yield":        T_YIELD,

	// Soft type names (tokenized as T_STRING in PHP < 8, but we treat as T_STRING)
	"int":      T_STRING,
	"float":    T_STRING,
	"bool":     T_STRING,
	"string":   T_STRING,
	"void":     T_STRING,
	"never":    T_STRING,
	"mixed":    T_STRING,
	"iterable": T_STRING,
	"self":     T_STRING,
	"parent":   T_STRING,
	"true":     T_TRUE,
	"false":    T_FALSE,
	"null":     T_NULL,
}

// TokenGetAllFunction implements PHP's token_get_all()
type TokenGetAllFunction struct{}

func NewTokenGetAllFunction() data.FuncStmt {
	return &TokenGetAllFunction{}
}

func (f *TokenGetAllFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	source, _ := ctx.GetIndexValue(0)
	if source == nil {
		return data.NewArrayValue(nil), nil
	}

	srcStr := source.AsString()
	tokens := tokenizePhp(srcStr)
	return data.NewArrayValue(tokens), nil
}

func (f *TokenGetAllFunction) GetName() string            { return "token_get_all" }
func (f *TokenGetAllFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *TokenGetAllFunction) GetIsStatic() bool          { return false }
func (f *TokenGetAllFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "source", 0, nil, data.String{}),
	}
}
func (f *TokenGetAllFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "source", 0, nil),
	}
}
func (f *TokenGetAllFunction) GetReturnType() data.Types { return data.NewBaseType("array") }

// tokenizePhp tokenizes PHP source code, matching PHP's token_get_all() behavior
func tokenizePhp(source string) []data.Value {
	var tokens []data.Value
	src := []rune(source)
	pos := 0
	line := 1
	length := len(src)

	// Helper to create a PHP token array [tokenId, content, line]
	makeToken := func(id int, content string) data.Value {
		return data.NewArrayValue([]data.Value{
			data.NewIntValue(id),
			data.NewStringValue(content),
			data.NewIntValue(line),
		})
	}

	// Helper to add a simple character token (single char like ';', '(', etc.)
	addCharToken := func(ch rune) {
		tokens = append(tokens, data.NewStringValue(string(ch)))
	}

	// Helper to add a PHP token array
	addToken := func(id int, content string) {
		if content != "" {
			tokens = append(tokens, makeToken(id, content))
		}
	}

	// Check if remaining source starts with a string
	startsWith := func(s string) bool {
		r := []rune(s)
		if pos+len(r) > length {
			return false
		}
		for i, c := range r {
			if src[pos+i] != c {
				return false
			}
		}
		return true
	}

	// Skip whitespace and return the whitespace content
	scanWhitespace := func() string {
		start := pos
		for pos < length && (src[pos] == ' ' || src[pos] == '\t' || src[pos] == '\n' || src[pos] == '\r') {
			if src[pos] == '\n' {
				line++
			}
			pos++
		}
		return string(src[start:pos])
	}

	// isIdentChar returns true if the rune is valid in an identifier
	isIdentChar := func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
	}

	// Main tokenizer for PHP code inside <?php ... ?>
	tokenizePhpCode := func() {
		for pos < length {
			// Check for ?> closing tag (but not ??>)
			if pos+1 < length && src[pos] == '?' && src[pos+1] == '>' {
				return
			}

			ch := src[pos]

			// Whitespace
			if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
				content := scanWhitespace()
				addToken(T_WHITESPACE, content)
				continue
			}

			// Line comment: // or #
			if (ch == '/' && pos+1 < length && src[pos+1] == '/') ||
				(ch == '#') {
				start := pos
				for pos < length && src[pos] != '\n' {
					pos++
				}
				addToken(T_COMMENT, string(src[start:pos]))
				continue
			}

			// Block comment: /* ... */ or /** ... */
			if ch == '/' && pos+1 < length && src[pos+1] == '*' {
				start := pos
				isDoc := (pos+2 < length && src[pos+2] == '*')
				pos += 2 // skip /*
				for pos+1 < length {
					if src[pos] == '*' && src[pos+1] == '/' {
						pos += 2
						break
					}
					if src[pos] == '\n' {
						line++
					}
					pos++
				}
				content := string(src[start:pos])
				if isDoc {
					addToken(T_DOC_COMMENT, content)
				} else {
					addToken(T_COMMENT, content)
				}
				continue
			}

			// Single-quoted string: '...'
			if ch == '\'' {
				start := pos
				pos++ // skip opening '
				for pos < length {
					if src[pos] == '\'' {
						if pos+1 < length && src[pos+1] == '\'' {
							// escaped single quote ''
							pos += 2
							continue
						}
						pos++ // skip closing '
						break
					}
					if src[pos] == '\\' {
						pos++ // skip backslash (escaped char)
						if pos < length {
							pos++
						}
						continue
					}
					if src[pos] == '\n' {
						line++
					}
					pos++
				}
				addToken(T_CONSTANT_ENCAPSED_STRING, string(src[start:pos]))
				continue
			}

			// Double-quoted string: "..."
			if ch == '"' {
				start := pos
				pos++ // skip opening "
				hasInterpolation := false
				for pos < length {
					if src[pos] == '"' {
						pos++ // skip closing "
						break
					}
					if src[pos] == '\\' {
						pos++ // skip backslash
						if pos < length {
							pos++
						}
						continue
					}
					if src[pos] == '$' {
						hasInterpolation = true
					}
					if src[pos] == '\n' {
						line++
					}
					pos++
				}
				if hasInterpolation {
					addToken(T_ENCAPSED_AND_WHITESPACE, string(src[start:pos]))
				} else {
					addToken(T_CONSTANT_ENCAPSED_STRING, string(src[start:pos]))
				}
				continue
			}

			// Variable: $identifier or ${...}
			if ch == '$' {
				start := pos
				pos++
				if pos < length && src[pos] == '{' {
					// ${...}
					pos++ // skip {
					depth := 1
					for pos < length && depth > 0 {
						if src[pos] == '{' {
							depth++
						} else if src[pos] == '}' {
							depth--
						}
						if src[pos] == '\n' {
							line++
						}
						pos++
					}
				} else if pos < length && (unicode.IsLetter(src[pos]) || src[pos] == '_') {
					// $identifier
					for pos < length && isIdentChar(src[pos]) {
						pos++
					}
				} else {
					// Just $ followed by non-identifier - include $ and next char
					if pos < length {
						pos++
					}
				}
				addToken(T_VARIABLE, string(src[start:pos]))
				continue
			}

			// Numeric literals
			if unicode.IsDigit(ch) || (ch == '.' && pos+1 < length && unicode.IsDigit(src[pos+1])) {
				start := pos
				isFloat := false

				// Check for hex, octal, binary
				if ch == '0' && pos+1 < length {
					next := src[pos+1]
					if next == 'x' || next == 'X' {
						pos += 2
						for pos < length && isIdentChar(src[pos]) {
							pos++
						}
						addToken(T_LNUMBER, string(src[start:pos]))
						continue
					} else if next == 'o' || next == 'O' {
						pos += 2
						for pos < length && unicode.IsDigit(src[pos]) {
							pos++
						}
						addToken(T_LNUMBER, string(src[start:pos]))
						continue
					} else if next == 'b' || next == 'B' {
						pos += 2
						for pos < length && (src[pos] == '0' || src[pos] == '1') {
							pos++
						}
						addToken(T_LNUMBER, string(src[start:pos]))
						continue
					}
				}

				// Decimal number
				for pos < length && unicode.IsDigit(src[pos]) {
					pos++
				}
				if pos < length && src[pos] == '.' {
					isFloat = true
					pos++
					for pos < length && unicode.IsDigit(src[pos]) {
						pos++
					}
				}
				if pos < length && (src[pos] == 'e' || src[pos] == 'E') {
					isFloat = true
					pos++
					if pos < length && (src[pos] == '+' || src[pos] == '-') {
						pos++
					}
					for pos < length && unicode.IsDigit(src[pos]) {
						pos++
					}
				}
				if isFloat {
					addToken(T_DNUMBER, string(src[start:pos]))
				} else {
					addToken(T_LNUMBER, string(src[start:pos]))
				}
				continue
			}

			// Identifiers and keywords
			if unicode.IsLetter(ch) || ch == '_' {
				start := pos
				for pos < length && isIdentChar(src[pos]) {
					pos++
				}
				word := string(src[start:pos])
				if tokId, ok := keywordMap[strings.ToLower(word)]; ok {
					addToken(tokId, word)
				} else {
					addToken(T_STRING, word)
				}
				continue
			}

			// Multi-character operators (sorted by length descending)
			if startsWith("?->") {
				addToken(T_NULLSAFE_OBJECT_OPERATOR, "?->")
				pos += 3
				continue
			}
			if startsWith("**=") {
				addToken(T_POW_EQUAL, "**=")
				pos += 3
				continue
			}
			if startsWith("<=>") {
				addToken(T_SPACESHIP, "<=>")
				pos += 3
				continue
			}
			if startsWith("...") {
				addToken(T_ELLIPSIS, "...")
				pos += 3
				continue
			}
			if startsWith("??=") {
				addCharToken('?')
				addCharToken('?')
				addCharToken('=')
				pos += 3
				continue
			}
			if startsWith("===") {
				addCharToken('=')
				addCharToken('=')
				addCharToken('=')
				pos += 3
				continue
			}
			if startsWith("!==") {
				addCharToken('!')
				addCharToken('=')
				addCharToken('=')
				pos += 3
				continue
			}
			if startsWith("->") {
				addToken(T_OBJECT_OPERATOR, "->")
				pos += 2
				continue
			}
			if startsWith("::") {
				addToken(T_PAAMAYIM_NEKUDOTAYIM, "::")
				pos += 2
				continue
			}
			if startsWith("=>") {
				addToken(T_DOUBLE_ARROW, "=>")
				pos += 2
				continue
			}
			if startsWith("??") {
				addToken(T_COALESCE, "??")
				pos += 2
				continue
			}
			if startsWith("**") {
				addToken(T_POW, "**")
				pos += 2
				continue
			}
			if startsWith("==") {
				addCharToken('=')
				addCharToken('=')
				pos += 2
				continue
			}
			if startsWith("!=") {
				addCharToken('!')
				addCharToken('=')
				pos += 2
				continue
			}
			if startsWith("<=") {
				addCharToken('<')
				addCharToken('=')
				pos += 2
				continue
			}
			if startsWith(">=") {
				addCharToken('>')
				addCharToken('=')
				pos += 2
				continue
			}
			if startsWith("&&") {
				addCharToken('&')
				addCharToken('&')
				pos += 2
				continue
			}
			if startsWith("||") {
				addCharToken('|')
				addCharToken('|')
				pos += 2
				continue
			}
			if startsWith("<<") {
				addCharToken('<')
				addCharToken('<')
				pos += 2
				continue
			}
			if startsWith(">>") {
				addCharToken('>')
				addCharToken('>')
				pos += 2
				continue
			}
			if startsWith("++") {
				addCharToken('+')
				addCharToken('+')
				pos += 2
				continue
			}
			if startsWith("--") {
				addCharToken('-')
				addCharToken('-')
				pos += 2
				continue
			}
			if startsWith(".=") {
				addCharToken('.')
				addCharToken('=')
				pos += 2
				continue
			}
			if startsWith("+=") {
				addCharToken('+')
				addCharToken('=')
				pos += 2
				continue
			}
			if startsWith("-=") {
				addCharToken('-')
				addCharToken('=')
				pos += 2
				continue
			}
			if startsWith("*=") {
				addCharToken('*')
				addCharToken('=')
				pos += 2
				continue
			}
			if startsWith("/=") {
				addCharToken('/')
				addCharToken('=')
				pos += 2
				continue
			}
			if startsWith("%=") {
				addCharToken('%')
				addCharToken('=')
				pos += 2
				continue
			}
			if startsWith("&=") {
				addCharToken('&')
				addCharToken('=')
				pos += 2
				continue
			}
			if startsWith("|=") {
				addCharToken('|')
				addCharToken('=')
				pos += 2
				continue
			}
			if startsWith("^=") {
				addCharToken('^')
				addCharToken('=')
				pos += 2
				continue
			}

			if startsWith("\\") {
				addToken(T_NS_SEPARATOR, "\\")
				pos += 1
				continue
			}
			if startsWith("::") {
				addToken(T_PAAMAYIM_NEKUDOTAYIM, "::")
				pos += 2
				continue
			}

			// Single character tokens
			switch ch {
			case ';', '(', ')', '{', '}', '[', ']', '=', '!', '<', '>',
				'+', '-', '*', '/', '%', '&', '|', '^', '~', ',', '.',
				':', '?', '@':
				addCharToken(ch)
				pos++
			default:
				// Unknown character - include it to preserve source
				addCharToken(ch)
				pos++
			}
		}
	}

	// Main tokenizer loop
	for pos < length {
		// Check for PHP close tags first (closing a previous PHP block)
		if startsWith("?>") {
			addToken(T_CLOSE_TAG, "?>")
			pos += 2
			continue
		}

		// Check for PHP open tags
		if startsWith("<?php") || startsWith("<?PHP") || startsWith("<?Php") {
			lower := strings.ToLower(string(src[pos : pos+5]))
			if lower == "<?php" {
				// Consume <?php
				openTag := string(src[pos : pos+5])
				pos += 5

				// Consume optional whitespace after <?php (PHP eats one whitespace char)
				if pos < length && (src[pos] == ' ' || src[pos] == '\t') {
					openTag += string(src[pos])
					pos++
				}

				addToken(T_OPEN_TAG, openTag)
				tokenizePhpCode()
				continue
			}
		}

		if startsWith("<?=") {
			addToken(T_OPEN_TAG_WITH_ECHO, "<?=")
			pos += 3
			tokenizePhpCode()
			continue
		}

		if startsWith("<?") && !startsWith("<?xml") {
			addToken(T_OPEN_TAG, "<?"+string(src[pos+2]))
			if pos+2 < length && (src[pos+2] == ' ' || src[pos+2] == '\t') {
				pos += 3
			} else {
				pos += 2
			}
			tokenizePhpCode()
			continue
		}

		// Inline HTML
		htmlStart := pos
		for pos < length {
			if startsWith("<?") || startsWith("?>") {
				break
			}
			if src[pos] == '\n' {
				line++
			}
			pos++
		}
		if pos > htmlStart {
			addToken(T_INLINE_HTML, string(src[htmlStart:pos]))
		}
	}

	return tokens
}

// InitTokenConstants 注册 PHP token 常量
func InitTokenConstants(vm data.VM) {
	consts := map[string]int{
		"T_REQUIRE_ONCE":             T_REQUIRE_ONCE,
		"T_LOGICAL_AND":              T_LOGICAL_AND,
		"T_LOGICAL_OR":               T_LOGICAL_OR,
		"T_LOGICAL_XOR":              T_LOGICAL_XOR,
		"T_CLONE":                    T_CLONE,
		"T_ELSEIF":                   T_ELSEIF,
		"T_ELSE":                     T_ELSE,
		"T_ENDIF":                    T_ENDIF,
		"T_ECHO":                     T_ECHO,
		"T_DO":                       T_DO,
		"T_WHILE":                    T_WHILE,
		"T_ENDWHILE":                 T_ENDWHILE,
		"T_FOR":                      T_FOR,
		"T_ENDFOR":                   T_ENDFOR,
		"T_FOREACH":                  T_FOREACH,
		"T_ENDFOREACH":               T_ENDFOREACH,
		"T_DECLARE":                  T_DECLARE,
		"T_ENDDECLARE":               T_ENDDECLARE,
		"T_AS":                       T_AS,
		"T_SWITCH":                   T_SWITCH,
		"T_ENDSWITCH":                T_ENDSWITCH,
		"T_CASE":                     T_CASE,
		"T_DEFAULT":                  T_DEFAULT,
		"T_BREAK":                    T_BREAK,
		"T_CONTINUE":                 T_CONTINUE,
		"T_GOTO":                     T_GOTO,
		"T_FUNCTION":                 T_FUNCTION,
		"T_CONST":                    T_CONST,
		"T_RETURN":                   T_RETURN,
		"T_YIELD":                    T_YIELD,
		"T_TRY":                      T_TRY,
		"T_CATCH":                    T_CATCH,
		"T_FINALLY":                  T_FINALLY,
		"T_THROW":                    T_THROW,
		"T_IF":                       T_IF,
		"T_INSTANCEOF":               T_INSTANCEOF,
		"T_NEW":                      T_NEW,
		"T_EXIT":                     T_EXIT,
		"T_EVAL":                     T_EVAL,
		"T_EMPTY":                    T_EMPTY,
		"T_INCLUDE":                  T_INCLUDE,
		"T_INCLUDE_ONCE":             T_INCLUDE_ONCE,
		"T_REQUIRE":                  T_REQUIRE,
		"T_USE":                      T_USE,
		"T_GLOBAL":                   T_GLOBAL,
		"T_ISSET":                    T_ISSET,
		"T_UNSET":                    T_UNSET,
		"T_LIST":                     T_LIST,
		"T_ARRAY":                    T_ARRAY,
		"T_PRINT":                    T_PRINT,
		"T_NAMESPACE":                T_NAMESPACE,
		"T_OBJECT_OPERATOR":          T_OBJECT_OPERATOR,
		"T_PAAMAYIM_NEKUDOTAYIM":     T_PAAMAYIM_NEKUDOTAYIM,
		"T_DOUBLE_ARROW":             T_DOUBLE_ARROW,
		"T_DOUBLE_COLON":             T_DOUBLE_COLON,
		"T_NS_SEPARATOR":             T_NS_SEPARATOR,
		"T_ELLIPSIS":                 T_ELLIPSIS,
		"T_COALESCE":                 T_COALESCE,
		"T_SPACESHIP":                T_SPACESHIP,
		"T_POW":                      T_POW,
		"T_POW_EQUAL":                T_POW_EQUAL,
		"T_INLINE_HTML":              T_INLINE_HTML,
		"T_OPEN_TAG":                 T_OPEN_TAG,
		"T_OPEN_TAG_WITH_ECHO":       T_OPEN_TAG_WITH_ECHO,
		"T_CLOSE_TAG":                T_CLOSE_TAG,
		"T_WHITESPACE":               T_WHITESPACE,
		"T_COMMENT":                  T_COMMENT,
		"T_DOC_COMMENT":              T_DOC_COMMENT,
		"T_STRING":                   T_STRING,
		"T_VARIABLE":                 T_VARIABLE,
		"T_LNUMBER":                  T_LNUMBER,
		"T_DNUMBER":                  T_DNUMBER,
		"T_NUM_STRING":               T_NUM_STRING,
		"T_CONSTANT_ENCAPSED_STRING": T_CONSTANT_ENCAPSED_STRING,
		"T_ENCAPSED_AND_WHITESPACE":  T_ENCAPSED_AND_WHITESPACE,
		"T_CHARACTER":                T_CHARACTER,
		"T_BAD_CHARACTER":            T_BAD_CHARACTER,
		"T_ABSTRACT":                 T_ABSTRACT,
		"T_FINAL":                    T_FINAL,
		"T_PRIVATE":                  T_PRIVATE,
		"T_PROTECTED":                T_PROTECTED,
		"T_PUBLIC":                   T_PUBLIC,
		"T_STATIC":                   T_STATIC,
		"T_TRAIT":                    T_TRAIT,
		"T_INTERFACE":                T_INTERFACE,
		"T_CLASS":                    T_CLASS,
		"T_CALLABLE":                 T_CALLABLE,
		"T_EXTENDS":                  T_EXTENDS,
		"T_IMPLEMENTS":               T_IMPLEMENTS,
		"T_VAR":                      T_VAR,
		"T_READONLY":                 T_READONLY,
		"T_ENUM":                     T_ENUM,
		"T_NAME_FULLY_QUALIFIED":     T_NAME_FULLY_QUALIFIED,
		"T_NAME_RELATIVE":            T_NAME_RELATIVE,
		"T_NAME_QUALIFIED":           T_NAME_QUALIFIED,
		"T_MATCH":                    T_MATCH,
		"T_ATTRIBUTE":                T_ATTRIBUTE,
		"T_NULLSAFE_OBJECT_OPERATOR": T_NULLSAFE_OBJECT_OPERATOR,
		"T_FN":                       T_FN,
	}
	for name, val := range consts {
		vm.SetConstant(name, data.NewIntValue(val))
	}

	// Boolean-like constants
	vm.SetConstant("T_TRUE", data.NewIntValue(T_TRUE))
	vm.SetConstant("T_FALSE", data.NewIntValue(T_FALSE))
	vm.SetConstant("T_NULL", data.NewIntValue(T_NULL))
}
