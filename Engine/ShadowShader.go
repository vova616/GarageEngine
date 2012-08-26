package Engine

import (
	"github.com/vova616/gl"
	//"gl/glu"
	//"log"
	//"image/png"
	//"image"
	//"os"
	//. "github.com/vova616/chipmunk/vect"
	//"fmt"
	//"glfw"
	"image/color"
)


var ShadowShaderProgram  gl.Program
var ShadowCalcProgram gl.Program


const vertexShadowShader = `
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


const fragmentShadowShader = `
#version 130
precision highp float; // needed only for version 1.30
 
out vec4 color;

in vec2 UV;
uniform sampler2D mytexture;
 
void main(void)
{
	vec4 c = texture2D(mytexture, UV);
	if (c.r > 0 || c.g > 0 || c.b > 0 || c.a > 0) {
		c.rgba = vec4(0,0,0,1);
	}
  	color = c;
}
`




const vertexShadowShader2 = `
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


const fragmentShadowShader2 = `
#version 130
precision highp float; // needed only for version 1.30
 
out vec4 color;

in vec2 UV;
uniform sampler2D mytexture;
 
void main(void)
{
	vec2 newUV = UV;
	newUV.y = 1-newUV.y;
	vec4 oc = texture2D(mytexture, newUV);

	
	vec2 px = vec2(1024*UV.x,768*UV.y);
	
	float ds = distance(px, vec2(512, 384));
	ds /= 2;
	
	vec2 d = (vec2(0.5,0.5)-newUV)/ds;
	
	//
	
	vec2 pos = newUV;
	
	for (int i=0;i<ds;i++) {
		pos += d;
		vec4 c = texture2D(mytexture, pos);
		
		float a = 0;
		if (ds > 300) {
			a = ds/300;
		} else {
			a = ds/300;
		}
		
		if (c.r == 0) {
			if (oc.r == 0) {
				c.rgba = vec4(0.2,0.2,0.2,a+0.2);	
			} else {
				c.rgba = vec4(0.2,0.2,0.2,a);	
			}
			color = c;
			return;
		}
	}
	if (oc.r == 0) {
		color = vec4(0,0,0,ds/300);
	} else {	
		color = vec4(0,0,0,1);	
	}
}
`
	
	
	
	//	Width  = 1024
	//	Height = 768
	
	func loadShadowShader() {
	
		program := gl.CreateProgram()
		vrt := gl.CreateShader(gl.VERTEX_SHADER)
		frg := gl.CreateShader(gl.FRAGMENT_SHADER)
		
		vrt.Source(vertexShadowShader)
		frg.Source(fragmentShadowShader)
		
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
		 
		ShadowShaderProgram = program 
	}
	
	func loadShadowCalc() {
	
		program := gl.CreateProgram()
		vrt := gl.CreateShader(gl.VERTEX_SHADER)
		frg := gl.CreateShader(gl.FRAGMENT_SHADER)
		
		vrt.Source(vertexShadowShader2)
		frg.Source(fragmentShadowShader2)
		
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
		 
		ShadowCalcProgram = program 
	}
	
	
	type ShadowShader struct {
		BaseComponent
		Texture 	*Texture
		FrameBuffer gl.Framebuffer
		Camera *Camera
		Sprite *Sprite
	}
	
	func NewShadowShader(c *Camera) *ShadowShader {
	
		return &ShadowShader{BaseComponent: NewComponent(), Camera:c}
	}
	
	func (s *ShadowShader) Start() {
		loadShadowShader()
		loadShadowCalc()
		
		frameBuffer := gl.GenFramebuffer()
		texture := NewTextureEmpty(Width, Height, color.RGBAModel)
		
		s.Texture = texture
		s.FrameBuffer = frameBuffer
	}
	
	func (s *ShadowShader) Draw() {
		return
		
		s.FrameBuffer.Bind()
		gl.ClearColor(255, 255, 255, 255)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0,gl.TEXTURE_2D, s.Texture.GLTexture(), 0)
		//gl.DrawBuffers(gl.COLOR_ATTACHMENT0) 
		
		t := TextureShader
		TextureShader = ShadowShaderProgram
		s.Camera.Render()
		TextureShader = t
		
	
		
		s.FrameBuffer.Unbind()
		
		//s.FrameBuffer.Bind()
	
		//gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0,gl.TEXTURE_2D, s.Texture.GLTexture(), 0)
		
		t = TextureShader
		TextureShader = ShadowCalcProgram
		sp := NewSprite(s.Texture)
		sp.DrawScreen()
		TextureShader = t
		//s.FrameBuffer.Unbind()
		
		//s.Sprite.Texture = s.Texture
	}
	