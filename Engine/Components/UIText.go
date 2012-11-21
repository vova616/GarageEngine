package Components

import (
	"github.com/vova616/gl"
	//"image"
	//"github.com/jteeuwen/glfw"
	//"gl/glu"
	//"log"

	//"bufio"
	//"image/png"
	//"os" 
	//"strconv"
	//"github.com/jteeuwen/glfw"
	"github.com/vova616/GarageEngine/Engine"
	//. "github.com/vova616/GarageEngine/Engine/Input"
	"github.com/vova616/chipmunk/vect"
)

type UIText struct {
	Engine.BaseComponent
	Font           *Engine.Font
	text           string
	buffer         gl.Buffer
	vertexCount    int
	texcoordsIndex int

	width  float32
	height float32
	align  Engine.AlignType

	hover bool
	red   bool
}

func NewUIText(font *Engine.Font, text string) *UIText {
	if font == nil {
		return nil
	}

	uitext := &UIText{Engine.NewComponent(), font, text, gl.GenBuffer(), 0, 0, 0, 0, Engine.AlignCenter, false, false}
	uitext.SetString(text)
	return uitext
}

func (ui *UIText) OnComponentBind(binded *Engine.GameObject) {
	ph := ui.GameObject().AddComponent(Engine.NewPhysics(false, 1, 1)).(*Engine.Physics)
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

		uv := Engine.IndexUV(ui.Font, rune)

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
	//if MousePress(glfw.MouseLeft) {
	ui.red = ui.hover
	//}
}

func (ui *UIText) OnCollisionEnter(arbiter *Engine.Arbiter) bool {

	return true
}

func (ui *UIText) OnCollisionExit(arbiter *Engine.Arbiter) {

}

func (ui *UIText) OnMouseEnter(arbiter *Engine.Arbiter) bool {
	ui.hover = true
	return true
}

func (ui *UIText) OnMouseExit(arbiter *Engine.Arbiter) {
	ui.hover = false
}

func (ui *UIText) UpdateCollider() {
	//if ui.GameObject().Physics.Body.Enabled {
	b := ui.GameObject().Physics.Box
	if b != nil {
		h := float64(ui.height) * float64(ui.GameObject().Transform().WorldScale().Y)
		w := float64(ui.width) * float64(ui.GameObject().Transform().WorldScale().X)
		if vect.Float(h) != b.Height || vect.Float(w) != b.Width {
			b.Width = vect.Float(w)
			b.Height = vect.Float(h)
			b.UpdatePoly()
		}
	}
	//log.Println(b.Height, b.Width, ui.GameObject().Transform().Scale().X, ui.GameObject().Name())
	//}
}

func (ui *UIText) Align() Engine.AlignType {
	return ui.align
}

func (ui *UIText) SetAlign(align Engine.AlignType) {
	ui.align = align
}

func (ui *UIText) Draw() {
	if ui.text == "" {
		return
	}

	v := Engine.Align(ui.align)
	v.X *= ui.width
	v.Y *= ui.height

	Engine.TextureMaterial.Begin(ui.GameObject())

	vert := Engine.TextureMaterial.Verts
	uv := Engine.TextureMaterial.UV
	mp := Engine.TextureMaterial.ProjMatrix
	mv := Engine.TextureMaterial.ViewMatrix
	mm := Engine.TextureMaterial.ModelMatrix
	tx := Engine.TextureMaterial.Texture
	color := Engine.TextureMaterial.AddColor

	vert.EnableArray()
	uv.EnableArray()

	ui.buffer.Bind(gl.ARRAY_BUFFER)

	vert.AttribPointerPtr(3, gl.FLOAT, false, 0, 0)
	uv.AttribPointerPtr(2, gl.FLOAT, false, 0, ui.texcoordsIndex)

	camera := Engine.GetScene().SceneBase().Camera

	view := camera.Transform().Matrix()
	view = view.Invert()
	model := Engine.NewIdentity()
	model.Translate(v.X, v.Y, 0)
	model.Mul(ui.GameObject().Transform().Matrix())

	/*
		view := camera.Transform().Matrix()
		view = view.Invert()
		model := ui.GameObject().Transform().Matrix()
	*/

	mv.Uniform4fv([]float32(view[:]))
	mp.Uniform4fv([]float32(camera.Projection[:]))
	mm.Uniform4fv([]float32(model[:]))

	ui.Font.Bind()
	gl.ActiveTexture(gl.TEXTURE0)
	tx.Uniform1i(0)

	if ui.red {
		color.Uniform4f(1, 0, 0, 1)
	} else {
		color.Uniform4f(1, 1, 1, 1)
	}

	gl.DrawArrays(gl.QUADS, 0, ui.vertexCount)

	ui.Font.Unbind()
	vert.DisableArray()
	uv.DisableArray()

}
