package engine

import (
	"github.com/vova616/gl"
)

var (
	defaultBuffer gl.Buffer
)

func initDefaultPlane() {
	defaultBuffer = gl.GenBuffer()

	//Triagles
	data := make([]float32, 20)
	data[0] = -0.5
	data[1] = -0.5
	data[2] = 1

	data[3] = 0.5
	data[4] = -0.5
	data[5] = 1

	data[6] = 0.5
	data[7] = 0.5
	data[8] = 1

	data[9] = -0.5
	data[10] = 0.5
	data[11] = 1

	// UV
	data[12] = 0
	data[13] = 1

	data[14] = 1
	data[15] = 1

	data[16] = 1
	data[17] = 0

	data[18] = 0
	data[19] = 0

	defaultBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, data, gl.STATIC_DRAW)

}

func InsideScreen(ratio float32, position Vector, scale Vector) bool {
	cameraPos := GetScene().SceneBase().Camera.Transform().WorldPosition()

	bigScale := scale.X * ratio
	if scale.Y > bigScale {
		bigScale = scale.Y
	}
	bigScale = -bigScale

	x := (position.X - cameraPos.X) + (bigScale / 2)
	y := (position.Y - cameraPos.Y) + (bigScale / 2)
	if x > float32(Width) || x < bigScale {
		return false
	}
	if y > float32(Height) || y < bigScale {
		return false
	}
	return true
}

func DrawSprite(tex *Texture, uv UV, position Vector, scale Vector, rotation float32, aling AlignType, color Vector) {
	if !InsideScreen(uv.Ratio, position, scale) {
		return
	}

	TextureMaterial.Begin(nil)

	vert := TextureMaterial.Verts
	uvb := TextureMaterial.UV
	mp := TextureMaterial.ProjMatrix
	mv := TextureMaterial.ViewMatrix
	mm := TextureMaterial.ModelMatrix
	tx := TextureMaterial.Texture
	ac := TextureMaterial.AddColor
	ti := TextureMaterial.Tiling
	of := TextureMaterial.Offset

	vert.EnableArray()
	uvb.EnableArray()

	defaultBuffer.Bind(gl.ARRAY_BUFFER)

	vert.AttribPointer(3, gl.FLOAT, false, 0, uintptr(0))
	uvb.AttribPointer(2, gl.FLOAT, false, 0, uintptr(12*4))

	v := Align(aling)
	v.X *= uv.Ratio

	camera := GetScene().SceneBase().Camera
	view := camera.InvertedMatrix()
	model := Identity()
	model.Translate(v.X, v.Y, 0)

	model.Scale(scale.X*uv.Ratio, scale.Y, scale.Z)
	model.Rotate(rotation, 0, 0, -1)
	model.Translate(position.X, position.Y, position.Z)

	mv.UniformMatrix4fv(false, view)
	mp.UniformMatrix4f(false, (*[16]float32)(camera.Projection))
	mm.UniformMatrix4fv(false, model)
	ac.Uniform4f(color.X, color.Y, color.Z, 1)
	ti.Uniform2f(uv.U2-uv.U1, uv.V2-uv.V1)
	of.Uniform2f(uv.U1, uv.V1)

	tx.Uniform1i(0)

	tex.Bind()

	gl.DrawArrays(gl.QUADS, 0, 4)

	TextureMaterial.End(nil)
}
