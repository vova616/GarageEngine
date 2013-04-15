package engine

import (
	"github.com/go-gl/gl"
)

var (
	lastVAO       VAO
	lastVBO       VBO
	lastVBOTarget gl.GLenum
)

type VBO gl.Buffer
type VAO gl.VertexArray

func GenBuffer() VBO {
	return VBO(gl.GenBuffer())
}

func GenVertexArray() VAO {
	return VAO(gl.GenVertexArray())
}

func (this VBO) Bind(target gl.GLenum) {
	if lastVBO != this || lastVBOTarget != target {
		lastVBO = this
		lastVBOTarget = target
		gl.Buffer.Bind(gl.Buffer(this), target)
	}
}

func (this VAO) Bind() {
	if lastVAO != this {
		lastVAO = this
		gl.VertexArray.Bind(gl.VertexArray(this))
	}
}

type Batch interface {
	Add(position, scale Vector, rotation float32, uv UV) (index int)
	Update(index int, position, scale Vector, rotation float32, uv UV)
	UpdateUV(uv UV, index int)
	Remove(index int)
	Render()
}

type StaticBatch struct {
	Tex          *Texture
	Vertices     VBO
	UVs          VBO
	Indecies     VBO
	tempPBuffer  []float32
	tempUVBuffer []float32
	tempIBuffer  []uint16
}

func NewStaticBatch(tex *Texture) *StaticBatch {
	return &StaticBatch{
		Tex:          tex,
		Vertices:     GenBuffer(),
		UVs:          GenBuffer(),
		Indecies:     GenBuffer(),
		tempPBuffer:  make([]float32, 4*3),
		tempUVBuffer: make([]float32, 4*2),
		tempIBuffer:  make([]uint16, 3*2),
	}
}

func (this *StaticBatch) Add(position, scale Vector, rotation float32, uv UV) (index int) {
	panic("Not implemented yet")
}

func (this *StaticBatch) Update(position, scale Vector, rotation float32, uv UV, index int) {
	panic("Not implemented yet")
}

func (this *StaticBatch) UpdateUV(uv UV, index int) {
	panic("Not implemented yet")
}

func (this *StaticBatch) Remove(index int) {
	panic("Not implemented yet")
}

func (this *StaticBatch) Render() {
	panic("Not implemented yet")
}
