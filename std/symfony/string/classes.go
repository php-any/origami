package sfstring

import (
	"crypto/rand"
	"math/bits"

	"github.com/php-any/origami/data"
)

func registerCommon(c *strClass) {
	defEmpty := data.NewStringValue("")
	defZero := data.NewIntValue(0)
	defOne := data.NewIntValue(1)
	defFalse := data.NewBoolValue(false)
	defTrue := data.NewBoolValue(true)
	defSpace := data.NewStringValue(" ")

	c.add("__construct", false, []data.GetValue{p("string", 0, defEmpty)}, []data.Variable{v("string", 0)}, methodConstruct)
	c.add("__clone", false, nil, nil, methodClone)
	c.add("__toString", false, nil, nil, methodToString)
	c.add("toString", false, nil, nil, methodToStringAlias)
	c.add("jsonSerialize", false, nil, nil, methodJsonSerialize)
	c.add("length", false, nil, nil, methodLength)
	c.add("isEmpty", false, nil, nil, methodIsEmpty)
	c.add("ignoreCase", false, nil, nil, methodIgnoreCase)
	c.add("append", false, []data.GetValue{variadic("suffix", 0)}, []data.Variable{v("suffix", 0)}, methodAppend)
	c.add("prepend", false, []data.GetValue{variadic("prefix", 0)}, []data.Variable{v("prefix", 0)}, methodPrepend)
	c.add("slice", false, []data.GetValue{p("start", 0, defZero), p("length", 1, data.NewNullValue())}, []data.Variable{v("start", 0), v("length", 1)}, methodSlice)
	c.add("splice", false, []data.GetValue{p("replacement", 0, nil), p("start", 1, defZero), p("length", 2, data.NewNullValue())}, []data.Variable{v("replacement", 0), v("start", 1), v("length", 2)}, methodSplice)
	c.add("chunk", false, []data.GetValue{p("length", 0, defOne)}, []data.Variable{v("length", 0)}, methodChunk)
	c.add("split", false, []data.GetValue{p("delimiter", 0, nil), p("limit", 1, data.NewNullValue()), p("flags", 2, data.NewNullValue())}, []data.Variable{v("delimiter", 0), v("limit", 1), v("flags", 2)}, methodSplit)
	c.add("indexOf", false, []data.GetValue{p("needle", 0, nil), p("offset", 1, defZero)}, []data.Variable{v("needle", 0), v("offset", 1)}, methodIndexOf)
	c.add("indexOfLast", false, []data.GetValue{p("needle", 0, nil), p("offset", 1, defZero)}, []data.Variable{v("needle", 0), v("offset", 1)}, methodIndexOfLast)
	c.add("startsWith", false, []data.GetValue{p("prefix", 0, nil)}, []data.Variable{v("prefix", 0)}, methodStartsWith)
	c.add("endsWith", false, []data.GetValue{p("suffix", 0, nil)}, []data.Variable{v("suffix", 0)}, methodEndsWith)
	c.add("equalsTo", false, []data.GetValue{p("string", 0, nil)}, []data.Variable{v("string", 0)}, methodEqualsTo)
	c.add("containsAny", false, []data.GetValue{p("needle", 0, nil)}, []data.Variable{v("needle", 0)}, methodContainsAny)
	c.add("contains", false, []data.GetValue{p("needle", 0, nil)}, []data.Variable{v("needle", 0)}, methodContainsAny)
	c.add("replace", false, []data.GetValue{p("from", 0, nil), p("to", 1, nil)}, []data.Variable{v("from", 0), v("to", 1)}, methodReplace)
	c.add("replaceMatches", false, []data.GetValue{p("fromRegexp", 0, nil), p("to", 1, nil)}, []data.Variable{v("fromRegexp", 0), v("to", 1)}, methodReplaceMatches)
	c.add("lower", false, nil, nil, methodLower)
	c.add("upper", false, nil, nil, methodUpper)
	c.add("folded", false, []data.GetValue{p("compat", 0, defTrue)}, []data.Variable{v("compat", 0)}, methodFolded)
	c.add("camel", false, nil, nil, methodCamel)
	c.add("snake", false, nil, nil, methodSnake)
	c.add("kebab", false, nil, nil, methodKebab)
	c.add("pascal", false, nil, nil, methodPascal)
	c.add("title", false, []data.GetValue{p("allWords", 0, defFalse)}, []data.Variable{v("allWords", 0)}, methodTitle)
	c.add("trim", false, []data.GetValue{p("chars", 0, nil)}, []data.Variable{v("chars", 0)}, methodTrim)
	c.add("trimStart", false, []data.GetValue{p("chars", 0, nil)}, []data.Variable{v("chars", 0)}, methodTrimStart)
	c.add("trimEnd", false, []data.GetValue{p("chars", 0, nil)}, []data.Variable{v("chars", 0)}, methodTrimEnd)
	c.add("trimPrefix", false, []data.GetValue{p("prefix", 0, nil)}, []data.Variable{v("prefix", 0)}, methodTrimPrefix)
	c.add("trimSuffix", false, []data.GetValue{p("suffix", 0, nil)}, []data.Variable{v("suffix", 0)}, methodTrimSuffix)
	c.add("reverse", false, nil, nil, methodReverse)
	c.add("repeat", false, []data.GetValue{p("multiplier", 0, nil)}, []data.Variable{v("multiplier", 0)}, methodRepeat)
	c.add("padBoth", false, []data.GetValue{p("length", 0, nil), p("padStr", 1, defSpace)}, []data.Variable{v("length", 0), v("padStr", 1)}, methodPadBoth)
	c.add("padEnd", false, []data.GetValue{p("length", 0, nil), p("padStr", 1, defSpace)}, []data.Variable{v("length", 0), v("padStr", 1)}, methodPadEnd)
	c.add("padStart", false, []data.GetValue{p("length", 0, nil), p("padStr", 1, defSpace)}, []data.Variable{v("length", 0), v("padStr", 1)}, methodPadStart)
	c.add("join", false, []data.GetValue{p("strings", 0, nil), p("lastGlue", 1, data.NewNullValue())}, []data.Variable{v("strings", 0), v("lastGlue", 1)}, methodJoin)
	c.add("bytesAt", false, []data.GetValue{p("offset", 0, nil)}, []data.Variable{v("offset", 0)}, methodBytesAt)
	c.add("toByteString", false, []data.GetValue{p("toEncoding", 0, data.NewNullValue())}, []data.Variable{v("toEncoding", 0)}, methodToByteString)
	c.add("toUnicodeString", false, []data.GetValue{p("fromEncoding", 0, data.NewNullValue())}, []data.Variable{v("fromEncoding", 0)}, methodToUnicodeString)
	c.add("ensureEnd", false, []data.GetValue{p("suffix", 0, nil)}, []data.Variable{v("suffix", 0)}, methodEnsureEnd)
	c.add("ensureStart", false, []data.GetValue{p("prefix", 0, nil)}, []data.Variable{v("prefix", 0)}, methodEnsureStart)
	c.add("after", false, []data.GetValue{p("needle", 0, nil), p("includeNeedle", 1, defFalse), p("offset", 2, defZero)}, []data.Variable{v("needle", 0), v("includeNeedle", 1), v("offset", 2)}, methodAfter)
	c.add("afterLast", false, []data.GetValue{p("needle", 0, nil), p("includeNeedle", 1, defFalse), p("offset", 2, defZero)}, []data.Variable{v("needle", 0), v("includeNeedle", 1), v("offset", 2)}, methodAfterLast)
	c.add("before", false, []data.GetValue{p("needle", 0, nil), p("includeNeedle", 1, defFalse), p("offset", 2, defZero)}, []data.Variable{v("needle", 0), v("includeNeedle", 1), v("offset", 2)}, methodBefore)
	c.add("beforeLast", false, []data.GetValue{p("needle", 0, nil), p("includeNeedle", 1, defFalse), p("offset", 2, defZero)}, []data.Variable{v("needle", 0), v("includeNeedle", 1), v("offset", 2)}, methodBeforeLast)
	c.add("truncate", false, []data.GetValue{p("length", 0, nil), p("ellipsis", 1, defEmpty), p("cut", 2, defTrue)}, []data.Variable{v("length", 0), v("ellipsis", 1), v("cut", 2)}, methodTruncate)
	c.add("wordwrap", false, []data.GetValue{p("width", 0, data.NewIntValue(75)), p("break", 1, data.NewStringValue("\n")), p("cut", 2, defFalse)}, []data.Variable{v("width", 0), v("break", 1), v("cut", 2)}, methodWordwrap)
	c.add("collapseWhitespace", false, nil, nil, methodCollapseWhitespace)
	c.add("match", false, []data.GetValue{p("regexp", 0, nil), p("flags", 1, defZero), p("offset", 2, defZero)}, []data.Variable{v("regexp", 0), v("flags", 1), v("offset", 2)}, methodMatch)
	c.add("wrap", true, []data.GetValue{p("values", 0, nil)}, []data.Variable{v("values", 0)}, methodWrap)
	c.add("unwrap", true, []data.GetValue{p("values", 0, nil)}, []data.Variable{v("values", 0)}, methodUnwrap)
}

func registerUnicode(c *strClass) {
	defTrue := data.NewBoolValue(true)
	defNFC := data.NewIntValue(4)
	c.add("ascii", false, []data.GetValue{p("rules", 0, nil)}, []data.Variable{v("rules", 0)}, methodAscii)
	c.add("normalize", false, []data.GetValue{p("form", 0, defNFC)}, []data.Variable{v("form", 0)}, methodNormalize)
	c.add("width", false, []data.GetValue{p("ignoreAnsiDecoration", 0, defTrue)}, []data.Variable{v("ignoreAnsiDecoration", 0)}, methodWidth)
	c.add("codePointsAt", false, []data.GetValue{p("offset", 0, nil)}, []data.Variable{v("offset", 0)}, methodCodePointsAt)
	c.add("toCodePointString", false, []data.GetValue{p("fromEncoding", 0, data.NewNullValue())}, []data.Variable{v("fromEncoding", 0)}, methodToCodePointString)
	c.add("fromCodePoints", true, []data.GetValue{variadic("codes", 0)}, []data.Variable{v("codes", 0)}, methodFromCodePoints)
}

func NewUnicodeStringClass() data.ClassStmt {
	if unicodeStringClass != nil && len(unicodeStringClass.methods) > 0 {
		return unicodeStringClass
	}
	c := &strClass{
		name:       unicodeStringName,
		parent:     abstractUnicodeName,
		kind:       kindGrapheme,
		methods:    map[string]data.Method{},
		implements: []string{"Stringable", "JsonSerializable"},
	}
	initStringProps(c)
	registerCommon(c)
	registerUnicode(c)
	unicodeStringClass = c
	return c
}

func NewByteStringClass() data.ClassStmt {
	if byteStringClass != nil && len(byteStringClass.methods) > 0 {
		return byteStringClass
	}
	c := &strClass{
		name:       byteStringName,
		parent:     abstractStringName,
		kind:       kindByte,
		methods:    map[string]data.Method{},
		implements: []string{"Stringable", "JsonSerializable"},
	}
	initStringProps(c)
	registerCommon(c)
	c.add("isUtf8", false, nil, nil, methodIsUtf8)
	c.add("fromRandom", true,
		[]data.GetValue{p("length", 0, data.NewIntValue(16)), p("alphabet", 1, data.NewNullValue())},
		[]data.Variable{v("length", 0), v("alphabet", 1)},
		methodFromRandom)
	c.add("toCodePointString", false, []data.GetValue{p("fromEncoding", 0, data.NewNullValue())}, []data.Variable{v("fromEncoding", 0)}, methodToCodePointString)
	byteStringClass = c
	return c
}

func NewCodePointStringClass() data.ClassStmt {
	if codePointStringClass != nil && len(codePointStringClass.methods) > 0 {
		return codePointStringClass
	}
	c := &strClass{
		name:       codePointStringName,
		parent:     abstractUnicodeName,
		kind:       kindCodePoint,
		methods:    map[string]data.Method{},
		implements: []string{"Stringable", "JsonSerializable"},
	}
	initStringProps(c)
	registerCommon(c)
	registerUnicode(c)
	codePointStringClass = c
	return c
}

func methodFromRandom(ctx data.Context) (data.GetValue, data.Control) {
	length := argInt(ctx, 0, 16)
	if length <= 0 {
		return nil, throwInvalid("A strictly positive length is expected, \"" + itoa(length) + "\" given.")
	}
	alphabet := argString(ctx, 1, byteAlphabet)
	if alphabet == "" {
		alphabet = byteAlphabet
	}
	alphaSize := len(alphabet)
	if alphaSize < 2 {
		return nil, throwInvalid("The length of the alphabet must in the [2^1, 2^56] range.")
	}
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	out := make([]byte, length)
	mask := byte(1<<bits.Len(uint(alphaSize-1)) - 1)
	if mask == 0 {
		mask = 0xff
	}
	for i := 0; i < length; i++ {
		out[i] = alphabet[int(buf[i]&mask)%alphaSize]
	}
	return newTyped(ctx, ensureByteClass(), string(out)), nil
}
