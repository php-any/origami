package fyne

import (
	"image/color"

	"github.com/php-any/origami/data"
)

// colorToGo 将 Fyne\Color 的 ClassValue 转换为 Go 的 color.Color
func colorToGo(v data.Value) color.Color {
	if cv, ok := v.(*data.ClassValue); ok {
		r, _ := cv.GetProperty("r")
		g, _ := cv.GetProperty("g")
		b, _ := cv.GetProperty("b")
		a, _ := cv.GetProperty("a")
		ri, _ := r.(data.AsInt).AsInt()
		gi, _ := g.(data.AsInt).AsInt()
		bi, _ := b.(data.AsInt).AsInt()
		ai, _ := a.(data.AsInt).AsInt()
		return color.NRGBA{R: uint8(ri), G: uint8(gi), B: uint8(bi), A: uint8(ai)}
	}
	return color.Black
}

// fyneColorToGo 将 data.Value 转换为 color.Color，支持字符串或 Color 对象
func fyneColorToGo(v data.Value) color.Color {
	if v == nil {
		return color.Black
	}
	switch c := v.(type) {
	case *data.ClassValue:
		return colorToGo(c)
	case *data.StringValue:
		return parseColorString(c.AsString())
	default:
		return color.Black
	}
}

// parseColorString 解析颜色字符串，支持 "#RRGGBB" 和 "#RRGGBBAA" 格式
func parseColorString(s string) color.Color {
	if len(s) == 0 {
		return color.Black
	}
	if s[0] == '#' {
		s = s[1:]
	}
	var r, g, b, a uint8 = 0, 0, 0, 255
	switch len(s) {
	case 6:
		r = hexToByte(s[0:2])
		g = hexToByte(s[2:4])
		b = hexToByte(s[4:6])
	case 8:
		r = hexToByte(s[0:2])
		g = hexToByte(s[2:4])
		b = hexToByte(s[4:6])
		a = hexToByte(s[6:8])
	}
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

func hexToByte(s string) uint8 {
	var v byte
	for _, c := range s {
		v *= 16
		if c >= '0' && c <= '9' {
			v += byte(c - '0')
		} else if c >= 'a' && c <= 'f' {
			v += byte(c-'a') + 10
		} else if c >= 'A' && c <= 'F' {
			v += byte(c-'A') + 10
		}
	}
	return v
}
