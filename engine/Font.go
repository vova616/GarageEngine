package engine

import (
	"github.com/vova616/freetype-go/freetype"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	//"github.com/go-gl/glfw"
	//"gl/glu"
	//"log"
	"fmt"
	"math"
	//"bufio"
	//"image/png"
	//"os"
)

type Font struct {
	*Texture
	lettersArray map[rune]*LetterInfo
	fontSize     float64
	dpi          int
	sdf          bool
}

type LetterInfo struct {
	Rect           image.Rectangle
	YOffset        float32
	XOffset        float32
	XAdvance       float32
	RelativeWidth  float32
	RelativeHeight float32
}

func (t *Font) Size() float64 {
	return t.fontSize
}

func (t *Font) IsSDF() bool {
	return t.sdf
}

func (t *Font) Index(runei ID) image.Rectangle {
	letter, ok := runei.(rune)
	if !ok {
		panic("runei is not rune")
	}
	l, exists := t.lettersArray[letter]
	if exists {
		return l.Rect
	}
	return image.Rectangle{}
}

func (t *Font) Group(id ID) []image.Rectangle {
	panic("font does not have groups")
}

func (t *Font) CheckText(s string) {

}

func (t *Font) LetterInfo(letter rune) *LetterInfo {
	l, exists := t.lettersArray[letter]
	if exists {
		return l
	}
	return nil
}

func NewFont(fontPath string, size float64) (*Font, error) {
	return NewFont2(fontPath, size, 72, false, 0, 255)
}

func NewSDFFont(fontPath string, size float64) (*Font, error) {
	return NewSDFFont2(fontPath, size, 512, 32)
}

func NewSDFFont2(fontPath string, size float64, sdfSize float64, scanRange int) (*Font, error) {
	return NewSDFFont3(fontPath, size, 72, false, 0, 255, sdfSize, scanRange)
}

func NewSDFFont3(fontPath string, size float64, dpi int, readonly bool, firstRune, lastRune rune, sdfSize float64, scanRange int) (*Font, error) {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	fontBytes = nil

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(sdfSize)
	c.SetSrc(image.White)

	ratio := sdfSize / size

	text := ""
	for i := firstRune; i < lastRune+1; i++ {
		text += string(i)
	}

	rects := make([]image.Rectangle, 0)
	LetterArray := make(map[rune]*LetterInfo)

	for _, r := range text {
		index := font.Index(r)
		mask, offset, err := c.Glyph(index, freetype.Pt(0, 0))
		if err != nil {
			fmt.Println("Rune generation error:", err)
			continue
		}
		bd := mask.Bounds()

		AdvanceWidth := c.FUnitToFix32(int(font.HMetric(index).AdvanceWidth)).Float()
		LeftSideBearing := c.FUnitToFix32(int(font.HMetric(index).LeftSideBearing)).Float()
		AdvanceWidth = AdvanceWidth / float32(sdfSize)
		LeftSideBearing = LeftSideBearing / float32(sdfSize)
		YOffset := (float32(-offset.Y) - float32(mask.Bounds().Max.Y)) / float32(sdfSize)
		relativeWidth := float32(bd.Dx()) / float32(sdfSize)
		relativeHeight := float32(bd.Dy()) / float32(sdfSize)

		sdfBounds := mask.Bounds()
		sdfBounds.Max.X = int(float64(sdfBounds.Max.X)/ratio) + 2
		sdfBounds.Max.Y = int(float64(sdfBounds.Max.Y)/ratio) + 2

		rects = append(rects, sdfBounds)

		LetterArray[r] = &LetterInfo{sdfBounds, YOffset, LeftSideBearing, AdvanceWidth, relativeWidth, relativeHeight}
	}

	ay, ax, e := FindOptimalSize(10, rects...)
	if e != nil {
		return nil, e
	}

	dst := image.NewRGBA(image.Rect(0, 0, int(ax), int(ay)))
	node := NewBin(int(ax), int(ay), Padding)

	rects, e = node.InsertArray(rects)
	if e != nil {
		return nil, e
	}

	rectIndex := 0
	for _, r := range text {
		index := font.Index(r)
		mask, _, err := c.Glyph(index, freetype.Pt(0, 0))
		if err != nil {
			fmt.Println("Rune generation error:", err)
			continue
		}

		rect := rects[rectIndex]
		newMask := image.NewAlpha(rect)

		//Note: this is slow we need to find better algorithm
		for xx := 0; xx < newMask.Bounds().Dx(); xx++ {
			for yy := 0; yy < newMask.Bounds().Dy(); yy++ {
				alpha := FindSDFAlpha(mask, int(float64(xx)*ratio), int(float64(yy)*ratio), scanRange)
				newMask.SetAlpha(xx, yy, color.Alpha{uint8(alpha)})
			}
		}

		draw.Draw(dst, rect, newMask, image.ZP, draw.Src)
		LetterArray[r].Rect = rect
		rectIndex++
	}

	texture, err := NewTexture(dst, dst.Pix)
	if err != nil {
		return nil, err
	}

	if readonly {
		texture.SetReadOnly()
	}

	texture.SetFiltering(Linear, Linear)

	return &Font{texture, LetterArray, size, dpi, true}, nil
}

func NextPowerOfTwo(x uint64) uint64 {
	power := uint64(1)
	for power < x {
		power *= 2
	}
	return power
}

func NewFont2(fontPath string, size float64, dpi int, readonly bool, firstRune, lastRune rune) (*Font, error) {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	fontBytes = nil

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(size)
	c.SetSrc(image.White)

	text := ""
	for i := firstRune; i < lastRune+1; i++ {
		text += string(i)
	}

	rects := make([]image.Rectangle, 0)
	LetterArray := make(map[rune]*LetterInfo)

	for _, r := range text {
		index := font.Index(r)
		mask, offset, err := c.Glyph(index, freetype.Pt(0, 0))
		if err != nil {
			fmt.Println("Rune generation error:", err)
			continue
		}
		bd := mask.Bounds()

		AdvanceWidth := c.FUnitToFix32(int(font.HMetric(index).AdvanceWidth)).Float()
		LeftSideBearing := c.FUnitToFix32(int(font.HMetric(index).LeftSideBearing)).Float()
		AdvanceWidth = AdvanceWidth / float32(size)
		LeftSideBearing = LeftSideBearing / float32(size)
		YOffset := (float32(-offset.Y) - float32(mask.Bounds().Max.Y)) / float32(size)
		relativeWidth := float32(bd.Dx()) / float32(size)
		relativeHeight := float32(bd.Dy()) / float32(size)

		rects = append(rects, mask.Bounds())

		LetterArray[r] = &LetterInfo{bd, YOffset, LeftSideBearing, AdvanceWidth, relativeWidth, relativeHeight}
	}

	ay, ax, e := FindOptimalSize(10, rects...)
	if e != nil {
		return nil, e
	}

	dst := image.NewRGBA(image.Rect(0, 0, int(ax), int(ay)))
	node := NewBin(int(ax), int(ay), Padding)
	rects, e = node.InsertArray(rects)
	if e != nil {
		return nil, e
	}

	rectIndex := 0
	for _, r := range text {
		index := font.Index(r)
		mask, _, err := c.Glyph(index, freetype.Pt(0, 0))
		if err != nil {
			fmt.Println("Rune generation error:", err)
			continue
		}

		rect := rects[rectIndex]

		draw.Draw(dst, rect, mask, image.ZP, draw.Src)
		LetterArray[r].Rect = rect
		rectIndex++
	}

	texture, err := NewTexture(dst, dst.Pix)
	if err != nil {
		return nil, err
	}

	if readonly {
		texture.SetReadOnly()
	}

	texture.SetFiltering(Linear, Linear)

	return &Font{texture, LetterArray, size, dpi, false}, nil

}

func FindSDF(img image.Image, x, y, maxRadius int) int {

	c := img.At(x, y)
	if c == nil {
		panic("this point does not exists.")
	}
	_, _, _, alpha := c.RGBA()

	distance := maxRadius*maxRadius + 1
	for radius := 1; (radius <= maxRadius) && (radius*radius < distance); radius++ {
		for line := -radius; line <= radius; line++ {
			nx, ny := x+line, y+radius
			if (image.Point{nx, ny}.In(img.Bounds())) {
				c = img.At(nx, ny)
				_, _, _, a := c.RGBA()
				//fmt.Println(line, x, ny, a, alpha)
				if a != alpha {
					nx = x - nx
					ny = y - ny
					d := (nx * nx) + (ny * ny)
					if d < distance {
						distance = d
					}
				}
			}
		}

		for line := -radius; line <= radius; line++ {
			nx, ny := x+line, y-radius
			if (image.Point{nx, ny}.In(img.Bounds())) {
				c = img.At(nx, ny)
				_, _, _, a := c.RGBA()
				if a != alpha {
					nx = x - nx
					ny = y - ny
					d := (nx * nx) + (ny * ny)
					if d < distance {
						distance = d
					}
				}
			}
		}

		for line := -radius; line <= radius; line++ {
			nx, ny := x+radius, y+line
			if (image.Point{nx, ny}.In(img.Bounds())) {
				c = img.At(nx, ny)
				_, _, _, a := c.RGBA()
				if a != alpha {
					nx = x - nx
					ny = y - ny
					d := (nx * nx) + (ny * ny)
					if d < distance {
						distance = d
					}
				}
			}
		}

		for line := -radius; line <= radius; line++ {
			nx, ny := x-radius, y+line
			if (image.Point{nx, ny}.In(img.Bounds())) {
				c = img.At(nx, ny)
				_, _, _, a := c.RGBA()
				if a != alpha {
					nx = x - nx
					ny = y - ny
					d := (nx * nx) + (ny * ny)
					if d < distance {
						distance = d
					}
				}
			}
		}
	}
	SDF := float32(math.Sqrt(float64(distance)))
	if alpha == 0 {
		SDF = -SDF
	}
	SDF *= 127.5 / float32(maxRadius)
	SDF += 127.5
	if SDF < 0 {
		SDF = 0
	} else if SDF > 255 {
		SDF = 255
	}
	return int(SDF + 0.5)
}

func FindSDFAlpha(img *image.Alpha, x, y, maxRadius int) int {

	w := img.Bounds().Dx()
	distance := maxRadius*maxRadius + 1
	alpha := uint8(0)
	if (image.Point{x, y}.In(img.Bounds())) {
		alpha = img.Pix[y*w+x]
	}

	a := uint8(0)
	for radius := 1; (radius <= maxRadius) && (radius*radius < distance); radius++ {
		for line := -radius; line <= radius; line++ {
			nx, ny := x+line, y+radius
			if (image.Point{nx, ny}.In(img.Bounds())) {
				a = img.Pix[ny*w+nx]
			} else {
				a = 0
			}
			//fmt.Println(line, x, ny, a, alpha)
			if a != alpha {
				nx = x - nx
				ny = y - ny
				d := (nx * nx) + (ny * ny)
				if d < distance {
					distance = d
				}
			}
			//}
		}

		for line := -radius; line <= radius; line++ {
			nx, ny := x+line, y-radius
			if (image.Point{nx, ny}.In(img.Bounds())) {
				a = img.Pix[ny*w+nx]
			} else {
				a = 0
			}
			if a != alpha {
				nx = x - nx
				ny = y - ny
				d := (nx * nx) + (ny * ny)
				if d < distance {
					distance = d
				}
			}
			//}
		}

		for line := -radius; line < radius; line++ {
			nx, ny := x+radius, y+line
			if (image.Point{nx, ny}.In(img.Bounds())) {
				a = img.Pix[ny*w+nx]
			} else {
				a = 0
			}
			if a != alpha {
				nx = x - nx
				ny = y - ny
				d := (nx * nx) + (ny * ny)
				if d < distance {
					distance = d
				}
			}
			//}
		}

		for line := -radius; line < radius; line++ {
			nx, ny := x-radius, y+line
			if (image.Point{nx, ny}.In(img.Bounds())) {
				a = img.Pix[ny*w+nx]
			} else {
				a = 0
			}
			if a != alpha {
				nx = x - nx
				ny = y - ny
				d := (nx * nx) + (ny * ny)
				if d < distance {
					distance = d
				}
			}
			//}
		}
	}
	SDF := float32(math.Sqrt(float64(distance)))
	if alpha == 0 {
		SDF = -SDF
	}
	SDF *= 127.5 / float32(maxRadius)
	SDF += 127.5
	if SDF < 0 {
		SDF = 0
	} else if SDF > 255 {
		SDF = 255
	}
	return int(SDF + 0.5)
}
