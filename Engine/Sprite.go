package Engine

import (
	"github.com/vova616/gl"
	//"gl/glu"
	//"log"
	//"image/png"
	//"image"
	//"os"
	//"fmt"

	"github.com/vova616/chipmunk/vect"
	//"glfw"
)

type OnAnimationEnd func(sprite *Sprite)

type Sprite struct {
	BaseComponent
	*Texture
	buffer               gl.Buffer
	AnimationSpeed       float32
	texcoordsIndex       int
	endAnimation         int
	startAnimation       int
	animation            float32
	UVs                  AnimatedUV
	animMap              map[interface{}][2]int
	currentAnim          interface{}
	AnimationEndCallback OnAnimationEnd

	Render bool

	Border     bool
	BorderSize float32
	Color      Vector

	align AlignType
}

func NewSprite(tex *Texture) *Sprite {
	return NewSprite3(tex, AnimatedUV{NewUV(0, 0, 1, 1, float32(tex.Width())/float32(tex.Height()))})
}

func NewSprite2(tex *Texture, uv UV) *Sprite {
	return NewSprite3(tex, AnimatedUV{uv})
}

func NewSprite3(tex *Texture, uv AnimatedUV) *Sprite {

	sp := &Sprite{
		BaseComponent:  NewComponent(),
		Texture:        tex,
		buffer:         gl.GenBuffer(),
		AnimationSpeed: 1,
		endAnimation:   len(uv),
		UVs:            uv,
		Render:         true,
		Color:          Vector{1, 1, 1},
		align:          AlignCenter,
	}
	sp.CreateVBO(uv...)

	return sp
}

func (p *Sprite) BindAnimations(animMap map[interface{}][2]int) {
	p.animMap = animMap
}

func (p *Sprite) SetAnimation(id interface{}) {
	a, e := p.animMap[id]
	if !e {
		panic("no such id")
	}
	p.currentAnim = id
	p.animation = float32(a[0])
	p.startAnimation = a[0]
	p.endAnimation = a[1]
}

func (p *Sprite) Align() AlignType {
	return p.align
}

func (p *Sprite) SetAlign(align AlignType) {
	p.align = align
}

func (p *Sprite) CurrentAnimation() interface{} {
	return p.currentAnim
}

func (sp *Sprite) OnComponentBind(binded *GameObject) {
	binded.Sprite = sp
}

func (sp *Sprite) CreateVBO(uvs ...UV) {

	l := len(uvs)

	if l == 0 {
		sp.buffer.Bind(gl.ARRAY_BUFFER)
		gl.BufferData(gl.ARRAY_BUFFER, 1, []byte{1}, gl.STATIC_DRAW)
		return
	}

	lt := l * 12
	_ = lt
	data := make([]float32, 20*l)

	vertexCount := 0
	texcoordsIndex := lt * 4
	for i, uv := range uvs {

		yratio := float32(1)
		xratio := uv.Ratio
		ygrid := float32(-0.5)
		xgrid := float32(-uv.Ratio / 2)

		vertexCount += 4

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

	sp.texcoordsIndex = texcoordsIndex
	sp.buffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, data, gl.STATIC_DRAW)
}

func (sp *Sprite) Start() {

}

func (sp *Sprite) SetAnimationIndex(index int) {
	if index <= 0 {
		sp.animation = 0
	} else {
		sp.animation = float32(index % sp.endAnimation)
	}
}

func (sp *Sprite) Update() {
	if int(sp.animation) < sp.endAnimation {
		sp.animation += sp.AnimationSpeed * DeltaTime()
	}
	if sp.animation >= float32(sp.endAnimation) {
		if sp.AnimationEndCallback != nil {
			sp.AnimationEndCallback(sp)
		}
		sp.animation = float32(sp.startAnimation)
	}

	sp.UpdateShape()
}

/*
Todo: make this an interface.
*/
func (sp *Sprite) UpdateShape() {
	if sp.GameObject().Physics != nil {
		ph := sp.GameObject().Physics
		box := ph.Box
		cir := ph.Shape.GetAsCircle()

		scale := sp.Transform().WorldScale()
		ratio := sp.UVs[int(sp.animation)].Ratio
		scale.X *= ratio

		if box != nil {
			update := false
			if vect.Float(scale.Y) != box.Height || vect.Float(scale.X) != box.Width {
				box.Height = vect.Float(scale.Y)
				box.Width = vect.Float(scale.X)
				//box.Position = Vect{box.Width/2, box.Height/2}
				update = true
			}

			c := Align(sp.align)
			center := vect.Vect{vect.Float(c.X), vect.Float(c.Y)}
			center.X *= vect.Float(scale.X)
			center.Y *= vect.Float(scale.Y)

			if box.Position.X != center.X || box.Position.Y != center.Y {
				update = true
				box.Position.X, box.Position.Y = center.X, center.Y
			}

			if update {
				box.UpdatePoly()
				if !ph.Body.MomentIsInf() && box.Height != 0 && box.Width != 0 {
					ph.Body.SetMoment(vect.Float(box.Moment(float32(ph.Body.Mass()))))
				}
			}
		} else if cir != nil {
			update := false
			s := float32(0)
			if scale.X > scale.Y {
				s = scale.X
			} else {
				s = scale.Y
			}
			if float32(cir.Radius) != s/2 {
				cir.Radius = vect.Float(s / 2)
				update = true
			}

			c := Align(sp.align)
			center := vect.Vect{vect.Float(c.X), vect.Float(c.Y)}
			center.X *= vect.Float(s)
			center.Y *= vect.Float(s)

			if cir.Position.X != center.X || cir.Position.Y != center.Y {
				update = true
				cir.Position.X, cir.Position.Y = center.X, center.Y
			}

			if update {
				sp.GameObject().Physics.Body.UpdateShapes()
				if !ph.Body.MomentIsInf() && cir.Radius != 0 {
					//log.Println(sp.gameObject.name, cir.Radius, cir.Moment(float32(ph.Body.Mass())), scale.X, scale.Y, ph.Body.Mass(), cir.Position)
					ph.Body.SetMoment(vect.Float(cir.Moment(float32(ph.Body.Mass()))))
				}
			}
		}
	}
}

func Abs(val float32) float32 {
	if val < 0 {
		return -val
	}
	return val
}

func (sp *Sprite) Draw() {
	if sp.Texture != nil && sp.Render {

		camera := GetScene().SceneBase().Camera
		cameraPos := camera.Transform().WorldPosition()
		pos := sp.Transform().WorldPosition()
		scale := sp.Transform().WorldScale()
		if Abs(pos.X-cameraPos.X)-scale.X-float32(Width)/2 > float32(Width) {
			return
		}
		if Abs(pos.Y-cameraPos.Y)-scale.Y-float32(Height)/2 > float32(Height) {
			return
		}

		TextureMaterial.Begin(sp.GameObject())

		vert := TextureMaterial.Verts
		uv := TextureMaterial.UV
		mp := TextureMaterial.ProjMatrix
		mv := TextureMaterial.ViewMatrix
		mm := TextureMaterial.ModelMatrix
		mc := TextureMaterial.BorderColor
		tx := TextureMaterial.Texture
		ac := TextureMaterial.AddColor

		vert.EnableArray()
		uv.EnableArray()

		sp.buffer.Bind(gl.ARRAY_BUFFER)

		vert.AttribPointer(3, gl.FLOAT, false, 0, uintptr(int(sp.animation)*12*4))
		uv.AttribPointer(2, gl.FLOAT, false, 0, uintptr(sp.texcoordsIndex+(int(sp.animation)*8*4)))

		v := Align(sp.align)
		v.X *= float32(sp.Texture.Width()) / float32(sp.Texture.Height())

		view := camera.Transform().Matrix()
		view = view.Invert()
		model := Identity()
		model.Translate(v.X, v.Y, 0)
		model.Mul(sp.GameObject().Transform().Matrix())

		mv.UniformMatrix4fv(false, view)
		mp.UniformMatrix4fv(false, *camera.Projection)
		mm.UniformMatrix4fv(false, model)

		sp.Bind()
		gl.ActiveTexture(gl.TEXTURE0)
		tx.Uniform1i(0)

		//ac.Uniform4f(1, 1, 1, 0) 
		ac.Uniform4f(sp.Color.X, sp.Color.Y, sp.Color.Z, 1)

		if sp.Border {
			mc.Uniform4f(1, 1, 1, 0)

			gl.DrawArrays(gl.QUADS, 0, 4)

			scale := sp.Transform().Scale()
			scalex := scale.Mul2(1 - (sp.BorderSize / 100))
			sp.Transform().SetScale(scalex)
			model = sp.GameObject().Transform().Matrix()

			mm.UniformMatrix4fv(false, model)
			sp.Transform().SetScale(scale)
		}

		mc.Uniform4f(0, 0, 0, 0)

		gl.DrawArrays(gl.QUADS, 0, 4)

		sp.Unbind()
		vert.DisableArray()
		uv.DisableArray()

		TextureMaterial.End(sp.GameObject())
	}
}

func (sp *Sprite) DrawScreen() {
	if sp.Texture != nil && sp.Render {

		camera := GetScene().SceneBase().Camera
		pos := sp.Transform().WorldPosition()
		scale := sp.Transform().WorldScale()

		TextureMaterial.Begin(sp.GameObject())

		vert := TextureMaterial.Verts
		uv := TextureMaterial.UV
		mp := TextureMaterial.ProjMatrix
		mv := TextureMaterial.ViewMatrix
		mm := TextureMaterial.ModelMatrix
		mc := TextureMaterial.BorderColor
		tx := TextureMaterial.Texture
		ac := TextureMaterial.AddColor

		vert.EnableArray()
		uv.EnableArray()

		sp.buffer.Bind(gl.ARRAY_BUFFER)

		vert.AttribPointer(3, gl.FLOAT, false, 0, uintptr(int(sp.animation)*12*4))
		uv.AttribPointer(2, gl.FLOAT, false, 0, uintptr(sp.texcoordsIndex+(int(sp.animation)*8*4)))

		proj := camera.Projection
		view := Identity()
		model := Identity()
		model.Scale(scale.X, scale.Y, 1)
		model.Translate((float32(Width)/2)+pos.X, (float32(Height)/2)+pos.Y, 1)

		mv.UniformMatrix4fv(false, view)
		mp.UniformMatrix4fv(false, *proj)
		mm.UniformMatrix4fv(false, model)

		sp.Bind()
		gl.ActiveTexture(gl.TEXTURE0)
		tx.Uniform1i(0)

		//ac.Uniform4f(1, 1, 1, 0) 
		ac.Uniform4f(1, 1, 1, 1)

		if sp.Border {
			mc.Uniform4f(1, 1, 1, 0)

			gl.DrawArrays(gl.QUADS, 0, 4)

			scale := sp.Transform().Scale()
			scalex := scale.Mul2(1 - (sp.BorderSize / 100))
			sp.Transform().SetScale(scalex)
			model = sp.GameObject().Transform().Matrix()

			mm.UniformMatrix4fv(false, model)
			sp.Transform().SetScale(scale)
		}

		mc.Uniform4f(0, 0, 0, 0)

		gl.DrawArrays(gl.QUADS, 0, 4)

		sp.Unbind()
		vert.DisableArray()
		uv.DisableArray()

		TextureMaterial.End(sp.GameObject())
	}
}
