package core

import (
	"sort"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PHP 8.4 的现代 DOM API（命名空间 `Dom`，注意与旧 `DOMDocument` 的大小写不同）。
//
// 存在的唯一理由：vendor/symfony/html-sanitizer 的 NativeParser 硬编码使用它，
// 且 HtmlSanitizer 的构造函数在没有显式传入 parser 时默认 `new NativeParser()`
// （没有 class_exists 兜底）：
//
//	$document = @\Dom\HTMLDocument::createFromString('<!DOCTYPE html><body>…</body>');
//	$element  = $document->getElementsByTagName('body')->item(0);
//	return $element->hasChildNodes() ? $element : null;
//
// 所以这里只覆盖 NativeParser + DomVisitor 真正读到的面：
// nodeName / nodeType / nodeValue / textContent / childNodes / attributes /
// hasChildNodes / getElementsByTagName / item / getIterator。
// 树是只读的，属性在构建时一次算好（textContent 预先拼平），不做惰性求值。

const (
	domNodeTypeElement = 1
	domNodeTypeText    = 3
	domNodeTypeAttr    = 2
	domNodeTypeDoc     = 9
)

// ---------- Dom\Node ----------

type DomNodeClass struct {
	node.Node
}

func NewDomNodeClass() *DomNodeClass { return &DomNodeClass{} }

func (c *DomNodeClass) GetName() string                                 { return "Dom\\Node" }
func (c *DomNodeClass) GetExtend() *string                              { return nil }
func (c *DomNodeClass) GetImplements() []string                         { return nil }
func (c *DomNodeClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *DomNodeClass) GetPropertyList() []data.Property                { return nil }
func (c *DomNodeClass) GetConstruct() data.Method                       { return nil }
func (c *DomNodeClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DomNodeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DomNodeClass) GetMethod(name string) (data.Method, bool) {
	if name == "hasChildNodes" {
		return &DomNodeHasChildNodesMethod{}, true
	}
	return nil, false
}
func (c *DomNodeClass) GetMethods() []data.Method {
	return []data.Method{&DomNodeHasChildNodesMethod{}}
}

type DomNodeHasChildNodesMethod struct{}

func (m *DomNodeHasChildNodesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		if children, _ := cmc.ObjectValue.GetProperty("childNodes"); children != nil {
			if arr, ok := children.(*data.ArrayValue); ok {
				return data.NewBoolValue(arr.Len() > 0), nil
			}
		}
	}
	return data.NewBoolValue(false), nil
}
func (m *DomNodeHasChildNodesMethod) GetName() string               { return "hasChildNodes" }
func (m *DomNodeHasChildNodesMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DomNodeHasChildNodesMethod) GetIsStatic() bool             { return false }
func (m *DomNodeHasChildNodesMethod) GetReturnType() data.Types     { return data.NewBaseType("bool") }
func (m *DomNodeHasChildNodesMethod) GetParams() []data.GetValue    { return nil }
func (m *DomNodeHasChildNodesMethod) GetVariables() []data.Variable { return nil }

// ---------- Dom\Element ----------

type DomElementClass struct {
	node.Node
}

func NewDomElementClass() *DomElementClass { return &DomElementClass{} }

func (c *DomElementClass) GetName() string           { return "Dom\\Element" }
func (c *DomElementClass) GetExtend() *string        { s := "Dom\\Node"; return &s }
func (c *DomElementClass) GetImplements() []string   { return nil }
func (c *DomElementClass) GetConstruct() data.Method { return nil }
func (c *DomElementClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *DomElementClass) GetPropertyList() []data.Property                { return nil }
func (c *DomElementClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DomElementClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DomElementClass) GetMethod(name string) (data.Method, bool) {
	if name == "getElementsByTagName" {
		return &DomGetElementsByTagNameMethod{}, true
	}
	return nil, false
}
func (c *DomElementClass) GetMethods() []data.Method {
	return []data.Method{&DomGetElementsByTagNameMethod{}}
}

// ---------- Dom\Text ----------

type DomTextClass struct {
	node.Node
}

func NewDomTextClass() *DomTextClass { return &DomTextClass{} }

func (c *DomTextClass) GetName() string                                 { return "Dom\\Text" }
func (c *DomTextClass) GetExtend() *string                              { s := "Dom\\Node"; return &s }
func (c *DomTextClass) GetImplements() []string                         { return nil }
func (c *DomTextClass) GetConstruct() data.Method                       { return nil }
func (c *DomTextClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *DomTextClass) GetPropertyList() []data.Property                { return nil }
func (c *DomTextClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DomTextClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DomTextClass) GetMethod(name string) (data.Method, bool) { return nil, false }
func (c *DomTextClass) GetMethods() []data.Method                 { return nil }

// ---------- Dom\ProcessingInstruction ----------

type DomProcessingInstructionClass struct {
	node.Node
}

func NewDomProcessingInstructionClass() *DomProcessingInstructionClass {
	return &DomProcessingInstructionClass{}
}

func (c *DomProcessingInstructionClass) GetName() string           { return "Dom\\ProcessingInstruction" }
func (c *DomProcessingInstructionClass) GetExtend() *string        { s := "Dom\\Node"; return &s }
func (c *DomProcessingInstructionClass) GetImplements() []string   { return nil }
func (c *DomProcessingInstructionClass) GetConstruct() data.Method { return nil }
func (c *DomProcessingInstructionClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *DomProcessingInstructionClass) GetPropertyList() []data.Property { return nil }
func (c *DomProcessingInstructionClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *DomProcessingInstructionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DomProcessingInstructionClass) GetMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *DomProcessingInstructionClass) GetMethods() []data.Method { return nil }

// ---------- Dom\Document ----------

type DomDocumentClass struct {
	node.Node
}

func NewDomDocumentClass() *DomDocumentClass { return &DomDocumentClass{} }

func (c *DomDocumentClass) GetName() string           { return "Dom\\Document" }
func (c *DomDocumentClass) GetExtend() *string        { s := "Dom\\Node"; return &s }
func (c *DomDocumentClass) GetImplements() []string   { return nil }
func (c *DomDocumentClass) GetConstruct() data.Method { return nil }
func (c *DomDocumentClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *DomDocumentClass) GetPropertyList() []data.Property                { return nil }
func (c *DomDocumentClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DomDocumentClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DomDocumentClass) GetMethod(name string) (data.Method, bool) {
	if name == "getElementsByTagName" {
		return &DomGetElementsByTagNameMethod{}, true
	}
	return nil, false
}
func (c *DomDocumentClass) GetMethods() []data.Method {
	return []data.Method{&DomGetElementsByTagNameMethod{}}
}

// ---------- Dom\HTMLDocument ----------

type DomHTMLDocumentClass struct {
	node.Node
}

func NewDomHTMLDocumentClass() *DomHTMLDocumentClass { return &DomHTMLDocumentClass{} }

func (c *DomHTMLDocumentClass) GetName() string           { return "Dom\\HTMLDocument" }
func (c *DomHTMLDocumentClass) GetExtend() *string        { s := "Dom\\Document"; return &s }
func (c *DomHTMLDocumentClass) GetImplements() []string   { return nil }
func (c *DomHTMLDocumentClass) GetConstruct() data.Method { return nil }
func (c *DomHTMLDocumentClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *DomHTMLDocumentClass) GetPropertyList() []data.Property { return nil }
func (c *DomHTMLDocumentClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DomHTMLDocumentClass) GetMethod(name string) (data.Method, bool) {
	if name == "getElementsByTagName" {
		return &DomGetElementsByTagNameMethod{}, true
	}
	return nil, false
}
func (c *DomHTMLDocumentClass) GetMethods() []data.Method {
	return []data.Method{&DomGetElementsByTagNameMethod{}, &DomHTMLDocumentCreateFromStringMethod{}}
}
func (c *DomHTMLDocumentClass) GetStaticMethod(name string) (data.Method, bool) {
	if name == "createFromString" {
		return &DomHTMLDocumentCreateFromStringMethod{}, true
	}
	return nil, false
}

type DomHTMLDocumentCreateFromStringMethod struct{}

func (m *DomHTMLDocumentCreateFromStringMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	sourceVal, _ := ctx.GetIndexValue(0)
	if sourceVal == nil {
		return data.NewNullValue(), nil
	}
	return buildDomHTMLDocument(sourceVal.AsString(), ctx), nil
}
func (m *DomHTMLDocumentCreateFromStringMethod) GetName() string { return "createFromString" }
func (m *DomHTMLDocumentCreateFromStringMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *DomHTMLDocumentCreateFromStringMethod) GetIsStatic() bool { return true }
func (m *DomHTMLDocumentCreateFromStringMethod) GetReturnType() data.Types {
	return nil
}
func (m *DomHTMLDocumentCreateFromStringMethod) GetParams() []data.GetValue {
	return domHTMLDocumentCreateFromStringMethodParams
}

var domHTMLDocumentCreateFromStringMethodParams = []data.GetValue{
	node.NewParameter(nil, "source", 0, nil, data.NewBaseType("string")),
	node.NewParameter(nil, "options", 1, node.NewIntLiteral(nil, "0"), data.NewBaseType("int")),
}

func (m *DomHTMLDocumentCreateFromStringMethod) GetVariables() []data.Variable {
	return domHTMLDocumentCreateFromStringMethodVariables
}

var domHTMLDocumentCreateFromStringMethodVariables = []data.Variable{
	node.NewVariable(nil, "source", 0, data.NewBaseType("string")),
	node.NewVariable(nil, "options", 1, data.NewBaseType("int")),
}

// ---------- getElementsByTagName（Document / Element 共用） ----------

type DomGetElementsByTagNameMethod struct{}

func (m *DomGetElementsByTagNameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cmc, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return newDOMNodeList(nil, ctx), nil
	}
	tagName := ""
	if nameVal, _ := ctx.GetIndexValue(0); nameVal != nil {
		tagName = strings.ToLower(nameVal.AsString())
	}
	if tagName == "" {
		return newDOMNodeList(nil, ctx), nil
	}

	var results []data.Value
	if children, _ := cmc.ObjectValue.GetProperty("childNodes"); children != nil {
		if arr, ok := children.(*data.ArrayValue); ok {
			for arraySlots116, arrayPosition116 := arr.View(), 0; arrayPosition116 < arraySlots116.Len(); arrayPosition116++ {
				zval := arraySlots116.At(arrayPosition116)
				if child, ok := zval.ReadValue().(*data.ClassValue); ok {
					collectElementsByTagName(child, tagName, &results)
				}
			}
		}
	}
	return newDOMNodeList(results, ctx), nil
}
func (m *DomGetElementsByTagNameMethod) GetName() string            { return "getElementsByTagName" }
func (m *DomGetElementsByTagNameMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DomGetElementsByTagNameMethod) GetIsStatic() bool          { return false }
func (m *DomGetElementsByTagNameMethod) GetReturnType() data.Types  { return nil }
func (m *DomGetElementsByTagNameMethod) GetParams() []data.GetValue {
	return domGetElementsByTagNameMethodParams
}

var domGetElementsByTagNameMethodParams = []data.GetValue{
	node.NewParameter(nil, "qualifiedName", 0, nil, data.NewBaseType("string")),
}

func (m *DomGetElementsByTagNameMethod) GetVariables() []data.Variable {
	return domGetElementsByTagNameMethodVariables
}

var domGetElementsByTagNameMethodVariables = []data.Variable{
	node.NewVariable(nil, "qualifiedName", 0, data.NewBaseType("string")),
}

// ---------- Dom\Attr ----------

type DomAttrClass struct {
	node.Node
}

func NewDomAttrClass() *DomAttrClass { return &DomAttrClass{} }

func (c *DomAttrClass) GetName() string                                 { return "Dom\\Attr" }
func (c *DomAttrClass) GetExtend() *string                              { s := "Dom\\Node"; return &s }
func (c *DomAttrClass) GetImplements() []string                         { return nil }
func (c *DomAttrClass) GetConstruct() data.Method                       { return nil }
func (c *DomAttrClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *DomAttrClass) GetPropertyList() []data.Property                { return nil }
func (c *DomAttrClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DomAttrClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DomAttrClass) GetMethod(name string) (data.Method, bool) { return nil, false }
func (c *DomAttrClass) GetMethods() []data.Method                 { return nil }

// ---------- Dom\NamedNodeMap ----------
//
// DomVisitor 只用到 `$node->attributes?->getIterator()`，返回可 foreach 的属性列表。

type DomNamedNodeMapClass struct {
	node.Node
}

func NewDomNamedNodeMapClass() *DomNamedNodeMapClass { return &DomNamedNodeMapClass{} }

func (c *DomNamedNodeMapClass) GetName() string           { return "Dom\\NamedNodeMap" }
func (c *DomNamedNodeMapClass) GetExtend() *string        { return nil }
func (c *DomNamedNodeMapClass) GetImplements() []string   { return nil }
func (c *DomNamedNodeMapClass) GetConstruct() data.Method { return nil }
func (c *DomNamedNodeMapClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *DomNamedNodeMapClass) GetPropertyList() []data.Property                { return nil }
func (c *DomNamedNodeMapClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DomNamedNodeMapClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DomNamedNodeMapClass) GetMethod(name string) (data.Method, bool) {
	if name == "getIterator" {
		return &DomNamedNodeMapGetIteratorMethod{}, true
	}
	if name == "count" {
		return &DomNamedNodeMapCountMethod{}, true
	}
	return nil, false
}
func (c *DomNamedNodeMapClass) GetMethods() []data.Method {
	return []data.Method{&DomNamedNodeMapGetIteratorMethod{}, &DomNamedNodeMapCountMethod{}}
}

type DomNamedNodeMapGetIteratorMethod struct{}

func (m *DomNamedNodeMapGetIteratorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		if items, _ := cmc.ObjectValue.GetProperty("_items"); items != nil {
			return items, nil
		}
		if items, _ := cmc.ObjectValue.GetProperty("_attributes"); items != nil {
			return items, nil
		}
	}
	return data.NewArrayValue(nil), nil
}
func (m *DomNamedNodeMapGetIteratorMethod) GetName() string               { return "getIterator" }
func (m *DomNamedNodeMapGetIteratorMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DomNamedNodeMapGetIteratorMethod) GetIsStatic() bool             { return false }
func (m *DomNamedNodeMapGetIteratorMethod) GetReturnType() data.Types     { return nil }
func (m *DomNamedNodeMapGetIteratorMethod) GetParams() []data.GetValue    { return nil }
func (m *DomNamedNodeMapGetIteratorMethod) GetVariables() []data.Variable { return nil }

type DomNamedNodeMapCountMethod struct{}

func (m *DomNamedNodeMapCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		if length, _ := cmc.ObjectValue.GetProperty("length"); length != nil {
			return length, nil
		}
	}
	return data.NewIntValue(0), nil
}
func (m *DomNamedNodeMapCountMethod) GetName() string               { return "count" }
func (m *DomNamedNodeMapCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DomNamedNodeMapCountMethod) GetIsStatic() bool             { return false }
func (m *DomNamedNodeMapCountMethod) GetReturnType() data.Types     { return data.NewBaseType("int") }
func (m *DomNamedNodeMapCountMethod) GetParams() []data.GetValue    { return nil }
func (m *DomNamedNodeMapCountMethod) GetVariables() []data.Variable { return nil }

// ---------- 树构建 ----------

// buildDomHTMLDocument 按 HTML 解析器的隐含结构建树：document > html > [head, body]。
// 顶层若是 html/head/body 包裹标签，其内容被并进对应元素（与浏览器解析一致），
// 其余顶层节点一律进 body —— 这样 `getElementsByTagName($context)` 对
// NativeParser 传入的任意 context（body/div/…）都能命中原包裹元素。
func buildDomHTMLDocument(source string, ctx data.Context) *data.ClassValue {
	base := ctx.CreateBaseContext()
	parsed := parseHTML(htmlStripped(source))

	head := &htmlNode{tag: "head"}
	body := &htmlNode{tag: "body"}
	html := &htmlNode{tag: "html", children: []*htmlNode{head, body}}

	for _, child := range parsed.children {
		if child.isText {
			body.children = append(body.children, child)
			continue
		}
		switch child.tag {
		case "html":
			for _, grand := range child.children {
				if !grand.isText && grand.tag == "head" {
					head.children = append(head.children, grand.children...)
				} else if !grand.isText && grand.tag == "body" {
					body.children = append(body.children, grand.children...)
				} else {
					body.children = append(body.children, grand)
				}
			}
		case "head":
			head.children = append(head.children, child.children...)
		case "body":
			body.children = append(body.children, child.children...)
		default:
			body.children = append(body.children, child)
		}
	}

	headVal := buildDomNode(head, base)
	bodyVal := buildDomNode(body, base)
	htmlVal := buildDomNode(html, base)

	doc := data.NewClassValue(NewDomHTMLDocumentClass(), base)
	doc.ObjectValue.SetProperty("nodeName", data.NewStringValue("#document"))
	doc.ObjectValue.SetProperty("nodeType", data.NewIntValue(domNodeTypeDoc))
	doc.ObjectValue.SetProperty("nodeValue", data.NewNullValue())
	doc.ObjectValue.SetProperty("textContent", data.NewStringValue(domTextContent(body)))
	doc.ObjectValue.SetProperty("childNodes", data.NewArrayValue([]data.Value{htmlVal}))
	doc.ObjectValue.SetProperty("attributes", data.NewNullValue())
	doc.ObjectValue.SetProperty("documentElement", htmlVal)
	doc.ObjectValue.SetProperty("head", headVal)
	doc.ObjectValue.SetProperty("body", bodyVal)
	return doc
}

// buildDomNode 把 htmlNode 转成 Dom\* 实例。属性一次算全，树只读。
func buildDomNode(n *htmlNode, ctx data.Context) *data.ClassValue {
	if n.isText {
		cv := data.NewClassValue(NewDomTextClass(), ctx)
		cv.ObjectValue.SetProperty("nodeName", data.NewStringValue("#text"))
		cv.ObjectValue.SetProperty("nodeType", data.NewIntValue(domNodeTypeText))
		cv.ObjectValue.SetProperty("nodeValue", data.NewStringValue(n.text))
		cv.ObjectValue.SetProperty("textContent", data.NewStringValue(n.text))
		cv.ObjectValue.SetProperty("childNodes", data.NewArrayValue(nil))
		cv.ObjectValue.SetProperty("attributes", data.NewNullValue())
		return cv
	}

	cv := data.NewClassValue(NewDomElementClass(), ctx)
	cv.ObjectValue.SetProperty("nodeName", data.NewStringValue(n.tag))
	cv.ObjectValue.SetProperty("nodeType", data.NewIntValue(domNodeTypeElement))
	cv.ObjectValue.SetProperty("nodeValue", data.NewNullValue())
	cv.ObjectValue.SetProperty("textContent", data.NewStringValue(domTextContent(n)))
	cv.ObjectValue.SetProperty("attributes", buildDomNamedNodeMap(n.attrs, ctx))

	children := make([]data.Value, 0, len(n.children))
	for _, child := range n.children {
		children = append(children, buildDomNode(child, ctx))
	}
	cv.ObjectValue.SetProperty("childNodes", data.NewArrayValue(children))
	return cv
}

// domTextContent 把子树文本预先拼平（PHP: Element::$textContent 即全部后代文本）。
func domTextContent(n *htmlNode) string {
	if n.isText {
		return n.text
	}
	if len(n.children) == 0 {
		return ""
	}
	var b strings.Builder
	var walk func(node *htmlNode)
	walk = func(node *htmlNode) {
		for _, child := range node.children {
			if child.isText {
				b.WriteString(child.text)
			} else {
				walk(child)
			}
		}
	}
	walk(n)
	return b.String()
}

// buildDomNamedNodeMap 属性名排序，保证同一份 HTML 每次渲染结果一致
// （parseHTML 的属性存在 map 里，直接遍历顺序随机）。
func buildDomNamedNodeMap(attrs map[string]string, ctx data.Context) *data.ClassValue {
	cv := data.NewClassValue(NewDomNamedNodeMapClass(), ctx)

	items := make([]data.Value, 0, len(attrs))
	if len(attrs) > 0 {
		names := make([]string, 0, len(attrs))
		for name := range attrs {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			value := attrs[name]
			attr := data.NewClassValue(NewDomAttrClass(), ctx)
			attr.ObjectValue.SetProperty("name", data.NewStringValue(name))
			attr.ObjectValue.SetProperty("value", data.NewStringValue(value))
			attr.ObjectValue.SetProperty("nodeName", data.NewStringValue(name))
			attr.ObjectValue.SetProperty("nodeType", data.NewIntValue(domNodeTypeAttr))
			attr.ObjectValue.SetProperty("nodeValue", data.NewStringValue(value))
			attr.ObjectValue.SetProperty("textContent", data.NewStringValue(value))
			items = append(items, attr)
		}
	}

	cv.ObjectValue.SetProperty("length", data.NewIntValue(len(items)))
	cv.ObjectValue.SetProperty("_items", data.NewArrayValue(items))
	return cv
}
