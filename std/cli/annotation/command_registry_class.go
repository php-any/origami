package annotation

import (
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// lastCliApplication 最近一次注册的 CLI 应用（用于 list / about 输出）
var lastCliApplication *CliApplication

// CommandRegistryClass 供 PHP 读取已注册命令与应用信息
type CommandRegistryClass struct {
	node.Node
}

func (c *CommandRegistryClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *CommandRegistryClass) GetName() string    { return "Cli\\Annotation\\CommandRegistry" }
func (c *CommandRegistryClass) GetExtend() *string { return nil }
func (c *CommandRegistryClass) GetImplements() []string {
	return nil
}
func (c *CommandRegistryClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (c *CommandRegistryClass) GetPropertyList() []data.Property           { return []data.Property{} }
func (c *CommandRegistryClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "getCommands":
		return &CommandRegistryGetCommandsMethod{}, true
	case "getAppName":
		return &CommandRegistryGetAppNameMethod{}, true
	case "getAppVersion":
		return &CommandRegistryGetAppVersionMethod{}, true
	case "getLongVersion":
		return &CommandRegistryGetLongVersionMethod{}, true
	}
	return nil, false
}
func (c *CommandRegistryClass) GetMethods() []data.Method {
	return []data.Method{
		&CommandRegistryGetCommandsMethod{},
		&CommandRegistryGetAppNameMethod{},
		&CommandRegistryGetAppVersionMethod{},
		&CommandRegistryGetLongVersionMethod{},
	}
}
func (c *CommandRegistryClass) GetConstruct() data.Method { return nil }

type CommandRegistryGetCommandsMethod struct{}

func (m *CommandRegistryGetCommandsMethod) GetName() string            { return "getCommands" }
func (m *CommandRegistryGetCommandsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *CommandRegistryGetCommandsMethod) GetIsStatic() bool          { return true }
func (m *CommandRegistryGetCommandsMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *CommandRegistryGetCommandsMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *CommandRegistryGetCommandsMethod) GetReturnType() data.Types { return data.NewBaseType("array") }
func (m *CommandRegistryGetCommandsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	names := make([]string, 0, len(registeredCommands))
	for name := range registeredCommands {
		names = append(names, name)
	}
	sort.Strings(names)

	entries := make([]data.Value, 0, len(names))
	for _, name := range names {
		cmd := registeredCommands[name]
		obj := data.NewObjectValue()
		obj.SetProperty("name", data.NewStringValue(name))
		obj.SetProperty("description", data.NewStringValue(cmd.GetDescription()))
		entries = append(entries, obj)
	}
	return data.NewArrayValue(entries), nil
}

type CommandRegistryGetAppNameMethod struct{}

func (m *CommandRegistryGetAppNameMethod) GetName() string            { return "getAppName" }
func (m *CommandRegistryGetAppNameMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *CommandRegistryGetAppNameMethod) GetIsStatic() bool          { return true }
func (m *CommandRegistryGetAppNameMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *CommandRegistryGetAppNameMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *CommandRegistryGetAppNameMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *CommandRegistryGetAppNameMethod) Call(_ data.Context) (data.GetValue, data.Control) {
	if lastCliApplication != nil && lastCliApplication.name != "" {
		return data.NewStringValue(lastCliApplication.name), nil
	}
	return data.NewStringValue("CLI"), nil
}

type CommandRegistryGetAppVersionMethod struct{}

func (m *CommandRegistryGetAppVersionMethod) GetName() string            { return "getAppVersion" }
func (m *CommandRegistryGetAppVersionMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *CommandRegistryGetAppVersionMethod) GetIsStatic() bool          { return true }
func (m *CommandRegistryGetAppVersionMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *CommandRegistryGetAppVersionMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *CommandRegistryGetAppVersionMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *CommandRegistryGetAppVersionMethod) Call(_ data.Context) (data.GetValue, data.Control) {
	if lastCliApplication != nil && lastCliApplication.version != "" {
		return data.NewStringValue(lastCliApplication.version), nil
	}
	return data.NewStringValue("1.0.0"), nil
}

type CommandRegistryGetLongVersionMethod struct{}

func (m *CommandRegistryGetLongVersionMethod) GetName() string            { return "getLongVersion" }
func (m *CommandRegistryGetLongVersionMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *CommandRegistryGetLongVersionMethod) GetIsStatic() bool          { return true }
func (m *CommandRegistryGetLongVersionMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *CommandRegistryGetLongVersionMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *CommandRegistryGetLongVersionMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *CommandRegistryGetLongVersionMethod) Call(_ data.Context) (data.GetValue, data.Control) {
	name := "Laravel Framework"
	version := "1.0.0"
	if lastCliApplication != nil {
		if lastCliApplication.name != "" {
			name = lastCliApplication.name
		}
		if lastCliApplication.version != "" {
			version = lastCliApplication.version
		}
	}
	return data.NewStringValue(name + " " + version), nil
}
