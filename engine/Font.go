package engine

import (
	"github.com/vova616/GarageEngine/freetype"
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
	return NewSDFFont2(fontPath, size, 16, 32)
}

func NewSDFFont2(fontPath string, size float64, scaler float64, scanRange int) (*Font, error) {
	return NewSDFFont3(fontPath, size, 72, false, 0, 255, scaler, scanRange)
}

func NewSDFFont3(fontPath string, size float64, dpi int, readonly bool, firstRune, lastRune rune, scaler float64, scanRange int) (*Font, error) {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	fontBytes = nil

	pRange := scanRange

	_, _ = scaler, pRange

	osize := size
	size = size * scaler

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(size)
	c.SetSrc(image.White)

	border := 2

	pt := freetype.Pt(0, border+int(size))

	x := pt.X
	mx := pt.X

	text := ""
	for i := firstRune; i < lastRune+1; i++ {
		text += string(i)
	}

	//Note: need to fix this
	for i, r := range text {
		if i%15 == 0 && i != 0 {
			pt.Y += c.PointToFix32(size + float64(border))
			if pt.X > mx {
				mx = pt.X
			}
			pt.X = x
		}
		mask, offset, _ := c.Glyph(font.Index(r), pt)
		if mask == nil {
			continue
		}

		bd := mask.Bounds().Add(offset)
		pt.X = c.PointToFix32(float64(bd.Max.X) + float64(border) + 8)
	}

	dst := image.NewRGBA(image.Rect(0, 0, ((int(mx/256))/int(scaler))+2+int(osize), (int(pt.Y/256)/int(scaler))+2+int(osize)))
	dstBounds := dst.Bounds()
	fmt.Println("atlas size", dstBounds)

	c.SetDst(dst)
	c.SetClip(dstBounds)

	LetterArray := make(map[rune]*LetterInfo)

	pt = freetype.Pt(0, border+int(osize))

	for i, r := range text {
		if i%15 == 0 && i != 0 {
			pt.Y += c.PointToFix32(osize + float64(border))
			pt.X = x
		}

		mask, offset, err := c.Glyph2(font.Index(r), pt)
		_ = offset
		if err != nil {
			fmt.Println("Rune generation error:", err)
			continue
		}

		newMask := image.NewAlpha(image.Rect(0, 0, int(2+float64(mask.Bounds().Dx())/(scaler)), int(2+float64(mask.Bounds().Dy())/(scaler))))

		//Note: this is slow we need to find better algorithm	
		for xx := 0; xx < newMask.Bounds().Dx(); xx++ {
			for yy := 0; yy < newMask.Bounds().Dy(); yy++ {
				alpha := FindSDFAlpha(mask, xx*int(scaler), yy*int(scaler), pRange)

				newMask.SetAlpha(xx, yy, color.Alpha{uint8(alpha)})

				//_, _, _, a := mask.At(xx*int(scaler), yy*int(scaler)).RGBA()

				//newMask.SetAlpha(xx, yy, color.Alpha{uint8(a / 257)})

			}
		}
		//panic("asd")	

		realoffy := -(float32(offset.Y) + float32(mask.Bounds().Max.Y)) / float32(size)
		planeW := float32(mask.Bounds().Dx()) / float32(size)
		planeH := float32(mask.Bounds().Dy()) / float32(size)
		mask = newMask
		offx := (float64(offset.X) / (scaler)) + float64(int(pt.X)/256)
		offy := (float64(offset.Y) / (scaler)) + float64(int(pt.Y)/256)

		bd := mask.Bounds().Add(image.Pt(int(offx+0.5), int(offy+0.5)))

		mp := image.Point{0, 0}
		draw.DrawMask(dst, bd, image.White, image.ZP, mask, mp, draw.Over)
		pt.X = c.PointToFix32(float64(bd.Max.X) + float64(border))

		adv := int(c.FUnitToFix32(int(font.HMetric(font.Index(r)).AdvanceWidth)))
		adv2 := int(c.FUnitToFix32(int(font.HMetric(font.Index(r)).LeftSideBearing)))
		LeftSideBearing := (float32(adv2/256) + float32(adv2%256/256)) / float32(size)
		realWidth := (float32(adv/256) + float32(adv%256/256)) / float32(size)

		LetterArray[r] = &LetterInfo{bd, realoffy, LeftSideBearing, realWidth, planeW, planeH}
	}

	texture, err := NewTexture(dst, dst.Pix)
	if err != nil {
		return nil, err
	}

	if readonly {
		texture.SetReadOnly()
	}

	texture.SetFiltering(Linear, Linear)

	return &Font{texture, LetterArray, osize, dpi, true}, nil

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

	border := 2

	pt := freetype.Pt(0, border+int(size))

	x := pt.X
	mx := pt.X

	text := ""
	for i := firstRune; i < lastRune+1; i++ {
		text += string(i)
	}

	for i, r := range text {
		if i%15 == 0 && i != 0 {
			pt.Y += c.PointToFix32(size + float64(border))
			if pt.X > mx {
				mx = pt.X
			}
			pt.X = x
		}
		mask, offset, _ := c.Glyph(font.Index(r), pt)
		if mask == nil {
			continue
		}
		bd := mask.Bounds().Add(offset)
		pt.X = c.PointToFix32(float64(bd.Max.X) + float64(border))
	}

	dst := image.NewRGBA(image.Rect(0, 0, int(mx/256)+2, int(pt.Y/256)+2+int(size)))
	dstBounds := dst.Bounds()

	c.SetDst(dst)
	c.SetClip(dstBounds)

	LetterArray := make(map[rune]*LetterInfo)

	pt = freetype.Pt(0, border+int(size))

	for i, r := range text {
		if i%15 == 0 && i != 0 {
			pt.Y += c.PointToFix32(size + float64(border))
			pt.X = x
		}

		mask, offset, err := c.Glyph(font.Index(r), pt)
		if err != nil {
			fmt.Println("Rune generation error:", err)
			continue
		}
		bd := mask.Bounds().Add(offset)

		mp := image.Point{0, 0}
		draw.DrawMask(dst, bd, image.White, image.ZP, mask, mp, draw.Over)
		pt.X = c.PointToFix32(float64(bd.Max.X) + float64(border))

		adv := c.FUnitToFix32(int(font.HMetric(font.Index(r)).AdvanceWidth))
		adv2 := c.FUnitToFix32(int(font.HMetric(font.Index(r)).LeftSideBearing))
		LeftSideBearing := (float32(adv2/256) + float32(adv2%256/256)) / float32(size)
		realWidth := (float32(adv/256) + float32(adv%256/256)) / float32(size)

		LetterArray[r] = &LetterInfo{bd, (float32(pt.Y/256) - float32(bd.Max.Y)) / float32(size), LeftSideBearing, realWidth, float32(bd.Dx()) / float32(size), float32(bd.Dy()) / float32(size)}
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
