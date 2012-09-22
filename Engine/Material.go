package Engine

import (
	"github.com/vova616/gl"
)

type Material interface {
	Load()
	Begin(gobj *GameObject)
	End(gobj *GameObject)
}

type BasicMaterial struct {
	Program        gl.Program
	vertexShader   string
	fragmentShader string

	ViewMatrix, ProjMatrix, ModelMatrix, BorderColor, AddColor, Texture gl.UniformLocation
	Verts, UV                                                           gl.AttribLocation
}

func NewBasicMaterial(vertexShader, fragmentShader string) *BasicMaterial {
	return &BasicMaterial{Program: gl.CreateProgram(), vertexShader: vertexShader, fragmentShader: fragmentShader}
}

func (b *BasicMaterial) Load() {
	program := b.Program
	vrt := gl.CreateShader(gl.VERTEX_SHADER)
	frg := gl.CreateShader(gl.FRAGMENT_SHADER)

	vrt.Source(vertexShader)
	frg.Source(fragmentShader)

	vrt.Compile()
	if vrt.Get(gl.COMPILE_STATUS) != 1 {
		println(vrt.GetInfoLog())
	}
	frg.Compile()
	if frg.Get(gl.COMPILE_STATUS) != 1 {
		println(frg.GetInfoLog())
	}

	program.AttachShader(vrt)
	program.AttachShader(frg)

	program.BindAttribLocation(0, "vertexPos")
	program.BindAttribLocation(1, "vertexUV")

	program.Link()

	b.Verts = program.GetAttribLocation("vertexPos")
	b.UV = program.GetAttribLocation("vertexUV")
	b.ViewMatrix = program.GetUniformLocation("MView")
	b.ProjMatrix = program.GetUniformLocation("MProj")
	b.ModelMatrix = program.GetUniformLocation("MModel")
	b.BorderColor = program.GetUniformLocation("bcolor")
	b.Texture = program.GetUniformLocation("mytexture")
	b.AddColor = program.GetUniformLocation("addcolor")

}

func (b *BasicMaterial) Begin(gobj *GameObject) {
	b.Program.Use()
}

func (b *BasicMaterial) End(gobj *GameObject) {

}

var TextureShader gl.Program
var TextureMaterial *BasicMaterial

const vertexShader = `
#version 130

uniform mat4 MProj;
uniform mat4 MView;
uniform mat4 MModel;

in  vec3 vertexPos;
in  vec2 vertexUV;
out vec2 UV;


 
void main(void)
{
	gl_Position = MProj * MView * MModel * vec4(vertexPos, 1.0);
	UV = vertexUV;
}
`

const fragmentShader = `
#version 130
precision highp float; // needed only for version 1.30
 
out vec4 color;

in vec2 UV; 
uniform sampler2D mytexture;
uniform vec4 bcolor;
uniform vec4 addcolor;
 
void main(void)
{ 
  	vec4 tcolor = texture2D(mytexture, UV);
	if (tcolor.a > 0) {
		tcolor += bcolor;
	}
	float a = tcolor.a;
	tcolor = mix(addcolor, tcolor*addcolor, 1); 
	
	color = tcolor;
}
`

func loadShader() {

	program := gl.CreateProgram()
	vrt := gl.CreateShader(gl.VERTEX_SHADER)
	frg := gl.CreateShader(gl.FRAGMENT_SHADER)

	vrt.Source(vertexShader)
	frg.Source(fragmentShader)

	vrt.Compile()
	if vrt.Get(gl.COMPILE_STATUS) != 1 {
		println(vrt.GetInfoLog())
	}
	frg.Compile()
	if frg.Get(gl.COMPILE_STATUS) != 1 {
		println(frg.GetInfoLog())
	}

	program.AttachShader(vrt)
	program.AttachShader(frg)

	program.BindAttribLocation(0, "vertexPos")
	program.BindAttribLocation(1, "vertexUV")

	program.Link()
	program.Use()

	TextureShader = program
}
