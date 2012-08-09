package Engine

import (
	"github.com/banthar/gl"
	//"gl/glu"
	//"log"
	//"image/png"
	//"image"
	//"os"
	. "github.com/vova616/chipmunk/vect"
	//"fmt"
	//"glfw"
)

type Sprite struct {
	BaseComponent
	*Texture
	buffer         gl.Buffer
	AnimationSpeed float32
	texcoordsIndex int
	endAnimation   int	
	startAnimation int
	animation	   float32	
	UVs		 	   AnimatedUV
	animMap 	   map[interface{}][2]int
	currentAnim	   interface{}
}

func NewSprite(tex *Texture) *Sprite {
	sp := &Sprite{NewComponent(), tex, gl.GenBuffer(), 1,0,1,0,0, AnimatedUV{NewUV(0, 0, 1, 1, float32(tex.Width())/float32(tex.Height()))}, nil,nil}
	sp.CreateVBO(sp.UVs...)

	return sp
}

func NewSprite2(tex *Texture, uv UV) *Sprite {
	sp := &Sprite{NewComponent(), tex, gl.GenBuffer(), 1,0,1,0,0, AnimatedUV{uv}, nil,nil}
	sp.CreateVBO(uv)

	return sp
}

func NewSprite3(tex *Texture, uv AnimatedUV) *Sprite {
	sp := &Sprite{NewComponent(), tex, gl.GenBuffer(), 1,0,len(uv),0,0, uv, nil,nil}
	sp.CreateVBO(uv...)

	return sp
}

func (p *Sprite)BindAnimations(animMap map[interface{}][2]int) {
	p.animMap = animMap
}

func (p *Sprite)SetAnimation(id interface{}) {
	a,e := p.animMap[id]
	if !e {
		panic("no such id")
	}
	p.currentAnim = id
	p.animation = float32(a[0])
	p.startAnimation = a[0]
	p.endAnimation = a[1]
} 

func (p *Sprite)CurrentAnimation() interface{} {
	return p.currentAnim
}

func (p *Sprite)OnComponentBind(binded *GameObject) {
	binded.Sprite = p	
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
		xgrid := float32(-uv.Ratio/2)

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

func (sp *Sprite) Update() {
	if int(sp.animation) < sp.endAnimation {
		sp.animation += sp.AnimationSpeed*DeltaTime()
	}
	if sp.animation >= float32(sp.endAnimation) {
		sp.animation = float32(sp.startAnimation)
	}

	if sp.GameObject().Physics != nil {
		box := sp.GameObject().Physics.Box
		cir := sp.GameObject().Physics.Shape.GetAsCircle()
		
		scale := sp.Transform().WorldScale()
		ratio := sp.UVs[int(sp.animation)].Ratio
		scale.X *= ratio
		
		if box != nil {
			if Float(scale.Y) != box.Height || Float(scale.X) != box.Width {
				box.Height = Float(scale.Y)
				box.Width = Float(scale.X)
				//box.Position = Vect{box.Width/2, box.Height/2}
				box.UpdatePoly()
			}
		} else if cir != nil {
			if float32(cir.Radius) != scale.X/2 {
				cir.Radius = Float(scale.X/2)
				sp.GameObject().Physics.Body.UpdateShapes()
			}
		}
	}
}

func (sp *Sprite) Draw() {
	if sp.Texture != nil {

		sp.Bind()

		gl.EnableClientState(gl.VERTEX_ARRAY)
		gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
		
		sp.buffer.Bind(gl.ARRAY_BUFFER)
		//gl.VertexPointer
		gl.VertexPointerVBO(3, gl.FLOAT, 0,  (int(sp.animation)*12*4))
		gl.TexCoordPointerVBO(2, gl.FLOAT, 0, (sp.texcoordsIndex+(int(sp.animation)*8*4)))
		
		gl.DrawArrays(gl.QUADS, 0, 4)
		
		gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)
		gl.DisableClientState(gl.VERTEX_ARRAY)
		
		sp.Unbind()
		//gl.PopMatrix()
	}
}


