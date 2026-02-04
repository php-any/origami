package data

type ASTValue struct {
	Node GetValue
	Ctx  Context
}

func NewASTValue(node GetValue, ctx Context) *ASTValue {
	return &ASTValue{Node: node, Ctx: ctx}
}

func (a *ASTValue) GetValue(ctx Context) (GetValue, Control) {
	return a, nil
}

func (a *ASTValue) AsString() string {
	return "ast"
}
