package components

import (
	"github.com/go-gl/gl"
	//"image"
	//"github.com/go-gl/glfw"
	//"gl/glu"
	//"log"

	//"bufio"
	//"image/png"
	//"os"
	//"strconv"
	//"github.com/go-gl/glfw"
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/input"
	"github.com/vova616/chipmunk/vect"
	//"runtime"
)

type UIText struct {
	engine.BaseComponent
	Font           *engine.Font
	text           string
	buffer         engine.VBO
	vertexCount    int
	texcoordsIndex int

	tabSize int

	width  float32
	height float32
	align  engine.Align

	focused    bool
	writeable  bool
	updateText bool

	autoFocus bool //This will go away

	Color engine.Color
}

func NewUIText(font *engine.Font, text string) *UIText {
	if font == nil {
		return nil
	}

	uitext := &UIText{BaseComponent: engine.NewComponent(),
		Font:      font,
		text:      text,
		buffer:    engine.GenBuffer(),
		align:     engine.AlignCenter,
		writeable: false,
		tabSize:   4,
		Color:     engine.Color_White}

	uitext.setString(text)
	input.AddCharCallback(func(rn rune) { uitext.charCallback(rn) })
	return uitext
}

func (ui *UIText) OnComponentAdd() {
	ui.GameObject().AddComponent(engine.NewPhysics(false))
	ph := ui.GameObject().Physics
	ph.Body.IgnoreGravity = true
	ph.Shape.IsSensor = true
}

func (ui *UIText) Width() float32 {
	return ui.width
}

func (ui *UIText) Height() float32 {
	return ui.height
}

func (ui *UIText) SetString(text string) {
	ui.text = text
	ui.updateText = true
}

func (ui *UIText) setString(text string) {
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

		yratio := atlasImage.RelativeHeight
		xratio := atlasImage.RelativeWidth

		//ygrid := -0.5 + (atlasImage.YGrid)
		//xgrid := (-0.5 + (atlasImage.XGrid)) + space
		ygrid := -(h / 2) + (atlasImage.YOffset)
		xgrid := -(w / 2) + (atlasImage.XOffset) + space
		space += atlasImage.XAdvance * spaceMult

		vertexCount += 4

		uv := engine.IndexUV(ui.Font, rune)

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
		width += atlasImage.XAdvance * spaceMult
		/*
			yratio := atlasImage.PlaneHeight
			ygrid := atlasImage.YGrid
			if yratio < 0 {
				yratio = -yratio
			}
			if ygrid < 0 {
				ygrid = -ygrid
			}

			if yratio+ygrid > height {
				height = yratio + ygrid
			}
		*/
		height = 1
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

//var speed = float32(30)

func (ui *UIText) Update() {
	ui.UpdateCollider()

	/*
		s := ui.Transform().Scale()

		if s.X < 20 && speed < 0 {
			speed = -speed
		} else if s.X > 200 && speed > 0 {
			speed = -speed
		}

		s.X += speed * engine.DeltaTime()
		s.Y += speed * engine.DeltaTime()

		ui.Transform().SetScale(s)
	*/

	//Handle Tab & Backspace
	if ui.focused && ui.writeable {
		if len(ui.text) > 0 && input.KeyPress(input.KeyBackspace) {
			ui.updateText = true
			ui.text = ui.text[:len(ui.text)-1]
		}
		if input.KeyPress(input.KeyTab) {
			ui.updateText = true
			ui.text += "\t"
		}
	}
}

func (ui *UIText) LateUpdate() {
	if ui.updateText {
		ui.updateText = false
		ui.setString(ui.text)
	}
}

/*
Todo: make this an interface.
*/
func (ui *UIText) UpdateCollider() {
	//if ui.GameObject().Physics.Body.Enabled {
	if ui.GameObject().Physics == nil {
		return
	}
	b := ui.GameObject().Physics.Box
	body := ui.GameObject().Physics.Body
	if b != nil {
		h := vect.Float(float64(ui.height) * float64(ui.GameObject().Transform().WorldScale().Y))
		w := vect.Float(float64(ui.width) * float64(ui.GameObject().Transform().WorldScale().X))
		update := false
		if h != b.Height || w != b.Width {
			b.Width = w
			b.Height = h
			update = true
		}

		c := ui.align.Vector()
		center := vect.Vect{vect.Float(c.X), vect.Float(c.Y)}
		center.X = (center.X * w)
		center.Y = (center.Y * h)

		if b.Position.X != center.X || b.Position.Y != center.Y {
			update = true
			b.Position.X, b.Position.Y = center.X, center.Y
		}

		if update {
			b.UpdatePoly()
			if !body.MomentIsInf() && h != 0 && w != 0 {
				body.SetMoment(vect.Float(b.Moment(float32(body.Mass()))))
			}
		}
	}
	//log.Println(b.Height, b.Width, ui.GameObject().Transform().Scale().X, ui.GameObject().Name())
	//}
}

func (ui *UIText) Align() engine.Align {
	return ui.align
}

func (ui *UIText) SetAlign(align engine.Align) {
	ui.align = align
}

func (ui *UIText) Draw() {
	if ui.text == "" {
		return
	}

	v := ui.align.Vector()
	v.X = (v.X * ui.width)
	v.Y = (v.Y * ui.height)

	mat := engine.TextureMaterial
	if ui.Font.IsSDF() {
		mat = engine.SDFMaterial
	}
	mat.Begin(ui.GameObject())

	mp := mat.ProjMatrix
	mv := mat.ViewMatrix
	mm := mat.ModelMatrix
	tx := mat.Texture
	ti := mat.Tiling
	of := mat.Offset
	color := mat.AddColor
	_ = color

	ui.buffer.Bind(gl.ARRAY_BUFFER)
	mat.Verts.EnableArray()
	mat.Verts.AttribPointer(3, gl.FLOAT, false, 0, uintptr(0))
	mat.UV.EnableArray()
	mat.UV.AttribPointer(2, gl.FLOAT, false, 0, uintptr(ui.texcoordsIndex))

	camera := engine.GetScene().SceneBase().Camera

	view := camera.InvertedMatrix()
	model := engine.Identity()
	model.Translate(v.X, v.Y, 0)
	model.Mul(ui.GameObject().Transform().Matrix())
	model.Translate(0.75, 0.75, 0)
	/*
		view := camera.Transform().Matrix()
		view = view.Invert()
		model := ui.GameObject().Transform().Matrix()
	*/

	mv.UniformMatrix4fv(false, view)
	mp.UniformMatrix4f(false, (*[16]float32)(camera.Projection))
	mm.UniformMatrix4fv(false, model)
	ti.Uniform2f(1, 1)
	of.Uniform2f(0, 0)

	ui.Font.Bind()
	tx.Uniform1i(0)

	color.Uniform4f(ui.Color.R, ui.Color.G, ui.Color.B, ui.Color.A)

	gl.DrawArrays(gl.QUADS, 0, ui.vertexCount)

}
