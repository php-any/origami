package parser

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/php-any/origami/data"

	"github.com/php-any/origami/lexer"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// Parser 表示解析器
type Parser struct {
	vm               data.VM
	source           *string
	lexer            *lexer.Lexer      // 词法分析器
	tokens           []lexer.Token     // 词法单元列表
	position         int               // 当前处理位置
	errors           []data.Control    // 错误列表
	scopeManager     *ScopeManager     // 作用域管理器
	expressionParser *ExpressionParser // 表达式解析器

	identTryString bool

	namespace        *node.Namespace
	uses             map[string]string // 类引用
	classPathManager ClassPathManager  // 类路径管理器
}

// NewParser 创建一个新的解析器
func NewParser() *Parser {
	p := &Parser{
		lexer:            lexer.NewLexer(),
		tokens:           make([]lexer.Token, 0),
		position:         0,
		errors:           make([]data.Control, 0),
		scopeManager:     NewScopeManager(),
		uses:             make(map[string]string),
		classPathManager: NewDefaultClassPathManager(),
	}
	p.AddScanNamespace("app", "./")
	p.expressionParser = NewExpressionParser(p)
	return p
}

// reset 重置解析器状态
func (p *Parser) reset() {
	p.tokens = make([]lexer.Token, 0)
	p.position = 0
	p.errors = make([]data.Control, 0)
	p.uses = make(map[string]string)
	p.namespace = nil
	p.scopeManager = NewScopeManager()
}

func (p *Parser) Clone() *Parser {
	// 创建新的解析器实例
	cloned := &Parser{
		vm:               p.vm,    // VM 是共享的，不需要克隆
		source:           nil,     // 字符串指针，共享即可
		lexer:            p.lexer, // 创建新的词法分析器
		tokens:           make([]lexer.Token, 0),
		position:         0,
		errors:           make([]data.Control, 0),
		scopeManager:     NewScopeManager(),  // 创建新的作用域管理器
		expressionParser: p.expressionParser, // 稍后设置
		identTryString:   p.identTryString,
		namespace:        nil, // 命名空间节点，共享即可
		uses:             make(map[string]string),
		classPathManager: p.classPathManager, // 类路径管理器是共享的
	}
	cloned.expressionParser = NewExpressionParser(cloned)
	return cloned
}

func (p *Parser) SetVM(vm data.VM) {
	p.vm = vm
}

// ParseFile 解析文件
func (p *Parser) ParseFile(filename string) (*node.Program, error) {
	// 读取文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 重置解析器状态
	p.reset()

	p.source = &filename
	// 进行分词
	p.tokens = p.lexer.Tokenize(string(content))

	// 解析程序
	program := p.parseProgram()

	// 检查是否有错误
	if len(p.errors) > 0 {
		return nil, errors.New(p.errors[0].AsString())
	}

	return program, nil
}

// parseProgram 解析程序
func (p *Parser) parseProgram() *node.Program {
	statements := make([]node.Statement, 0)

	last := 0
	// 解析所有语句
	for !p.isEOF() {
		stmt, acl := p.parseStatement()
		if acl != nil {
			p.addControl(acl)
		}
		if stmt != nil {
			if n, ok := stmt.(*node.Namespace); ok {
				p.namespace = n
				statements = append(statements, stmt)
			} else {
				if n, ok := stmt.(*node.UseStatement); ok {
					p.uses[n.Alias] = n.Namespace
					continue
				}

				if p.namespace != nil {
					p.namespace.Statements = append(p.namespace.Statements, stmt)
				} else {
					statements = append(statements, stmt)
				}
			}
		} else if p.position != last {
			last = p.position
		} else {
			p.addControl(data.NewErrorThrow(p.newFrom(), errors.New("无法识别语句")))
			return nil
		}
	}

	return node.NewProgram(nil, statements)
}

// current 返回当前词法单元
func (p *Parser) current() lexer.Token {
	if p.position >= len(p.tokens) {
		return lexer.Token{Type: token.EOF}
	}
	return p.tokens[p.position]
}

// peek 向前查看指定位置的token
func (p *Parser) peek(offset int) lexer.Token {
	pos := p.position + offset
	if pos >= len(p.tokens) {
		return lexer.Token{Type: token.EOF}
	}
	return p.tokens[pos]
}

// 检查后续单词的可能类型
func (p *Parser) checkPositionIs(position int, checks ...token.TokenType) bool {
	if len(checks) == 1 && checks[0] == token.EOF {
		if p.position+position >= len(p.tokens) {
			return true
		}
		for _, check := range checks {
			if p.tokens[p.position+position].Type == check {
				return true
			}
		}
	} else {
		if p.position+position >= len(p.tokens) {
			return false
		}
	}

	for _, check := range checks {
		if p.tokens[p.position+position].Type == check {
			return true
		}
	}
	return false
}

func (p *Parser) currentIsTypeOrEOF(check token.TokenType) bool {
	if p.position >= len(p.tokens) {
		return true
	}
	return p.tokens[p.position].Type == check
}

// 打印剩余的
func (p *Parser) printRemaining() {
	for !p.isEOF() {
		fmt.Println(p.current().Literal)
		p.next()
	}
}

func (p *Parser) GetVariables() []data.Variable {
	return p.scopeManager.CurrentScope().GetVariables()
}

// next 移动到下一个词法单元
func (p *Parser) next() {
	p.position++
}

func (p *Parser) nextAndCheck(t token.TokenType) {
	if p.current().Type != t {
		p.addError("检查符号不一致")
	}
	p.position++
}

func (p *Parser) nextAndCheckStip(t token.TokenType) {
	if p.current().Type == t {
		p.position++
	}
}

// isEOF 检查是否到达文件末尾
func (p *Parser) isEOF() bool {
	return p.position >= len(p.tokens)
}

// addError
func (p *Parser) addError(err string) {
	from := node.NewTokenFrom(p.source, p.current().Start, p.current().End, p.current().Line, p.current().Pos)
	p.errors = append(p.errors, data.NewErrorThrow(from, errors.New(err)))

	// 打印详细的错误信息
	p.printDetailedError(err, from)

	panic(err)
}

func (p *Parser) addControl(acl data.Control) {
	p.vm.ThrowControl(acl)
}

func (p *Parser) ShowControl(acl data.Control) {
	err := acl.AsString()

	if acl, ok := acl.(node.GetFrom); ok {
		from := acl.GetFrom()
		p.errors = append(p.errors, data.NewErrorThrow(from, errors.New(err)))
		// 打印详细的错误信息
		p.printDetailedError(err, from)
	} else {
		from := node.NewTokenFrom(p.source, p.current().Start, p.current().End, p.current().Line, p.current().Pos)
		p.errors = append(p.errors, data.NewErrorThrow(from, errors.New(err)))
		// 打印详细的错误信息
		p.printDetailedError(err, from)
	}
}

func (p *Parser) GetStart() int {
	return p.current().Start
}

// Deprecated: 使用 NewFromBuilder() 或其他新方法替代
func (p *Parser) NewTokenFrom(start int) *node.TokenFrom {
	return node.NewTokenFrom(p.source, start, p.current().End, p.current().Line, p.current().Pos)
}

// StartPosition 开始位置跟踪，返回当前位置
func (p *Parser) StartPosition() int {
	return p.position
}

// EndPosition 结束位置跟踪，返回当前位置
func (p *Parser) EndPosition() int {
	return p.position
}

// FromPositionRange 从位置范围创建From信息
func (p *Parser) FromPositionRange(startPos, endPos int) *node.TokenFrom {
	if startPos >= len(p.tokens) || endPos >= len(p.tokens) {
		return p.FromCurrentToken()
	}

	startToken := p.tokens[startPos]
	endToken := p.tokens[endPos]

	// 创建 TokenFrom 并设置结束位置
	tf := node.NewTokenFrom(p.source, startToken.Start, endToken.End, startToken.Line, startToken.Pos)

	// 总是设置结束位置，确保位置信息完整
	// 即使 startPos == endPos，我们也需要正确的结束位置信息
	tf.SetEndPosition(endToken.Line, endToken.Pos)

	return tf
}

// isTokensAdjacent 检查两个 token 是否相邻（没有空白字符或其他分隔符）
func (p *Parser) isTokensAdjacent(token1, token2 lexer.Token) bool {
	// 如果第一个 token 的结束位置等于第二个 token 的开始位置，说明它们是相邻的
	return token1.End == token2.Start
}

// 获取类的完整路径, 类定义自己不用, 但是继承、实现接口需要调用
func (p *Parser) getClassName(try bool) (string, data.Control) {
	// 获取完整的类名路径
	var parts []string
	for {
		parts = append(parts, p.current().Literal)
		p.next()

		if p.current().Type != token.NAMESPACE_SEPARATOR {
			break
		}
		p.next()
	}
	className := ""
	for _, part := range parts {
		if className == "" {
			className = part
		} else {
			className = className + "\\" + part
		}
	}
	if len(parts) == 1 {
		// 如果只有一个单词, 则认为可能是别名
		if full, ok := p.uses[className]; ok {
			return full, nil
		}
		// 也有可能是同一个包内的类
		if try {
			if full, ok := p.findFullClassNameByNamespace(className); ok {
				return full, p.tryLoadClass(full)
			}
		}
	}

	return className, nil
}

// parseStatement 解析语句
func (p *Parser) parseStatement() (data.GetValue, data.Control) {
	return p.expressionParser.Parse()
}

// 只会获取单个值, 不会有表达式, 并且必须有值, 没有就是错误
func (p *Parser) parseValue() (data.GetValue, bool) {
	switch p.current().Type {
	case token.INT:
		value := p.current().Literal
		from := p.FromCurrentToken()
		p.next()
		return node.NewIntLiteral(from, value), true
	case token.FLOAT:
		value := p.current().Literal
		from := p.FromCurrentToken()
		p.next()
		return node.NewFloatLiteral(from, value), true
	case token.STRING:
		value := p.current().Literal
		from := p.FromCurrentToken()
		p.next()
		return node.NewStringLiteral(from, value), true
	case token.TRUE:
		from := p.FromCurrentToken()
		p.next()
		return node.NewBooleanLiteral(from, true), true
	case token.FALSE:
		from := p.FromCurrentToken()
		p.next()
		return node.NewBooleanLiteral(from, false), true
	case token.NULL:
		from := p.FromCurrentToken()
		p.next()
		return node.NewNullLiteral(from), true
	case token.THIS:
		stmt, acl := NewThisParser(p).Parse()
		_ = acl
		return stmt, true
	case token.VARIABLE:
		vp := &VariableParser{p}
		return vp.parseVariable(), true
	case token.IDENTIFIER:
		vp := &VariableParser{p}
		return vp.parseVariable(), true
	default:
		return nil, false
	}
}

// parseBlock 解析语句块
func (p *Parser) parseBlock() []data.GetValue {
	statements := make([]data.GetValue, 0)

	// 检查是否是语句块开始
	if p.current().Type != token.LBRACE {
		// 如果不是语句块，则解析单个语句
		stmt, acl := p.parseStatement()
		if acl != nil {
			p.addControl(acl)
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
		return statements
	}

	// 跳过左花括号
	p.next()

	for p.checkPositionIs(0, token.SEMICOLON) {
		p.next()
	}

	// 解析语句块中的所有语句
	for !p.isEOF() && p.current().Type != token.RBRACE {
		stmt, acl := p.parseStatement()
		if acl != nil {
			p.addControl(acl)
		}
		for p.checkPositionIs(0, token.SEMICOLON) {
			p.next()
		}
		if stmt != nil {
			statements = append(statements, stmt)
		} else {
			p.addControl(data.NewErrorThrow(p.newFrom(), errors.New("语法块无法识别")))
		}
	}

	// 跳过右花括号
	p.nextAndCheck(token.RBRACE)

	return statements
}

func (p *Parser) AddScanNamespace(namespace string, path string) {
	// 使用类路径管理器添加命名空间
	p.classPathManager.AddNamespace(namespace, path)
}

// 默认 try = true
func (p *Parser) findFullClassNameByNamespace(name string, try ...bool) (string, bool) {
	tryName := name

	if full, ok := p.uses[name]; ok {
		return full, true
	}

	if p.namespace != nil {
		tryName = p.namespace.GetName() + "\\" + name
	}
	// 本包
	if stmt, ok := p.vm.GetClass(tryName); ok {
		return stmt.GetName(), true
	}
	if stmt, ok := p.vm.GetInterface(tryName); ok {
		return stmt.GetName(), true
	}
	// 顶命名
	if stmt, ok := p.vm.GetClass(name); ok {
		return stmt.GetName(), true
	}
	if stmt, ok := p.vm.GetInterface(name); ok {
		return stmt.GetName(), true
	}

	if len(try) > 0 && !try[0] {
		return "", false
	}

	// 尝试加载同目录的同名文件
	p.tryLoadClass(tryName)

	// 加载成功，再次尝试查找类
	if full, ok := p.findFullClassNameByNamespace(name, false); ok {
		return full, true
	}
	return "", false
}

func (p *Parser) findFullFunNameByNamespace(name string) (string, bool) {
	if full, ok := p.uses[name]; ok {
		return full, true
	}
	tryName := name
	if p.namespace != nil && strings.Index(name, "\\") == -1 {
		tryName = p.namespace.GetName() + "\\" + name
	}
	if stmt, ok := p.vm.GetFunc(tryName); ok {
		return stmt.GetName(), true
	}

	return "", false
}

func (p *Parser) newFrom() data.From {
	return p.FromCurrentToken()
}

// 尝试加载类
func (p *Parser) tryLoadClass(full string) data.Control {
	if _, ok := p.vm.GetClass(full); ok {
		return nil
	}
	if _, ok := p.vm.GetInterface(full); ok {
		return nil
	}

	// 使用类路径管理器加载类
	return p.classPathManager.LoadClass(full, p)
}

// SetClassPathManager 设置类路径管理器
func (p *Parser) SetClassPathManager(manager ClassPathManager) {
	p.classPathManager = manager
}

// GetClassPathManager 获取类路径管理器
func (p *Parser) GetClassPathManager() ClassPathManager {
	return p.classPathManager
}

// ParseExpressionFromString 从字符串解析表达式
func (p *Parser) ParseExpressionFromString(exprStr string) (data.GetValue, data.Control) {
	// 保存当前状态
	originalTokens := p.tokens
	originalPosition := p.position

	// 重置解析器状态
	p.tokens = make([]lexer.Token, 0)
	p.position = 0

	// 对表达式字符串进行分词
	p.tokens = p.lexer.Tokenize(exprStr)

	// 使用表达式解析器解析
	exprParser := NewExpressionParser(p)
	result, ctl := exprParser.Parse()

	// 恢复原始状态
	p.tokens = originalTokens
	p.position = originalPosition

	return result, ctl
}

// ParseString 从字符串解析程序
func (p *Parser) ParseString(content string, filePath string) (*node.Program, error) {
	// 保存当前状态
	originalTokens := p.tokens
	originalPosition := p.position
	originalSource := p.source

	// 重置解析器状态
	p.reset()

	// 设置源文件路径，确保符号位置信息正确
	p.source = &filePath

	// 进行分词
	p.tokens = p.lexer.Tokenize(content)

	// 解析程序
	program := p.parseProgram()

	// 检查是否有错误
	if len(p.errors) > 0 {
		return nil, errors.New(p.errors[0].AsString())
	}

	// 恢复原始状态
	p.tokens = originalTokens
	p.position = originalPosition
	p.source = originalSource

	return program, nil
}
