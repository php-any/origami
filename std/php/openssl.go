package php

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// OpenSSLCipherIVLengthFunction 实现 openssl_cipher_iv_length
type OpenSSLCipherIVLengthFunction struct{}

func NewOpenSSLCipherIVLengthFunction() data.FuncStmt {
	return &OpenSSLCipherIVLengthFunction{}
}

func (f *OpenSSLCipherIVLengthFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cipherName := ""
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		cipherName = v.AsString()
	}
	cipherName = strings.ToLower(cipherName)

	switch cipherName {
	case "aes-128-cbc", "aes-256-cbc", "aes-128-gcm", "aes-256-gcm":
		return data.NewIntValue(aes.BlockSize), nil
	default:
		return data.NewBoolValue(false), nil
	}
}

func (f *OpenSSLCipherIVLengthFunction) GetName() string { return "openssl_cipher_iv_length" }
var openSSLCipherIVLengthFunctionGetParams = []data.GetValue{node.NewParameter(nil, "cipher_algo", 0, nil, nil)}

func (f *OpenSSLCipherIVLengthFunction) GetParams() []data.GetValue {
	return openSSLCipherIVLengthFunctionGetParams
}
var openSSLCipherIVLengthFunctionGetVariables = []data.Variable{node.NewVariable(nil, "cipher_algo", 0, nil)}

func (f *OpenSSLCipherIVLengthFunction) GetVariables() []data.Variable {
	return openSSLCipherIVLengthFunctionGetVariables
}

// OpenSSLEncryptFunction 实现 openssl_encrypt
type OpenSSLEncryptFunction struct{}

func NewOpenSSLEncryptFunction() data.FuncStmt {
	return &OpenSSLEncryptFunction{}
}

func (f *OpenSSLEncryptFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	plain := ""
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		plain = v.AsString()
	}
	cipherName := ""
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		cipherName = v.AsString()
	}
	cipherName = strings.ToLower(cipherName)
	keyStr := ""
	if v, ok := ctx.GetIndexValue(2); ok && v != nil {
		keyStr = v.AsString()
	}
	options := 0
	if v, ok := ctx.GetIndexValue(3); ok && v != nil {
		if n, err := asIntVal(v); err == nil {
			options = n
		}
	}
	ivStr := ""
	if v, ok := ctx.GetIndexValue(4); ok && v != nil {
		ivStr = v.AsString()
	}

	key := []byte(keyStr)
	iv := []byte(ivStr)

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
	case "aes-128-gcm", "aes-256-gcm":
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		sealed := gcm.Seal(nil, iv, []byte(plain), nil)
		if options&1 == 1 {
			return data.NewStringValue(string(sealed)), nil
		}
		return data.NewStringValue(base64.StdEncoding.EncodeToString(sealed)), nil
	default:
		return data.NewBoolValue(false), nil
	}
}

func (f *OpenSSLEncryptFunction) GetName() string { return "openssl_encrypt" }
var openSSLEncryptFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "data", 0, nil, nil),
	node.NewParameter(nil, "cipher_algo", 1, nil, nil),
	node.NewParameter(nil, "passphrase", 2, nil, nil),
	node.NewParameter(nil, "options", 3, node.NewIntLiteral(nil, "0"), nil),
	node.NewParameter(nil, "iv", 4, node.NewStringLiteral(nil, ""), nil),
	node.NewParameter(nil, "tag", 5, nil, nil),
}

func (f *OpenSSLEncryptFunction) GetParams() []data.GetValue {
	return openSSLEncryptFunctionGetParams
}
var openSSLEncryptFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "data", 0, nil),
	node.NewVariable(nil, "cipher_algo", 1, nil),
	node.NewVariable(nil, "passphrase", 2, nil),
	node.NewVariable(nil, "options", 3, nil),
	node.NewVariable(nil, "iv", 4, nil),
	node.NewVariable(nil, "tag", 5, nil),
}

func (f *OpenSSLEncryptFunction) GetVariables() []data.Variable {
	return openSSLEncryptFunctionGetVariables
}

// OpenSSLDecryptFunction 实现 openssl_decrypt
type OpenSSLDecryptFunction struct{}

func NewOpenSSLDecryptFunction() data.FuncStmt {
	return &OpenSSLDecryptFunction{}
}

func (f *OpenSSLDecryptFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	raw := ""
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		raw = v.AsString()
	}
	cipherName := ""
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		cipherName = v.AsString()
	}
	cipherName = strings.ToLower(cipherName)
	keyStr := ""
	if v, ok := ctx.GetIndexValue(2); ok && v != nil {
		keyStr = v.AsString()
	}
	options := 0
	if v, ok := ctx.GetIndexValue(3); ok && v != nil {
		if n, err := asIntVal(v); err == nil {
			options = n
		}
	}
	ivStr := ""
	if v, ok := ctx.GetIndexValue(4); ok && v != nil {
		ivStr = v.AsString()
	}

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

	key := []byte(keyStr)
	iv := []byte(ivStr)

	block, err := aes.NewCipher(key)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	if len(iv) != aes.BlockSize {
		return data.NewBoolValue(false), nil
	}

	switch cipherName {
	case "aes-128-cbc", "aes-256-cbc":
		if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
			return data.NewBoolValue(false), nil
		}
		dst := make([]byte, len(ciphertext))
		cipher.NewCBCDecrypter(block, iv).CryptBlocks(dst, ciphertext)
		unpadded, err := pkcs7Unpad(dst)
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		return data.NewStringValue(string(unpadded)), nil
	case "aes-128-gcm", "aes-256-gcm":
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		if len(ciphertext) < gcm.Overhead() {
			return data.NewBoolValue(false), nil
		}
		// Need separate tag parameter for GCM
		tagStr := ""
		if v, ok := ctx.GetIndexValue(5); ok && v != nil {
			tagStr = v.AsString()
		}
		var tag []byte
		if options&1 == 1 {
			tag = []byte(tagStr)
		} else {
			tag, err = base64.StdEncoding.DecodeString(tagStr)
			if err != nil {
				return data.NewBoolValue(false), nil
			}
		}
		combined := make([]byte, 0, len(ciphertext)+len(tag))
		combined = append(combined, ciphertext...)
		combined = append(combined, tag...)
		plain, err := gcm.Open(nil, iv, combined, nil)
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		return data.NewStringValue(string(plain)), nil
	default:
		return data.NewBoolValue(false), nil
	}
}

func (f *OpenSSLDecryptFunction) GetName() string { return "openssl_decrypt" }
var openSSLDecryptFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "data", 0, nil, nil),
	node.NewParameter(nil, "cipher_algo", 1, nil, nil),
	node.NewParameter(nil, "passphrase", 2, nil, nil),
	node.NewParameter(nil, "options", 3, node.NewIntLiteral(nil, "0"), nil),
	node.NewParameter(nil, "iv", 4, node.NewStringLiteral(nil, ""), nil),
	node.NewParameter(nil, "tag", 5, node.NewStringLiteral(nil, ""), nil),
}

func (f *OpenSSLDecryptFunction) GetParams() []data.GetValue {
	return openSSLDecryptFunctionGetParams
}
var openSSLDecryptFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "data", 0, nil),
	node.NewVariable(nil, "cipher_algo", 1, nil),
	node.NewVariable(nil, "passphrase", 2, nil),
	node.NewVariable(nil, "options", 3, nil),
	node.NewVariable(nil, "iv", 4, nil),
	node.NewVariable(nil, "tag", 5, nil),
}

func (f *OpenSSLDecryptFunction) GetVariables() []data.Variable {
	return openSSLDecryptFunctionGetVariables
}

func asIntVal(v data.Value) (int, error) {
	if asInt, ok := v.(data.AsInt); ok {
		return asInt.AsInt()
	}
	if str, ok := v.(data.AsString); ok {
		return strconv.Atoi(str.AsString())
	}
	return 0, fmt.Errorf("not an int")
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
