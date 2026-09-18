package captcha

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	defaultCharset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	defaultLength  = 4
	defaultWidth   = 160
	defaultHeight  = 50
	defaultNoise   = 5
	defaultExpire  = 300
	defaultKey     = "captcha"
	modeText       = "text"
	modeMath       = "math"
)

func cryptoInt(max int) (int, error) {
	if max <= 0 {
		return 0, nil
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func randomPhrase(charset string, length int) (string, error) {
	runes := []rune(charset)
	if len(runes) == 0 {
		runes = []rune(defaultCharset)
	}
	length = clampInt(length, 1, 16)
	out := make([]rune, length)
	for i := 0; i < length; i++ {
		idx, err := cryptoInt(len(runes))
		if err != nil {
			return "", err
		}
		out[i] = runes[idx]
	}
	return string(out), nil
}

func randomMath() (display, phrase string, err error) {
	a, err := cryptoInt(9)
	if err != nil {
		return "", "", err
	}
	b, err := cryptoInt(9)
	if err != nil {
		return "", "", err
	}
	a++
	b++
	op, err := cryptoInt(2)
	if err != nil {
		return "", "", err
	}
	if op == 0 {
		return fmt.Sprintf("%d+%d", a, b), fmt.Sprintf("%d", a+b), nil
	}
	if a < b {
		a, b = b, a
	}
	return fmt.Sprintf("%d-%d", a, b), fmt.Sprintf("%d", a-b), nil
}

func matchPhrase(stored, input string, ignoreCase bool) bool {
	stored = strings.TrimSpace(stored)
	input = strings.TrimSpace(input)
	if stored == "" {
		return false
	}
	if ignoreCase {
		return strings.EqualFold(stored, input)
	}
	return stored == input
}
