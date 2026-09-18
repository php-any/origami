package captcha

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"math/rand/v2"
)

// RenderOptions 控制验证码图像外观。
type RenderOptions struct {
	Phrase string
	Width  int
	Height int
	Noise  int
	Rng    *rand.Rand
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodeJPEG(img image.Image, quality int) ([]byte, error) {
	quality = clampInt(quality, 1, 100)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func renderImage(opt RenderOptions) *image.NRGBA {
	w := clampInt(opt.Width, 40, 800)
	h := clampInt(opt.Height, 20, 400)
	noise := clampInt(opt.Noise, 0, 10)
	rng := opt.Rng
	if rng == nil {
		rng = rand.New(rand.NewPCG(1, 2))
	}

	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	bg := color.NRGBA{
		R: uint8(230 + rng.IntN(20)),
		G: uint8(230 + rng.IntN(20)),
		B: uint8(225 + rng.IntN(20)),
		A: 255,
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)

	if noise > 0 {
		dotN := noise * 18
		for i := 0; i < dotN; i++ {
			c := color.NRGBA{
				R: uint8(160 + rng.IntN(80)),
				G: uint8(160 + rng.IntN(80)),
				B: uint8(160 + rng.IntN(80)),
				A: 255,
			}
			img.SetNRGBA(rng.IntN(w), rng.IntN(h), c)
		}
	}

	runes := []rune(opt.Phrase)
	if len(runes) == 0 {
		return img
	}
	drawPhrase(img, runes, rng)

	if noise > 0 {
		lineN := 1 + noise/2
		for i := 0; i < lineN; i++ {
			c := color.NRGBA{
				R: uint8(80 + rng.IntN(120)),
				G: uint8(80 + rng.IntN(120)),
				B: uint8(80 + rng.IntN(120)),
				A: 180,
			}
			drawLine(img,
				rng.IntN(w), rng.IntN(h),
				rng.IntN(w), rng.IntN(h),
				c)
		}
		amp := 1.2 + float64(noise)*0.35
		img = waveDistort(img, amp, rng)
	}
	return img
}

func drawPhrase(img *image.NRGBA, runes []rune, rng *rand.Rand) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	n := len(runes)
	cell := w / n
	maxScaleW := (cell - 4) / 5
	maxScaleH := (h - 8) / 7
	scale := maxScaleW
	if maxScaleH < scale {
		scale = maxScaleH
	}
	if scale < 2 {
		scale = 2
	}

	for i, r := range runes {
		g, ok := lookupGlyph(r)
		if !ok {
			g, _ = lookupGlyph('?')
		}
		gw, gh := 5*scale, 7*scale
		pad := 1
		if scale >= 4 {
			pad = 1
		}
		tmp := image.NewNRGBA(image.Rect(0, 0, gw+pad*2, gh+pad*2))
		col := color.NRGBA{
			R: uint8(20 + rng.IntN(70)),
			G: uint8(20 + rng.IntN(60)),
			B: uint8(40 + rng.IntN(80)),
			A: 255,
		}
		paintGlyph(tmp, g, pad, pad, scale, col)
		angle := (rng.Float64() - 0.5) * 0.7
		cx := i*cell + cell/2 + rng.IntN(5) - 2
		cy := h/2 + rng.IntN(7) - 3
		blitRotated(img, tmp, cx, cy, angle)
	}
}

func paintGlyph(img *image.NRGBA, g glyph5x7, x, y, scale int, col color.NRGBA) {
	for row := 0; row < 7; row++ {
		bits := g[row]
		for colx := 0; colx < 5; colx++ {
			if bits&(1<<uint(4-colx)) == 0 {
				continue
			}
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.SetNRGBA(x+colx*scale+dx, y+row*scale+dy, col)
				}
			}
		}
	}
}

func blitRotated(dst *image.NRGBA, src *image.NRGBA, cx, cy int, angle float64) {
	sb := src.Bounds()
	sw, sh := sb.Dx(), sb.Dy()
	sinA, cosA := math.Sin(angle), math.Cos(angle)
	hw, hh := float64(sw)/2, float64(sh)/2
	bound := int(math.Ceil(math.Hypot(hw, hh))) + 1
	db := dst.Bounds()
	for y := -bound; y <= bound; y++ {
		for x := -bound; x <= bound; x++ {
			// 逆旋转取源像素
			sx := cosA*float64(x) + sinA*float64(y) + hw
			sy := -sinA*float64(x) + cosA*float64(y) + hh
			ix, iy := int(math.Floor(sx)), int(math.Floor(sy))
			if ix < 0 || iy < 0 || ix >= sw || iy >= sh {
				continue
			}
			c := src.NRGBAAt(ix, iy)
			if c.A == 0 {
				continue
			}
			dx, dy := cx+x, cy+y
			if dx < db.Min.X || dy < db.Min.Y || dx >= db.Max.X || dy >= db.Max.Y {
				continue
			}
			dst.SetNRGBA(dx, dy, c)
		}
	}
}

func drawLine(img *image.NRGBA, x0, y0, x1, y1 int, c color.NRGBA) {
	dx := absInt(x1 - x0)
	dy := -absInt(y1 - y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy
	b := img.Bounds()
	for {
		if x0 >= b.Min.X && y0 >= b.Min.Y && x0 < b.Max.X && y0 < b.Max.Y {
			img.SetNRGBA(x0, y0, c)
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func waveDistort(src *image.NRGBA, amp float64, rng *rand.Rand) *image.NRGBA {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	dst := image.NewNRGBA(b)
	fx := 0.04 + rng.Float64()*0.03
	fy := 0.05 + rng.Float64()*0.04
	px := rng.Float64() * 6.28
	py := rng.Float64() * 6.28
	ampX := amp
	ampY := amp * 0.6
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			sx := x + int(ampX*math.Sin(fy*float64(y)+px))
			sy := y + int(ampY*math.Sin(fx*float64(x)+py))
			if sx < 0 {
				sx = 0
			} else if sx >= w {
				sx = w - 1
			}
			if sy < 0 {
				sy = 0
			} else if sy >= h {
				sy = h - 1
			}
			dst.SetNRGBA(x, y, src.NRGBAAt(sx, sy))
		}
	}
	return dst
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func decodePNGBytes(bin []byte) (image.Image, error) {
	return png.Decode(bytes.NewReader(bin))
}

func newVisualRNG(seed [32]byte) *rand.Rand {
	s0 := binary.LittleEndian.Uint64(seed[0:8])
	s1 := binary.LittleEndian.Uint64(seed[8:16])
	if s0 == 0 && s1 == 0 {
		s0 = 1
	}
	return rand.New(rand.NewPCG(s0, s1))
}
