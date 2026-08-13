package openai

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ============================================================
// SystemMessage — 对应 openai-go 的 openai.SystemMessage(content)
// ============================================================

type SystemMessageClass struct{}

func NewSystemMessageClass() data.ClassStmt { return &SystemMessageClass{} }

func (c *SystemMessageClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *SystemMessageClass) GetFrom() data.From                              { return nil }
func (c *SystemMessageClass) GetName() string                                 { return "OpenAI\\SystemMessage" }
func (c *SystemMessageClass) GetExtend() *string                              { return nil }
func (c *SystemMessageClass) GetImplements() []string                         { return nil }
func (c *SystemMessageClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *SystemMessageClass) GetPropertyList() []data.Property                { return nil }
func (c *SystemMessageClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *SystemMessageClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *SystemMessageClass) GetMethods() []data.Method                       { return nil }
func (c *SystemMessageClass) GetConstruct() data.Method {
	return newSingleArgConstruct("OpenAI\\SystemMessage", "content")
}

// ============================================================
// UserMessage — 对应 openai-go 的 openai.UserMessage(content)
// ============================================================

type UserMessageClass struct{}

func NewUserMessageClass() data.ClassStmt { return &UserMessageClass{} }

func (c *UserMessageClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *UserMessageClass) GetFrom() data.From                              { return nil }
func (c *UserMessageClass) GetName() string                                 { return "OpenAI\\UserMessage" }
func (c *UserMessageClass) GetExtend() *string                              { return nil }
func (c *UserMessageClass) GetImplements() []string                         { return nil }
func (c *UserMessageClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *UserMessageClass) GetPropertyList() []data.Property                { return nil }
func (c *UserMessageClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *UserMessageClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *UserMessageClass) GetMethods() []data.Method                       { return nil }
func (c *UserMessageClass) GetConstruct() data.Method {
	return newSingleArgConstruct("OpenAI\\UserMessage", "content")
}

// ============================================================
// AssistantMessage — 对应 openai-go 的 openai.AssistantMessage(content)
// ============================================================

type AssistantMessageClass struct{}

func NewAssistantMessageClass() data.ClassStmt { return &AssistantMessageClass{} }

func (c *AssistantMessageClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *AssistantMessageClass) GetFrom() data.From                              { return nil }
func (c *AssistantMessageClass) GetName() string                                 { return "OpenAI\\AssistantMessage" }
func (c *AssistantMessageClass) GetExtend() *string                              { return nil }
func (c *AssistantMessageClass) GetImplements() []string                         { return nil }
func (c *AssistantMessageClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *AssistantMessageClass) GetPropertyList() []data.Property                { return nil }
func (c *AssistantMessageClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *AssistantMessageClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *AssistantMessageClass) GetMethods() []data.Method                       { return nil }
func (c *AssistantMessageClass) GetConstruct() data.Method {
	return newSingleArgConstruct("OpenAI\\AssistantMessage", "content")
}

// ============================================================
// DeveloperMessage — 对应 openai-go 的 openai.DeveloperMessage(content)
// ============================================================

type DeveloperMessageClass struct{}

func NewDeveloperMessageClass() data.ClassStmt { return &DeveloperMessageClass{} }

func (c *DeveloperMessageClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *DeveloperMessageClass) GetFrom() data.From                              { return nil }
func (c *DeveloperMessageClass) GetName() string                                 { return "OpenAI\\DeveloperMessage" }
func (c *DeveloperMessageClass) GetExtend() *string                              { return nil }
func (c *DeveloperMessageClass) GetImplements() []string                         { return nil }
func (c *DeveloperMessageClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *DeveloperMessageClass) GetPropertyList() []data.Property                { return nil }
func (c *DeveloperMessageClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *DeveloperMessageClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DeveloperMessageClass) GetMethods() []data.Method                       { return nil }
func (c *DeveloperMessageClass) GetConstruct() data.Method {
	return newSingleArgConstruct("OpenAI\\DeveloperMessage", "content")
}

// ============================================================
// ToolMessage — 对应 openai-go 的 openai.ToolMessage(toolCallID, content)
// ============================================================

type ToolMessageClass struct{}

func NewToolMessageClass() data.ClassStmt { return &ToolMessageClass{} }

func (c *ToolMessageClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *ToolMessageClass) GetFrom() data.From                              { return nil }
func (c *ToolMessageClass) GetName() string                                 { return "OpenAI\\ToolMessage" }
func (c *ToolMessageClass) GetExtend() *string                              { return nil }
func (c *ToolMessageClass) GetImplements() []string                         { return nil }
func (c *ToolMessageClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *ToolMessageClass) GetPropertyList() []data.Property                { return nil }
func (c *ToolMessageClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *ToolMessageClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *ToolMessageClass) GetMethods() []data.Method                       { return nil }
func (c *ToolMessageClass) GetConstruct() data.Method {
	return &toolMessageConstruct{}
}

// ============================================================
// 构造函数实现
// ============================================================

// singleArgConstruct 使用 PromotedParameter，运行时自动将参数值赋到对象属性
type singleArgConstruct struct {
	className    string
	propertyName string
}

func newSingleArgConstruct(className, propertyName string) data.Method {
	return &singleArgConstruct{className: className, propertyName: propertyName}
}

func (m *singleArgConstruct) GetName() string            { return "__construct" }
func (m *singleArgConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *singleArgConstruct) GetIsStatic() bool          { return false }
func (m *singleArgConstruct) GetReturnType() data.Types  { return nil }

func (m *singleArgConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewPromotedParameter(nil, m.propertyName, 0, nil, data.NewBaseType("string")),
	}
}

func (m *singleArgConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, m.propertyName, 0, data.NewBaseType("string")),
	}
}

func (m *singleArgConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil // PromotedParameter 自动处理属性赋值
}

// toolMessageConstruct 是 ToolMessage 的构造函数，接受 content 和 tool_call_id
type toolMessageConstruct struct{}

func (m *toolMessageConstruct) GetName() string            { return "__construct" }
func (m *toolMessageConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *toolMessageConstruct) GetIsStatic() bool          { return false }
func (m *toolMessageConstruct) GetReturnType() data.Types  { return nil }

func (m *toolMessageConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewPromotedParameter(nil, "content", 0, nil, data.NewBaseType("string")),
		node.NewPromotedParameter(nil, "tool_call_id", 1, data.NewStringValue(""), data.NewBaseType("string")),
	}
}

func (m *toolMessageConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "content", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "tool_call_id", 1, data.NewBaseType("string")),
	}
}

func (m *toolMessageConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil // PromotedParameter 自动处理属性赋值
}
