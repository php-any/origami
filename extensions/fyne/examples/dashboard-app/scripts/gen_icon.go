//go:build ignore

package main

import (
	"encoding/binary"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

func main() {
	size := 256
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Background: rounded rectangle (dark purple)
	bgColor := color.RGBA{R: 88, G: 56, B: 156, A: 255}
	radius := size / 4
	drawRoundedRect(img, image.Rect(0, 0, size, size), radius, bgColor)

	// Music note color: white
	noteColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Draw music note in center
	cx, cy := size/2, size/2
	noteW, noteH := size/5, size/4

	// Note head (filled ellipse)
	headX := cx - noteW/2
	headY := cy + noteH/4
	drawFilledEllipse(img, headX, headY, noteW, noteH/2, noteColor)

	// Note stem (vertical bar, right side)
	stemX := cx + noteW/2 - noteW/10
	stemTop := cy - noteH
	stemBottom := cy + noteH/4
	for y := stemTop; y <= stemBottom; y++ {
		for x := stemX; x < stemX+noteW/10; x++ {
			if x >= 0 && x < size && y >= 0 && y < size {
				img.Set(x, y, noteColor)
			}
		}
	}

	// Note flag (curved line at top)
	flagStartX := stemX + noteW/10
	for i := 0; i < noteH/3; i++ {
		yy := stemTop + i
		xx := flagStartX + int(float64(i)*2.0)
		for dx := 0; dx < 4; dx++ {
			x := xx + dx
			if x >= 0 && x < size && yy >= 0 && yy < size {
				img.Set(x, yy, noteColor)
			}
		}
	}

	// Save PNG
	pngF, err := os.Create("Icon.png")
	if err != nil {
		panic(err)
	}
	png.Encode(pngF, img)
	pngF.Close()

	// Generate ICO (wrap PNG in ICO format)
	icoF, err := os.Create("Icon.ico")
	if err != nil {
		panic(err)
	}
	// Re-encode PNG to get the data
	pngData := encodePNG(img)
	writeICO(icoF, pngData, size)
	icoF.Close()

	println("Icon.png + Icon.ico generated")
}

func drawRoundedRect(img *image.RGBA, rect image.Rectangle, radius int, c color.RGBA) {
	w, h := rect.Dx(), rect.Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Check if in rounded corner
			inCorner := false
			if x < radius && y < radius {
				dx, dy := float64(radius-x), float64(radius-y)
				inCorner = math.Sqrt(dx*dx+dy*dy) > float64(radius)
			} else if x >= w-radius && y < radius {
				dx, dy := float64(x-(w-radius)), float64(radius-y)
				inCorner = math.Sqrt(dx*dx+dy*dy) > float64(radius)
			} else if x < radius && y >= h-radius {
				dx, dy := float64(radius-x), float64(y-(h-radius))
				inCorner = math.Sqrt(dx*dx+dy*dy) > float64(radius)
			} else if x >= w-radius && y >= h-radius {
				dx, dy := float64(x-(w-radius)), float64(y-(h-radius))
				inCorner = math.Sqrt(dx*dx+dy*dy) > float64(radius)
			}
			if !inCorner {
				img.Set(x, y, c)
			}
		}
	}
}

func drawFilledEllipse(img *image.RGBA, cx, cy, w, h int, c color.RGBA) {
	for y := cy - h; y <= cy+h; y++ {
		for x := cx - w; x <= cx+w; x++ {
			if x < 0 || x >= img.Bounds().Dx() || y < 0 || y >= img.Bounds().Dy() {
				continue
			}
			dx := float64(x-cx) / float64(w)
			dy := float64(y-cy) / float64(h)
			if dx*dx+dy*dy <= 1.0 {
				img.Set(x, y, c)
			}
		}
	}
}

func encodePNG(img image.Image) []byte {
	// Use a buffer approach
	tmp, _ := os.CreateTemp("", "icon-*.png")
	tmpName := tmp.Name()
	tmp.Close()
	pngF, _ := os.Create(tmpName)
	png.Encode(pngF, img)
	pngF.Close()
	data, _ := os.ReadFile(tmpName)
	os.Remove(tmpName)
	return data
}

func writeICO(f *os.File, pngData []byte, size int) {
	// ICO header
	binary.Write(f, binary.LittleEndian, uint16(0)) // reserved
	binary.Write(f, binary.LittleEndian, uint16(1)) // ICO type
	binary.Write(f, binary.LittleEndian, uint16(1)) // 1 image

	// ICO directory entry
	binary.Write(f, binary.LittleEndian, uint8(size))          // width
	binary.Write(f, binary.LittleEndian, uint8(size))          // height
	binary.Write(f, binary.LittleEndian, uint8(0))             // palette
	binary.Write(f, binary.LittleEndian, uint8(0))             // reserved
	binary.Write(f, binary.LittleEndian, uint16(1))            // color planes
	binary.Write(f, binary.LittleEndian, uint16(32))           // bits per pixel
	binary.Write(f, binary.LittleEndian, uint32(len(pngData))) // image size
	binary.Write(f, binary.LittleEndian, uint32(22))           // image offset (6+16)

	// PNG image data
	f.Write(pngData)
}

// Ensure draw is imported (used by image.NewRGBA)
var _ = draw.Draw
