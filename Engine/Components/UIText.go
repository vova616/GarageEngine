package Components

import (
	"github.com/banthar/gl"
	//"image"
	//"github.com/jteeuwen/glfw"
	//"gl/glu"
	//"log"

	//"bufio"
	//"image/png"
	//"os" 
	//"strconv"
	. "github.com/vova616/GarageEngine/Engine"
	. "github.com/vova616/GarageEngine/Engine/Input"
	"github.com/jteeuwen/glfw"
	. "github.com/vova616/chipmunk/vect"
)

type UIText struct {
	BaseComponent
	Font           *Font
	text           string
	buffer         gl.Buffer
	vertexCount    int
	texcoordsIndex int

	width  float32
	height float32
	align  AlignType

	hover bool
	red   bool
}

func NewUIText(font *Font, text string) *UIText {
	if font == nil {
		return nil
	}

	uitext := &UIText{NewComponent(), font, text, gl.GenBuffer(), 0, 0, 0, 0, AlignCenter, false, false}
	uitext.SetString(text)
	return uitext
}

func (ui *UIText) OnComponentBind(binded *GameObject) {
	h := (ui.height) * (ui.GameObject().Transform().WorldScale().Y)
	w := (ui.width) * (ui.GameObject().Transform().WorldScale().X)
	ph := binded.AddComponent(NewPhysics(true,w,h)).(*Physics)
	_ = ph 
	ph.Body.IgnoreGravity = true
	ph.Shape.IsSensor = true
}

func (ui *UIText) SetString(text string) {
	ui.text = text

	if text == "" {
		ui.buffer.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, 1, []byte{1}, gl.STATIC_DRAW)
		return
	}

	l := len(text)
	lt := l * 12
	_ = lt
	data := make([]float32, 20*l)

	vertexCount := 0
	texcoordsIndex := lt * 4
	space := float32(0)
	height := float32(0)
	for i, rune := range text {
		atlasImage := ui.Font.LetterInfo(rune)

		if atlasImage == nil {
			continue
		}

		yratio := atlasImage.PlaneHeight
		xratio := atlasImage.PlaneWidth

		//ygrid := -0.5 + (atlasImage.YGrid)
		//xgrid := (-0.5 + (atlasImage.XGrid)) + space
		ygrid := (0 + atlasImage.YGrid)
		xgrid := (0 + (atlasImage.XGrid)) + space
		space += atlasImage.RealWidth

		if yratio+ygrid > height {
			height = (yratio) + ygrid
		}

		vertexCount += 4

		uv := IndexUV(ui.Font, rune)

		data[(i*12)+0] = xgrid
		data[(i*12)+1] = ygrid
		data[(i*12)+2] = 1
		data[(i*12)+3] = (xratio) + xgrid
		data[(i*12)+4] = ygrid
		data[(i*12)+5] = 1
		data[(i*12)+6] = (xratio) + xgrid
		data[(i*12)+7] = (yratio) + ygrid
		data[(i*12)+8] = 1
		data[(i*12)+9] = xgrid
		data[(i*12)+10] = (yratio) + ygrid
		data[(i*12)+11] = 1

		data[lt+(i*8)+0] = uv.U1
		data[lt+(i*8)+1] = uv.V2
		data[lt+(i*8)+2] = uv.U2
		data[lt+(i*8)+3] = uv.V2

		data[lt+(i*8)+4] = uv.U2
		data[lt+(i*8)+5] = uv.V1
		data[lt+(i*8)+6] = uv.U1
		data[lt+(i*8)+7] = uv.V1
	}

	ui.width = space
	ui.height = height

	ui.vertexCount = vertexCount
	ui.texcoordsIndex = texcoordsIndex
	ui.buffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, data, gl.STATIC_DRAW)
}

func (ui *UIText) String() string {
	return ui.text
}

func (ui *UIText) Start() {

}

func (ui *UIText) Update() {
	ui.UpdateCollider()
	if MousePress(glfw.MouseLeft) {
		if ui.hover {
			ui.red = !ui.red
		}
	}
}

func (ui *UIText) OnCollisionEnter(collision Collision) {
	ui.red = !ui.red
}

func (ui *UIText) OnCollisionExit(collision Collision) {
	ui.red = !ui.red
}

func (ui *UIText) OnMouseHover() {
	ui.hover = true
}

func (ui *UIText) OnMouseExit() {
	ui.hover = false
}

func (ui *UIText) UpdateCollider() {
	//if ui.GameObject().Physics.Body.Enabled {
		b := ui.GameObject().Physics.Box
		h := float64(ui.height) * float64(ui.GameObject().Transform().WorldScale().Y)
		w := float64(ui.width) * float64(ui.GameObject().Transform().WorldScale().X)
		if Float(h) != b.Height || Float(w) != b.Width {
			b.Width = Float(w)
			b.Height = Float(h)
			b.UpdatePoly()
		}
		//log.Println(b.Height, b.Width, ui.GameObject().Transform().Scale().X, ui.GameObject().Name())
	//}
}

func (ui *UIText) Align() AlignType {
	return ui.align
}

func (ui *UIText) SetAlign(align AlignType) {
	ui.align = align
}

func (ui *UIText) Draw() {
	if ui.text == "" {
		return
	}
	ui.Font.Bind()
	gl.PushMatrix()

	v := Align(ui.align)
	v.X *= ui.width
	v.Y *= ui.height
	gl.Translatef(v.X, v.Y, v.Z)

	if ui.red {
		gl.Color3ub(255, 0, 0)
	} else {
		gl.Color3ub(255, 255, 255)
	}

	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)

	ui.buffer.Bind(gl.ARRAY_BUFFER)
	gl.VertexPointerVBO(3, gl.FLOAT, 0, 0)
	gl.TexCoordPointerVBO(2, gl.FLOAT, 0, (ui.texcoordsIndex))

	gl.DrawArrays(gl.QUADS, 0, ui.vertexCount)

	gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.DisableClientState(gl.VERTEX_ARRAY)

	gl.Color3ub(255, 255, 255)

	gl.PopMatrix()
	ui.Font.Unbind()
}
