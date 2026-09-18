// Package captcha 是 Origami 原生验证码标准库。
// 提供 Captcha 类与 captcha_img / captcha_src / captcha_check 函数：
// 生成短语或算术题、绘制 PNG/JPEG、写入 $_SESSION 并一次性校验。
package captcha

import (
	"encoding/base64"
	"errors"
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func Load(vm data.VM) {
	vm.AddClass(NewCaptchaClass())
	for _, fn := range []data.FuncStmt{
		NewCaptchaImgFunction(),
		NewCaptchaSrcFunction(),
		NewCaptchaCheckFunction(),
	} {
		vm.AddFunc(fn)
	}
}

func encodeB64(bin string) string {
	return base64.StdEncoding.EncodeToString([]byte(bin))
}

type CaptchaImgFunction struct{}

func NewCaptchaImgFunction() data.FuncStmt { return &CaptchaImgFunction{} }

func (f *CaptchaImgFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := createFromArgs(ctx)
	if acl != nil {
		return nil, acl
	}
	bin, acl := ensurePNG(cv)
	if acl != nil {
		return nil, acl
	}
	w := propInt(cv, "width", defaultWidth)
	h := propInt(cv, "height", defaultHeight)
	html := `<img src="data:image/png;base64,` + encodeB64(bin) + `" alt="captcha" width="` +
		strconv.Itoa(w) + `" height="` + strconv.Itoa(h) + `">`
	return data.NewStringValue(html), nil
}

func (f *CaptchaImgFunction) GetName() string { return "captcha_img" }
func (f *CaptchaImgFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "options", 0, data.NewNullValue(), nil)}
}
func (f *CaptchaImgFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "options", 0, nil)}
}

type CaptchaSrcFunction struct{}

func NewCaptchaSrcFunction() data.FuncStmt { return &CaptchaSrcFunction{} }

func (f *CaptchaSrcFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv, acl := createFromArgs(ctx)
	if acl != nil {
		return nil, acl
	}
	bin, acl := ensurePNG(cv)
	if acl != nil {
		return nil, acl
	}
	return data.NewStringValue("data:image/png;base64," + encodeB64(bin)), nil
}

func (f *CaptchaSrcFunction) GetName() string { return "captcha_src" }
func (f *CaptchaSrcFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "options", 0, data.NewNullValue(), nil)}
}
func (f *CaptchaSrcFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "options", 0, nil)}
}

type CaptchaCheckFunction struct{}

func NewCaptchaCheckFunction() data.FuncStmt { return &CaptchaCheckFunction{} }

func (f *CaptchaCheckFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return sessionCheck(ctx, true)
}

func (f *CaptchaCheckFunction) GetName() string { return "captcha_check" }
func (f *CaptchaCheckFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "input", 0, nil, nil),
		node.NewParameter(nil, "key", 1, data.NewNullValue(), nil),
	}
}
func (f *CaptchaCheckFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "input", 0, nil),
		node.NewVariable(nil, "key", 1, nil),
	}
}

func createFromArgs(ctx data.Context) (*data.ClassValue, data.Control) {
	gv, acl := capCreate(ctx)
	if acl != nil {
		return nil, acl
	}
	if cv, ok := gv.(*data.ClassValue); ok {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, errors.New("captcha_img/captcha_src 未能创建 Captcha 实例"))
}
