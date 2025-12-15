package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrrposFunction 实现 strrpos 函数
type StrrposFunction struct{}

func NewStrrposFunction() data.FuncStmt {
	return &StrrposFunction{}
}

func (f *StrrposFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	haystackValue, _ := ctx.GetIndexValue(0)
	needleValue, _ := ctx.GetIndexValue(1)
	offsetValue, _ := ctx.GetIndexValue(2)

	haystack := haystackValue.AsString()
	needle := needleValue.AsString()
	offset := 0
	// fmt.Printf("Debug strrpos input: haystack=%s, needle=%s, offsetValue=%T, %v\n", haystack, needle, offsetValue, offsetValue)

	if offsetValue != nil {
		if i, ok := offsetValue.(data.AsInt); ok {
			if val, err := i.AsInt(); err == nil {
				offset = val
			}
		}
	}

	if offset >= 0 {
		if offset > len(haystack) {
			return data.NewBoolValue(false), nil
		}
		searchSpace := haystack[offset:]
		idx := strings.LastIndex(searchSpace, needle)
		if idx == -1 {
			return data.NewBoolValue(false), nil
		}
		return data.NewIntValue(idx + offset), nil
	} else {
		// Negative offset
		// Search starts at len(haystack) + offset from the end, moving backwards.
		// Effectively, we search in the substring haystack[:len(haystack)+offset+len(needle)]?
		// Actually, "search starting at X" for backward search means we look at X, X-1, ...
		// So the match must start at or before X.

		endPos := len(haystack) + offset
		if endPos < 0 {
			return data.NewBoolValue(false), nil
		}

		// We want to find the last occurrence where the match starts at index <= endPos.
		// So we can search in haystack[:endPos+len(needle)].
		// But we must ensure we don't go out of bounds.

		limit := endPos + len(needle)
		if limit > len(haystack) {
			limit = len(haystack)
		}

		// fmt.Printf("Debug strrpos: offset=%d, endPos=%d, limit=%d, len=%d\n", offset, endPos, limit, len(haystack))

		searchSpace := haystack[:limit]
		idx := strings.LastIndex(searchSpace, needle)
		if idx == -1 {
			return data.NewBoolValue(false), nil
		}

		// If the match starts AFTER endPos, it's invalid (should not happen if we cut correctly).
		if idx > endPos {
			// This implies we included too much.
			// But since we used LastIndex, and we want match <= endPos.
			// If we cut at endPos + len(needle), we allow match starting at endPos.
			// So it should be fine.
			// However, if limit was clamped to len(haystack), we might find a match starting > endPos?
			// Example: haystack "ABC", needle "C", offset -1. endPos = 2 ('C').
			// limit = 2 + 1 = 3. searchSpace "ABC".
			// LastIndex("ABC", "C") -> 2.
			// 2 <= 2. OK.

			// Example: haystack "ABC", needle "C", offset -2. endPos = 1 ('B').
			// limit = 1 + 1 = 2. searchSpace "AB".
			// LastIndex("AB", "C") -> -1. OK.

			// Example: haystack "ABC", needle "BC", offset -2. endPos = 1 ('B').
			// limit = 1 + 2 = 3. searchSpace "ABC".
			// LastIndex("ABC", "BC") -> 1.
			// 1 <= 1. OK.
		}

		return data.NewIntValue(idx), nil
	}
}

func (f *StrrposFunction) GetName() string {
	return "strrpos"
}

func (f *StrrposFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, nil),
		node.NewParameter(nil, "needle", 1, nil, nil),
		node.NewParameter(nil, "offset", 2, node.NewIntLiteral(nil, "0"), nil),
	}
}

func (f *StrrposFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "needle", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "offset", 2, data.NewBaseType("int")),
	}
}

// StrriposFunction 实现 strripos 函数 (case-insensitive)
type StrriposFunction struct{}

func NewStrriposFunction() data.FuncStmt {
	return &StrriposFunction{}
}

func (f *StrriposFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	haystackValue, _ := ctx.GetIndexValue(0)
	needleValue, _ := ctx.GetIndexValue(1)
	offsetValue, _ := ctx.GetIndexValue(2)

	haystack := strings.ToLower(haystackValue.AsString())
	needle := strings.ToLower(needleValue.AsString())

	// Reuse logic from Strrpos but with lowercased strings
	// Copy-paste logic for now or delegate?
	// Since I can't easily delegate to another function struct without creating it, I'll duplicate logic.

	offset := 0
	if offsetValue != nil {
		if i, ok := offsetValue.(data.AsInt); ok {
			if val, err := i.AsInt(); err == nil {
				offset = val
			}
		}
	}

	if offset >= 0 {
		if offset > len(haystack) {
			return data.NewBoolValue(false), nil
		}
		searchSpace := haystack[offset:]
		idx := strings.LastIndex(searchSpace, needle)
		if idx == -1 {
			return data.NewBoolValue(false), nil
		}
		return data.NewIntValue(idx + offset), nil
	} else {
		endPos := len(haystack) + offset
		if endPos < 0 {
			return data.NewBoolValue(false), nil
		}
		limit := endPos + len(needle)
		if limit > len(haystack) {
			limit = len(haystack)
		}
		searchSpace := haystack[:limit]
		idx := strings.LastIndex(searchSpace, needle)
		if idx == -1 {
			return data.NewBoolValue(false), nil
		}
		return data.NewIntValue(idx), nil
	}
}

func (f *StrriposFunction) GetName() string {
	return "strripos"
}

func (f *StrriposFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, nil),
		node.NewParameter(nil, "needle", 1, nil, nil),
		node.NewParameter(nil, "offset", 2, node.NewIntLiteral(nil, "0"), nil),
	}
}

func (f *StrriposFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "needle", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "offset", 2, data.NewBaseType("int")),
	}
}
