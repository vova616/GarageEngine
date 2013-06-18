package login

import (
	"fmt"
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/audio"
	"github.com/vova616/GarageEngine/engine/audio/ibxm"
	"github.com/vova616/GarageEngine/engine/components"
	"github.com/vova616/GarageEngine/engine/cr"
	//"github.com/vova616/GarageEngine/engine/components/tween"
	"github.com/vova616/GarageEngine/spaceCookies/game"
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

	//"fmt"
)

type LoginScene struct {
	*engine.SceneData
}

var (
	LoginSceneGeneral *LoginScene

	backgroundTexture *engine.Texture
	button            *engine.Texture
	ArialFont         *engine.Font
	ArialFont2        *engine.Font
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
	ArialFont, e = engine.NewFont("./data/Fonts/arial.ttf", 24)
	CheckError(e)
	ArialFont.Texture.SetReadOnly()

	ArialFont2, e = engine.NewFont("./data/Fonts/arial.ttf", 24)
	CheckError(e)
	ArialFont2.Texture.SetReadOnly()

	backgroundTexture, e = engine.LoadTexture("./data/spaceCookies/background.png")
	CheckError(e)

	button, e = engine.LoadTexture("./data/spaceCookies/button.png")
	CheckError(e)
}

func (s *LoginScene) Load() {
	engine.SetTitle("Space Cookies")
	LoadTextures()

	rand.Seed(time.Now().UnixNano())

	LoginSceneGeneral = s

	s.Camera = engine.NewCamera()

	cam := engine.NewGameObject("Camera")
	cam.AddComponent(s.Camera)
	cam.Transform().SetPosition(engine.NewVector3(0, 0, -50))
	cam.Transform().SetScalef(1, 1)

	background := engine.NewGameObject("Background")
	background.AddComponent(engine.NewSprite(backgroundTexture))
	//background.Transform().SetScalef(float32(backgroung.Height()), float32(backgroung.Height()), 1)
	background.Transform().SetScalef(800, 800)
	background.Transform().SetPositionf(0, 0)
	background.Transform().SetDepth(-1)

	gui := engine.NewGameObject("GUI")
	gui.Transform().SetParent2(cam)

	mouse := engine.NewGameObject("Mouse")
	mouse.AddComponent(engine.NewMouse())
	mouse.Transform().SetParent2(gui)

	FPSDrawer := engine.NewGameObject("FPS")
	FPSDrawer.Transform().SetParent2(gui)
	txt := FPSDrawer.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
	fps := FPSDrawer.AddComponent(engine.NewFPS()).(*engine.FPS)
	fps.SetAction(func(fps float64) {
		txt.SetString("FPS: " + strconv.FormatFloat(fps, 'f', 2, 32))
	})
	txt.SetAlign(engine.AlignLeft)

	FPSDrawer.Transform().SetPositionf((float32(-engine.Width)/2)+20, (float32(engine.Height)/2)-20)
	FPSDrawer.Transform().SetScalef(20, 20)

	/*
		tween.Create(&tween.Tween{Target: FPSDrawer, From: []float32{1}, To: []float32{100},
			Algo: tween.Linear, Type: tween.Scale, Time: time.Second * 3, Loop: tween.PingPong})

			tween.Create(&tween.Tween{Target: FPSDrawer, From: []float32{400}, To: []float32{500},
				Algo: tween.Linear, Type: tween.Position, Time: time.Second * 3, Loop: tween.PingPong, Format: "y"})

			tween.Create(&tween.Tween{Target: FPSDrawer, From: []float32{0}, To: []float32{180},
				Algo: tween.Linear, Type: tween.Rotation, Time: time.Second * 6, Loop: tween.PingPong})

		txt.SetAlign(engine.AlignCenter)
	*/
	/*
		{
			FPSDrawer := engine.NewGameObject("FPS")
			FPSDrawer.Transform().SetParent2(gui)
			txt := FPSDrawer.AddComponent(components.NewUIText(ArialFont, "")).(*components.UIText)
			fps := FPSDrawer.AddComponent(engine.NewFPS()).(*engine.FPS)
			fps.SetAction(func(fps float32) {
				txt.SetString("FPS: " + strconv.FormatFloat(float64(fps), 'f', 2, 32))
			})
			txt.SetAlign(engine.AlignLeft)

			FPSDrawer.Transform().SetPositionf(20, float32(engine.Height)-500)
			FPSDrawer.Transform().SetScalef(20, 20)
		}
	*/

	//
	tBox := engine.NewGameObject("TextBox")
	tBox.Transform().SetParent2(gui)

	txt2 := tBox.AddComponent(components.NewUIText(ArialFont2, "Type your name: ")).(*components.UIText)
	txt2.SetFocus(false)
	txt2.SetWritable(false)
	txt2.SetAlign(engine.AlignLeft)
	txt2.Transform().SetDepth(1)

	tBox.Transform().SetPositionf(-txt2.Width()*20, 0)
	tBox.Transform().SetScalef(20, 20)
	//
	input := engine.NewGameObject("TextBoxInput")
	input.Transform().SetParent2(gui)
	p := tBox.Transform().Position()
	p.X += txt2.Width() * 20
	input.Transform().SetPosition(p)
	input.Transform().SetScalef(20, 20)
	input.Transform().SetDepth(1)

	name := input.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
	name.SetFocus(true)
	name.SetWritable(true)
	name.SetAlign(engine.AlignLeft)
	//
	/*
		{
			input := engine.NewGameObject("TextBoxInput")
			input.Transform().SetParent2(gui)
			p := tBox.Transform().Position()
			p.X += txt2.Width() * 20
			input.Transform().SetPosition(p)
			input.Transform().SetScalef(20, 20)

			name := input.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
			name.SetFocus(true)
			name.SetWritable(true)
			name.SetAlign(engine.AlignTopCenter)
		}
		{
			input := engine.NewGameObject("TextBoxInput")
			input.Transform().SetParent2(gui)
			p := tBox.Transform().Position()
			p.X += txt2.Width() * 20
			input.Transform().SetPosition(p)
			input.Transform().SetScalef(20, 20)

			name := input.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
			name.SetFocus(true)
			name.SetWritable(true)
			name.SetAlign(engine.AlignBottomRight)
		}
	*/
	//
	errLabel := engine.NewGameObject("TextBoxInput")
	errLabel.Transform().SetParent2(gui)
	errLabel.Transform().SetPositionf(0, -100)
	errLabel.Transform().SetScalef(24, 24)

	errLabelTxt := errLabel.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
	errLabelTxt.SetFocus(false)
	errLabelTxt.SetWritable(false)
	errLabelTxt.SetAlign(engine.AlignCenter)
	errLabelTxt.Color = engine.Color_White
	//
	LoginButton := engine.NewGameObject("LoginButton")
	LoginButton.Transform().SetParent2(cam)
	LoginButton.Transform().SetPositionf(0, -50)
	LoginButton.AddComponent(engine.NewSprite(button))
	LoginButton.AddComponent(engine.NewPhysics(false))
	LoginButton.Physics.Shape.IsSensor = true
	LoginButton.Transform().SetScalef(50, 50)
	LoginButton.Sprite.Color = engine.Color{0.5, 0.5, 0.5, 1}
	LoginButton.Transform().SetDepth(0)

	/*
		{
			LoginButton := engine.NewGameObject("LoginButton")
			LoginButton.Transform().SetParent2(cam)
			LoginButton.Transform().SetPositionf(float32(engine.Width)/2, float32(engine.Height)/2-50)
			LoginButton.AddComponent(engine.NewSprite(button))
			LoginButton.AddComponent(engine.NewPhysics(false, 1, 1))
			LoginButton.Physics.Shape.IsSensor = true
			LoginButton.Transform().SetScalef(50, 50)
			LoginButton.Sprite.Color = engine.Vector{0.5, 0.5, 0.5}
			LoginButton.Sprite.SetAlign(engine.AlignTopLeft)
			LoginButton.AddComponent(components.NewUIButton(nil, func(enter bool) {
				if enter {
					LoginButton.Sprite.Color = engine.Vector{0.4, 0.4, 0.4}
				} else {
					LoginButton.Sprite.Color = engine.Vector{0.5, 0.5, 0.5}
				}
			}))
		}
		{
			LoginButton := engine.NewGameObject("LoginButton")
			LoginButton.Transform().SetParent2(cam)
			LoginButton.Transform().SetPositionf(float32(engine.Width)/2, float32(engine.Height)/2-50)
			LoginButton.AddComponent(engine.NewSprite(button))
			LoginButton.AddComponent(engine.NewPhysics(false, 1, 1))
			LoginButton.Physics.Shape.IsSensor = true
			LoginButton.Transform().SetScalef(50, 50)
			LoginButton.Sprite.Color = engine.Vector{0.5, 0.5, 0.5}
			LoginButton.Sprite.SetAlign(engine.AlignBottomRight)
			LoginButton.AddComponent(components.NewUIButton(nil, func(enter bool) {
				if enter {
					LoginButton.Sprite.Color = engine.Vector{0.4, 0.4, 0.4}
				} else {
					LoginButton.Sprite.Color = engine.Vector{0.5, 0.5, 0.5}
				}
			}))
		}
	*/
	loginText := engine.NewGameObject("LoginButtonText")
	loginText.Transform().SetParent2(LoginButton)
	loginText.Transform().SetWorldScalef(24, 24)
	loginText.Transform().SetPositionf(0, 0.1)
	loginText.Transform().SetDepth(1)

	if game.MyClient != nil {
		game.MyClient.Socket.Close()
		game.MyClient = nil
	}

	var errChan chan error
	LoginButton.AddComponent(components.NewUIButton(func() {
		if errChan == nil && game.MyClient == nil {
			go game.Connect(name.String(), &errChan)
			errLabelTxt.SetString("Connecting...")
		}
	}, func(enter bool) {
		if enter {
			LoginButton.Sprite.Color = engine.Color{0.4, 0.4, 0.4, 1}
		} else {
			LoginButton.Sprite.Color = engine.Color{0.5, 0.5, 0.5, 1}
		}
	}))

	txt2 = loginText.AddComponent(components.NewUIText(ArialFont2, "Log in")).(*components.UIText)
	txt2.SetFocus(false)
	txt2.SetWritable(false)
	txt2.SetAlign(engine.AlignCenter)
	txt2.Color = engine.Color{1, 1, 1, 1}
	//

	cr.Start(func() {
		for {

			if errChan == nil {
				cr.YieldSkip()
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
			cr.YieldSkip()
		}
	})

	//SPACCCEEEEE
	engine.Space.Gravity.Y = 0
	engine.Space.Iterations = 1

	//cam.Transform().SetPositionf(float32(engine.Width)/2, float32(engine.Height)/2)
	//
	cam.AddComponent(audio.NewAudioListener())
	clip, e := ibxm.NewClip("./data/LoginSong.xm")
	if e != nil {
		panic(e)
	}
	music := engine.NewGameObject("Music")
	as := audio.NewAudioSource(clip)
	music.AddComponent(as)

	s.AddGameObject(cam)
	s.AddGameObject(background)
	s.AddGameObject(music)

	fmt.Println("LoginScene loaded")
}

func (s *LoginScene) New() engine.Scene {
	gs := new(LoginScene)
	gs.SceneData = engine.NewScene("LoginScene")
	return gs
}
