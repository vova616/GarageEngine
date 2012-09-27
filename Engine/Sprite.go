package Engine

import (
	"github.com/vova616/gl"
	//"gl/glu"
	//"log"
	//"image/png"
	//"image"
	//"os"
	//"fmt"
	. "github.com/vova616/chipmunk/vect"
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

	Border     bool
	BorderSize float32
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
		UVs:            uv}
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

func (sp *Sprite) UpdateShape() {
	if sp.GameObject().Physics != nil {
		ph := sp.GameObject().Physics
		box := ph.Box
		cir := ph.Shape.GetAsCircle()

		scale := sp.Transform().WorldScale()
		ratio := sp.UVs[int(sp.animation)].Ratio
		scale.X *= ratio

		if box != nil {
			if Float(scale.Y) != box.Height || Float(scale.X) != box.Width {
				box.Height = Float(scale.Y)
				box.Width = Float(scale.X)
				if !ph.Body.MomentIsInf() {
					ph.Body.SetMoment(Float(box.Moment(float32(ph.Body.Mass()))))
				}
				//box.Position = Vect{box.Width/2, box.Height/2}
				box.UpdatePoly()
			}
		} else if cir != nil {
			if float32(cir.Radius) != scale.X/2 {
				cir.Radius = Float(scale.X / 2)
				if !ph.Body.MomentIsInf() {
					ph.Body.SetMoment(Float(cir.Moment(float32(ph.Body.Mass()))))
				}
				sp.GameObject().Physics.Body.UpdateShapes()
			}
		}
	}
}

func (sp *Sprite) Draw() {
	if sp.Texture != nil {

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

		vert.AttribPointerPtr(3, gl.FLOAT, false, 0, int(sp.animation)*12*4)
		uv.AttribPointerPtr(2, gl.FLOAT, false, 0, sp.texcoordsIndex+(int(sp.animation)*8*4))

		camera := GetScene().SceneBase().Camera

		view := camera.Transform().Matrix()
		view = view.Invert()
		model := sp.GameObject().Transform().Matrix()

		mv.Uniform4fv([]float32(view[:]))
		mp.Uniform4fv([]float32(camera.Projection[:]))
		mm.Uniform4fv([]float32(model[:]))

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

			mm.Uniform4fv([]float32(model[:]))
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
	if sp.Texture != nil {

		program := TextureShader
		program.Use()

		vert := program.GetAttribLocation("vectexPos")
		uv := program.GetAttribLocation("vertexUV")

		vert.EnableArray()
		uv.EnableArray()

		sp.buffer.Bind(gl.ARRAY_BUFFER)

		vert.AttribPointerPtr(3, gl.FLOAT, false, 0, int(sp.animation)*12*4)
		uv.AttribPointerPtr(2, gl.FLOAT, false, 0, sp.texcoordsIndex+(int(sp.animation)*8*4))

		camera := GetScene().SceneBase().Camera
		proj := NewIdentity()
		proj = camera.Projection
		view := NewIdentity()
		model := NewIdentity()
		model.Scale(float32(Height), float32(Height), 1)
		model.Translate(float32(Width)/2, float32(Height)/2, 1)

		mv := program.GetUniformLocation("MView")
		mv.Uniform4fv([]float32(view[:]))
		mp := program.GetUniformLocation("MProj")
		mp.Uniform4fv([]float32(proj[:]))
		mm := program.GetUniformLocation("MModel")
		mm.Uniform4fv([]float32(model[:]))

		sp.Bind()
		gl.ActiveTexture(gl.TEXTURE0)
		tx := program.GetUniformLocation("mytexture")
		tx.Uniform1i(0)

		gl.DrawArrays(gl.QUADS, 0, 4)

		sp.Unbind()
		vert.DisableArray()
		uv.DisableArray()
	}

}
