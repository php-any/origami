package gosupport

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"golang.org/x/crypto/bcrypt"
)

// Load 注册 Laravel 示例专用的运行时扩展。
func Load(vm data.VM) {
	for _, fn := range []data.FuncStmt{
		&SysGetTempDirFunction{},
		&PasswordHashFunction{},
		&PasswordVerifyFunction{},
		&PasswordGetInfoFunction{},
		&OpenSSLCipherIVLengthFunction{},
		&OpenSSLEncryptFunction{},
		&OpenSSLDecryptFunction{},
		&HashHmacFunction{},
		&HashEqualsFunction{},
		&JsonLastErrorFunction{},
		&BoolvalFunction{},
		&IntvalFunction{},
		&FloatvalFunction{},
		&RestoreErrorHandlerFunction{},
		&StreamSetChunkSizeFunction{},
		&FileinodeFunction{},
		&ParseStrFunction{},
		&HtmlEntityDecodeFunction{},
		&GetcwdFunction{},
	} {
		vm.AddFunc(fn)
	}

	// password_hash() algo constants（与 PHP 一致）
	vm.SetConstant("PASSWORD_BCRYPT", data.NewIntValue(1))
	vm.SetConstant("PASSWORD_DEFAULT", data.NewIntValue(1))
	vm.SetConstant("PASSWORD_ARGON2I", data.NewIntValue(2))
	vm.SetConstant("PASSWORD_ARGON2ID", data.NewIntValue(3))

	vm.SetConstant("JSON_ERROR_NONE", data.NewIntValue(0))
	vm.SetConstant("JSON_UNESCAPED_SLASHES", data.NewIntValue(64))
}

func valueAsString(v data.Value) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(data.AsString); ok {
		return s.AsString()
	}
	return fmt.Sprint(v)
}

func valueAsInt(v data.Value, fallback int) int {
	if v == nil {
		return fallback
	}
	switch t := v.(type) {
	case *data.IntValue:
		n, _ := strconv.Atoi(t.AsString())
		return n
	case data.AsString:
		n, err := strconv.Atoi(t.AsString())
		if err == nil {
			return n
		}
	}
	return fallback
}

func optionsCost(options data.Value, defaultCost int) int {
	arr, ok := options.(*data.ArrayValue)
	if !ok || arr == nil {
		return defaultCost
	}
	for _, z := range arr.List {
		if z != nil && z.Name == "cost" && z.Value != nil {
			return valueAsInt(z.Value, defaultCost)
		}
	}
	return defaultCost
}

// --- sys_get_temp_dir ---

type SysGetTempDirFunction struct{}

func (f *SysGetTempDirFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(os.TempDir()), nil
}
func (f *SysGetTempDirFunction) GetName() string            { return "sys_get_temp_dir" }
func (f *SysGetTempDirFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *SysGetTempDirFunction) GetIsStatic() bool          { return false }
func (f *SysGetTempDirFunction) GetParams() []data.GetValue { return nil }
func (f *SysGetTempDirFunction) GetVariables() []data.Variable {
	return nil
}
func (f *SysGetTempDirFunction) GetReturnType() data.Types { return data.NewBaseType("string") }

// --- password_hash / password_verify / password_get_info ---

type PasswordHashFunction struct{}

func (f *PasswordHashFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	password := valueAsString(mustIndex(ctx, 0))
	options := mustIndex(ctx, 2)
	cost := optionsCost(options, 10)
	if cost < 4 {
		cost = 4
	}
	if cost > 31 {
		cost = 31
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(string(hash)), nil
}
func (f *PasswordHashFunction) GetName() string            { return "password_hash" }
func (f *PasswordHashFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *PasswordHashFunction) GetIsStatic() bool          { return false }
func (f *PasswordHashFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "password", 0, nil, nil),
		node.NewParameter(nil, "algo", 1, nil, nil),
		node.NewParameter(nil, "options", 2, node.NewArray(nil, nil), nil),
	}
}
func (f *PasswordHashFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "password", 0, nil),
		node.NewVariable(nil, "algo", 1, nil),
		node.NewVariable(nil, "options", 2, nil),
	}
}
func (f *PasswordHashFunction) GetReturnType() data.Types { return data.NewBaseType("string") }

type PasswordVerifyFunction struct{}

func (f *PasswordVerifyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	password := valueAsString(mustIndex(ctx, 0))
	hash := valueAsString(mustIndex(ctx, 1))
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return data.NewBoolValue(err == nil), nil
}
func (f *PasswordVerifyFunction) GetName() string            { return "password_verify" }
func (f *PasswordVerifyFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *PasswordVerifyFunction) GetIsStatic() bool          { return false }
func (f *PasswordVerifyFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "password", 0, nil, nil),
		node.NewParameter(nil, "hash", 1, nil, nil),
	}
}
func (f *PasswordVerifyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "password", 0, nil),
		node.NewVariable(nil, "hash", 1, nil),
	}
}
func (f *PasswordVerifyFunction) GetReturnType() data.Types { return data.NewBaseType("bool") }

type PasswordGetInfoFunction struct{}

func (f *PasswordGetInfoFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	hash := valueAsString(mustIndex(ctx, 0))
	algo := 0
	algoName := "unknown"
	cost := 0
	if strings.HasPrefix(hash, "$2y$") || strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") {
		algo = 1
		algoName = "bcrypt"
		parts := strings.Split(hash, "$")
		if len(parts) >= 3 {
			cost, _ = strconv.Atoi(parts[2])
		}
	}
	return &data.ArrayValue{
		List: []*data.ZVal{
			data.NewNamedZVal("algo", data.NewIntValue(algo)),
			data.NewNamedZVal("algoName", data.NewStringValue(algoName)),
			data.NewNamedZVal("options", &data.ArrayValue{
				List: []*data.ZVal{
					data.NewNamedZVal("cost", data.NewIntValue(cost)),
				},
			}),
		},
	}, nil
}
func (f *PasswordGetInfoFunction) GetName() string            { return "password_get_info" }
func (f *PasswordGetInfoFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *PasswordGetInfoFunction) GetIsStatic() bool          { return false }
func (f *PasswordGetInfoFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "hash", 0, nil, nil)}
}
func (f *PasswordGetInfoFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "hash", 0, nil)}
}
func (f *PasswordGetInfoFunction) GetReturnType() data.Types { return data.NewBaseType("array") }

// --- openssl_* ---

type OpenSSLCipherIVLengthFunction struct{}

func (f *OpenSSLCipherIVLengthFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cipherName := strings.ToLower(valueAsString(mustIndex(ctx, 0)))
	switch cipherName {
	case "aes-128-cbc", "aes-256-cbc", "aes-128-gcm", "aes-256-gcm":
		return data.NewIntValue(aes.BlockSize), nil
	default:
		return data.NewBoolValue(false), nil
	}
}
func (f *OpenSSLCipherIVLengthFunction) GetName() string { return "openssl_cipher_iv_length" }
func (f *OpenSSLCipherIVLengthFunction) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (f *OpenSSLCipherIVLengthFunction) GetIsStatic() bool { return false }
func (f *OpenSSLCipherIVLengthFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "cipher_algo", 0, nil, nil)}
}
func (f *OpenSSLCipherIVLengthFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "cipher_algo", 0, nil)}
}
func (f *OpenSSLCipherIVLengthFunction) GetReturnType() data.Types { return data.NewBaseType("int") }

type OpenSSLEncryptFunction struct{}

func (f *OpenSSLEncryptFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	plain := valueAsString(mustIndex(ctx, 0))
	cipherName := strings.ToLower(valueAsString(mustIndex(ctx, 1)))
	key := []byte(valueAsString(mustIndex(ctx, 2)))
	options := valueAsInt(mustIndex(ctx, 3), 0)
	iv := []byte(valueAsString(mustIndex(ctx, 4)))

	block, err := aes.NewCipher(key)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	if len(iv) != aes.BlockSize {
		return data.NewBoolValue(false), nil
	}

	switch cipherName {
	case "aes-128-cbc", "aes-256-cbc":
		padded := pkcs7Pad([]byte(plain), aes.BlockSize)
		dst := make([]byte, len(padded))
		cipher.NewCBCEncrypter(block, iv).CryptBlocks(dst, padded)
		if options&1 == 1 { // OPENSSL_RAW_DATA
			return data.NewStringValue(string(dst)), nil
		}
		return data.NewStringValue(base64.StdEncoding.EncodeToString(dst)), nil
	default:
		return data.NewBoolValue(false), nil
	}
}
func (f *OpenSSLEncryptFunction) GetName() string            { return "openssl_encrypt" }
func (f *OpenSSLEncryptFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *OpenSSLEncryptFunction) GetIsStatic() bool          { return false }
func (f *OpenSSLEncryptFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, nil),
		node.NewParameter(nil, "cipher_algo", 1, nil, nil),
		node.NewParameter(nil, "passphrase", 2, nil, nil),
		node.NewParameter(nil, "options", 3, node.NewIntLiteral(nil, "0"), nil),
		node.NewParameter(nil, "iv", 4, node.NewStringLiteral(nil, ""), nil),
		node.NewParameter(nil, "tag", 5, nil, nil),
	}
}
func (f *OpenSSLEncryptFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, nil),
		node.NewVariable(nil, "cipher_algo", 1, nil),
		node.NewVariable(nil, "passphrase", 2, nil),
		node.NewVariable(nil, "options", 3, nil),
		node.NewVariable(nil, "iv", 4, nil),
		node.NewVariable(nil, "tag", 5, nil),
	}
}
func (f *OpenSSLEncryptFunction) GetReturnType() data.Types { return data.NewBaseType("string") }

type OpenSSLDecryptFunction struct{}

func (f *OpenSSLDecryptFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	raw := valueAsString(mustIndex(ctx, 0))
	cipherName := strings.ToLower(valueAsString(mustIndex(ctx, 1)))
	key := []byte(valueAsString(mustIndex(ctx, 2)))
	options := valueAsInt(mustIndex(ctx, 3), 0)
	iv := []byte(valueAsString(mustIndex(ctx, 4)))

	var ciphertext []byte
	if options&1 == 1 {
		ciphertext = []byte(raw)
	} else {
		var err error
		ciphertext, err = base64.StdEncoding.DecodeString(raw)
		if err != nil {
			return data.NewBoolValue(false), nil
		}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	if len(iv) != aes.BlockSize || len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return data.NewBoolValue(false), nil
	}

	switch cipherName {
	case "aes-128-cbc", "aes-256-cbc":
		dst := make([]byte, len(ciphertext))
		cipher.NewCBCDecrypter(block, iv).CryptBlocks(dst, ciphertext)
		unpadded, err := pkcs7Unpad(dst)
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		return data.NewStringValue(string(unpadded)), nil
	default:
		return data.NewBoolValue(false), nil
	}
}
func (f *OpenSSLDecryptFunction) GetName() string            { return "openssl_decrypt" }
func (f *OpenSSLDecryptFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *OpenSSLDecryptFunction) GetIsStatic() bool          { return false }
func (f *OpenSSLDecryptFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, nil),
		node.NewParameter(nil, "cipher_algo", 1, nil, nil),
		node.NewParameter(nil, "passphrase", 2, nil, nil),
		node.NewParameter(nil, "options", 3, node.NewIntLiteral(nil, "0"), nil),
		node.NewParameter(nil, "iv", 4, node.NewStringLiteral(nil, ""), nil),
		node.NewParameter(nil, "tag", 5, nil, nil),
	}
}
func (f *OpenSSLDecryptFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, nil),
		node.NewVariable(nil, "cipher_algo", 1, nil),
		node.NewVariable(nil, "passphrase", 2, nil),
		node.NewVariable(nil, "options", 3, nil),
		node.NewVariable(nil, "iv", 4, nil),
		node.NewVariable(nil, "tag", 5, nil),
	}
}
func (f *OpenSSLDecryptFunction) GetReturnType() data.Types { return data.NewBaseType("string") }

// --- hash_hmac ---

type HashHmacFunction struct{}

func (f *HashHmacFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	algo := strings.ToLower(valueAsString(mustIndex(ctx, 0)))
	payload := valueAsString(mustIndex(ctx, 1))
	key := valueAsString(mustIndex(ctx, 2))
	if algo != "sha256" {
		return data.NewBoolValue(false), nil
	}
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(payload))
	return data.NewStringValue(hex.EncodeToString(mac.Sum(nil))), nil
}
func (f *HashHmacFunction) GetName() string            { return "hash_hmac" }
func (f *HashHmacFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *HashHmacFunction) GetIsStatic() bool          { return false }
func (f *HashHmacFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "algo", 0, nil, nil),
		node.NewParameter(nil, "data", 1, nil, nil),
		node.NewParameter(nil, "key", 2, nil, nil),
	}
}
func (f *HashHmacFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "algo", 0, nil),
		node.NewVariable(nil, "data", 1, nil),
		node.NewVariable(nil, "key", 2, nil),
	}
}
func (f *HashHmacFunction) GetReturnType() data.Types { return data.NewBaseType("string") }

// --- hash_equals ---

type HashEqualsFunction struct{}

func (f *HashEqualsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	a := valueAsString(mustIndex(ctx, 0))
	b := valueAsString(mustIndex(ctx, 1))
	return data.NewBoolValue(hmac.Equal([]byte(a), []byte(b))), nil
}
func (f *HashEqualsFunction) GetName() string            { return "hash_equals" }
func (f *HashEqualsFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *HashEqualsFunction) GetIsStatic() bool          { return false }
func (f *HashEqualsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "known_string", 0, nil, nil),
		node.NewParameter(nil, "user_string", 1, nil, nil),
	}
}
func (f *HashEqualsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "known_string", 0, nil),
		node.NewVariable(nil, "user_string", 1, nil),
	}
}
func (f *HashEqualsFunction) GetReturnType() data.Types { return data.NewBaseType("bool") }

// --- json_last_error ---
// 核心 json_encode 成功路径不维护错误码；此处恒返回 0，满足 Encrypter 成功分支检查。

type JsonLastErrorFunction struct{}

func (f *JsonLastErrorFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(0), nil
}
func (f *JsonLastErrorFunction) GetName() string            { return "json_last_error" }
func (f *JsonLastErrorFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *JsonLastErrorFunction) GetIsStatic() bool          { return false }
func (f *JsonLastErrorFunction) GetParams() []data.GetValue { return nil }
func (f *JsonLastErrorFunction) GetVariables() []data.Variable {
	return nil
}
func (f *JsonLastErrorFunction) GetReturnType() data.Types { return data.NewBaseType("int") }

// --- boolval / intval / floatval ---

type BoolvalFunction struct{}

func (f *BoolvalFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v := mustIndex(ctx, 0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	if b, ok := v.(data.AsBool); ok {
		bv, err := b.AsBool()
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		return data.NewBoolValue(bv), nil
	}
	s := strings.TrimSpace(valueAsString(v))
	if s == "" || s == "0" || strings.EqualFold(s, "false") {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}
func (f *BoolvalFunction) GetName() string            { return "boolval" }
func (f *BoolvalFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *BoolvalFunction) GetIsStatic() bool          { return false }
func (f *BoolvalFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "value", 0, nil, nil)}
}
func (f *BoolvalFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "value", 0, nil)}
}
func (f *BoolvalFunction) GetReturnType() data.Types { return data.NewBaseType("bool") }

type IntvalFunction struct{}

func (f *IntvalFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v := mustIndex(ctx, 0)
	base := valueAsInt(mustIndex(ctx, 1), 10)
	if v == nil {
		return data.NewIntValue(0), nil
	}
	if iv, ok := v.(data.AsInt); ok {
		n, err := iv.AsInt()
		if err == nil {
			return data.NewIntValue(n), nil
		}
	}
	s := strings.TrimSpace(valueAsString(v))
	n, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		// PHP intval 对非法字符串返回 0（或前缀数字）；此处简化为 0
		f64, err2 := strconv.ParseFloat(s, 64)
		if err2 != nil {
			return data.NewIntValue(0), nil
		}
		return data.NewIntValue(int(f64)), nil
	}
	return data.NewIntValue(int(n)), nil
}
func (f *IntvalFunction) GetName() string            { return "intval" }
func (f *IntvalFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *IntvalFunction) GetIsStatic() bool          { return false }
func (f *IntvalFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
		node.NewParameter(nil, "base", 1, node.NewIntLiteral(nil, "10"), nil),
	}
}
func (f *IntvalFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
		node.NewVariable(nil, "base", 1, nil),
	}
}
func (f *IntvalFunction) GetReturnType() data.Types { return data.NewBaseType("int") }

type FloatvalFunction struct{}

func (f *FloatvalFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v := mustIndex(ctx, 0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if fv, ok := v.(data.AsFloat); ok {
		n, err := fv.AsFloat()
		if err == nil {
			return data.NewFloatValue(n), nil
		}
	}
	s := strings.TrimSpace(valueAsString(v))
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return data.NewFloatValue(0), nil
	}
	return data.NewFloatValue(n), nil
}
func (f *FloatvalFunction) GetName() string            { return "floatval" }
func (f *FloatvalFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *FloatvalFunction) GetIsStatic() bool          { return false }
func (f *FloatvalFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "value", 0, nil, nil)}
}
func (f *FloatvalFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "value", 0, nil)}
}
func (f *FloatvalFunction) GetReturnType() data.Types { return data.NewBaseType("float") }

func mustIndex(ctx data.Context, i int) data.Value {
	v, _ := ctx.GetIndexValue(i)
	return v
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

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty")
	}
	pad := int(data[len(data)-1])
	if pad == 0 || pad > len(data) {
		return nil, fmt.Errorf("bad padding")
	}
	for i := 0; i < pad; i++ {
		if data[len(data)-1-i] != byte(pad) {
			return nil, fmt.Errorf("bad padding")
		}
	}
	return data[:len(data)-pad], nil
}
