package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const encrypterClassName = "Illuminate\\Encryption\\Encrypter"

var supportedCiphers = map[string]int{
	"aes-128-cbc": 16,
	"aes-256-cbc": 32,
}

type EncrypterClass struct {
	node.Node
	methods map[string]data.Method
}

func NewEncrypterClass() data.ClassStmt {
	c := &EncrypterClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *EncrypterClass) GetName() string                          { return encrypterClassName }
func (c *EncrypterClass) GetExtend() *string                       { return nil }
func (c *EncrypterClass) GetImplements() []string                  {
	return []string{"Illuminate\\Contracts\\Encryption\\Encrypter", "Illuminate\\Contracts\\Encryption\\StringEncrypter"}
}
func (c *EncrypterClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "key", "cipher", "previousKeys":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *EncrypterClass) GetPropertyList() []data.Property {
	out := make([]data.Property, 0, 3)
	for _, n := range []string{"key", "cipher", "previousKeys"} {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *EncrypterClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *EncrypterClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *EncrypterClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *EncrypterClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *EncrypterClass) GetStaticMethod(name string) (data.Method, bool) {
	lower := strings.ToLower(name)
	switch lower {
	case "supported", "generatekey", "appearsencrypted":
		if m, ok := c.methods[lower]; ok {
			return m, true
		}
	}
	return c.GetMethod(name)
}

func (c *EncrypterClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"key", "cipher"}, 1, encConstruct)
	c.methods["supported"] = kit.StaticMethod("supported", []string{"key", "cipher"}, -1, encSupported)
	c.methods["generatekey"] = kit.StaticMethod("generateKey", []string{"cipher"}, -1, encGenerateKey)
	c.methods["appearsencrypted"] = kit.StaticMethod("appearsEncrypted", []string{"value"}, -1, encAppearsEncrypted)
	c.methods["encrypt"] = kit.InstanceMethodOpt("encrypt", []string{"value", "serialize"}, 1, encEncrypt)
	c.methods["encryptstring"] = kit.InstanceMethod("encryptString", []string{"value"}, encEncryptString)
	c.methods["decrypt"] = kit.InstanceMethodOpt("decrypt", []string{"payload", "unserialize"}, 1, encDecrypt)
	c.methods["decryptstring"] = kit.InstanceMethod("decryptString", []string{"payload"}, encDecryptString)
	c.methods["getkey"] = kit.InstanceMethod("getKey", nil, encGetKey)
	c.methods["getallkeys"] = kit.InstanceMethod("getAllKeys", nil, encGetAllKeys)
	c.methods["getpreviouskeys"] = kit.InstanceMethod("getPreviousKeys", nil, encGetPreviousKeys)
	c.methods["previouskeys"] = kit.InstanceMethod("previousKeys", []string{"keys"}, encPreviousKeys)
}

func encRecv(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Encrypter missing $this"))
}

func encKeyCipher(cv *data.ClassValue) (key, cipher string, ctl data.Control) {
	kv, _ := cv.GetProperty("key")
	cvProp, _ := cv.GetProperty("cipher")
	if kv != nil {
		key = kv.AsString()
	}
	if cvProp != nil {
		cipher = strings.ToLower(cvProp.AsString())
	}
	if key == "" || supportedCiphers[cipher] == 0 || len(key) != supportedCiphers[cipher] {
		return "", "", data.NewErrorThrow(nil, fmt.Errorf("Unsupported cipher or incorrect key length."))
	}
	return key, cipher, nil
}

func encConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := encRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := ""
	if v := kit.Arg(ctx, 0); v != nil {
		key = v.AsString()
	}
	cipherName := "aes-128-cbc"
	if v := kit.Arg(ctx, 1); v != nil && !kit.IsNull(v) {
		cipherName = strings.ToLower(v.AsString())
	}
	if sz, ok := supportedCiphers[cipherName]; !ok || len(key) != sz {
		ciphers := "aes-128-cbc, aes-256-cbc"
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Unsupported cipher or incorrect key length. Supported ciphers are: %s.", ciphers))
	}
	_ = cv.SetProperty("key", data.NewStringValue(key))
	_ = cv.SetProperty("cipher", data.NewStringValue(cipherName))
	_ = cv.SetProperty("previousKeys", data.NewArrayValue(nil))
	return cv, nil
}

func encSupported(ctx data.Context) (data.GetValue, data.Control) {
	key := ""
	cipherName := "aes-128-cbc"
	if v := kit.Arg(ctx, 0); v != nil {
		key = v.AsString()
	}
	if v := kit.Arg(ctx, 1); v != nil {
		cipherName = strings.ToLower(v.AsString())
	}
	sz, ok := supportedCiphers[cipherName]
	return data.NewBoolValue(ok && len(key) == sz), nil
}

func encGenerateKey(ctx data.Context) (data.GetValue, data.Control) {
	cipherName := "aes-128-cbc"
	if v := kit.Arg(ctx, 0); v != nil && !kit.IsNull(v) {
		cipherName = strings.ToLower(v.AsString())
	}
	sz, ok := supportedCiphers[cipherName]
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Unsupported cipher or incorrect key length."))
	}
	buf := make([]byte, sz)
	if _, err := rand.Read(buf); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewStringValue(string(buf)), nil
}

func encAppearsEncrypted(ctx data.Context) (data.GetValue, data.Control) {
	v := kit.Unwrap(kit.Arg(ctx, 0))
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	if _, ok := v.(*data.StringValue); !ok {
		return data.NewBoolValue(false), nil
	}
	raw, err := base64.StdEncoding.DecodeString(v.AsString())
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return data.NewBoolValue(false), nil
	}
	_, hasIV := payload["iv"]
	_, hasValue := payload["value"]
	_, hasMAC := payload["mac"]
	return data.NewBoolValue(hasIV && hasValue && hasMAC), nil
}

func encGetKey(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := encRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	v, _ := cv.GetProperty("key")
	if v == nil {
		return data.NewStringValue(""), nil
	}
	return v, nil
}

func encGetPreviousKeys(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := encRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	v, _ := cv.GetProperty("previousKeys")
	if v == nil {
		return data.NewArrayValue(nil), nil
	}
	return v, nil
}

func encGetAllKeys(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := encRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	arr := data.NewArrayValue(nil).(*data.ArrayValue)
	if kv, _ := cv.GetProperty("key"); kv != nil {
		arr.SetIntKey(0, kv)
	}
	if pv, _ := cv.GetProperty("previousKeys"); pv != nil {
		if av, ok := kit.Unwrap(pv).(*data.ArrayValue); ok {
			i := len(arr.List)
			for _, e := range av.List {
				if e == nil {
					continue
				}
				arr.SetIntKey(i, e.Value)
				i++
			}
		}
	}
	return arr, nil
}

func encPreviousKeys(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := encRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cipherName := "aes-128-cbc"
	if cp, _ := cv.GetProperty("cipher"); cp != nil {
		cipherName = strings.ToLower(cp.AsString())
	}
	keysArg := kit.Unwrap(kit.Arg(ctx, 0))
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	i := 0
	appendKey := func(k data.Value) data.Control {
		if k == nil {
			return nil
		}
		key := k.AsString()
		sz, ok := supportedCiphers[cipherName]
		if !ok || len(key) != sz {
			return data.NewErrorThrow(nil, fmt.Errorf("Unsupported cipher or incorrect key length. Supported ciphers are: aes-128-cbc, aes-256-cbc."))
		}
		out.SetIntKey(i, k)
		i++
		return nil
	}
	switch t := keysArg.(type) {
	case *data.ArrayValue:
		for _, e := range t.List {
			if e == nil {
				continue
			}
			if ctl := appendKey(e.Value); ctl != nil {
				return nil, ctl
			}
		}
	case *data.ObjectValue:
		var ctl data.Control
		t.RangeProperties(func(_ string, val data.Value) bool {
			ctl = appendKey(val)
			return ctl == nil
		})
		if ctl != nil {
			return nil, ctl
		}
	}
	_ = cv.SetProperty("previousKeys", out)
	return cv, nil
}

func encEncryptString(ctx data.Context) (data.GetValue, data.Control) {
	return encEncryptPayload(ctx, false)
}

func encEncrypt(ctx data.Context) (data.GetValue, data.Control) {
	serialize := true
	if v := kit.Arg(ctx, 1); v != nil {
		if bv, ok := v.(data.AsBool); ok {
			serialize, _ = bv.AsBool()
		}
	}
	return encEncryptPayload(ctx, serialize)
}

func encEncryptPayload(ctx data.Context, serialize bool) (data.GetValue, data.Control) {
	cv, ctl := encRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, cipherName, ctl := encKeyCipher(cv)
	if ctl != nil {
		return nil, ctl
	}
	plain := ""
	if v := kit.Arg(ctx, 0); v != nil {
		if serialize {
			if fn, ok := ctx.GetVM().GetFunc("serialize"); ok {
				sctx := ctx.CreateContext(fn.GetVariables())
				_ = sctx.SetVariableValue(fn.GetVariables()[0], v)
				if out, ctl := fn.Call(sctx); ctl == nil && out != nil {
					if sv, ok := out.(data.Value); ok {
						plain = sv.AsString()
					}
				}
			}
		} else {
			plain = v.AsString()
		}
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Could not encrypt the data."))
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Could not encrypt the data."))
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	padded := pkcs7Pad([]byte(plain), aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	mode.CryptBlocks(ciphertext, padded)
	encVal := base64.StdEncoding.EncodeToString(ciphertext)
	ivB64 := base64.StdEncoding.EncodeToString(iv)
	mac := encMac(ivB64, encVal, key)
	payload, _ := json.Marshal(map[string]string{
		"iv":    ivB64,
		"value": encVal,
		"mac":   mac,
		"tag":   "",
	})
	_ = cipherName
	return data.NewStringValue(base64.StdEncoding.EncodeToString(payload)), nil
}

func encDecryptString(ctx data.Context) (data.GetValue, data.Control) {
	return encDecryptPayload(ctx, false)
}

func encDecrypt(ctx data.Context) (data.GetValue, data.Control) {
	unserialize := true
	if v := kit.Arg(ctx, 1); v != nil {
		if bv, ok := v.(data.AsBool); ok {
			unserialize, _ = bv.AsBool()
		}
	}
	return encDecryptPayload(ctx, unserialize)
}

func encDecryptPayload(ctx data.Context, unserialize bool) (data.GetValue, data.Control) {
	cv, ctl := encRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, cipherName, ctl := encKeyCipher(cv)
	if ctl != nil {
		return nil, ctl
	}
	payloadStr := kit.Arg(ctx, 0).AsString()
	raw, err := base64.StdEncoding.DecodeString(payloadStr)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("The payload is invalid."))
	}
	var payload map[string]string
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("The payload is invalid."))
	}
	if payload["iv"] == "" || payload["value"] == "" || payload["mac"] == "" {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("The payload is invalid."))
	}
	if !hmac.Equal([]byte(encMac(payload["iv"], payload["value"], key)), []byte(payload["mac"])) {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("The MAC is invalid."))
	}
	iv, _ := base64.StdEncoding.DecodeString(payload["iv"])
	ct, _ := base64.StdEncoding.DecodeString(payload["value"])
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Could not decrypt the data."))
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	if len(ct)%aes.BlockSize != 0 {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Could not decrypt the data."))
	}
	plain := make([]byte, len(ct))
	mode.CryptBlocks(plain, ct)
	plain, err = pkcs7Unpad(plain, aes.BlockSize)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Could not decrypt the data."))
	}
	_ = cipherName
	if unserialize {
		if fn, ok := ctx.GetVM().GetFunc("unserialize"); ok {
			sctx := ctx.CreateContext(fn.GetVariables())
			_ = sctx.SetVariableValue(fn.GetVariables()[0], data.NewStringValue(string(plain)))
			return fn.Call(sctx)
		}
	}
	return data.NewStringValue(string(plain)), nil
}

func encMac(iv, value, key string) string {
	m := hmac.New(sha256.New, []byte(key))
	m.Write([]byte(iv + value))
	return fmt.Sprintf("%x", m.Sum(nil))
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	pad := blockSize - len(data)%blockSize
	out := make([]byte, len(data)+pad)
	copy(out, data)
	for i := len(data); i < len(out); i++ {
		out[i] = byte(pad)
	}
	return out
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	pad := int(data[len(data)-1])
	if pad <= 0 || pad > blockSize {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:len(data)-pad], nil
}
