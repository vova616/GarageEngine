package Engine

import (
	"github.com/vova616/GarageEngine/freetype"
	"image"
	"image/draw"
	"io/ioutil"
	//"github.com/jteeuwen/glfw"
	//"gl/glu"
	//"log"

	//"bufio"
	//"image/png"
	//"os" 
)

type Font struct {
	*Texture
	lettersArray []*LetterInfo
	firstRune    rune
	lastRune     rune
	fontSize     float64
	dpi          int
}

type LetterInfo struct {
	Rect        image.Rectangle
	YGrid       float32
	XGrid       float32
	RealWidth   float32
	PlaneWidth  float32
	PlaneHeight float32
}

func (t *Font) Size() float64 {
	return t.fontSize
}

func (t *Font) Index(runei interface{}) image.Rectangle {
	letter, ok := runei.(rune)
	if !ok {
		panic("runei is not rune")
	}
	if letter >= t.firstRune && letter <= t.lastRune {
		return t.lettersArray[letter+t.firstRune].Rect
	}
	return image.Rectangle{}
}

func (t *Font) Group(id interface{}) []image.Rectangle {
	panic("font does not have groups")
}

func (t *Font) LetterInfo(letter rune) *LetterInfo {
	if letter >= t.firstRune && letter <= t.lastRune {
		return t.lettersArray[letter+t.firstRune]
	}
	return nil
}

func NewFont(fontPath string, size float64) (*Font, error) {
	return NewFont2(fontPath, size, 72, false, 0, 255)
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
		bd := mask.Bounds().Add(offset)
		pt.X = c.PointToFix32(float64(bd.Max.X) + float64(border))
	}

	dst := image.NewRGBA(image.Rect(0, 0, int(mx/256)+2, int(pt.Y/256)+2+int(size)))
	dstBounds := dst.Bounds()

	c.SetDst(dst)
	c.SetClip(dstBounds)

	LetterArray := make([]*LetterInfo, len(text))

	pt = freetype.Pt(0, border+int(size))

	for i, r := range text {
		if i%15 == 0 && i != 0 {
			pt.Y += c.PointToFix32(size + float64(border))
			pt.X = x
		}

		mask, offset, _ := c.Glyph(font.Index(r), pt)
		bd := mask.Bounds().Add(offset)

		mp := image.Point{0, 0}
		draw.DrawMask(dst, bd, image.White, image.ZP, mask, mp, draw.Over)
		pt.X = c.PointToFix32(float64(bd.Max.X) + float64(border))

		adv := c.FUnitToFix32(int(font.HMetric(font.Index(r)).AdvanceWidth))
		adv2 := c.FUnitToFix32(int(font.HMetric(font.Index(r)).LeftSideBearing))
		LeftSideBearing := (float32(adv2/256) + float32(adv2%256/256)) / float32(size)
		realWidth := (float32(adv/256) + float32(adv%256/256)) / float32(size)

		LetterArray[int(r)-int(firstRune)] = &LetterInfo{bd, (float32(pt.Y/256) - float32(bd.Max.Y)) / float32(size), LeftSideBearing, realWidth, float32(bd.Dx()) / float32(size), float32(bd.Dy()) / float32(size)}
	}

	texture, err := NewTexture(dst, dst.Pix)
	if err != nil {
		return nil, err
	}

	if readonly {
		texture.SetReadOnly()
	}

	texture.SetFiltering(MinFilter, Linear)
	texture.SetFiltering(MagFilter, Linear)

	return &Font{texture, LetterArray, firstRune, lastRune, size, dpi}, nil

}
