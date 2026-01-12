package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// 枚举反 desugar 相关的内部约定集中在这里，避免散落硬编码
const (
	// enum 默认的底层类型（目前只实现 string backed enum）
	defaultEnumBackingType = "string"
	// enum 继承的基础类名，用于支持 instanceof BackedEnum / ->value 语义
	enumBaseClassName = "BackedEnum"
)

// EnumParser 解析 PHP 8.1 enum 声明，并将其反 desugar 为继承 \BackedEnum 的普通类。
//
// 目前支持的子集：
//
//	enum Status: string {
//	    case OPEN = 'open';
//	    case CLOSED = 'closed';
//	}
//
// 反 desugar 逻辑大致等价于：
//
//	class Status extends \BackedEnum {
//	    // BackedEnum 提供 public string $value 和 __construct(string $value)
//	    public static $OPEN = new Status('open');
//	    public static $CLOSED = new Status('closed');
//	}
//
// 因此：
//   - Status::OPEN 得到一个 Status 实例
//   - $value instanceof \BackedEnum 成立
//   - $value->value 为底层值
type EnumParser struct {
	*Parser
}

func NewEnumParser(p *Parser) StatementParser {
	return &EnumParser{Parser: p}
}

func (p *EnumParser) Parse() (data.GetValue, data.Control) {
	// 跳过 enum 关键字
	p.next()
	tracker := p.StartTracking()

	// 解析枚举名
	if p.current().Type() != token.IDENTIFIER {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum 后缺少名称"))
	}
	enumName := p.current().Literal()
	p.next()

	// 加上命名空间前缀
	if p.namespace != nil {
		enumName = p.namespace.GetName() + "\\" + enumName
	}
	p.vm.SetClassPathCache(enumName, *p.source)

	// 解析可选的底层类型: enum Status: string
	// 当前解析器接受类型标注，但不再依赖其具体值做行为分支（避免 string-only 硬编码）
	if p.current().Type() == token.COLON {
		p.next()
		if p.current().Type() != token.IDENTIFIER {
			return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum 底层类型缺失或非法"))
		}
		// 忽略具体类型名称，只做语法校验
		_ = p.current().Literal()
		p.next()
	}

	// 解析枚举体
	if p.current().Type() != token.LBRACE {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum 声明后缺少 '{'"))
	}
	p.next()

	type enumCase struct {
		name  string
		value data.GetValue // 底层值表达式
	}
	var cases []enumCase
	methods := map[string]data.Method{}
	staticMethods := map[string]data.Method{}

	for !p.currentIsTypeOrEOF(token.RBRACE) {
		if p.current().Type() == token.SEMICOLON {
			p.next()
			continue
		}

		// 先尝试解析注解
		var memberAnnotations []*node.Annotation
		cp := &ClassParser{
			Parser:               p.Parser,
			FunctionParserCommon: NewFunctionParserCommon(p.Parser),
		}
		for p.checkPositionIs(0, token.AT, token.HASH) {
			ann, acl := cp.parseAnnotation()
			if acl != nil {
				return nil, acl
			}
			if ann != nil {
				memberAnnotations = append(memberAnnotations, ann)
			}
		}

		// 解析访问修饰符（方法需要）
		modifier := cp.parseModifier()

		// 解析 static 关键字
		isStatic := false
		if p.current().Type() == token.STATIC {
			isStatic = true
			p.next()
		}

		// 检查是否是 case 声明
		if p.current().Type() == token.CASE {
			p.next()

			if p.current().Type() != token.IDENTIFIER {
				return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum case 缺少名称"))
			}
			caseName := p.current().Literal()
			p.next()

			var val data.GetValue

			// 支持: case OPEN = 'open';
			if p.current().Type() == token.ASSIGN {
				p.next()
				exprParser := NewExpressionParser(p.Parser)
				var acl data.Control
				val, acl = exprParser.Parse()
				if acl != nil {
					return nil, acl
				}
			} else {
				// 未显式赋值时，默认使用 case 名称的**字符串**作为枚举底层值。
				// 底层值类型本身不再被这里限制，BackedEnum 也允许任意类型。
				from := tracker.EndBefore()
				val = node.NewStringLiteralByAst(from, caseName)
			}

			// 跳过可选分号
			if p.current().Type() == token.SEMICOLON {
				p.next()
			}

			cases = append(cases, enumCase{name: caseName, value: val})
		} else if p.current().Type() == token.FUNC {
			// 解析方法
			method, _, acl := cp.parseMethodWithAnnotations(modifier, isStatic, false, memberAnnotations, nil, nil)
			if acl != nil {
				return nil, acl
			}
			if method != nil {
				if isStatic {
					staticMethods[method.GetName()] = method
				} else {
					methods[method.GetName()] = method
				}
			}
		} else {
			return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum 体内只支持 case 声明和方法声明"))
		}
	}

	// 跳过枚举体结束右括号
	p.next()

	// 不直接定义属性：
	//   - value 属性与构造函数都由 BackedEnum 提供
	//   - 这里保持 properties 为空
	properties := []data.Property{}

	// enum 反 desugar 成：
	//   class EnumName extends BackedEnum { ... }
	// 注意：这里不加前导反斜杠，保持与 BackedEnumClass.GetName() 一致
	extends := enumBaseClassName
	classStmt := node.NewClassStatement(
		tracker.EndBefore(),
		enumName,
		extends,
		nil,
		properties,
		methods,
	)
	classStmt.StaticMethods = staticMethods

	// 继承 BackedEnum 的构造函数：确保 new Status('open') 会运行 BackedEnum::__construct
	if parent, ok := p.vm.GetClass(extends); ok {
		if ctor := parent.GetConstruct(); ctor != nil {
			classStmt.Construct = ctor
		}
	}

	// 将枚举作为类注册到 VM
	if acl := p.vm.AddClass(classStmt); acl != nil {
		return nil, acl
	}

	// 为每个 case 注入一个静态属性：public static $CASE = new EnumName(<value>);
	// 此时类已注册到 VM，可以安全地构造枚举实例
	for _, ccase := range cases {
		if ccase.value == nil {
			continue
		}
		from := tracker.EndBefore()
		newExpr := node.NewNewExpression(from, enumName, []data.GetValue{ccase.value})

		ctx := p.vm.CreateContext(nil)
		v, acl := newExpr.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if val, ok := v.(data.Value); ok {
			classStmt.StaticProperty.Store(ccase.name, val)
		}
	}

	return classStmt, nil
}
