package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionClassGetFileNameMethod 实现 ReflectionClass::getFileName(): string|false
type ReflectionClassGetFileNameMethod struct{}

func (m *ReflectionClassGetFileNameMethod) GetName() string               { return "getFileName" }
func (m *ReflectionClassGetFileNameMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ReflectionClassGetFileNameMethod) GetIsStatic() bool             { return false }
func (m *ReflectionClassGetFileNameMethod) GetParams() []data.GetValue    { return nil }
func (m *ReflectionClassGetFileNameMethod) GetVariables() []data.Variable { return nil }
func (m *ReflectionClassGetFileNameMethod) GetReturnType() data.Types     { return data.Mixed{} }

func (m *ReflectionClassGetFileNameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewBoolValue(false), nil
	}
	return sourceFileFrom(classStmt), nil
}

func sourceFileFrom(v any) data.Value {
	if gf, ok := v.(interface{ GetFrom() data.From }); ok {
		if from := gf.GetFrom(); from != nil {
			if src := from.GetSource(); src != "" {
				return data.NewStringValue(src)
			}
		}
	}
	// 内部/Go 实现的类无源文件，对齐 PHP 返回 false
	return data.NewBoolValue(false)
}

func sourceStartLineFrom(v any) data.Value {
	if gf, ok := v.(interface{ GetFrom() data.From }); ok {
		if from := gf.GetFrom(); from != nil {
			line, _ := from.GetStartPosition()
			// TokenFrom / lexer 行号为 0-based，PHP getStartLine 为 1-based
			return data.NewIntValue(line + 1)
		}
	}
	return data.NewBoolValue(false)
}
