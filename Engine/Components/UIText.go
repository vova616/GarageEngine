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
	"github.com/vova616/GarageEngine/Engine/Input"
	"github.com/vova616/chipmunk/vect"
)

type UIText struct {
	Engine.BaseComponent
	Font           *Engine.Font
	text           string
	buffer         gl.Buffer
	vertexCount    int
	texcoordsIndex int

	tabSize int

	width  float32
	height float32
	align  Engine.AlignType

	hover      bool
	focused    bool
	writeable  bool
	updateText bool

	autoFocus bool //This will go away

	Color Engine.Vector

	onHoverCallback func(bool)
}

func NewUIText(font *Engine.Font, text string) *UIText {
	if font == nil {
		return nil
	}

	uitext := &UIText{BaseComponent: Engine.NewComponent(),
		Font:      font,
		text:      text,
		buffer:    gl.GenBuffer(),
		align:     Engine.AlignCenter,
		writeable: false,
		tabSize:   4,
		Color:     Engine.Vector{1, 1, 1}}

	uitext.SetString(text)
	Input.AddCharCallback(func(rn rune) { uitext.charCallback(rn) })
	return uitext
}

func (ui *UIText) OnComponentBind(binded *Engine.GameObject) {
	ph := binded.AddComponent(Engine.NewPhysics(false, 1, 1)).(*Engine.Physics)
	_ = ph
	ph.Body.IgnoreGravity = true
	ph.Shape.IsSensor = true
}

func (ui *UIText) SetHoverCallBack(callback func(bool)) {
	ui.onHoverCallback = callback
}

func (ui *UIText) Width() float32 {
	return ui.width
}

func (ui *UIText) Height() float32 {
	return ui.height
}

func (ui *UIText) SetText(text string) {
	ui.text = text
	ui.updateText = true
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
	w, h := ui.GetPixelSize(text)

	index := 0
	for _, rune := range text {
		spaceMult := float32(1)
		if rune == '\t' {
			rune = ' '
			spaceMult = float32(index % ui.tabSize)
		}
		atlasImage := ui.Font.LetterInfo(rune)

		if atlasImage == nil {
			continue
		}

		yratio := atlasImage.PlaneHeight
		xratio := atlasImage.PlaneWidth

		//ygrid := -0.5 + (atlasImage.YGrid)
		//xgrid := (-0.5 + (atlasImage.XGrid)) + space
		ygrid := -0.5 + atlasImage.YGrid
		xgrid := -(w / 2) + (atlasImage.XGrid) + space
		space += atlasImage.RealWidth * spaceMult

		vertexCount += 4

		uv := Engine.IndexUV(ui.Font, rune)

		data[(index*12)+0] = xgrid
		data[(index*12)+1] = ygrid
		data[(index*12)+2] = 1
		data[(index*12)+3] = (xratio) + xgrid
		data[(index*12)+4] = ygrid
		data[(index*12)+5] = 1
		data[(index*12)+6] = (xratio) + xgrid
		data[(index*12)+7] = (yratio) + ygrid
		data[(index*12)+8] = 1
		data[(index*12)+9] = xgrid
		data[(index*12)+10] = (yratio) + ygrid
		data[(index*12)+11] = 1

		data[lt+(index*8)+0] = uv.U1
		data[lt+(index*8)+1] = uv.V2
		data[lt+(index*8)+2] = uv.U2
		data[lt+(index*8)+3] = uv.V2

		data[lt+(index*8)+4] = uv.U2
		data[lt+(index*8)+5] = uv.V1
		data[lt+(index*8)+6] = uv.U1
		data[lt+(index*8)+7] = uv.V1
		index++
	}

	ui.width = w
	ui.height = h

	ui.vertexCount = vertexCount
	ui.texcoordsIndex = texcoordsIndex
	ui.buffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, data, gl.STATIC_DRAW)
}

func (ui *UIText) GetPixelSize(text string) (width float32, height float32) {
	index := 0
	for _, rune := range text {
		spaceMult := float32(1)
		if rune == '\t' {
			rune = ' '
			spaceMult = float32(index % ui.tabSize)
		}
		atlasImage := ui.Font.LetterInfo(rune)

		if atlasImage == nil {
			continue
		}

		yratio := atlasImage.PlaneHeight
		width += atlasImage.RealWidth * spaceMult

		if yratio+atlasImage.YGrid > height {
			height = (yratio) + atlasImage.YGrid
		}
	}
	return
}

func (ui *UIText) String() string {
	return ui.text
}

func (ui *UIText) Start() {

}

func (ui *UIText) SetFocus(b bool) {
	ui.focused = b
}
func (ui *UIText) SetWritable(b bool) {
	ui.writeable = b
}

func (ui *UIText) charCallback(rn rune) {
	if ui.focused && ui.writeable {
		ui.text += string(rn)
		ui.updateText = true
	}
}

func (ui *UIText) Update() {
	ui.UpdateCollider()

	//Handle Tab & Backspace
	if ui.focused && ui.writeable {
		if len(ui.text) > 0 && Input.KeyPress(Input.KeyBackspace) {
			ui.updateText = true
			ui.text = ui.text[:len(ui.text)-1]
		}
		if Input.KeyPress(Input.KeyTab) {
			ui.updateText = true
			ui.text += "\t"
		}
	}
	if ui.updateText {
		ui.updateText = false
		ui.SetString(ui.text)
	}

	/*
		Just for tests
	*/
	if ui.autoFocus {
		mousePressed := Input.MousePress(Input.Mouse1)
		if mousePressed {
			if ui.hover && ui.writeable {
				ui.focused = true
			} else {
				ui.focused = false
			}
		}
	}
}

func (ui *UIText) OnCollisionEnter(arbiter *Engine.Arbiter) bool {

	return true
}

func (ui *UIText) OnCollisionExit(arbiter *Engine.Arbiter) {

}

func (ui *UIText) OnMouseEnter(arbiter *Engine.Arbiter) bool {
	ui.hover = true
	if ui.onHoverCallback != nil {
		ui.onHoverCallback(true)
	}
	return true
}

func (ui *UIText) OnMouseExit(arbiter *Engine.Arbiter) {
	ui.hover = false
	if ui.onHoverCallback != nil {
		ui.onHoverCallback(false)
	}
}

func (ui *UIText) UpdateCollider() {
	//if ui.GameObject().Physics.Body.Enabled {
	if ui.GameObject().Physics == nil {
		return
	}
	b := ui.GameObject().Physics.Box
	if b != nil {
		h := vect.Float(float64(ui.height) * float64(ui.GameObject().Transform().WorldScale().Y))
		w := vect.Float(float64(ui.width) * float64(ui.GameObject().Transform().WorldScale().X))
		update := false
		if h != b.Height || w != b.Width {
			b.Width = w
			b.Height = h
			update = true
		}

		c := Engine.Align(ui.align)
		center := vect.Vect{vect.Float(c.X), vect.Float(c.Y)}
		center.X = (center.X * w)
		center.Y = (center.Y * h)

		if b.Position.X != center.X || b.Position.Y != center.Y {
			update = true
			b.Position.X, b.Position.Y = center.X, center.Y
		}

		if update {
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
	v.X = (v.X * ui.width)
	v.Y = (v.Y * ui.height)

	Engine.TextureMaterial.Begin(ui.GameObject())

	vert := Engine.TextureMaterial.Verts
	uv := Engine.TextureMaterial.UV
	mp := Engine.TextureMaterial.ProjMatrix
	mv := Engine.TextureMaterial.ViewMatrix
	mm := Engine.TextureMaterial.ModelMatrix
	tx := Engine.TextureMaterial.Texture
	color := Engine.TextureMaterial.AddColor
	_ = color

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

	color.Uniform4f(ui.Color.X, ui.Color.Y, ui.Color.Z, 1)

	gl.DrawArrays(gl.QUADS, 0, ui.vertexCount)

	ui.Font.Unbind()
	vert.DisableArray()
	uv.DisableArray()

}
