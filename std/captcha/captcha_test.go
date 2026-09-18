package captcha

import (
	"bytes"
	"fmt"
	"image/png"
	"testing"
)

func TestMatchPhrase(t *testing.T) {
	if !matchPhrase("Ab3D", "ab3d", true) {
		t.Fatal("ignoreCase 应为 true")
	}
	if matchPhrase("Ab3D", "ab3d", false) {
		t.Fatal("大小写敏感时不应匹配")
	}
	if matchPhrase("ABCD", "XXXX", true) {
		t.Fatal("错误输入不应通过")
	}
	if matchPhrase("", "A", true) {
		t.Fatal("空短语不应通过")
	}
}

func TestRandomPhraseCharset(t *testing.T) {
	s, err := randomPhrase("A", 4)
	if err != nil {
		t.Fatal(err)
	}
	if s != "AAAA" {
		t.Fatalf("charset=A length=4 得到 %q", s)
	}
}

func TestRandomMath(t *testing.T) {
	display, phrase, err := randomMath()
	if err != nil {
		t.Fatal(err)
	}
	var a, b int
	var op byte
	if _, err := fmt.Sscanf(display, "%d%c%d", &a, &op, &b); err != nil {
		t.Fatalf("无法解析算式 %q: %v", display, err)
	}
	got := a + b
	if op == '-' {
		got = a - b
	}
	if phrase != fmt.Sprintf("%d", got) {
		t.Fatalf("算式 %s 答案应为 %d 实际 %s", display, got, phrase)
	}
}

func TestRenderPNG(t *testing.T) {
	img := renderImage(RenderOptions{
		Phrase: "TEST",
		Width:  160,
		Height: 50,
		Noise:  3,
	})
	bin, err := encodePNG(img)
	if err != nil {
		t.Fatal(err)
	}
	if len(bin) < 32 || !bytes.HasPrefix(bin, []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}) {
		t.Fatal("PNG 魔数不正确")
	}
	decoded, err := png.Decode(bytes.NewReader(bin))
	if err != nil {
		t.Fatal(err)
	}
	b := decoded.Bounds()
	if b.Dx() != 160 || b.Dy() != 50 {
		t.Fatalf("尺寸 %dx%d", b.Dx(), b.Dy())
	}
}
