package parser

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// UseParser 表示use语句解析器
type UseParser struct {
	*Parser
}

// NewUseParser 创建一个新的use语句解析器
func NewUseParser(parser *Parser) StatementParser {
	return &UseParser{
		parser,
	}
}

// Parse 解析use语句，支持普通use、函数use、常量use、分组use
func (p *UseParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过 use 关键字
	p.next()

	// 解析 use 类型（function/const/class）
	useType := "" // "", "function", "const"
	if p.current().Type() != token.IDENTIFIER {
		if p.checkPositionIs(0, token.FUNC, token.CONST) {
			useType = p.current().Literal()
			p.next()
		} else {
			return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("use 命名空间名称不能为空"))
		}
	}

	// 获取完整的命名空间路径
	namespace := p.current().Literal()
	p.next()

	// 分组 use 语句：use Foo\{Bar, Baz as Qux};
	//                    use function Foo\{bar, baz};
	if p.current().Type() == token.LBRACE {
		return p.parseGroupUse(tracker, namespace, useType)
	}

	// 普通 use 语句
	return p.parseSingleUse(tracker, namespace)
}

// parseSingleUse 解析单个 use 语句：use Foo\Bar; 或 use Foo\Bar as Baz; 或 use function Foo\bar;
func (p *UseParser) parseSingleUse(tracker *PositionTracker, namespace string) (data.GetValue, data.Control) {
	// 检查是否有 as 关键字
	var alias string
	if p.current().Type() == token.AS {
		p.next()
		if p.current().Type() != token.IDENTIFIER {
			return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("as 关键字后需要变量名"))
		}
		alias = p.current().Literal()
		p.next()
	} else {
		// 如果没有 as 关键字，从完整路径中提取最后一个部分作为 alias
		// 例如：从 "a\b\c" 中提取 "c"
		parts := strings.Split(namespace, "\\")
		alias = parts[len(parts)-1]
	}

	// 检查分号
	if p.current().Type() != token.SEMICOLON {
		return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("use 语句缺少分号"))
	}
	p.next()

	from := tracker.EndBefore()
	// 创建 use 语句节点
	return node.NewUseStatement(
		from,
		namespace,
		alias,
	), nil
}

// parseGroupUse 解析分组 use 语句：
//
//	use Foo\{Bar, Baz as Qux};            // 类分组
//	use function Foo\{bar, baz as qux};   // 函数分组
//	use const Foo\{BAR, BAZ};             // 常量分组
//
// PHP 语义：分组 use 相当于把分组内的每个名称展开为独立的 use 语句。
// 分组内每个名称与分组前缀拼接成完整名称，别名默认取最后一个名称段。
func (p *UseParser) parseGroupUse(tracker *PositionTracker, namespace string, useType string) (data.GetValue, data.Control) {
	// 分组前缀必须以反斜杠结尾（PHP 语法要求：use Foo\{Bar};）
	prefix := namespace
	if !strings.HasSuffix(prefix, "\\") {
		prefix += "\\"
	}

	// 跳过 {
	p.next()

	var lastStmt *node.UseStatement
	var lastFrom data.From

	// 解析分组内的导入列表
	for {
		// 检查是否到达 }
		if p.current().Type() == token.RBRACE {
			break
		}

		if p.current().Type() != token.IDENTIFIER {
			return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("分组 use 语句中导入名称不能为空"))
		}

		// 获取导入的名称（可能是 "Bar" 或 "Bar\Baz"）
		importName := p.current().Literal()
		p.next()

		// 检查是否有 as 别名
		alias := ""
		if p.current().Type() == token.AS {
			p.next()
			if p.current().Type() != token.IDENTIFIER {
				return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("as 关键字后需要变量名"))
			}
			alias = p.current().Literal()
			p.next()
		} else {
			// 从导入名称中提取最后一个部分作为默认别名
			parts := strings.Split(importName, "\\")
			alias = parts[len(parts)-1]
		}

		// 拼接完整名称
		fullName := prefix + importName

		// 注册到 use 映射
		p.uses[alias] = fullName

		from := tracker.EndBefore()
		lastFrom = from
		lastStmt = node.NewUseStatement(from, fullName, alias)

		// 检查分隔符
		if p.current().Type() == token.COMMA {
			p.next()
			continue
		}
		if p.current().Type() == token.RBRACE {
			break
		}
		return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("分组 use 语句格式错误，应为 , 或 }"))
	}

	// 跳过 }
	if p.current().Type() != token.RBRACE {
		return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("分组 use 语句缺少 }"))
	}
	p.next()

	// 检查分号
	if p.current().Type() != token.SEMICOLON {
		return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("use 语句缺少分号"))
	}
	p.next()

	if lastStmt == nil {
		return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("分组 use 语句为空"))
	}

	// 返回最后一个 use 语句作为节点（所有别名已在上面注册到 p.uses）
	return node.NewUseStatement(
		lastFrom,
		lastStmt.GetNamespace(),
		lastStmt.GetAlias(),
	), nil
}
