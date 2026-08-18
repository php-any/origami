package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionClassIsInternalMethod 实现 ReflectionClass::isInternal(): bool
// 判断被反射的类是否由 PHP 内核/扩展提供（内部类），而非用户定义的 PHP 类
type ReflectionClassIsInternalMethod struct{}

func (m *ReflectionClassIsInternalMethod) GetName() string               { return "isInternal" }
func (m *ReflectionClassIsInternalMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ReflectionClassIsInternalMethod) GetIsStatic() bool             { return false }
func (m *ReflectionClassIsInternalMethod) GetParams() []data.GetValue    { return nil }
func (m *ReflectionClassIsInternalMethod) GetVariables() []data.Variable { return nil }
func (m *ReflectionClassIsInternalMethod) GetReturnType() data.Types     { return data.Bool{} }

func (m *ReflectionClassIsInternalMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, classStmt := getReflectionClassInfo(ctx)
	// 内部类（Go 实现的类）没有 PHP 源文件；无法获取类语句时按内部类处理
	return data.NewBoolValue(!hasSourceFile(classStmt)), nil
}

// hasSourceFile 判断类语句是否来自 PHP 源文件
func hasSourceFile(v any) bool {
	if gf, ok := v.(interface{ GetFrom() data.From }); ok {
		if from := gf.GetFrom(); from != nil {
			if src := from.GetSource(); src != "" {
				return true
			}
		}
	}
	return false
}
