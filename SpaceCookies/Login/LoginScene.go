package Login

import (
	"fmt"
	"github.com/vova616/GarageEngine/Engine"
	"github.com/vova616/GarageEngine/Engine/Components"
	//"github.com/vova616/GarageEngine/Engine/Components/Tween"
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
	//"fmt"
)

type LoginScene struct {
	*Engine.SceneData
}

var (
	LoginSceneGeneral *LoginScene

	backgroundTexture *Engine.Texture
	button            *Engine.Texture
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
	ArialFont, e = Engine.NewFont("./data/Fonts/arial.ttf", 24)
	CheckError(e)
	ArialFont.Texture.SetReadOnly()

	ArialFont2, e = Engine.NewFont("./data/Fonts/arial.ttf", 24)
	CheckError(e)
	ArialFont2.Texture.SetReadOnly()

	backgroundTexture, e = Engine.LoadTexture("./data/SpaceCookies/background.png")
	CheckError(e)

	button, e = Engine.LoadTexture("./data/SpaceCookies/button.png")
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
	gui.Transform().SetParent2(cam)

	mouse := Engine.NewGameObject("Mouse")
	mouse.AddComponent(Engine.NewMouse())
	mouse.Transform().SetParent2(gui)

	FPSDrawer := Engine.NewGameObject("FPS")
	FPSDrawer.Transform().SetParent2(gui)
	txt := FPSDrawer.AddComponent(Components.NewUIText(ArialFont2, "")).(*Components.UIText)
	fps := FPSDrawer.AddComponent(Engine.NewFPS()).(*Engine.FPS)
	fps.SetAction(func(fps float32) {
		txt.SetString("FPS: " + strconv.FormatFloat(float64(fps), 'f', 2, 32))
	})
	txt.SetAlign(Engine.AlignLeft)

	FPSDrawer.Transform().SetPositionf(20, float32(Engine.Height)-20)
	FPSDrawer.Transform().SetScalef(20, 20)

	//Tween.CreateTween(&Tween.Tween{Target: FPSDrawer, From: []float32{1, 1}, To: []float32{100, 100},
	//	Algo: Tween.EaseInBounce, Type: Tween.Scale, Time: time.Second * 5, LoopF: Tween.None})
	/*
		{
			FPSDrawer := Engine.NewGameObject("FPS")
			FPSDrawer.Transform().SetParent2(gui)
			txt := FPSDrawer.AddComponent(Components.NewUIText(ArialFont, "")).(*Components.UIText)
			fps := FPSDrawer.AddComponent(Engine.NewFPS()).(*Engine.FPS)
			fps.SetAction(func(fps float32) {
				txt.SetString("FPS: " + strconv.FormatFloat(float64(fps), 'f', 2, 32))
			})
			txt.SetAlign(Engine.AlignLeft)

			FPSDrawer.Transform().SetPositionf(20, float32(Engine.Height)-500)
			FPSDrawer.Transform().SetScalef(20, 20)
		}
	*/

	//
	tBox := Engine.NewGameObject("TextBox")
	tBox.Transform().SetParent2(gui)

	txt2 := tBox.AddComponent(Components.NewUIText(ArialFont2, "Type your name: ")).(*Components.UIText)
	txt2.SetFocus(false)
	txt2.SetWritable(false)
	txt2.SetAlign(Engine.AlignLeft)

	tBox.Transform().SetPositionf(float32(Engine.Width)/2-txt2.Width()*20, float32(Engine.Height)/2)
	tBox.Transform().SetScalef(20, 20)
	//
	input := Engine.NewGameObject("TextBoxInput")
	input.Transform().SetParent2(gui)
	p := tBox.Transform().Position()
	p.X += txt2.Width() * 20
	input.Transform().SetPosition(p)
	input.Transform().SetScalef(20, 20)

	name := input.AddComponent(Components.NewUIText(ArialFont2, "")).(*Components.UIText)
	name.SetFocus(true)
	name.SetWritable(true)
	name.SetAlign(Engine.AlignLeft)
	//
	errLabel := Engine.NewGameObject("TextBoxInput")
	errLabel.Transform().SetParent2(gui)
	errLabel.Transform().SetPositionf(float32(Engine.Width)/2, float32(Engine.Height)/2-100)
	errLabel.Transform().SetScalef(24, 24)

	errLabelTxt := errLabel.AddComponent(Components.NewUIText(ArialFont2, "")).(*Components.UIText)
	errLabelTxt.SetFocus(false)
	errLabelTxt.SetWritable(false)
	errLabelTxt.SetAlign(Engine.AlignCenter)
	errLabelTxt.Color = Engine.Vector{1, 1, 1}
	//
	LoginButton := Engine.NewGameObject("LoginButton")
	LoginButton.Transform().SetParent2(cam)
	LoginButton.Transform().SetPositionf(float32(Engine.Width)/2, float32(Engine.Height)/2-50)
	LoginButton.AddComponent(Engine.NewSprite(button))
	LoginButton.AddComponent(Engine.NewPhysics(false, 1, 1))
	LoginButton.Physics.Shape.IsSensor = true
	LoginButton.Transform().SetScalef(50, 50)
	LoginButton.Sprite.Color = Engine.Vector{0.5, 0.5, 0.5}

	loginText := Engine.NewGameObject("LoginButtonText")
	loginText.Transform().SetParent2(LoginButton)
	loginText.Transform().SetWorldScalef(24, 24)
	loginText.Transform().SetPositionf(0, 0.1)

	var errChan chan error
	LoginButton.AddComponent(Components.NewUIButton(func() {
		if errChan == nil && Game.MyClient == nil {
			go Game.Connect(name.String(), &errChan)
			errLabelTxt.SetString("Connecting...")
		}
	}, func(enter bool) {
		if enter {
			LoginButton.Sprite.Color = Engine.Vector{0.4, 0.4, 0.4}
		} else {
			LoginButton.Sprite.Color = Engine.Vector{0.5, 0.5, 0.5}
		}
	}))

	txt2 = loginText.AddComponent(Components.NewUIText(ArialFont2, "Log in")).(*Components.UIText)
	txt2.SetFocus(false)
	txt2.SetWritable(false)
	txt2.SetAlign(Engine.AlignCenter)
	txt2.Color = Engine.Vector{1, 1, 1}
	//	

	Engine.StartCoroutine(func() {
		for {

			if errChan == nil {
				Engine.CoYieldSkip()
				continue
			}
			select {
			case loginErr := <-errChan:
				if loginErr != nil {
					errLabelTxt.SetString(loginErr.Error())
					errChan = nil
				}
			default:

			}
			Engine.CoYieldSkip()
		}
	})

	//SPACCCEEEEE
	Engine.Space.Gravity.Y = 0
	Engine.Space.Iterations = 1

	s.AddGameObject(cam)
	s.AddGameObject(background)

	fmt.Println("LoginScene loaded")
}

func (s *LoginScene) SceneBase() *Engine.SceneData {
	return s.SceneData
}

func (s *LoginScene) New() Engine.Scene {
	gs := new(LoginScene)
	gs.SceneData = Engine.NewScene("LoginScene")
	return gs
}
