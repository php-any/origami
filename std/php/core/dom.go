package core

import (
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ---------- DOMNode ----------

type DOMNodeClass struct {
	node.Node
}

func NewDOMNodeClass() *DOMNodeClass { return &DOMNodeClass{} }

func (c *DOMNodeClass) GetName() string                                 { return "DOMNode" }
func (c *DOMNodeClass) GetExtend() *string                              { return nil }
func (c *DOMNodeClass) GetImplements() []string                         { return nil }
func (c *DOMNodeClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *DOMNodeClass) GetPropertyList() []data.Property                { return nil }
func (c *DOMNodeClass) GetConstruct() data.Method                       { return nil }
func (c *DOMNodeClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DOMNodeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DOMNodeClass) GetMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *DOMNodeClass) GetMethods() []data.Method { return nil }

// ---------- DOMDocument ----------

type DOMDocumentClass struct {
	node.Node
}

func NewDOMDocumentClass() *DOMDocumentClass { return &DOMDocumentClass{} }

func (c *DOMDocumentClass) GetName() string           { return "DOMDocument" }
func (c *DOMDocumentClass) GetExtend() *string        { s := "DOMNode"; return &s }
func (c *DOMDocumentClass) GetImplements() []string   { return nil }
func (c *DOMDocumentClass) GetConstruct() data.Method { return &DOMDocumentConstructMethod{} }

func (c *DOMDocumentClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *DOMDocumentClass) GetPropertyList() []data.Property                { return nil }
func (c *DOMDocumentClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DOMDocumentClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *DOMDocumentClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &DOMDocumentConstructMethod{}, true
	case "loadHTML":
		return &DOMDocumentLoadHTMLMethod{}, true
	case "getElementsByTagName":
		return &DOMDocumentGetElementsByTagNameMethod{}, true
	case "saveXML":
		return &DOMDocumentSaveXMLMethod{}, true
	}
	return nil, false
}

func (c *DOMDocumentClass) GetMethods() []data.Method {
	return []data.Method{
		&DOMDocumentConstructMethod{},
		&DOMDocumentLoadHTMLMethod{},
		&DOMDocumentGetElementsByTagNameMethod{},
		&DOMDocumentSaveXMLMethod{},
	}
}

// ---------- DOMDocument::__construct ----------

type DOMDocumentConstructMethod struct{}

func (m *DOMDocumentConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}
func (m *DOMDocumentConstructMethod) GetName() string               { return "__construct" }
func (m *DOMDocumentConstructMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DOMDocumentConstructMethod) GetIsStatic() bool             { return false }
func (m *DOMDocumentConstructMethod) GetReturnType() data.Types     { return nil }
func (m *DOMDocumentConstructMethod) GetParams() []data.GetValue    { return nil }
func (m *DOMDocumentConstructMethod) GetVariables() []data.Variable { return nil }

// ---------- DOMDocument::loadHTML ----------

type DOMDocumentLoadHTMLMethod struct{}

func (m *DOMDocumentLoadHTMLMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	htmlVal, _ := ctx.GetIndexValue(0)
	if htmlVal == nil {
		return data.NewBoolValue(false), nil
	}
	htmlStr := htmlStripped(htmlVal.AsString())

	// Parse the HTML
	parsed := parseHTML(htmlStr)

	// Build DOM tree starting from parsed root
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		docNode := buildDOMNode(parsed, cmc.ClassValue, ctx)
		cmc.ObjectValue.SetProperty("documentElement", docNode)
	}

	return data.NewBoolValue(true), nil
}
func (m *DOMDocumentLoadHTMLMethod) GetName() string            { return "loadHTML" }
func (m *DOMDocumentLoadHTMLMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DOMDocumentLoadHTMLMethod) GetIsStatic() bool          { return false }
func (m *DOMDocumentLoadHTMLMethod) GetReturnType() data.Types  { return data.NewBaseType("bool") }
func (m *DOMDocumentLoadHTMLMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "source", 0, nil, nil),
		node.NewParameter(nil, "options", 1, node.NewIntLiteral(nil, "0"), nil),
	}
}
func (m *DOMDocumentLoadHTMLMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "source", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "options", 1, data.NewBaseType("int")),
	}
}

// ---------- DOMDocument::getElementsByTagName ----------

type DOMDocumentGetElementsByTagNameMethod struct{}

func (m *DOMDocumentGetElementsByTagNameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	nameVal, _ := ctx.GetIndexValue(0)
	if nameVal == nil {
		return newDOMNodeList(nil, ctx), nil
	}
	tagName := strings.ToLower(nameVal.AsString())

	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		docElem, _ := cmc.ObjectValue.GetProperty("documentElement")
		var results []data.Value
		if docElem != nil {
			if cv, ok := docElem.(*data.ClassValue); ok {
				collectElementsByTagName(cv, tagName, &results)
			}
		}
		return newDOMNodeList(results, ctx), nil
	}
	return newDOMNodeList(nil, ctx), nil
}
func (m *DOMDocumentGetElementsByTagNameMethod) GetName() string { return "getElementsByTagName" }
func (m *DOMDocumentGetElementsByTagNameMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *DOMDocumentGetElementsByTagNameMethod) GetIsStatic() bool         { return false }
func (m *DOMDocumentGetElementsByTagNameMethod) GetReturnType() data.Types { return nil }
func (m *DOMDocumentGetElementsByTagNameMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, nil),
	}
}
func (m *DOMDocumentGetElementsByTagNameMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.NewBaseType("string")),
	}
}

// ---------- DOMDocument::saveXML ----------

type DOMDocumentSaveXMLMethod struct{}

func (m *DOMDocumentSaveXMLMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	nodeVal, _ := ctx.GetIndexValue(0)
	if nodeVal == nil {
		// Save entire document
		if cmc, ok := ctx.(*data.ClassMethodContext); ok {
			docElem, _ := cmc.ObjectValue.GetProperty("documentElement")
			if docElem != nil {
				xml := nodeToXML(docElem)
				return data.NewStringValue(xml), nil
			}
		}
		return data.NewStringValue(""), nil
	}
	xml := nodeToXML(nodeVal)
	return data.NewStringValue(xml), nil
}
func (m *DOMDocumentSaveXMLMethod) GetName() string            { return "saveXML" }
func (m *DOMDocumentSaveXMLMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DOMDocumentSaveXMLMethod) GetIsStatic() bool          { return false }
func (m *DOMDocumentSaveXMLMethod) GetReturnType() data.Types  { return data.NewBaseType("string") }
func (m *DOMDocumentSaveXMLMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "node", 0, node.NewNullLiteral(nil), nil),
	}
}
func (m *DOMDocumentSaveXMLMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "node", 0, data.NewNullableType(nil)),
	}
}

// ---------- DOMNodeList ----------

type DOMNodeListClass struct {
	node.Node
}

func NewDOMNodeListClass() *DOMNodeListClass { return &DOMNodeListClass{} }

func (c *DOMNodeListClass) GetName() string                                 { return "DOMNodeList" }
func (c *DOMNodeListClass) GetExtend() *string                              { return nil }
func (c *DOMNodeListClass) GetImplements() []string                         { return nil }
func (c *DOMNodeListClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *DOMNodeListClass) GetPropertyList() []data.Property                { return nil }
func (c *DOMNodeListClass) GetConstruct() data.Method                       { return nil }
func (c *DOMNodeListClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DOMNodeListClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DOMNodeListClass) GetMethod(name string) (data.Method, bool) {
	if name == "item" {
		return &DOMNodeListItemMethod{}, true
	}
	if name == "count" {
		return &DOMNodeListCountMethod{}, true
	}
	return nil, false
}
func (c *DOMNodeListClass) GetMethods() []data.Method {
	return []data.Method{&DOMNodeListItemMethod{}, &DOMNodeListCountMethod{}}
}

func newDOMNodeList(items []data.Value, ctx data.Context) *data.ClassValue {
	cls := NewDOMNodeListClass()
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	arr := data.NewArrayValue(items)
	cv.ObjectValue.SetProperty("length", data.NewIntValue(len(items)))
	cv.ObjectValue.SetProperty("_items", arr)
	return cv
}

type DOMNodeListItemMethod struct{}

func (m *DOMNodeListItemMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	idxVal, _ := ctx.GetIndexValue(0)
	if idxVal == nil {
		return data.NewNullValue(), nil
	}
	idx := 0
	if iv, ok := idxVal.(*data.IntValue); ok {
		idx = iv.Value
	}
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		items, _ := cmc.ObjectValue.GetProperty("_items")
		if arr, ok := items.(*data.ArrayValue); ok {
			if idx >= 0 && idx < len(arr.List) {
				return arr.List[idx].Value, nil
			}
		}
	}
	return data.NewNullValue(), nil
}
func (m *DOMNodeListItemMethod) GetName() string            { return "item" }
func (m *DOMNodeListItemMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DOMNodeListItemMethod) GetIsStatic() bool          { return false }
func (m *DOMNodeListItemMethod) GetReturnType() data.Types  { return nil }
func (m *DOMNodeListItemMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "index", 0, nil, nil)}
}
func (m *DOMNodeListItemMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "index", 0, data.NewBaseType("int"))}
}

type DOMNodeListCountMethod struct{}

func (m *DOMNodeListCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		length, _ := cmc.ObjectValue.GetProperty("length")
		return length, nil
	}
	return data.NewIntValue(0), nil
}
func (m *DOMNodeListCountMethod) GetName() string               { return "count" }
func (m *DOMNodeListCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DOMNodeListCountMethod) GetIsStatic() bool             { return false }
func (m *DOMNodeListCountMethod) GetReturnType() data.Types     { return data.NewBaseType("int") }
func (m *DOMNodeListCountMethod) GetParams() []data.GetValue    { return nil }
func (m *DOMNodeListCountMethod) GetVariables() []data.Variable { return nil }

// ---------- DOMElement ----------

type DOMElementClass struct {
	node.Node
}

func NewDOMElementClass() *DOMElementClass { return &DOMElementClass{} }

func (c *DOMElementClass) GetName() string                                 { return "DOMElement" }
func (c *DOMElementClass) GetExtend() *string                              { s := "DOMNode"; return &s }
func (c *DOMElementClass) GetImplements() []string                         { return nil }
func (c *DOMElementClass) GetConstruct() data.Method                       { return nil }
func (c *DOMElementClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *DOMElementClass) GetPropertyList() []data.Property                { return nil }
func (c *DOMElementClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DOMElementClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DOMElementClass) GetMethod(name string) (data.Method, bool) {
	if name == "getAttribute" {
		return &DOMElementGetAttributeMethod{}, true
	}
	return nil, false
}
func (c *DOMElementClass) GetMethods() []data.Method {
	return []data.Method{&DOMElementGetAttributeMethod{}}
}

type DOMElementGetAttributeMethod struct{}

func (m *DOMElementGetAttributeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	nameVal, _ := ctx.GetIndexValue(0)
	if nameVal == nil {
		return data.NewStringValue(""), nil
	}
	attrName := nameVal.AsString()
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		attrs, _ := cmc.ObjectValue.GetProperty("_attributes")
		if obj, ok := attrs.(*data.ObjectValue); ok {
			val, _ := obj.GetProperty(attrName)
			if val != nil {
				return val, nil
			}
		}
	}
	return data.NewStringValue(""), nil
}
func (m *DOMElementGetAttributeMethod) GetName() string            { return "getAttribute" }
func (m *DOMElementGetAttributeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DOMElementGetAttributeMethod) GetIsStatic() bool          { return false }
func (m *DOMElementGetAttributeMethod) GetReturnType() data.Types  { return data.NewBaseType("string") }
func (m *DOMElementGetAttributeMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "name", 0, nil, nil)}
}
func (m *DOMElementGetAttributeMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "name", 0, data.NewBaseType("string"))}
}

// ---------- DOMText ----------

type DOMTextClass struct {
	node.Node
}

func NewDOMTextClass() *DOMTextClass { return &DOMTextClass{} }

func (c *DOMTextClass) GetName() string                                 { return "DOMText" }
func (c *DOMTextClass) GetExtend() *string                              { s := "DOMNode"; return &s }
func (c *DOMTextClass) GetImplements() []string                         { return nil }
func (c *DOMTextClass) GetConstruct() data.Method                       { return nil }
func (c *DOMTextClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *DOMTextClass) GetPropertyList() []data.Property                { return nil }
func (c *DOMTextClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DOMTextClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DOMTextClass) GetMethod(name string) (data.Method, bool) { return nil, false }
func (c *DOMTextClass) GetMethods() []data.Method                 { return nil }

// ---------- DOMComment ----------

type DOMCommentClass struct {
	node.Node
}

func NewDOMCommentClass() *DOMCommentClass { return &DOMCommentClass{} }

func (c *DOMCommentClass) GetName() string                                 { return "DOMComment" }
func (c *DOMCommentClass) GetExtend() *string                              { s := "DOMNode"; return &s }
func (c *DOMCommentClass) GetImplements() []string                         { return nil }
func (c *DOMCommentClass) GetConstruct() data.Method                       { return nil }
func (c *DOMCommentClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *DOMCommentClass) GetPropertyList() []data.Property                { return nil }
func (c *DOMCommentClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *DOMCommentClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DOMCommentClass) GetMethod(name string) (data.Method, bool) { return nil, false }
func (c *DOMCommentClass) GetMethods() []data.Method                 { return nil }

// ---------- DOMNodeList traversal (foreach support) ----------

// The foreach over DOMNodeList works because we store _items as an ArrayValue
// and Origami's foreach iterates ArrayValue directly.

// ---------- HTML Parser ----------

type htmlNode struct {
	tag       string
	attrs     map[string]string
	text      string
	children  []*htmlNode
	isText    bool
	isComment bool
}

func parseHTML(src string) *htmlNode {
	root := &htmlNode{tag: "body"}
	p := &htmlParser{src: src, pos: 0}
	p.skipWhitespace()

	for p.pos < len(p.src) {
		if p.peek() == '<' {
			if p.pos+1 < len(p.src) && p.src[p.pos+1] == '!' {
				// Comment or doctype - skip
				p.skipComment()
				continue
			}
			n := p.parseElement()
			if n != nil {
				root.children = append(root.children, n)
			}
		} else {
			text := p.parseText()
			if text != "" {
				root.children = append(root.children, &htmlNode{text: text, isText: true})
			}
		}
	}
	return root
}

type htmlParser struct {
	src string
	pos int
}

func (p *htmlParser) peek() byte {
	if p.pos >= len(p.src) {
		return 0
	}
	return p.src[p.pos]
}

func (p *htmlParser) advance() byte {
	if p.pos >= len(p.src) {
		return 0
	}
	c := p.src[p.pos]
	p.pos++
	return c
}

func (p *htmlParser) skipWhitespace() {
	for p.pos < len(p.src) && unicode.IsSpace(rune(p.src[p.pos])) {
		p.pos++
	}
}

func (p *htmlParser) skipComment() {
	// Skip <!-- ... --> or <!DOCTYPE ... >
	for p.pos < len(p.src) && p.src[p.pos] != '>' {
		p.pos++
	}
	if p.pos < len(p.src) {
		p.pos++ // skip '>'
	}
}

func (p *htmlParser) parseElement() *htmlNode {
	if p.peek() != '<' {
		return nil
	}
	p.advance() // skip '<'

	// Check for closing tag
	if p.peek() == '/' {
		// Skip closing tag
		for p.pos < len(p.src) && p.src[p.pos] != '>' {
			p.pos++
		}
		if p.pos < len(p.src) {
			p.pos++
		}
		return nil
	}

	// Parse tag name
	tagName := p.parseTagName()
	if tagName == "" {
		return nil
	}
	tagName = strings.ToLower(tagName)

	// Self-closing tags
	selfClosing := tagName == "br" || tagName == "hr" || tagName == "img" || tagName == "input"

	// Parse attributes
	attrs := make(map[string]string)
	for p.pos < len(p.src) {
		p.skipWhitespace()
		if p.peek() == '>' {
			p.advance()
			break
		}
		if p.peek() == '/' {
			p.advance()
			selfClosing = true
			if p.peek() == '>' {
				p.advance()
			}
			break
		}
		name, value := p.parseAttribute()
		if name != "" {
			attrs[name] = value
		}
	}

	node := &htmlNode{tag: tagName, attrs: attrs}

	if selfClosing {
		return node
	}

	// Parse children until closing tag
	for p.pos < len(p.src) {
		p.skipWhitespace()
		if p.pos >= len(p.src) {
			break
		}

		// Check for closing tag
		if p.peek() == '<' {
			if p.pos+1 < len(p.src) && p.src[p.pos+1] == '/' {
				// Find the matching closing tag
				savePos := p.pos
				p.advance() // <
				p.advance() // /
				closeTag := ""
				for p.pos < len(p.src) && p.src[p.pos] != '>' {
					closeTag += string(p.src[p.pos])
					p.pos++
				}
				if p.pos < len(p.src) {
					p.pos++ // skip '>'
				}
				closeTag = strings.ToLower(closeTag)
				if closeTag == tagName {
					break
				}
				// Not matching - treat as text
				p.pos = savePos
				text := p.parseText()
				if text != "" {
					node.children = append(node.children, &htmlNode{text: text, isText: true})
				}
				continue
			}
			if p.pos+1 < len(p.src) && p.src[p.pos+1] == '!' {
				p.skipComment()
				continue
			}
			child := p.parseElement()
			if child != nil {
				node.children = append(node.children, child)
			}
		} else {
			text := p.parseText()
			if text != "" {
				node.children = append(node.children, &htmlNode{text: text, isText: true})
			}
		}
	}

	return node
}

func (p *htmlParser) parseTagName() string {
	start := p.pos
	for p.pos < len(p.src) && (p.src[p.pos] == ':' || p.src[p.pos] == '-' || p.src[p.pos] == '_' ||
		p.src[p.pos] == '.' || unicode.IsLetter(rune(p.src[p.pos])) || unicode.IsDigit(rune(p.src[p.pos]))) {
		p.pos++
	}
	return p.src[start:p.pos]
}

func (p *htmlParser) parseAttribute() (string, string) {
	// Parse attribute name
	start := p.pos
	for p.pos < len(p.src) && p.src[p.pos] != '=' && p.src[p.pos] != '>' && p.src[p.pos] != '/' &&
		!unicode.IsSpace(rune(p.src[p.pos])) {
		p.pos++
	}
	name := p.src[start:p.pos]
	if name == "" {
		return "", ""
	}

	p.skipWhitespace()
	if p.peek() != '=' {
		return name, ""
	}
	p.advance() // skip '='
	p.skipWhitespace()

	// Parse attribute value
	quote := p.peek()
	if quote == '"' || quote == '\'' {
		p.advance()
		start = p.pos
		for p.pos < len(p.src) && p.src[p.pos] != quote {
			if p.src[p.pos] == '\\' && p.pos+1 < len(p.src) {
				p.pos += 2
			} else {
				p.pos++
			}
		}
		value := p.src[start:p.pos]
		if p.pos < len(p.src) {
			p.pos++ // skip closing quote
		}
		return name, value
	}

	// Unquoted value
	start = p.pos
	for p.pos < len(p.src) && !unicode.IsSpace(rune(p.src[p.pos])) && p.src[p.pos] != '>' {
		p.pos++
	}
	return name, p.src[start:p.pos]
}

func (p *htmlParser) parseText() string {
	start := p.pos
	for p.pos < len(p.src) && p.src[p.pos] != '<' {
		p.pos++
	}
	text := p.src[start:p.pos]
	return htmlEntityDecode(text)
}

func htmlEntityDecode(s string) string {
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&apos;", "'")
	// Decode numeric entities like &#39;
	for i := 0; i < len(s); i++ {
		if s[i] == '&' && i+2 < len(s) && s[i+1] == '#' {
			end := strings.IndexByte(s[i:], ';')
			if end > 0 {
				numStr := s[i+2 : i+end]
				var code int
				if numStr[0] == 'x' || numStr[0] == 'X' {
					for _, c := range numStr[1:] {
						code = code * 16
						if c >= '0' && c <= '9' {
							code += int(c - '0')
						} else if c >= 'a' && c <= 'f' {
							code += int(c - 'a' + 10)
						} else if c >= 'A' && c <= 'F' {
							code += int(c - 'A' + 10)
						}
					}
				} else {
					for _, c := range numStr {
						code = code*10 + int(c-'0')
					}
				}
				s = s[:i] + string(rune(code)) + s[i+end+1:]
			}
		}
	}
	return s
}

// htmlStripped strips PHP XML encoding header and whitespace
func htmlStripped(s string) string {
	if idx := strings.Index(s, "?>"); idx >= 0 {
		s = s[idx+2:]
	}
	return strings.TrimSpace(s)
}

// ---------- DOM Tree Builder ----------

func buildDOMNode(n *htmlNode, parentClass *data.ClassValue, ctx data.Context) data.Value {
	var cls data.ClassStmt
	var nodeType int

	if n.isText {
		cls = NewDOMTextClass()
		nodeType = 3 // XML_TEXT_NODE
	} else if n.isComment {
		cls = NewDOMCommentClass()
		nodeType = 8 // XML_COMMENT_NODE
	} else {
		cls = NewDOMElementClass()
		nodeType = 1 // XML_ELEMENT_NODE
	}

	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	// Set common properties
	cv.ObjectValue.SetProperty("nodeName", data.NewStringValue(n.tag))
	cv.ObjectValue.SetProperty("nodeValue", data.NewStringValue(n.text))
	cv.ObjectValue.SetProperty("nodeType", data.NewIntValue(nodeType))

	// For elements, set attributes
	if !n.isText && !n.isComment {
		attrObj := data.NewObjectValue()
		for k, v := range n.attrs {
			attrObj.SetProperty(k, data.NewStringValue(v))
		}
		cv.ObjectValue.SetProperty("_attributes", attrObj)
	}

	// Build children
	childNodes := make([]data.Value, 0, len(n.children))
	for _, child := range n.children {
		childVal := buildDOMNode(child, cv, ctx)
		if childVal != nil {
			childNodes = append(childNodes, childVal)
		}
	}
	cv.ObjectValue.SetProperty("childNodes", data.NewArrayValue(childNodes))

	return cv
}

// collectElementsByTagName recursively collects elements by tag name
func collectElementsByTagName(node *data.ClassValue, tagName string, results *[]data.Value) {
	nodeName, _ := node.GetProperty("nodeName")
	if nodeName != nil && strings.ToLower(nodeName.AsString()) == tagName {
		*results = append(*results, node)
	}

	children, _ := node.GetProperty("childNodes")
	if arr, ok := children.(*data.ArrayValue); ok {
		for _, zval := range arr.List {
			if child, ok := zval.Value.(*data.ClassValue); ok {
				collectElementsByTagName(child, tagName, results)
			}
		}
	}
}

// nodeToXML converts a DOM node back to XML string
func nodeToXML(val data.Value) string {
	cv, ok := val.(*data.ClassValue)
	if !ok {
		return val.AsString()
	}

	nodeName, _ := cv.GetProperty("nodeName")
	nodeValue, _ := cv.GetProperty("nodeValue")
	name := ""
	text := ""
	if nodeName != nil {
		name = nodeName.AsString()
	}
	if nodeValue != nil {
		text = nodeValue.AsString()
	}

	// Check if text node
	nodeType, _ := cv.GetProperty("nodeType")
	if nt, ok := nodeType.(*data.IntValue); ok && nt.Value == 3 {
		return text
	}

	// Element node
	xml := "<" + name

	// Attributes
	attrs, _ := cv.GetProperty("_attributes")
	if obj, ok := attrs.(*data.ObjectValue); ok {
		for k, v := range obj.GetProperties() {
			xml += " " + k + `="` + v.AsString() + `"`
		}
	}

	children, _ := cv.GetProperty("childNodes")
	if arr, ok := children.(*data.ArrayValue); ok && len(arr.List) > 0 {
		xml += ">"
		for _, zval := range arr.List {
			xml += nodeToXML(zval.Value)
		}
		xml += "</" + name + ">"
	} else {
		xml += "/>"
	}

	return xml
}
