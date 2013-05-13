package engine

import (
	"errors"
	"github.com/go-gl/gl"
	"image"
	"image/color"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	//"log"
	"os"
	"reflect"
	"unsafe"
)

type Filter int

type WrapType int
type Wrap int

const (
	Nearest              = Filter(gl.NEAREST)
	Linear               = Filter(gl.LINEAR)
	MipMapLinearNearest  = Filter(gl.LINEAR_MIPMAP_NEAREST)
	MipMapLinearLinear   = Filter(gl.LINEAR_MIPMAP_LINEAR)
	MipMapNearestLinear  = Filter(gl.NEAREST_MIPMAP_LINEAR)
	MipMapNearestNearest = Filter(gl.NEAREST_MIPMAP_NEAREST)

	WrapS = WrapType(gl.TEXTURE_WRAP_S)
	WrapR = WrapType(gl.TEXTURE_WRAP_R)
	WrapT = WrapType(gl.TEXTURE_WRAP_T)

	ClampToEdge    = Wrap(gl.CLAMP_TO_EDGE)
	ClampToBorder  = Wrap(gl.CLAMP_TO_BORDER)
	MirroredRepeat = Wrap(gl.MIRRORED_REPEAT)
	Repeat         = Wrap(gl.REPEAT)
	//gl.CLAMP
)

/*

*/

var (
	CustomColorModels = make(map[color.Model]*GLColorModel)
	lastBindedTexture gl.Texture
)

type Align byte

const (
	AlignLeft   = Align(1)
	AlignCenter = Align(2)
	AlignRight  = Align(4)

	AlignTopLeft   = Align(8 | AlignLeft)
	AlignTopCenter = Align(8 | AlignCenter)
	AlignTopRight  = Align(8 | AlignRight)

	AlignBottomLeft   = Align(16 | AlignLeft)
	AlignBottomCenter = Align(16 | AlignCenter)
	AlignBottomRight  = Align(16 | AlignRight)
)

func (typ Align) Vector() Vector {
	vect := NewVector2(0, 0)
	switch {
	case typ&AlignLeft != 0:
		vect.X = 0.5
	case typ&AlignCenter != 0:
		vect.X = 0
	case typ&AlignRight != 0:
		vect.X = -0.5
	}
	switch {
	case typ&8 != 0:
		vect.Y = 0.5
	case typ&16 != 0:
		vect.Y = -0.5
	}
	return vect
}

type EngineColorModel interface {
	color.Model
	Data() interface{}
}

type GLColorModel struct {
	InternalFormat int
	Type           gl.GLenum
	Format         gl.GLenum
	Target         gl.GLenum
	PixelBytesSize int
	Model          EngineColorModel
}

type GLTexture interface {
	GLTexture() gl.Texture
	Height() int
	Width() int
	PixelSize() int
	Bind()
}

type Texture struct {
	handle         gl.Texture
	readOnly       bool
	data           interface{}
	format         gl.GLenum
	typ            gl.GLenum
	internalFormat int
	target         gl.GLenum
	width          int
	height         int
}

func (t *Texture) GLTexture() gl.Texture {
	return t.handle
}

func (t *Texture) Height() int {
	return t.height
}

func (t *Texture) Width() int {
	return t.width
}

func LoadTexture(path string) (tex *Texture, err error) {
	img, e := LoadImage(path)
	if e != nil {
		return nil, e
	}
	return LoadTextureFromImage(img)
}

func LoadImage(path string) (img image.Image, err error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	img, _, e = image.Decode(f)
	if e != nil {
		return nil, e
	}
	return img, nil
}

func LoadGIF(path string) (imgs []image.Image, err error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	GIF, e := gif.DecodeAll(f)
	if e != nil {
		return nil, e
	}

	imgs = make([]image.Image, len(GIF.Image))
	for i, img := range GIF.Image {
		imgs[i] = img
	}

	return imgs, nil
}

func LoadImageQuiet(path string) (img image.Image) {
	panic("Deprecated")
}

func LoadTextureFromImage(image image.Image) (tex *Texture, err error) {
	/*
		val := reflect.ValueOf(image)
		pixs := val.FieldByName("Pix")
		if pixs.IsValid() {
			return
		} else {

		}
		return nil, false, nil
	*/
	internalFormat, typ, format, target, e := ColorModelToGLTypes(image.ColorModel())
	if e != nil {
		return nil, nil
	}

	//
	w := image.Bounds().Dx()
	h := image.Bounds().Dy()
	model := image.ColorModel()
	var data []byte

	switch model.(type) {
	case color.Palette:
		memHandle := Allocate(4 * h * w)
		data = memHandle.Bytes()
		defer memHandle.Release()
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				offset := (x + (y * w)) * 4
				r, g, b, a := image.At(x, y).RGBA()
				data[offset] = byte(r / 257)
				data[offset+1] = byte(g / 257)
				data[offset+2] = byte(b / 257)
				data[offset+3] = byte(a / 257)
			}
		}
	}

	switch model {
	case color.YCbCrModel:
		memHandle := Allocate(3 * h * w)
		data = memHandle.Bytes()
		defer memHandle.Release()
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				offset := (x + y*w) * 3
				r, g, b, _ := image.At(x, y).RGBA()
				data[offset] = byte(r / 257)
				data[offset+1] = byte(g / 257)
				data[offset+2] = byte(b / 257)
			}
		}
	case color.RGBAModel, color.NRGBAModel:
		memHandle := Allocate(4 * h * w)
		data = memHandle.Bytes()
		defer memHandle.Release()
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				offset := (x + (y * w)) * 4
				r, g, b, a := image.At(x, y).RGBA()
				data[offset] = byte(r / 257)
				data[offset+1] = byte(g / 257)
				data[offset+2] = byte(b / 257)
				data[offset+3] = byte(a / 257)
			}
		}
	case color.RGBA64Model, color.NRGBA64Model:
		memHandle := Allocate(4 * h * w)
		data = memHandle.Bytes()
		defer memHandle.Release()
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				offset := (x + y*w) * 4
				r, g, b, a := image.At(x, y).RGBA()
				data[offset] = byte(r / 257)
				data[offset+1] = byte(g / 257)
				data[offset+2] = byte(b / 257)
				data[offset+3] = byte(a / 257)
			}
		}
	default:
		m, e := CustomColorModels[model]
		if e {
			return NewTexture2(m.Model.Data(), image.Bounds().Dx(), image.Bounds().Dy(), target, internalFormat, typ, format), nil
		} else {
			return nil, errors.New("unsupported format")
		}
	}

	return NewTexture2(data, image.Bounds().Dx(), image.Bounds().Dy(), target, internalFormat, typ, format), nil
}

func ColorModelToGLTypes(model color.Model) (internalFormat int, typ gl.GLenum, format gl.GLenum, target gl.GLenum, err error) {

	switch model.(type) {
	case color.Palette:
		return gl.RGBA8, gl.RGBA, gl.UNSIGNED_BYTE, gl.TEXTURE_2D, nil
	}

	switch model {
	case color.RGBAModel, color.NRGBAModel:
		return gl.RGBA8, gl.RGBA, gl.UNSIGNED_BYTE, gl.TEXTURE_2D, nil
	case color.RGBA64Model, color.NRGBAModel:
		return gl.RGBA16, gl.RGBA, gl.UNSIGNED_SHORT, gl.TEXTURE_2D, nil
	case color.AlphaModel:
		return gl.ALPHA, gl.ALPHA, gl.UNSIGNED_BYTE, gl.TEXTURE_2D, nil
	case color.Alpha16Model:
		return gl.ALPHA16, gl.ALPHA, gl.UNSIGNED_SHORT, gl.TEXTURE_2D, nil
	case color.GrayModel:
		return gl.LUMINANCE, gl.LUMINANCE, gl.UNSIGNED_BYTE, gl.TEXTURE_2D, nil
	case color.Gray16Model:
		return gl.LUMINANCE16, gl.LUMINANCE, gl.UNSIGNED_SHORT, gl.TEXTURE_2D, nil
	case color.YCbCrModel:
		return gl.RGB8, gl.RGB, gl.UNSIGNED_BYTE, gl.TEXTURE_2D, nil
	default:
		m, e := CustomColorModels[model]
		if e {
			return m.InternalFormat, m.Type, m.Format, m.Target, nil
		}
		break
	}
	return 0, 0, 0, 0, errors.New("unsupported format")
}

func NewRGBTexture(rgbData interface{}, width int, height int) *Texture {
	return NewTexture2(rgbData, width, height, gl.TEXTURE_2D, gl.RGB8, gl.RGB, gl.UNSIGNED_BYTE)
}

func NewRGBATexture(rgbaData interface{}, width int, height int) *Texture {
	return NewTexture2(rgbaData, width, height, gl.TEXTURE_2D, gl.RGBA8, gl.RGBA, gl.UNSIGNED_BYTE)
}

func NewTexture(image image.Image, data interface{}) (texture *Texture, err error) {
	iF, typ, f, t, e := ColorModelToGLTypes(image.ColorModel())
	if e != nil {
		return nil, nil
	}
	return NewTexture2(data, image.Bounds().Dx(), image.Bounds().Dy(), t, iF, typ, f), nil
}

func NewTexture2(data interface{}, width int, height int, target gl.GLenum, internalFormat int, typ gl.GLenum, format gl.GLenum) *Texture {
	a := gl.GenTexture()
	a.Bind(target)
	gl.TexImage2D(target, 0, internalFormat, width, height, 0, typ, format, data)

	t := &Texture{a, false, data, format, typ, internalFormat, target, width, height}

	t.SetWraping(WrapS, ClampToEdge)
	t.SetWraping(WrapT, ClampToEdge)
	t.SetFiltering(Nearest, Nearest)

	//ansi := []float32{0}
	//gl.GetFloatv(gl.MAX_TEXTURE_MAX_ANISOTROPY_EXT, ansi)
	//gl.TexParameterf(target, gl.TEXTURE_MAX_ANISOTROPY_EXT, ansi[0])

	t.PreloadRender() //Forcing texture to go to VRAM and prevent shuttering
	t.data = nil
	data = nil

	ResourceManager.Add(t)

	return t
}

func NewTextureEmpty(width int, height int, model color.Model) *Texture {
	internalFormat, typ, format, target, e := ColorModelToGLTypes(model)
	if e != nil {
		return nil
	}
	a := gl.GenTexture()
	a.Bind(target)
	gl.TexImage2D(target, 0, internalFormat, width, height, 0, typ, format, nil)

	t := &Texture{a, false, nil, format, typ, internalFormat, target, width, height}

	t.SetWraping(WrapS, ClampToEdge)
	t.SetWraping(WrapT, ClampToEdge)
	t.SetFiltering(Nearest, Nearest)

	ResourceManager.Add(t)

	return t
}

func (t *Texture) Options(filter, clamp int) {
	t.Bind()
	gl.TexParameteri(t.target, gl.TEXTURE_MIN_FILTER, filter)
	gl.TexParameteri(t.target, gl.TEXTURE_MAG_FILTER, filter)
	gl.TexParameteri(t.target, gl.TEXTURE_WRAP_S, clamp)
	gl.TexParameteri(t.target, gl.TEXTURE_WRAP_T, clamp)
}

func (t *Texture) Param(filter, value int) {
	t.Bind()
	gl.TexParameteri(t.target, gl.GLenum(filter), value)
}

func (t *Texture) Paramf(filter int, value float32) {
	t.Bind()
	gl.TexParameterf(t.target, gl.GLenum(filter), value)
}

func (t *Texture) SetFiltering(minFilter Filter, magFilter Filter) {
	t.Bind()
	gl.TexParameteri(t.target, gl.TEXTURE_MAG_FILTER, int(magFilter))
	gl.TexParameteri(t.target, gl.TEXTURE_MIN_FILTER, int(minFilter))
}

func (t *Texture) SetWraping(wrapType WrapType, wrap Wrap) {
	t.Bind()
	gl.TexParameteri(t.target, gl.GLenum(wrapType), int(wrap))
}

func (t *Texture) BuildMipmaps() {
	t.Bind()
	gl.GenerateMipmap(t.target)
}

func (t *Texture) PixelSize() int {
	return 4
}

func (t *Texture) Image() image.Image {
	if !t.readOnly {
		return nil
	}
	return nil
}

func (t *Texture) ReadTextureFromGPU() []byte {
	t.Bind()
	b := gl.GenBuffer()
	b.Bind(gl.PIXEL_UNPACK_BUFFER)
	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, t.Width()*t.Height()*t.PixelSize(), 0, gl.STREAM_DRAW)
	//gl.GetTexImage(t.target, 0, t.format, buffer)
	b.Bind(gl.PIXEL_UNPACK_BUFFER)

	gl.TexSubImage2D(t.target, 0, 0, 0, t.Width(), t.Height(), t.format, t.typ, unsafe.Pointer(uintptr(0)))
	b.Bind(gl.PIXEL_UNPACK_BUFFER)

	l := t.Width() * t.Height() * t.PixelSize()

	gl.BufferData(gl.PIXEL_UNPACK_BUFFER, t.Width()*t.Height()*t.PixelSize(), 0, gl.STREAM_DRAW)
	ptr := gl.MapBuffer(gl.PIXEL_UNPACK_BUFFER, gl.WRITE_ONLY)

	var x []byte
	s := (*reflect.SliceHeader)(unsafe.Pointer(&x))
	s.Data = uintptr(ptr)
	s.Len = l
	s.Cap = l

	gl.UnmapBuffer(gl.PIXEL_UNPACK_BUFFER)

	return x
}

func (t *Texture) SetData(data interface{}) {
	panic("not implemented")
	//t.Bind()
	//gl.TexSubImage2D(t.target, 0, 0, 0, t.width, t.height, t.format, data)
}

func (t *Texture) SetReadOnly() {
	if t.readOnly {
		return
	}

	t.data = nil
	t.readOnly = true
}

func (t *Texture) Bind() {
	if t.handle != lastBindedTexture {
		t.handle.Bind(t.target)
		lastBindedTexture = t.handle
	}

}

func (t *Texture) Unbind() {
	t.handle.Unbind(t.target)
}

func (t *Texture) Render() {
	t.Bind()
	xratio := float32(t.width) / float32(t.height)
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-0.5, -0.5, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f((xratio)-0.5, -0.5, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f((xratio)-0.5, 0.5, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-0.5, 0.5, 1)
	gl.End()
}

func (t *Texture) PreloadRender() {
	t.Bind()
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(0, 0, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(0, 0, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(0, 0, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(0, 0, 1)
	gl.End()
}

func (t *Texture) Release() {
	t.data = nil
	if t.handle != 0 {
		t.handle.Delete()
		t.handle = 0
	}
}
