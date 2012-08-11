package Engine

import (
	"github.com/banthar/gl"
)

var TextureShader gl.Program


const vertexShader = `
#version 130

uniform mat4 MProj;
uniform mat4 MView;
uniform mat4 MModel;

in  vec3 vectexPos;
in  vec2 vertexUV;
out vec2 UV;

 
void main(void)
{
	gl_Position = MProj * MView * MModel * vec4(vectexPos, 1.0);
	UV = vertexUV;
}
`


const fragmentShader = `
#version 130
precision highp float; // needed only for version 1.30
 
out vec4 color;

in vec2 UV;
uniform sampler2D mytexture;
 
void main(void)
{
  	color = texture2D(mytexture, UV);
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
