package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NormalizerClass 模拟 intl 扩展提供的全局 Normalizer 类，只提供常量。
//
// 常量值参考 PHP 文档，精确数值在本运行时中并不关键，
// 主要用于 Symfony polyfill / string 组件中的等值比较与 switch 分支。
type NormalizerClass struct {
	node.Node
}

func (n *NormalizerClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(n, ctx.CreateBaseContext()), nil
}

func (n *NormalizerClass) GetFrom() data.From { return n.Node.GetFrom() }

func (n *NormalizerClass) GetName() string { return "Normalizer" }

func (n *NormalizerClass) GetExtend() *string { return nil }

func (n *NormalizerClass) GetImplements() []string { return nil }

// GetProperty 普通属性：这里只提供常量，不支持实例属性
func (n *NormalizerClass) GetProperty(name string) (data.Property, bool) { return nil, false }

// GetStaticProperty 提供 Normalizer 相关常量，作为“静态属性/常量”暴露给 PHP 层。
//
// 参考 PHP:
//   - Normalizer::FORM_D  = 1
//   - Normalizer::FORM_KD = 2
//   - Normalizer::FORM_C  = 4
//   - Normalizer::FORM_KC = 5
//   - Normalizer::NFD  = 1
//   - Normalizer::NFKD = 2
//   - Normalizer::NFC  = 4
//   - Normalizer::NFKC = 5
//   - Normalizer::NONE = 0
func (n *NormalizerClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "NONE":
		return data.NewIntValue(0), true
	case "FORM_D", "NFD":
		return data.NewIntValue(1), true
	case "FORM_KD", "NFKD":
		return data.NewIntValue(2), true
	case "FORM_C", "NFC":
		return data.NewIntValue(4), true
	case "FORM_KC", "NFKC":
		return data.NewIntValue(5), true
	default:
		return nil, false
	}
}

func (n *NormalizerClass) GetPropertyList() []data.Property { return nil }

func (n *NormalizerClass) GetMethod(name string) (data.Method, bool) { return nil, false }

func (n *NormalizerClass) GetMethods() []data.Method { return nil }

func (n *NormalizerClass) GetConstruct() data.Method { return nil }
