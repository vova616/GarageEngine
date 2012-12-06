package Login

import (
	"fmt"
	"github.com/vova616/GarageEngine/Engine"
	"github.com/vova616/GarageEngine/Engine/Components"
	"github.com/vova616/GarageEngine/SpaceCookies/Game"
	_ "image/jpeg"
	_ "image/png"
	//"gl"  
	"strconv"
	"time"
	//"strings"
	//"math"
	//"github.com/vova616/chipmunk"
	//"github.com/vova616/chipmunk/vect"
	//"image"
	//"image/color"
	//"encoding/json"
	"math/rand"
	//"os"
)

type LoginScene struct {
	*Engine.SceneData
}

var (
	LoginSceneGeneral *LoginScene

	backgroundTexture *Engine.Texture
	ArialFont         *Engine.Font
	ArialFont2        *Engine.Font
)

func CheckError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func LoadTextures() {
	var e error
	ArialFont, e = Engine.NewFont("./data/Fonts/arial.ttf", 48)
	CheckError(e)
	ArialFont.Texture.SetReadOnly()

	ArialFont2, e = Engine.NewFont("./data/Fonts/arial.ttf", 24)
	CheckError(e)
	ArialFont2.Texture.SetReadOnly()

	backgroundTexture, e = Engine.LoadTexture("./data/SpaceCookies/background.png")
	CheckError(e)
}

func init() {
	Engine.Title = "Space Cookies"
}

func (s *LoginScene) Load() {

	LoadTextures()

	rand.Seed(time.Now().UnixNano())

	LoginSceneGeneral = s

	s.Camera = Engine.NewCamera()

	cam := Engine.NewGameObject("Camera")
	cam.AddComponent(s.Camera)

	cam.Transform().SetScalef(1, 1)

	background := Engine.NewGameObject("Background")
	background.AddComponent(Engine.NewSprite(backgroundTexture))
	background.AddComponent(Game.NewBackground(background.Sprite))
	background.Sprite.Render = false
	//background.Transform().SetScalef(float32(backgroung.Height()), float32(backgroung.Height()), 1)
	background.Transform().SetScalef(800, 800)
	background.Transform().SetPositionf(0, 0)

	gui := Engine.NewGameObject("GUI")

	mouse := Engine.NewGameObject("Mouse")
	mouse.AddComponent(Engine.NewMouse())
	mouse.Transform().SetParent2(cam)

	FPSDrawer := Engine.NewGameObject("FPS")
	FPSDrawer.Transform().SetParent2(cam)
	txt := FPSDrawer.AddComponent(Components.NewUIText(ArialFont2, "")).(*Components.UIText)
	fps := FPSDrawer.AddComponent(Engine.NewFPS()).(*Engine.FPS)
	fps.SetAction(func(fps float32) {
		txt.SetString("FPS: " + strconv.FormatFloat(float64(fps), 'f', 2, 32))
	})
	txt.SetAlign(Engine.AlignLeft)

	FPSDrawer.Transform().SetPositionf(20, float32(Engine.Height)-20)
	FPSDrawer.Transform().SetScalef(20, 20)

	//
	tBox := Engine.NewGameObject("TextBox")
	tBox.Transform().SetParent2(cam)

	txt2 := tBox.AddComponent(Components.NewUIText(ArialFont2, "Type your name: ")).(*Components.UIText)
	txt2.SetFocus(false)
	txt2.SetWritable(false)
	txt2.SetAlign(Engine.AlignLeft)

	tBox.Transform().SetPositionf(float32(Engine.Width)/2-txt2.Width()*20, float32(Engine.Height)/2)
	tBox.Transform().SetScalef(20, 20)
	//
	input := Engine.NewGameObject("TextBoxInput")
	input.Transform().SetParent2(cam)
	p := tBox.Transform().Position()
	p.X += txt2.Width() * 20
	input.Transform().SetPosition(p)
	input.Transform().SetScalef(20, 20)

	txt2 = input.AddComponent(Components.NewUIText(ArialFont2, "")).(*Components.UIText)
	txt2.SetFocus(true)
	txt2.SetWritable(true)
	txt2.SetAlign(Engine.AlignLeft)
	//
	login := Engine.NewGameObject("TextBoxInput")
	login.Transform().SetParent2(cam)
	login.Transform().SetPositionf(float32(Engine.Width)/2, float32(Engine.Height)/2-50)
	login.Transform().SetScalef(24, 24)

	login.AddComponent(Components.NewUIButton(func() {
		Engine.LoadScene(Game.GameSceneGeneral)
	}))

	txt2 = login.AddComponent(Components.NewUIText(ArialFont2, "Log in")).(*Components.UIText)
	txt2.SetFocus(false)
	txt2.SetWritable(false)
	txt2.SetAlign(Engine.AlignCenter)
	txt2.SetHoverCallBack(func(enter bool) {
		if enter {
			txt2.Color = Engine.Vector{1, 0, 0}
		} else {
			txt2.Color = Engine.Vector{0.5, 0, 0}
		}
	})
	txt2.Color = Engine.Vector{0.5, 0, 0}
	//

	//SPACCCEEEEE
	Engine.Space.Gravity.Y = 0
	Engine.Space.Iterations = 1

	s.AddGameObject(cam)
	s.AddGameObject(gui)
	s.AddGameObject(background)

	fmt.Println("Scene loaded")
}

func (s *LoginScene) SceneBase() *Engine.SceneData {
	return s.SceneData
}

func (s *LoginScene) New() Engine.Scene {
	gs := new(LoginScene)
	gs.SceneData = Engine.NewScene("LoginScene")
	return gs
}
