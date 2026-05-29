package core

import (
	"encoding/base64"
	"os"
	"sync"
)

var (
	phptInputMu  sync.RWMutex
	phptInputBody string
)

// SetPhptInputBody 设置当前请求的 php://input 内容（由 PHPT 运行器注入）。
func SetPhptInputBody(body string) {
	phptInputMu.Lock()
	phptInputBody = body
	phptInputMu.Unlock()
}

// InitPhptInputFromEnv 从 ORIGAMI_PHPT_INPUT（base64）加载请求体。
func InitPhptInputFromEnv() {
	raw := os.Getenv("ORIGAMI_PHPT_INPUT")
	if raw == "" {
		return
	}
	if decoded, err := base64.StdEncoding.DecodeString(raw); err == nil {
		SetPhptInputBody(string(decoded))
		return
	}
	SetPhptInputBody(raw)
}

// PhptInputBody 返回 php://input 可读内容。
func PhptInputBody() string {
	phptInputMu.RLock()
	defer phptInputMu.RUnlock()
	return phptInputBody
}
