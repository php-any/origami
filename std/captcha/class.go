package captcha

import (
	crand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

const className = "Captcha"

type sessionProvider interface {
	EnsureSessionArray() *data.ObjectValue
}

// CaptchaClass Origami 原生验证码：生成短语、绘制 PNG/JPEG、写入 $_SESSION 并校验。
type CaptchaClass struct {
	node.Node
	methods map[string]data.Method
	props   []data.Property
}

func NewCaptchaClass() data.ClassStmt {
	c := &CaptchaClass{methods: map[string]data.Method{}}
	c.props = []data.Property{
		node.NewProperty(nil, "phrase", "protected", false, data.NewStringValue("")),
		node.NewProperty(nil, "display", "protected", false, data.NewStringValue("")),
		node.NewProperty(nil, "charset", "protected", false, data.NewStringValue(defaultCharset)),
		node.NewProperty(nil, "mode", "protected", false, data.NewStringValue(modeText)),
		node.NewProperty(nil, "length", "protected", false, data.NewIntValue(defaultLength)),
		node.NewProperty(nil, "width", "protected", false, data.NewIntValue(defaultWidth)),
		node.NewProperty(nil, "height", "protected", false, data.NewIntValue(defaultHeight)),
		node.NewProperty(nil, "noise", "protected", false, data.NewIntValue(defaultNoise)),
		node.NewProperty(nil, "ignoreCase", "protected", false, data.NewBoolValue(true)),
		node.NewProperty(nil, "expire", "protected", false, data.NewIntValue(defaultExpire)),
		node.NewProperty(nil, "key", "protected", false, data.NewStringValue(defaultKey)),
		node.NewProperty(nil, "png", "protected", false, data.NewStringValue("")),
	}
	c.methods["__construct"] = capMethod("__construct", []string{"phrase"}, false, capConstruct)
	c.methods["setlength"] = capMethod("setLength", []string{"length"}, false, capSetLength)
	c.methods["setcharset"] = capMethod("setCharset", []string{"charset"}, false, capSetCharset)
	c.methods["setsize"] = capMethod("setSize", []string{"width", "height"}, false, capSetSize)
	c.methods["setnoise"] = capMethod("setNoise", []string{"level"}, false, capSetNoise)
	c.methods["setignorecase"] = capMethod("setIgnoreCase", []string{"ignore"}, false, capSetIgnoreCase)
	c.methods["setexpire"] = capMethod("setExpire", []string{"seconds"}, false, capSetExpire)
	c.methods["setmode"] = capMethod("setMode", []string{"mode"}, false, capSetMode)
	c.methods["setkey"] = capMethod("setKey", []string{"key"}, false, capSetKey)
	c.methods["setphrase"] = capMethod("setPhrase", []string{"phrase"}, false, capSetPhrase)
	c.methods["build"] = capMethod("build", []string{"width", "height"}, false, capBuild)
	c.methods["getphrase"] = capMethod("getPhrase", nil, false, capGetPhrase)
	c.methods["getdisplay"] = capMethod("getDisplay", nil, false, capGetDisplay)
	c.methods["png"] = capMethod("png", nil, false, capPNG)
	c.methods["jpeg"] = capMethod("jpeg", []string{"quality"}, false, capJPEG)
	c.methods["base64"] = capMethod("base64", nil, false, capBase64)
	c.methods["inline"] = capMethod("inline", nil, false, capInline)
	c.methods["src"] = capMethod("src", nil, false, capInline)
	c.methods["img"] = capMethod("img", []string{"alt"}, false, capImg)
	c.methods["save"] = capMethod("save", []string{"path"}, false, capSave)
	c.methods["output"] = capMethod("output", []string{"format"}, false, capOutput)
	c.methods["test"] = capMethod("test", []string{"input"}, false, capTest)
	c.methods["verify"] = capMethod("verify", []string{"input"}, false, capTest)
	c.methods["store"] = capMethod("store", []string{"key"}, false, capStore)
	c.methods["__tostring"] = capMethod("__toString", nil, false, capInline)
	c.methods["create"] = capMethod("create", []string{"options"}, true, capCreate)
	c.methods["check"] = capMethod("check", []string{"input", "key"}, true, capCheck)
	c.methods["checkkeep"] = capMethod("checkKeep", []string{"input", "key"}, true, capCheckKeep)
	return c
}

func (c *CaptchaClass) GetName() string         { return className }
func (c *CaptchaClass) GetExtend() *string      { return nil }
func (c *CaptchaClass) GetImplements() []string { return []string{"Stringable"} }
func (c *CaptchaClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.props {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *CaptchaClass) GetPropertyList() []data.Property { return c.props }
func (c *CaptchaClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *CaptchaClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	base := ctx
	if ctx != nil {
		base = ctx.CreateBaseContext()
	}
	return data.NewClassValue(c, base), nil
}
func (c *CaptchaClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *CaptchaClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *CaptchaClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}
func (c *CaptchaClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "MODE_TEXT":
		return data.NewStringValue(modeText), true
	case "MODE_MATH":
		return data.NewStringValue(modeMath), true
	case "DEFAULT_KEY":
		return data.NewStringValue(defaultKey), true
	}
	return nil, false
}

type captchaMethod struct {
	name   string
	params []string
	static bool
	fn     func(data.Context) (data.GetValue, data.Control)
}

func capMethod(name string, params []string, static bool, fn func(data.Context) (data.GetValue, data.Control)) data.Method {
	return &captchaMethod{name: name, params: params, static: static, fn: fn}
}

func (m *captchaMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return m.fn(ctx)
}
func (m *captchaMethod) GetName() string { return m.name }
func (m *captchaMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *captchaMethod) GetIsStatic() bool { return m.static }
func (m *captchaMethod) GetReturnType() data.Types {
	return nil
}
func (m *captchaMethod) GetParams() []data.GetValue {
	out := make([]data.GetValue, len(m.params))
	for i, p := range m.params {
		var def data.GetValue
		switch p {
		case "phrase", "charset", "mode", "key", "input", "options", "alt", "path", "format", "width", "height":
			def = data.NewNullValue()
		case "quality":
			def = data.NewIntValue(80)
		case "length":
			def = data.NewIntValue(defaultLength)
		case "level":
			def = data.NewIntValue(defaultNoise)
		case "seconds":
			def = data.NewIntValue(defaultExpire)
		case "ignore":
			def = data.NewBoolValue(true)
		}
		out[i] = node.NewParameter(nil, p, i, def, nil)
	}
	return out
}
func (m *captchaMethod) GetVariables() []data.Variable {
	out := make([]data.Variable, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewVariable(nil, p, i, nil)
	}
	return out
}

func captchaSelf(ctx data.Context) *data.ClassValue {
	if c, ok := ctx.(*data.ClassMethodContext); ok {
		return c.ClassValue
	}
	if c, ok := ctx.(*data.ClassValue); ok {
		return c
	}
	return nil
}

func returnThis(ctx data.Context) (data.GetValue, data.Control) {
	if cv := captchaSelf(ctx); cv != nil {
		return data.NewThisValue(cv), nil
	}
	return data.NewNullValue(), nil
}

func needSelf(ctx data.Context) (*data.ClassValue, data.Control) {
	cv := captchaSelf(ctx)
	if cv == nil {
		return nil, utils.NewThrow(errors.New("Captcha 实例方法需要对象"))
	}
	return cv, nil
}

func capConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	if v, ok := ctx.GetIndexValue(0); ok && !isNullish(v) {
		s := strings.TrimSpace(v.AsString())
		if s != "" {
			if acl := cv.SetProperty("phrase", data.NewStringValue(s)); acl != nil {
				return nil, acl
			}
			if acl := cv.SetProperty("display", data.NewStringValue(s)); acl != nil {
				return nil, acl
			}
		}
	}
	return data.NewNullValue(), nil
}

func capSetLength(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	n := argInt(ctx, 0, defaultLength)
	if acl := cv.SetProperty("length", data.NewIntValue(clampInt(n, 1, 16))); acl != nil {
		return nil, acl
	}
	return invalidateImage(ctx)
}

func capSetCharset(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	s := argString(ctx, 0, defaultCharset)
	if s == "" {
		s = defaultCharset
	}
	if acl := cv.SetProperty("charset", data.NewStringValue(s)); acl != nil {
		return nil, acl
	}
	return invalidateImage(ctx)
}

func capSetSize(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	w := clampInt(argInt(ctx, 0, defaultWidth), 40, 800)
	h := clampInt(argInt(ctx, 1, defaultHeight), 20, 400)
	if acl := cv.SetProperty("width", data.NewIntValue(w)); acl != nil {
		return nil, acl
	}
	if acl := cv.SetProperty("height", data.NewIntValue(h)); acl != nil {
		return nil, acl
	}
	return invalidateImage(ctx)
}

func capSetNoise(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	if acl := cv.SetProperty("noise", data.NewIntValue(clampInt(argInt(ctx, 0, defaultNoise), 0, 10))); acl != nil {
		return nil, acl
	}
	return invalidateImage(ctx)
}

func capSetIgnoreCase(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	if acl := cv.SetProperty("ignoreCase", data.NewBoolValue(argBool(ctx, 0, true))); acl != nil {
		return nil, acl
	}
	return returnThis(ctx)
}

func capSetExpire(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	n := argInt(ctx, 0, defaultExpire)
	if n < 0 {
		n = 0
	}
	if acl := cv.SetProperty("expire", data.NewIntValue(n)); acl != nil {
		return nil, acl
	}
	return returnThis(ctx)
}

func capSetMode(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	mode := strings.ToLower(argString(ctx, 0, modeText))
	if mode != modeMath {
		mode = modeText
	}
	if acl := cv.SetProperty("mode", data.NewStringValue(mode)); acl != nil {
		return nil, acl
	}
	return invalidateImage(ctx)
}

func capSetKey(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	key := argString(ctx, 0, defaultKey)
	if key == "" {
		key = defaultKey
	}
	if acl := cv.SetProperty("key", data.NewStringValue(key)); acl != nil {
		return nil, acl
	}
	return returnThis(ctx)
}

func capSetPhrase(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	s := argString(ctx, 0, "")
	if acl := cv.SetProperty("phrase", data.NewStringValue(s)); acl != nil {
		return nil, acl
	}
	if acl := cv.SetProperty("display", data.NewStringValue(s)); acl != nil {
		return nil, acl
	}
	return invalidateImage(ctx)
}

func capBuild(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	if v, ok := ctx.GetIndexValue(0); ok && !isNullish(v) {
		if acl := cv.SetProperty("width", data.NewIntValue(clampInt(valueInt(v, defaultWidth), 40, 800))); acl != nil {
			return nil, acl
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok && !isNullish(v) {
		if acl := cv.SetProperty("height", data.NewIntValue(clampInt(valueInt(v, defaultHeight), 20, 400))); acl != nil {
			return nil, acl
		}
	}
	if acl := doBuild(cv); acl != nil {
		return nil, acl
	}
	return returnThis(ctx)
}

func capGetPhrase(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	return data.NewStringValue(propString(cv, "phrase", "")), nil
}

func capGetDisplay(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	d := propString(cv, "display", "")
	if d == "" {
		d = propString(cv, "phrase", "")
	}
	return data.NewStringValue(d), nil
}

func capPNG(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	bin, acl := ensurePNG(cv)
	if acl != nil {
		return nil, acl
	}
	return data.NewStringValue(bin), nil
}

func capJPEG(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	quality := argInt(ctx, 0, 80)
	bin, acl := ensurePNG(cv)
	if acl != nil {
		return nil, acl
	}
	img, err := decodePNGBytes([]byte(bin))
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	out, err := encodeJPEG(img, quality)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewStringValue(string(out)), nil
}

func capBase64(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	bin, acl := ensurePNG(cv)
	if acl != nil {
		return nil, acl
	}
	return data.NewStringValue(base64.StdEncoding.EncodeToString([]byte(bin))), nil
}

func capInline(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	bin, acl := ensurePNG(cv)
	if acl != nil {
		return nil, acl
	}
	uri := "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(bin))
	return data.NewStringValue(uri), nil
}

func capImg(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	bin, acl := ensurePNG(cv)
	if acl != nil {
		return nil, acl
	}
	alt := argString(ctx, 0, "captcha")
	if alt == "" {
		alt = "captcha"
	}
	w := propInt(cv, "width", defaultWidth)
	h := propInt(cv, "height", defaultHeight)
	uri := "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(bin))
	html := fmt.Sprintf(`<img src="%s" alt="%s" width="%d" height="%d">`, uri, htmlEscape(alt), w, h)
	return data.NewStringValue(html), nil
}

func capSave(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	path := argString(ctx, 0, "")
	if path == "" {
		return data.NewBoolValue(false), nil
	}
	lower := strings.ToLower(path)
	var blob []byte
	if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") {
		pngBin, acl := ensurePNG(cv)
		if acl != nil {
			return nil, acl
		}
		img, err := decodePNGBytes([]byte(pngBin))
		if err != nil {
			return nil, utils.NewThrow(err)
		}
		blob, err = encodeJPEG(img, 80)
		if err != nil {
			return nil, utils.NewThrow(err)
		}
	} else {
		pngBin, acl := ensurePNG(cv)
		if acl != nil {
			return nil, acl
		}
		blob = []byte(pngBin)
	}
	if err := os.WriteFile(path, blob, 0644); err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func capOutput(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	format := strings.ToLower(argString(ctx, 0, "png"))
	var body string
	ctype := "image/png"
	if format == "jpg" || format == "jpeg" {
		ctype = "image/jpeg"
		gv, acl := capJPEG(ctx)
		if acl != nil {
			return nil, acl
		}
		body = valueAsString(gv)
	} else {
		bin, acl := ensurePNG(cv)
		if acl != nil {
			return nil, acl
		}
		body = bin
	}
	setImageHeader(ctx, ctype)
	if acl := data.EmitOutput(ctx, body); acl != nil {
		return nil, acl
	}
	return data.NewNullValue(), nil
}

func capTest(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	input := argString(ctx, 0, "")
	ok := matchPhrase(propString(cv, "phrase", ""), input, propBool(cv, "ignoreCase", true))
	return data.NewBoolValue(ok), nil
}

func capStore(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := needSelf(ctx)
	if acl != nil {
		return nil, acl
	}
	if v, ok := ctx.GetIndexValue(0); ok && !isNullish(v) {
		key := strings.TrimSpace(v.AsString())
		if key != "" {
			if acl := cv.SetProperty("key", data.NewStringValue(key)); acl != nil {
				return nil, acl
			}
		}
	}
	if propString(cv, "phrase", "") == "" {
		if acl := doBuild(cv); acl != nil {
			return nil, acl
		}
	}
	if acl := storePhrase(ctx, cv); acl != nil {
		return nil, acl
	}
	return returnThis(ctx)
}

func capCreate(ctx data.Context) (data.GetValue, data.Control) {
	cls := NewCaptchaClass()
	gv, acl := cls.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	cv, ok := gv.(*data.ClassValue)
	if !ok {
		return gv, nil
	}
	opts, _ := ctx.GetIndexValue(0)
	if acl := applyOptions(cv, opts); acl != nil {
		return nil, acl
	}
	if acl := doBuild(cv); acl != nil {
		return nil, acl
	}
	store := true
	if opts != nil {
		if v, ok := optionValue(opts, "store"); ok {
			store = valueBool(v, true)
		}
	}
	if store {
		if acl := storePhrase(ctx, cv); acl != nil {
			return nil, acl
		}
	}
	return cv, nil
}

func capCheck(ctx data.Context) (data.GetValue, data.Control) {
	return sessionCheck(ctx, true)
}

func capCheckKeep(ctx data.Context) (data.GetValue, data.Control) {
	return sessionCheck(ctx, false)
}

func invalidateImage(ctx data.Context) (data.GetValue, data.Control) {
	if cv := captchaSelf(ctx); cv != nil {
		if acl := cv.SetProperty("png", data.NewStringValue("")); acl != nil {
			return nil, acl
		}
	}
	return returnThis(ctx)
}

func doBuild(cv *data.ClassValue) data.Control {
	mode := strings.ToLower(propString(cv, "mode", modeText))
	phrase := propString(cv, "phrase", "")
	display := propString(cv, "display", "")
	if phrase == "" {
		var err error
		if mode == modeMath {
			display, phrase, err = randomMath()
		} else {
			phrase, err = randomPhrase(propString(cv, "charset", defaultCharset), propInt(cv, "length", defaultLength))
			display = phrase
		}
		if err != nil {
			return utils.NewThrow(err)
		}
		if acl := cv.SetProperty("phrase", data.NewStringValue(phrase)); acl != nil {
			return acl
		}
		if acl := cv.SetProperty("display", data.NewStringValue(display)); acl != nil {
			return acl
		}
	} else if display == "" {
		display = phrase
		if acl := cv.SetProperty("display", data.NewStringValue(display)); acl != nil {
			return acl
		}
	}
	seed := [32]byte{}
	if _, err := crand.Read(seed[:]); err != nil {
		return utils.NewThrow(err)
	}
	rng := newVisualRNG(seed)
	img := renderImage(RenderOptions{
		Phrase: display,
		Width:  propInt(cv, "width", defaultWidth),
		Height: propInt(cv, "height", defaultHeight),
		Noise:  propInt(cv, "noise", defaultNoise),
		Rng:    rng,
	})
	bin, err := encodePNG(img)
	if err != nil {
		return utils.NewThrow(err)
	}
	return cv.SetProperty("png", data.NewStringValue(string(bin)))
}

func ensurePNG(cv *data.ClassValue) (string, data.Control) {
	pngBin := propString(cv, "png", "")
	if pngBin != "" {
		return pngBin, nil
	}
	if acl := doBuild(cv); acl != nil {
		return "", acl
	}
	return propString(cv, "png", ""), nil
}

func storePhrase(ctx data.Context, cv *data.ClassValue) data.Control {
	sess := sessionArray(ctx)
	if sess == nil {
		return utils.NewThrow(errors.New("无法访问 $_SESSION"))
	}
	key := propString(cv, "key", defaultKey)
	if key == "" {
		key = defaultKey
	}
	ttl := propInt(cv, "expire", defaultExpire)
	expireAt := 0
	if ttl > 0 {
		expireAt = int(time.Now().Unix()) + ttl
	}
	arr := data.NewArrayValue(nil).(*data.ArrayValue)
	arr.List = []*data.ZVal{
		data.NewNamedZVal("phrase", data.NewStringValue(propString(cv, "phrase", ""))),
		data.NewNamedZVal("expire", data.NewIntValue(expireAt)),
	}
	return sess.SetProperty(key, arr)
}

func sessionCheck(ctx data.Context, consume bool) (data.GetValue, data.Control) {
	input := argString(ctx, 0, "")
	key := argString(ctx, 1, defaultKey)
	if key == "" {
		key = defaultKey
	}
	sess := sessionArray(ctx)
	if sess == nil || !sess.HasProperty(key) {
		return data.NewBoolValue(false), nil
	}
	raw, acl := sess.GetProperty(key)
	if acl != nil {
		return nil, acl
	}
	phrase, expireAt, ignoreCase := readStored(raw)
	if phrase == "" {
		return data.NewBoolValue(false), nil
	}
	if expireAt > 0 && int(time.Now().Unix()) > expireAt {
		sess.UnsetProperty(key)
		return data.NewBoolValue(false), nil
	}
	ok := matchPhrase(phrase, input, ignoreCase)
	if ok && consume {
		sess.UnsetProperty(key)
	}
	return data.NewBoolValue(ok), nil
}

func readStored(v data.Value) (phrase string, expireAt int, ignoreCase bool) {
	ignoreCase = true
	switch t := v.(type) {
	case *data.StringValue:
		return t.AsString(), 0, true
	case *data.ArrayValue:
		for _, z := range t.List {
			if z == nil || z.Value == nil {
				continue
			}
			switch z.Name {
			case "phrase":
				phrase = z.Value.AsString()
			case "expire":
				expireAt = valueInt(z.Value, 0)
			case "ignoreCase":
				ignoreCase = valueBool(z.Value, true)
			}
		}
		return phrase, expireAt, ignoreCase
	case *data.ObjectValue:
		if p, _ := t.GetProperty("phrase"); p != nil {
			phrase = p.AsString()
		}
		if e, _ := t.GetProperty("expire"); e != nil {
			expireAt = valueInt(e, 0)
		}
		return phrase, expireAt, true
	default:
		if v != nil {
			return v.AsString(), 0, true
		}
	}
	return "", 0, true
}

func sessionArray(ctx data.Context) *data.ObjectValue {
	if ctx == nil {
		return nil
	}
	if p, ok := ctx.GetVM().(sessionProvider); ok {
		return p.EnsureSessionArray()
	}
	return nil
}

func applyOptions(cv *data.ClassValue, opts data.Value) data.Control {
	if opts == nil || isNullish(opts) {
		return nil
	}
	set := func(name string, val data.Value) data.Control {
		return cv.SetProperty(name, val)
	}
	pairs := optionMap(opts)
	if v, ok := pairs["length"]; ok {
		if acl := set("length", data.NewIntValue(clampInt(valueInt(v, defaultLength), 1, 16))); acl != nil {
			return acl
		}
	}
	if v, ok := pairs["width"]; ok {
		if acl := set("width", data.NewIntValue(clampInt(valueInt(v, defaultWidth), 40, 800))); acl != nil {
			return acl
		}
	}
	if v, ok := pairs["height"]; ok {
		if acl := set("height", data.NewIntValue(clampInt(valueInt(v, defaultHeight), 20, 400))); acl != nil {
			return acl
		}
	}
	if v, ok := pairs["noise"]; ok {
		if acl := set("noise", data.NewIntValue(clampInt(valueInt(v, defaultNoise), 0, 10))); acl != nil {
			return acl
		}
	}
	if v, ok := pairs["charset"]; ok {
		s := v.AsString()
		if s == "" {
			s = defaultCharset
		}
		if acl := set("charset", data.NewStringValue(s)); acl != nil {
			return acl
		}
	}
	if v, ok := pairs["mode"]; ok {
		mode := strings.ToLower(v.AsString())
		if mode != modeMath {
			mode = modeText
		}
		if acl := set("mode", data.NewStringValue(mode)); acl != nil {
			return acl
		}
	}
	if v, ok := pairs["ignoreCase"]; ok {
		if acl := set("ignoreCase", data.NewBoolValue(valueBool(v, true))); acl != nil {
			return acl
		}
	}
	if v, ok := pairs["expire"]; ok {
		n := valueInt(v, defaultExpire)
		if n < 0 {
			n = 0
		}
		if acl := set("expire", data.NewIntValue(n)); acl != nil {
			return acl
		}
	}
	if v, ok := pairs["key"]; ok {
		key := v.AsString()
		if key == "" {
			key = defaultKey
		}
		if acl := set("key", data.NewStringValue(key)); acl != nil {
			return acl
		}
	}
	if v, ok := pairs["phrase"]; ok {
		s := v.AsString()
		if acl := set("phrase", data.NewStringValue(s)); acl != nil {
			return acl
		}
		if acl := set("display", data.NewStringValue(s)); acl != nil {
			return acl
		}
	}
	return nil
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		`&`, "&amp;",
		`<`, "&lt;",
		`>`, "&gt;",
		`"`, "&quot;",
	)
	return replacer.Replace(s)
}
