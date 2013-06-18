package engine

import (
	"github.com/go-gl/gl"
)

var (
	defaultBuffer      VBO
	defaultIndexBuffer VBO
)

func initDefaultPlane() {
	defaultBuffer = GenBuffer()

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
	return CurrentCamera().InsideScreen(ratio, position, scale)
}

func DrawSprite(tex *Texture, uv UV, position Vector, scale Vector, rotation float32, aling Align, color Color) {
	if !InsideScreen(uv.Ratio, position, scale) {
		return
	}

	internalMaterial.Begin(nil)

	mp := internalMaterial.ProjMatrix
	mv := internalMaterial.ViewMatrix
	mm := internalMaterial.ModelMatrix
	tx := internalMaterial.Texture
	ac := internalMaterial.AddColor
	ti := internalMaterial.Tiling
	of := internalMaterial.Offset

	defaultBuffer.Bind(gl.ARRAY_BUFFER)
	internalMaterial.Verts.EnableArray()
	internalMaterial.Verts.AttribPointer(3, gl.FLOAT, false, 0, uintptr(0))
	internalMaterial.UV.EnableArray()
	internalMaterial.UV.AttribPointer(2, gl.FLOAT, false, 0, uintptr(12*4))

	v := aling.Vector()
	v.X *= uv.Ratio

	camera := GetScene().SceneBase().Camera
	view := camera.InvertedMatrix()
	model := Identity()
	model.Translate(v.X, v.Y, 0)

	model.Scale(scale.X*uv.Ratio, scale.Y, scale.Z)
	model.RotateZ(rotation, -1)
	model.Translate(position.X+0.75, position.Y+0.75, position.Z)

	mv.UniformMatrix4fv(false, view)
	mp.UniformMatrix4f(false, (*[16]float32)(camera.Projection))
	mm.UniformMatrix4fv(false, model)
	ac.Uniform4i(int(color.R), int(color.G), int(color.B), int(color.A))
	ti.Uniform2f(uv.U2-uv.U1, uv.V2-uv.V1)
	of.Uniform2f(uv.U1, uv.V1)

	tx.Uniform1i(0)

	tex.Bind()

	gl.DrawArrays(gl.QUADS, 0, 4)

	internalMaterial.End(nil)
}

func DrawSprites(tex *Texture, uvs []UV, positions []Vector, scales []Vector, rotations []float32, alings []Align, colors []Color) {

	internalMaterial.Begin(nil)

	mp := internalMaterial.ProjMatrix
	mv := internalMaterial.ViewMatrix
	mm := internalMaterial.ModelMatrix
	tx := internalMaterial.Texture
	ac := internalMaterial.AddColor
	ti := internalMaterial.Tiling
	of := internalMaterial.Offset

	defaultBuffer.Bind(gl.ARRAY_BUFFER)
	internalMaterial.Verts.EnableArray()
	internalMaterial.Verts.AttribPointer(3, gl.FLOAT, false, 0, uintptr(0))
	internalMaterial.UV.EnableArray()
	internalMaterial.UV.AttribPointer(2, gl.FLOAT, false, 0, uintptr(12*4))

	camera := GetScene().SceneBase().Camera
	view := camera.InvertedMatrix()
	mv.UniformMatrix4fv(false, view)
	mp.UniformMatrix4f(false, (*[16]float32)(camera.Projection))

	tex.Bind()
	tx.Uniform1i(0)

	for i := 0; i < len(uvs); i++ {

		uv, position, scale := uvs[i], positions[i], scales[i]

		if !InsideScreen(uv.Ratio, position, scale) {
			continue
		}

		rotation, aling, color := rotations[i], alings[i], colors[i]

		v := aling.Vector()
		v.X *= uv.Ratio

		model := Identity()
		model.Translate(v.X, v.Y, 0)

		model.Scale((scale.X * uv.Ratio), scale.Y, scale.Z)
		model.RotateZ(rotation, -1)
		model.Translate(position.X+0.75, position.Y+0.75, position.Z)

		mm.UniformMatrix4fv(false, model)
		ac.Uniform4f(color.R, color.G, color.B, color.A)
		ti.Uniform2f(uv.U2-uv.U1, uv.V2-uv.V1)
		of.Uniform2f(uv.U1, uv.V1)

		gl.DrawArrays(gl.QUADS, 0, 4)
	}

	internalMaterial.End(nil)
}
