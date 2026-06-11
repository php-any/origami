package compile

import (
	"reflect"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/net/annotation"
)

type specialHandler func(g *Generator, v data.GetValue) error

var specialHandlers map[reflect.Type]specialHandler

func init() {
	specialHandlers = map[reflect.Type]specialHandler{
		reflect.TypeOf((*node.CallExpression)(nil)):           emitCallExpression,
		reflect.TypeOf((*node.CallMethod)(nil)):               emitCallMethod,
		reflect.TypeOf((*node.CallStaticMethod)(nil)):         emitCallStaticMethod,
		reflect.TypeOf((*node.CallStaticProperty)(nil)):       emitCallStaticProperty,
		reflect.TypeOf((*node.CallStaticMethodLater)(nil)):    emitCallStaticMethodLater,
		reflect.TypeOf((*node.CallStaticPropertyLater)(nil)):  emitCallStaticPropertyLater,
		reflect.TypeOf((*node.CallLater)(nil)):                emitCallLater,
		reflect.TypeOf((*node.LambdaExpression)(nil)):         emitLambdaExpression,
		reflect.TypeOf((*node.ClassStatement)(nil)):           emitClassStatement,
		reflect.TypeOf((*node.AbstractClassStatement)(nil)):   emitAbstractClassStatement,
		reflect.TypeOf((*node.FunctionStatement)(nil)):        emitFunctionStatement,
		reflect.TypeOf((*node.InterfaceStatement)(nil)):       emitInterfaceStatement,
		reflect.TypeOf((*node.VarFastAssign)(nil)):            emitVarFastAssign,
		reflect.TypeOf((*node.VarPostIncr)(nil)):              emitVarPostIncr,
		reflect.TypeOf((*node.VarStmtIncr)(nil)):              emitVarStmtIncr,
		reflect.TypeOf((*node.VarIntLe)(nil)):                 emitVarIntLe,
		reflect.TypeOf((*data.ClassValue)(nil)):               emitClassValue,
		reflect.TypeOf((*node.Array)(nil)):                    emitArray,
		reflect.TypeOf((*node.Namespace)(nil)):                emitNamespace,
		reflect.TypeOf((*node.NewExpression)(nil)):            emitNewExpression,
		reflect.TypeOf((*node.NewVariableExpression)(nil)):    emitNewVariableExpression,
		reflect.TypeOf((*node.NewExpressionDynamic)(nil)):     emitNewExpressionDynamic,
		reflect.TypeOf((*node.NewSelfExpression)(nil)):        emitNewSelfExpression,
		reflect.TypeOf((*node.NewStaticExpression)(nil)):      emitNewStaticExpression,
		reflect.TypeOf((*node.InitClass)(nil)):                emitInitClass,
		reflect.TypeOf((*node.Kv)(nil)):                       emitKv,
		reflect.TypeOf((*node.Range)(nil)):                    emitRange,
		reflect.TypeOf((*node.IncludeStatement)(nil)):         emitIncludeStatement,
		reflect.TypeOf((*node.ConstStatement)(nil)):           emitConstStatement,
		reflect.TypeOf((*node.BinaryAssignVariable)(nil)):     emitBinaryAssignVariable,
		reflect.TypeOf((*node.BinaryAssignVariableList)(nil)): emitBinaryAssignVariableList,
	}
}

func emitCallExpression(g *Generator, v data.GetValue) error {
	n := v.(*node.CallExpression)
	g.printf("node.NewCallTodo(&node.CallExpression{Node: node.NewNode(from), FunName: %q, Args: []data.GetValue{\n", n.FunName)
	g.indent++
	for _, arg := range n.Args {
		if err := g.Emit(arg); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("}}, %q)", g.namespace)
	return nil
}

func emitCallMethod(g *Generator, v data.GetValue) error {
	n := v.(*node.CallMethod)
	g.printf("node.NewCallMethod(from,\n")
	g.indent++
	if err := g.Emit(n.Method); err != nil {
		return err
	}
	g.printf(",\n")
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Args {
		if err := g.Emit(arg); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
	return nil
}

func emitCallStaticMethod(g *Generator, v data.GetValue) error {
	n := v.(*node.CallStaticMethod)
	className := staticCallClassName(n.GetStmt())
	g.printf("node.NewCallStaticMethodLater(from, %q, %q, %q)", className, n.Method, g.namespace)
	return nil
}

func emitCallStaticProperty(g *Generator, v data.GetValue) error {
	n := v.(*node.CallStaticProperty)
	className := staticCallClassName(n.Stmt)
	g.printf("node.NewCallStaticPropertyLater(from, %q, %q, %q)", className, n.Property, g.namespace)
	return nil
}

func emitCallStaticMethodLater(g *Generator, v data.GetValue) error {
	n := v.(*node.CallStaticMethodLater)
	rv := reflect.ValueOf(n).Elem()
	className := rv.FieldByName("className").String()
	method := rv.FieldByName("method").String()
	namespace := rv.FieldByName("namespace").String()
	g.printf("node.NewCallStaticMethodLater(from, %q, %q, %q)", className, method, namespace)
	return nil
}

func emitCallStaticPropertyLater(g *Generator, v data.GetValue) error {
	n := v.(*node.CallStaticPropertyLater)
	rv := reflect.ValueOf(n).Elem()
	className := rv.FieldByName("className").String()
	property := rv.FieldByName("property").String()
	namespace := rv.FieldByName("namespace").String()
	g.printf("node.NewCallStaticPropertyLater(from, %q, %q, %q)", className, property, namespace)
	return nil
}

func emitCallLater(g *Generator, v data.GetValue) error {
	n := v.(*node.CallLater)
	rv := reflect.ValueOf(n).Elem()
	namespace := rv.FieldByName("namespace").String()
	if n.CallExpression == nil {
		return newEmitError(g.file, v, "CallLater missing CallExpression")
	}
	g.printf("node.NewCallTodo(&node.CallExpression{Node: node.NewNode(from), FunName: %q, Args: []data.GetValue{\n", n.FunName)
	g.indent++
	for _, arg := range n.Args {
		if err := g.Emit(arg); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("}}, %q)", namespace)
	return nil
}

func emitLambdaExpression(g *Generator, v data.GetValue) error {
	n := v.(*node.LambdaExpression)
	fs := n.FunctionStatement
	g.printf("node.NewLambdaExpression(from,\n")
	g.indent++
	if err := g.emitParamList(fs.Params); err != nil {
		return err
	}
	g.printf(",\n")
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, stmt := range fs.Body {
		if err := g.Emit(stmt); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.genMethodVars(fs.GetVariables())
	g.printf(",\n")
	g.genParentMap(n.GetParentBindings())
	g.printf(",\n")
	g.indent--
	g.printf(")")
	return nil
}

func emitClassStatement(g *Generator, v data.GetValue) error {
	n := v.(*node.ClassStatement)
	g.printf("func() data.GetValue {\n")
	g.indent++
	if err := g.emitClassStatementInit(n); err != nil {
		return err
	}
	g.printf("return cls\n")
	g.indent--
	g.printf("}()")
	return nil
}

func emitAbstractClassStatement(g *Generator, v data.GetValue) error {
	n := v.(*node.AbstractClassStatement)
	g.printf("node.NewAbstractClassStatement(func() *node.ClassStatement {\n")
	g.indent++
	if err := g.emitClassStatementInit(n.ClassStatement); err != nil {
		return err
	}
	g.printf("return cls\n")
	g.indent--
	g.printf("}())")
	return nil
}

func (g *Generator) emitClassStatementInit(n *node.ClassStatement) error {
	extends := ""
	if n.Extends != nil {
		extends = *n.Extends
	}
	g.printf("cls := node.NewClassStatement(from, %q, %q, ", n.Name, extends)
	g.genStringSlice(n.Implements)
	g.printf(", []data.Property{\n")
	g.indent++
	for _, name := range n.PropertiesIndex {
		if err := g.emitClassProperty(n.Properties[name]); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("}, map[string]data.Method{\n")
	g.indent++
	for name, method := range n.Methods {
		g.printf("%q: ", name)
		if err := g.emitClassMethod(method); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("})\n")
	if n.IsAbstract {
		g.printf("cls.IsAbstract = true\n")
	}
	if len(n.StaticMethods) > 0 {
		g.printf("cls.StaticMethods = map[string]data.Method{\n")
		g.indent++
		for name, method := range n.StaticMethods {
			g.printf("%q: ", name)
			if err := g.emitClassMethod(method); err != nil {
				return err
			}
			g.printf(",\n")
		}
		g.indent--
		g.printf("}\n")
	}
	for _, ann := range n.Annotations {
		g.printf("cls.AddAnnotations(")
		if err := g.emitClassAnnotation(ann); err != nil {
			return err
		}
		g.printf(")\n")
	}
	return nil
}

func emitFunctionStatement(g *Generator, v data.GetValue) error {
	n := v.(*node.FunctionStatement)
	g.printf("node.NewFunctionStatement(from, %q, ", n.Name)
	if err := g.emitParamList(n.Params); err != nil {
		return err
	}
	g.printf(", []data.GetValue{\n")
	g.indent++
	for _, stmt := range n.Body {
		if err := g.Emit(stmt); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("}, ")
	g.genMethodVars(n.GetVariables())
	g.printf(", ")
	g.genTypes(n.Ret)
	g.printf(", %v)", n.ReturnsReference)
	return nil
}

func emitInterfaceStatement(g *Generator, v data.GetValue) error {
	n := v.(*node.InterfaceStatement)
	g.printf("node.NewInterfaceStatement(from, %q, ", n.Name)
	g.genStringSlice(n.Extends)
	g.printf(", []data.Method{\n")
	g.indent++
	for _, method := range n.Methods {
		if err := g.emitMethod(method); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("})")
	return nil
}

func (g *Generator) emitInterfaceMethod(im *node.InterfaceMethod) error {
	mod := modifierName(im.Modifier)
	g.printf("node.NewInterfaceMethod(from, %q, %q, ", im.Name, mod)
	if err := g.emitParamList(im.Params); err != nil {
		return err
	}
	g.printf(", ")
	g.genTypes(im.ReturnType)
	g.printf(")")
	return nil
}

func emitVarFastAssign(g *Generator, v data.GetValue) error {
	n := v.(*node.VarFastAssign)
	rv := reflect.ValueOf(n).Elem()
	op := rv.FieldByName("op").Uint()
	g.printf("node.NewVarFastAssignCompiled(from,\n")
	g.indent++
	if err := g.Emit(n.Dst); err != nil {
		return err
	}
	g.printf(",\n")
	g.printf("%d, %d, %d, %d, %d, byte(%d),\n", n.DstIdx, n.LhsIdx, n.RhsIdx, n.LhsLit, n.RhsLit, op)
	if err := g.Emit(n.Slow); err != nil {
		return err
	}
	g.printf(",\n")
	g.indent--
	g.printf(")")
	return nil
}

func emitVarPostIncr(g *Generator, v data.GetValue) error {
	n := v.(*node.VarPostIncr)
	g.printf("&node.VarPostIncr{\n")
	g.indent++
	g.printf("Node: node.NewNode(from),\n")
	g.printf("VarIdx: %d,\n", n.VarIdx)
	g.printf("Var: ")
	if err := g.Emit(n.Var); err != nil {
		return err
	}
	g.printf(",\n")
	g.printf("Fallback: ")
	if err := g.Emit(n.Fallback); err != nil {
		return err
	}
	g.printf(",\n")
	g.indent--
	g.printf("}")
	return nil
}

func emitVarStmtIncr(g *Generator, v data.GetValue) error {
	n := v.(*node.VarStmtIncr)
	g.printf("&node.VarStmtIncr{\n")
	g.indent++
	g.printf("Node: node.NewNode(from),\n")
	g.printf("VarIdx: %d,\n", n.VarIdx)
	g.printf("Fallback: ")
	if err := g.Emit(n.Fallback); err != nil {
		return err
	}
	g.printf(",\n")
	g.indent--
	g.printf("}")
	return nil
}

func emitVarIntLe(g *Generator, v data.GetValue) error {
	n := v.(*node.VarIntLe)
	g.printf("&node.VarIntLe{\n")
	g.indent++
	g.printf("Node: node.NewNode(from),\n")
	g.printf("VarIdx: %d,\n", n.VarIdx)
	g.printf("Lit: %d,\n", n.Lit)
	g.printf("Le: ")
	if err := g.Emit(n.Le); err != nil {
		return err
	}
	g.printf(",\n")
	g.indent--
	g.printf("}")
	return nil
}

func emitClassValue(g *Generator, v data.GetValue) error {
	return g.emitClassAnnotation(v.(*data.ClassValue))
}

func (g *Generator) emitClassAnnotation(cv *data.ClassValue) error {
	if cv == nil {
		g.printf("nil")
		return nil
	}
	switch c := cv.Class.(type) {
	case *annotation.RouteClass:
		g.needAnnotationImport()
		g.printf("annotation.CompiledRouteValue(%q)", c.Prefix())
	case *annotation.ControllerClass:
		g.needAnnotationImport()
		g.printf("annotation.CompiledControllerValue(%q)", "")
	case *annotation.GetMappingClass:
		g.needAnnotationImport()
		g.printf("annotation.CompiledGetMappingValue(%q)", c.Path())
	case *annotation.PostMappingClass:
		g.needAnnotationImport()
		g.printf("annotation.CompiledPostMappingValue(%q)", c.Path())
	case *annotation.PutMappingClass:
		g.needAnnotationImport()
		g.printf("annotation.CompiledPutMappingValue(%q)", c.Path())
	case *annotation.DeleteMappingClass:
		g.needAnnotationImport()
		g.printf("annotation.CompiledDeleteMappingValue(%q)", c.Path())
	case *annotation.MiddlewareClass:
		g.needAnnotationImport()
		g.printf("annotation.CompiledMiddlewareValue(%q)", c.ClassName())
	default:
		return newEmitError(g.file, cv, "unsupported annotation type "+reflect.TypeOf(c).String())
	}
	return nil
}

func emitArray(g *Generator, v data.GetValue) error {
	n := v.(*node.Array)
	if len(n.Keys) > 0 {
		g.printf("node.NewArrayWithKeys(from,\n")
		g.indent++
		g.printf("[]data.GetValue{\n")
		g.indent++
		for _, elem := range n.V {
			if err := g.Emit(elem); err != nil {
				return err
			}
			g.printf(",\n")
		}
		g.indent--
		g.printf("},\n")
		g.printf("[]node.KvPair{\n")
		g.indent++
		for _, kv := range n.Keys {
			g.printf("{\n")
			g.indent++
			g.printf("Key: ")
			if err := g.Emit(kv.Key); err != nil {
				return err
			}
			g.printf(",\n")
			g.printf("Value: ")
			if err := g.Emit(kv.Value); err != nil {
				return err
			}
			g.printf(",\n")
			g.indent--
			g.printf("},\n")
		}
		g.indent--
		g.printf("},\n")
		g.indent--
		g.printf(")")
		return nil
	}
	g.printf("node.NewArray(from, []data.GetValue{\n")
	g.indent++
	for _, elem := range n.V {
		if err := g.Emit(elem); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("})")
	return nil
}

func emitNamespace(g *Generator, v data.GetValue) error {
	n := v.(*node.Namespace)
	g.printf("node.NewNamespace(from, %q, []data.GetValue{\n", n.Name)
	g.indent++
	for _, stmt := range n.Statements {
		if err := g.Emit(stmt); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("})")
	return nil
}

func emitNewExpression(g *Generator, v data.GetValue) error {
	n := v.(*node.NewExpression)
	g.printf("node.NewNewExpression(from,\n")
	g.indent++
	g.printf("%q,\n", n.ClassName)
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Arguments {
		if err := g.Emit(arg); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
	return nil
}

func emitNewVariableExpression(g *Generator, v data.GetValue) error {
	n := v.(*node.NewVariableExpression)
	g.printf("node.NewNewVariableExpression(from,\n")
	g.indent++
	if err := g.Emit(n.ClassNameExpr); err != nil {
		return err
	}
	g.printf(",\n")
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Arguments {
		if err := g.Emit(arg); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
	return nil
}

func emitNewExpressionDynamic(g *Generator, v data.GetValue) error {
	n := v.(*node.NewExpressionDynamic)
	g.printf("node.NewNewExpressionDynamic(from,\n")
	g.indent++
	if err := g.Emit(n.ClassExpr); err != nil {
		return err
	}
	g.printf(",\n")
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Arguments {
		if err := g.Emit(arg); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
	return nil
}

func emitNewSelfExpression(g *Generator, v data.GetValue) error {
	n := v.(*node.NewSelfExpression)
	g.printf("node.NewNewSelfExpression(from,\n")
	g.indent++
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Arguments {
		if err := g.Emit(arg); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
	return nil
}

func emitNewStaticExpression(g *Generator, v data.GetValue) error {
	n := v.(*node.NewStaticExpression)
	g.printf("node.NewNewStaticExpression(from,\n")
	g.indent++
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Arguments {
		if err := g.Emit(arg); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
	return nil
}

func emitInitClass(g *Generator, v data.GetValue) error {
	n := v.(*node.InitClass)
	g.printf("node.NewInitClass(from,\n")
	g.indent++
	g.printf("%q,\n", n.ClassName)
	g.printf("map[string]data.GetValue{\n")
	g.indent++
	for k, val := range n.KV {
		g.printf("%q: ", k)
		if err := g.Emit(val); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
	return nil
}

func emitKv(g *Generator, v data.GetValue) error {
	n := v.(*node.Kv)
	g.printf("node.NewKv(from, []node.KvPair{\n")
	g.indent++
	for _, kv := range n.V {
		g.printf("{\n")
		g.indent++
		g.printf("Key: ")
		if err := g.Emit(kv.Key); err != nil {
			return err
		}
		g.printf(",\n")
		g.printf("Value: ")
		if err := g.Emit(kv.Value); err != nil {
			return err
		}
		g.printf(",\n")
		g.indent--
		g.printf("},\n")
	}
	g.indent--
	g.printf("})")
	return nil
}

func emitRange(g *Generator, v data.GetValue) error {
	n := v.(*node.Range)
	g.printf("node.NewRange(from,\n")
	g.indent++
	if n.Array != nil {
		if err := g.Emit(n.Array); err != nil {
			return err
		}
	} else {
		g.printf("nil")
	}
	g.printf(",\n")
	if n.Start != nil {
		if err := g.Emit(n.Start); err != nil {
			return err
		}
	} else {
		g.printf("nil")
	}
	g.printf(",\n")
	if n.Stop != nil {
		if err := g.Emit(n.Stop); err != nil {
			return err
		}
	} else {
		g.printf("nil")
	}
	g.printf(",\n")
	g.indent--
	g.printf(")")
	return nil
}

func emitIncludeStatement(g *Generator, v data.GetValue) error {
	n := v.(*node.IncludeStatement)
	g.printf("node.NewIncludeStatement(from,\n")
	g.indent++
	if err := g.Emit(n.Expr); err != nil {
		return err
	}
	g.printf(",\n")
	g.printf("%v,\n", n.Once)
	g.printf("%v,\n", n.Required)
	g.indent--
	g.printf(")")
	return nil
}

func emitConstStatement(g *Generator, v data.GetValue) error {
	n := v.(*node.ConstStatement)
	g.printf("node.NewConstStatement(from,\n")
	g.indent++
	if n.Val != nil {
		g.printf("node.NewVariable(from, %q, %d, nil),\n", n.Val.GetName(), n.Val.GetIndex())
	} else {
		g.printf("nil,\n")
	}
	if err := g.Emit(n.Initializer); err != nil {
		return err
	}
	g.printf(",\n")
	g.indent--
	g.printf(")")
	return nil
}

func emitBinaryAssignVariable(g *Generator, v data.GetValue) error {
	n := v.(*node.BinaryAssignVariable)
	g.printf("&node.BinaryAssignVariable{Node: node.NewNode(from), Left: ")
	if err := g.emitVariable(n.Left); err != nil {
		return err
	}
	g.printf(", Right: ")
	if err := g.Emit(n.Right); err != nil {
		return err
	}
	g.printf("}")
	return nil
}

func emitBinaryAssignVariableList(g *Generator, v data.GetValue) error {
	n := v.(*node.BinaryAssignVariableList)
	g.printf("&node.BinaryAssignVariableList{Node: node.NewNode(from), Left: ")
	if err := g.emitVariableList(n.Left); err != nil {
		return err
	}
	g.printf(", Right: ")
	if err := g.Emit(n.Right); err != nil {
		return err
	}
	g.printf("}")
	return nil
}

func (g *Generator) emitParamList(params []data.GetValue) error {
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, p := range params {
		if err := g.Emit(p); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("}")
	return nil
}

func (g *Generator) emitClassProperty(p data.Property) error {
	cp, ok := p.(*node.ClassProperty)
	if !ok {
		return newEmitError(g.file, nil, "unsupported property type "+reflect.TypeOf(p).String())
	}
	mod := modifierName(cp.GetModifier())
	if cp.IsPromoted {
		g.printf("node.NewPropertyWithPromoted(from, %q, %q, %v, %v, %v, ", cp.Name, mod, cp.IsStatic, cp.IsReadonly, cp.IsPromoted)
	} else if cp.IsReadonly {
		g.printf("node.NewPropertyWithReadonly(from, %q, %q, %v, %v, ", cp.Name, mod, cp.IsStatic, cp.IsReadonly)
	} else {
		g.printf("node.NewProperty(from, %q, %q, %v, ", cp.Name, mod, cp.IsStatic)
	}
	if cp.DefaultValue != nil {
		if err := g.Emit(cp.DefaultValue); err != nil {
			return err
		}
	} else {
		g.printf("nil")
	}
	if cp.Type != nil {
		g.printf(", ")
		g.genTypes(cp.Type)
	}
	g.printf(")")
	return nil
}

func (g *Generator) emitClassMethod(method data.Method) error {
	if am, ok := method.(*node.AbstractMethod); ok {
		g.printf("node.NewAbstractMethod(")
		if err := g.emitClassMethodBody(am.ClassMethod); err != nil {
			return err
		}
		g.printf(")")
		return nil
	}
	cm, ok := method.(*node.ClassMethod)
	if !ok {
		return newEmitError(g.file, nil, "unsupported class method type "+reflect.TypeOf(method).String())
	}
	if len(cm.Annotations) == 0 {
		return g.emitClassMethodBody(cm)
	}
	g.printf("func() data.Method {\n")
	g.indent++
	g.printf("m := ")
	if err := g.emitClassMethodBody(cm); err != nil {
		return err
	}
	g.printf(".(*node.ClassMethod)\n")
	for _, ann := range cm.Annotations {
		g.printf("m.AddAnnotations(")
		if err := g.emitClassAnnotation(ann); err != nil {
			return err
		}
		g.printf(")\n")
	}
	g.printf("return m\n")
	g.indent--
	g.printf("}()")
	return nil
}

func (g *Generator) emitClassMethodBody(cm *node.ClassMethod) error {
	mod := modifierName(cm.GetModifier())
	g.printf("node.NewMethod(from, %q, %q, %v, ", cm.Name, mod, cm.IsStatic)
	if err := g.emitParamList(cm.Params); err != nil {
		return err
	}
	g.printf(", []data.GetValue{\n")
	g.indent++
	for _, stmt := range cm.Body {
		if err := g.Emit(stmt); err != nil {
			return err
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("}, ")
	g.genMethodVars(cm.GetVariables())
	g.printf(", ")
	g.genTypes(cm.Ret)
	g.printf(")")
	return nil
}

func (g *Generator) emitVariableList(n *node.VariableList) error {
	g.printf("node.NewVariableList([]*node.VariableExpression{\n")
	g.indent++
	for _, v := range n.Vars {
		g.printf("&node.VariableExpression{Node: node.NewNode(from), Name: %q, Index: %d},\n", v.Name, v.Index)
	}
	g.indent--
	g.printf("})")
	return nil
}
