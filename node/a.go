package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// TryGetCallClassName 尽可能尝试获取名称
func TryGetCallClassName(call data.GetValue) string {
	if call == nil {
		return "<nil>"
	}
	switch c := call.(type) {
	case *VariableExpression:
		return fmt.Sprintf("$%s", c.Name)
	case *data.ClassValue:
		return TryGetCallClassName(c.Class)
	case data.ClassStmt:
		return c.GetName()
	case *data.ThisValue:
		return TryGetCallClassName(c.ClassValue)
	case *data.NullValue:
		return "null"
	case *data.BoolValue:
		if c.Value {
			return "true"
		}
		return "false"
	case *data.IntValue:
		return fmt.Sprintf("%d", c.Value)
	case *data.StringValue:
		return fmt.Sprintf("%q", c.Value)
	case *data.FloatValue:
		return fmt.Sprintf("%g", c.Value)
	case *CallStaticMethod:
		return fmt.Sprintf("%s::%s", TryGetCallClassName(c.stmt), c.Method)
	case *This:
		return "this"
	case *IndexExpression:
		return fmt.Sprintf("%s[%s]", TryGetCallClassName(c.Array), TryGetCallClassName(c.Index))
	case *CallObjectProperty:
		return fmt.Sprintf("%s->%s", TryGetCallClassName(c.Object), c.Property)
	case *ReturnStatement:
		return fmt.Sprintf("return %s", TryGetCallClassName(c.Value))
	case *TernaryExpression:
		return fmt.Sprintf("%s ? %s : %s", TryGetCallClassName(c.Condition), TryGetCallClassName(c.TrueValue), TryGetCallClassName(c.FalseValue))
	case *CallObjectMethod:
		return fmt.Sprintf("%s->%s(%s)", TryGetCallClassName(c.Object), c.Method, tryArgsStr(c.Args))
	case *BinaryAssignVariable:
		return fmt.Sprintf("%s = %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryAssign:
		return fmt.Sprintf("%s = %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *CallExpression:
		return fmt.Sprintf("%s(%s)", c.FunName, tryArgsStr(c.Args))
	case *CallMethod:
		return fmt.Sprintf("call(%s)", TryGetCallClassName(c.Method))
	case *CallParentMethod:
		return fmt.Sprintf("parent::%s", c.Method)
	case *CallStaticProperty:
		return fmt.Sprintf("%s::$%s", TryGetCallClassName(c.Stmt), c.Property)
	case *NewExpression:
		return fmt.Sprintf("new %s", c.ClassName)
	case *NewVariableExpression:
		return fmt.Sprintf("new %s", TryGetCallClassName(c.ClassNameExpr))
	case *InstanceOfExpression:
		return fmt.Sprintf("%s instanceof %s", TryGetCallClassName(c.Object), TryGetCallClassName(c.ClassName))
	case *UnaryExpression:
		return fmt.Sprintf("%s%s", c.Operator, TryGetCallClassName(c.Right))
	case *UnaryIncr:
		return fmt.Sprintf("++%s", TryGetCallClassName(c.Right))
	case *UnaryDecr:
		return fmt.Sprintf("--%s", TryGetCallClassName(c.Right))
	case *NullCoalesceExpression:
		return fmt.Sprintf("%s ?? %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *StringLiteral:
		return fmt.Sprintf("%q", c.Value)
	case *IntLiteral:
		return fmt.Sprintf("%v", c.V)
	case *BooleanLiteral:
		if c.Value {
			return "true"
		}
		return "false"
	case *Todo:
		return "_todo_"
	// 二元运算符
	case *BinaryAdd:
		return fmt.Sprintf("%s + %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryAssignVariableList:
		return fmt.Sprintf("list(...) = %s", TryGetCallClassName(c.Right))
	case *BinaryDot:
		return fmt.Sprintf("%s . %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryEq:
		return fmt.Sprintf("%s == %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryEqStrict:
		return fmt.Sprintf("%s === %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryNe:
		return fmt.Sprintf("%s != %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryNeStrict:
		return fmt.Sprintf("%s !== %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryLt:
		return fmt.Sprintf("%s < %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryLe:
		return fmt.Sprintf("%s <= %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryGt:
		return fmt.Sprintf("%s > %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryGe:
		return fmt.Sprintf("%s >= %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryLand:
		return fmt.Sprintf("%s && %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryLor:
		return fmt.Sprintf("%s || %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryMul:
		return fmt.Sprintf("%s * %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryQuo:
		return fmt.Sprintf("%s / %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryRem:
		return fmt.Sprintf("%s %% %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinarySub:
		return fmt.Sprintf("%s - %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryBitAnd:
		return fmt.Sprintf("%s & %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryBitOr:
		return fmt.Sprintf("%s | %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryBitXor:
		return fmt.Sprintf("%s ^ %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryShl:
		return fmt.Sprintf("%s << %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryShr:
		return fmt.Sprintf("%s >> %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	case *BinaryLink:
		return fmt.Sprintf("%s . %s", TryGetCallClassName(c.Left), TryGetCallClassName(c.Right))
	// 静态方法/属性
	case *CallStaticMethodLater:
		return fmt.Sprintf("%s::%s()", c.className, c.method)
	case *CallSelfMethod:
		return fmt.Sprintf("self::%s()", c.Method)
	case *CallStaticKeywordMethod:
		return fmt.Sprintf("static::%s()", c.Method)
	case *CallStaticKeywordProperty:
		return fmt.Sprintf("static::$%s", c.Property)
	case *CallSelfProperty:
		return fmt.Sprintf("self::$%s", c.Property)
	// 后缀自增/自减
	case *PostfixIncr:
		return fmt.Sprintf("%s++", TryGetCallClassName(c.Left))
	case *PostfixDecr:
		return fmt.Sprintf("%s--", TryGetCallClassName(c.Left))
	// 其他表达式
	case *CloneExpression:
		return fmt.Sprintf("clone %s", TryGetCallClassName(c.Target))
	case *LambdaExpression:
		return "lambda::()"
	case *SpreadArgument:
		return fmt.Sprintf("...%s", TryGetCallClassName(c.Expr))
	case *MatchStatement:
		return fmt.Sprintf("match(%s){...}", TryGetCallClassName(c.Condition))
	case *NullsafeCall:
		return fmt.Sprintf("%s?->...", TryGetCallClassName(c.Object))
	case *data.FuncValue:
		return "func()"
	case *data.ArrayValue:
		return "array(...)"
	case *data.ObjectValue:
		return "object{...}"
	}
	// 未实现的打印内容，附带具体类型
	return fmt.Sprintf("?(%T)", call)
}

// tryArgsStr 将参数列表转换为字符串表示
func tryArgsStr(args []data.GetValue) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, a := range args {
		parts = append(parts, TryGetCallClassName(a))
	}
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ", "
		}
		result += p
	}
	return result
}
