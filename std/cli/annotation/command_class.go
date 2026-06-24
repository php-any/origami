package annotation

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// CommandClass 命令注解（类似 Symfony Console Command）
// 标注命令类，注册到 CLI 应用中。
//
// 用法示例：
//
//	#[Command(name: "greet", description: "Greet someone")]
//	class GreetCommand {
//	    public function execute(): void {
//	        echo "Hello!";
//	    }
//	}
type CommandClass struct {
	node.Node
	construct data.Method
}

func (c *CommandClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newCommand()

	return data.NewClassValue(&CommandClass{
		construct: &CommandConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (c *CommandClass) GetName() string    { return "Cli\\Annotation\\Command" }
func (c *CommandClass) GetExtend() *string { return nil }
func (c *CommandClass) GetImplements() []string {
	return []string{node.TypeFeature, node.TypeTargetClass}
}
func (c *CommandClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (c *CommandClass) GetPropertyList() []data.Property           { return []data.Property{} }
func (c *CommandClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return c.construct, true
	}
	return nil, false
}
func (c *CommandClass) GetMethods() []data.Method {
	return []data.Method{c.construct}
}
func (c *CommandClass) GetConstruct() data.Method { return c.construct }

// Command 命令元信息
type Command struct {
	name        string
	description string
	target      any // 被注解的命令类
}

// GetName 获取命令名称
func (c *Command) GetName() string {
	return c.name
}

// GetDescription 获取命令描述
func (c *Command) GetDescription() string {
	return c.description
}

// registeredCommands 存储已注册的命令
var registeredCommands = make(map[string]*Command)

func newCommand() *Command { return &Command{} }

type CommandConstructMethod struct{ cmd *Command }

func (m *CommandConstructMethod) GetName() string            { return "__construct" }
func (m *CommandConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *CommandConstructMethod) GetIsStatic() bool          { return false }
func (m *CommandConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "description", 1, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewAnnotationTargetParameter(nil, 2),
	}
}
func (m *CommandConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, nil),
		node.NewVariable(nil, "description", 1, nil),
		node.NewAnnotationTargetVariable(nil, 2),
	}
}
func (m *CommandConstructMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *CommandConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	name, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}
	description, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, nil
	}

	if v, ok := name.(*data.StringValue); ok {
		m.cmd.name = v.AsString()
	}
	if v, ok := description.(*data.StringValue); ok {
		m.cmd.description = v.AsString()
	}

	if target, ok := ctx.GetIndexValue(2); ok {
		if anyT, ok := target.(*data.AnyValue); ok {
			m.cmd.target = anyT.Value
		}
	}

	// 编译模式下只解析参数
	if data.CompileMode {
		return nil, nil
	}

	// 注册命令到全局命令注册表
	if m.cmd.name != "" {
		registeredCommands[m.cmd.name] = m.cmd
	}

	return nil, nil
}

// GetCommandName 获取命令名称
func (m *CommandConstructMethod) GetCommandName() string {
	return m.cmd.name
}

// GetCommandDescription 获取命令描述
func (m *CommandConstructMethod) GetCommandDescription() string {
	return m.cmd.description
}

// GetRegisteredCommands 获取所有已注册的命令
func GetRegisteredCommands() map[string]*Command {
	return registeredCommands
}

// ExecuteCommand 执行指定命令
func ExecuteCommand(ctx data.Context, commandName string) data.Control {
	cmd, exists := registeredCommands[commandName]
	if !exists {
		return data.NewErrorThrow(nil, errors.New("命令未找到: "+commandName))
	}

	cls, ok := cmd.target.(*node.ClassStatement)
	if !ok {
		return nil
	}

	baseCtx := ctx.CreateBaseContext()
	cv := data.NewClassValue(cls, baseCtx)

	// 查找 execute 方法
	method, has := cls.GetMethod("execute")
	if !has {
		return data.NewErrorThrow(nil, errors.New("命令 "+commandName+" 缺少 execute 方法"))
	}

	fnCtx := cv.CreateContext(method.GetVariables())
	_, acl := method.Call(fnCtx)
	return acl
}
