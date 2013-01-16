package Engine

import (
	"fmt"
	"github.com/go-gl/gl"
)

type Material interface {
	Load() error
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

func (b *BasicMaterial) Load() error {
	program := b.Program
	vrt := gl.CreateShader(gl.VERTEX_SHADER)
	frg := gl.CreateShader(gl.FRAGMENT_SHADER)

	vrt.Source(b.vertexShader)
	frg.Source(b.fragmentShader)

	vrt.Compile()
	if vrt.Get(gl.COMPILE_STATUS) != 1 {
		return fmt.Errorf("Error in Compiling Vertex Shader:%s\n", vrt.GetInfoLog())
	}
	frg.Compile()
	if frg.Get(gl.COMPILE_STATUS) != 1 {
		return fmt.Errorf("Error in Compiling Fragment Shader:%s\n", frg.GetInfoLog())
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
	return nil
}

func (b *BasicMaterial) Begin(gobj *GameObject) {
	b.Program.Use()
}

func (b *BasicMaterial) End(gobj *GameObject) {

}

var TextureShader gl.Program
var TextureMaterial *BasicMaterial

var SDFMaterial *BasicMaterial

const spriteVertexShader = `
#version 110

uniform mat4 MProj;
uniform mat4 MView;
uniform mat4 MModel;

attribute  vec3 vertexPos;
attribute  vec2 vertexUV; 
varying vec2 UV;


 
void main(void)
{
	gl_Position = MProj * MView * MModel * vec4(vertexPos, 1.0);
	UV = vertexUV;
}
`

const spriteFragmentShader = `
#version 110

varying vec2 UV; 
uniform sampler2D mytexture;
uniform vec4 bcolor;
uniform vec4 addcolor;

void main(void)
{ 
  	vec4 tcolor = texture2D(mytexture, UV);
	if (tcolor.a > 0.0) {
		tcolor += bcolor;
	}
	tcolor = tcolor*addcolor;

	//nice alpha detection
	//vec4 t = addcolor;
	//t.a = 0;
	//tcolor = mix(tcolor, t, tcolor.a);

	gl_FragColor = tcolor;
}
`

const sdfVertexShader = `
#version 110

uniform mat4 MProj;
uniform mat4 MView;
uniform mat4 MModel;

attribute  vec3 vertexPos;
attribute  vec2 vertexUV; 
varying vec2 UV;


 
void main(void)
{
	gl_Position = MProj * MView * MModel * vec4(vertexPos, 1.0);
	UV = vertexUV;
}
`

//Note: This shader needs to get better and get outline/shadow/glow support
const sdfFragmentShader = `
#version 110

varying vec2 UV; 
uniform sampler2D mytexture;
uniform vec4 bcolor;
uniform vec4 addcolor;
	
void main(void)
{ 
  	// retrieve distance from texture
	float sdf = texture2D( mytexture, UV).a;

	float smoothness = 32.0;
	float gamma = 2.2;

	vec4 basecolor = addcolor;

	
	bool outline = false;
	float outline_min = 0.4;
	float outline_min1 = 0.5;
	float outline_max = 0.4;
	float outline_max1 = 0.6;
	vec4 outline_color = vec4(1,0,0,1);

	bool softEdges = true;
	float softEdgeMin = 0.4;
	float softEdgeMax = 0.6;

	

	if (outline && sdf >= outline_min && sdf <= outline_max1 ) {
		float ofactor = 1.0;
		if (sdf <= outline_min1) {
			 ofactor = smoothstep(outline_min, outline_min1, sdf);
		} else {
			 ofactor = smoothstep(outline_max1, outline_max, sdf);
		}
		//lerp
		basecolor = addcolor + (outline_color - addcolor) * ofactor;
	}
	/*
	if (softEdges) {
		basecolor.a *= smoothstep(softEdgeMin, softEdgeMax, sdf);
	} else {
		if (sdf >= 0.5) {
			basecolor.a = 1.0;
		} else {
			basecolor.a = 0.0;
		}
	}	
	gl_FragColor = basecolor;
	*/
	


	// perform adaptive anti-aliasing of the edges\n
	float w = clamp( smoothness * (abs(dFdx(UV.x)) + abs(dFdy(UV.y))), 0.0, 0.5);
	float a = smoothstep(0.5-w, 0.5+w, sdf);

	a *= basecolor.a;

	// gamma correction for linear attenuation
	a = pow(a, 1.0/gamma);


	gl_FragColor.rgb = basecolor.rgb;
	gl_FragColor.a = a;

}
`
