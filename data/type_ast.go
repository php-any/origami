package data

// AST 表示注解 target 参数可接收的 AST 节点值类型。
type AST struct{}

func (AST) Is(value Value) bool {
	switch value.(type) {
	case *NullValue, *ASTValue, *AnyValue:
		return true
	default:
		return false
	}
}

func (AST) String() string {
	return "AstNode"
}
