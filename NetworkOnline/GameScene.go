package NetworkOnline

import (
	. "../Engine"
	. "../Engine/Components"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	//"gl"  
	"strconv"
	//"time" 
	//"strings"
	//"math"
	c "chipmunk"
	. "chipmunk/vect"
)

type GameScene struct {
	*SceneData
	Camera *GameObject
	Layer1 *GameObject
}

var (
	GameSceneGeneral *GameScene
)

func (s *GameScene) Load() {
	ArialFont, err := NewFont("./data/Fonts/arial.ttf", 48)
	if err != nil {
		panic(err)
	}

	ArialFont2, err := NewFont("./data/Fonts/arial.ttf", 24)
	if err != nil {
		panic(err)
	}

	GameSceneGeneral = s
	s.Camera = NewGameObject("Camera")
	//Camera.AddComponent(NewController()) 
	gui := NewGameObject("GUI")
	gui.AddComponent(NewGUI())
	gui.Transform().SetParent2(s.Camera)
	Layer1 := NewGameObject("Layer1")
	Layer2 := NewGameObject("Layer2")
	Layer3 := NewGameObject("Layer3")
	Layer1.Transform().SetParent2(s.Camera)
	Layer2.Transform().SetParent2(s.Camera)
	Layer3.Transform().SetParent2(s.Camera)
	
	//Layer2.Transform().SetScale(NewVector3(0.5,0.5,0.5))
	//s.Camera.Transform().Translate2(100,0,0)

	s.Layer1 = Layer1
 
	mouse := NewGameObject("Mouse")
	mouse.AddComponent(NewMouse())
	mouse.Transform().SetParent2(gui)

	FPSDrawer := NewGameObject("FPS")
	txt := FPSDrawer.AddComponent(NewUIText(ArialFont2, "")).(*UIText)
	fps := FPSDrawer.AddComponent(NewFPS()).(*FPS)
	fps.SetAction(func(fps float32) {
		txt.SetString("FPS: " + strconv.FormatFloat(float64(fps), 'f', 2, 32))
	})
	FPSDrawer.Transform().SetParent2(gui)
	FPSDrawer.Transform().SetPosition(NewVector2(60, 460))
	FPSDrawer.Transform().SetScale(NewVector2(20, 20))
	
	Texts := NewGameObject("Fuck")
	//Texts.AddComponent(NewRotator(Texts))
	Texts.Transform().SetPosition(NewVector2(100, 400))
	//Texts.Transform().SetParent2(Layer1)
	Texts.AddComponent(NewUIText(ArialFont, "Engine test"))
	Texts.Transform().SetScale(NewVector2(48, 48))

	ArialUIText := NewGameObject("Hello")
	ArialUIText.AddComponent(NewUIText(ArialFont, ArialUIText.Name()))
	//ArialUIText.AddComponent(NewUIText(ArialUIText,ArialFont, "  Marrr"))	
	//ArialUIText.AddComponent(NewRotator())
	//ArialUIText.Transform().SetParent(Texts.Transform())
	ArialUIText.Transform().SetPosition(NewVector2(300, 300))
	ArialUIText.Transform().SetScale(NewVector2(48, 48))
	//ArialUIText.Transform().SetParent2(Layer1)

	clone2 := NewGameObject("Sprite")

	//clone2.Transform().SetParent(Texts.Transform())
	//clone2.AddComponent(NewRotator(clone2))
	
	//GRAAVVIITTYYY
	Space.Gravity.Y = -300
	
	
	//Space.Gravity.X = 0 

	atlas := NewManagedAtlas(1024, 512)

	atlas.AddImage(LoadImageQuiet("./data/n0.png"), 0)
	atlas.AddImage(LoadImageQuiet("./data/n1.png"), 1)
	atlas.AddImage(LoadImageQuiet("./data/n2.png"), 2)
	atlas.AddImage(LoadImageQuiet("./data/n3.png"), 3)
	atlas.AddImage(LoadImageQuiet("./data/n4.png"), 4)
	atlas.AddImage(LoadImageQuiet("./data/n5.png"), 5)
	atlas.AddImage(LoadImageQuiet("./data/n6.png"), 6)
	atlas.AddImage(LoadImageQuiet("./data/n7.png"), 7)
	
	

	atlas.AddImage(LoadImageQuiet("./data/fire0.png"), 10)
	atlas.AddImage(LoadImageQuiet("./data/fire1.png"), 11)
	atlas.AddImage(LoadImageQuiet("./data/fire2.png"), 12)
	atlas.AddImage(LoadImageQuiet("./data/fire3.png"), 13)
	atlas.AddImage(LoadImageQuiet("./data/fire4.png"), 14)
	atlas.AddImage(LoadImageQuiet("./data/fire5.png"), 15)
	atlas.AddImage(LoadImageQuiet("./data/fire6.png"), 16)
	atlas.AddImage(LoadImageQuiet("./data/fire7.png"), 17)

	atlas.BuildAtlas()
	//atlas.Options(gl.NEAREST, gl.CLAMP_TO_EDGE)
	uvs := AnimatedUVs(atlas, 0, 1, 2, 3, 4, 5, 6, 7)
	_ = uvs
	uvs2 := AnimatedUVs(atlas, 10, 11, 12, 13, 14, 15, 16, 17)
	_ = uvs2

	s2 := clone2.AddComponent(NewSprite3(atlas.Texture, uvs2)).(*Sprite)
	s2.AnimationSpeed = 6
	//clone2.AddComponent(NewSprite2(clone2, atlas.Texture, NewUV(0,0,1,1,1)))
	clone2.Transform().SetPosition(NewVector2(600, 300))
	clone2.Transform().SetScale(NewVector2(58, 58))
	clone2.Transform().SetParent2(Layer1)

	sp, _ := LoadTexture("./data/rect.png")
	//spc := NewSprite(sp)
	
	cir, _ := LoadTexture("./data/circle.png")
	
/*
	sprite2 := NewGameObject("Sprite")
	sprite2.AddComponent(NewSprite(sp))
	//sprite2.AddComponent(NewController(sprite2))
	sprite2.Transform().SetParent2(Layer1)
	sprite2.Transform().SetPosition(NewVector2(100, 300))
	sprite2.Transform().SetScale(NewVector2(40, 40))
	//sprite2.AddComponent(NewPhysics(true))
/*
	spritet := NewGameObject("EngineText")
	spritet.AddComponent(NewUIText(ArialFont, "Engine Test"))
	spritet.Transform().SetParent2(Layer1)
	spritet.Transform().SetPosition(NewVector2(201, 200))
	spritet.Transform().SetScale(NewVector2(40, 40))
	//spritet.AddComponent(NewRotator())
*/
	for i := 0; i < 0; i++ {
		sprite3 := NewGameObject("Sprite" + fmt.Sprint(i))
		sprite3.AddComponent(NewSprite(sp))
		sprite3.Transform().SetParent2(Layer2)
		sprite3.Transform().SetRotation(NewVector3(0, 0, 180))
		sprite3.Transform().SetPosition(NewVector2(200+float32(i/4)*35, 45+float32(i%4)*30))
		sprite3.Transform().SetScale(NewVector2(30, 30))
		
		phx := sprite3.AddComponent(NewPhysics(false)).(*Physics)
		phx.Shape.SetFriction(1)
		phx.Shape.SetElasticity(0.5)
	}

	for i := 0; i < 3; i++ {
		sprite3 := NewGameObject("Sprite" + fmt.Sprint(i))
		sprite3.AddComponent(NewSprite(cir))
		sprite3.Transform().SetParent2(Layer2)
		sprite3.Transform().SetPosition(NewVector2(200+float32(int(i/20))*35, 120+float32(i%20)*40))
		sprite3.Transform().SetScale(NewVector2(30, 30))
		phx := sprite3.AddComponent(NewPhysics2(false, c.NewCircle(Vect{0,0},Float(15)))).(*Physics)
		phx.Body.SetMoment(Inf)
		phx.Shape.SetFriction(1)
		phx.Shape.SetElasticity(0)
	}

	//Layer2.Transform().Position.Y += 200

	atlas, e := AtlasLoadDirectory("./data/Charecter")
	if atlas == nil {
		panic(e)
	}
	if e != nil {
		fmt.Println(e)
	}
	atlas.BuildAtlas()

	uvs2, ind := AnimatedGroupUVs(atlas, "stand", "walk")
	_ = uvs2 
	_ = ind
	//	/fmt.Println(ind)
	/*
		sprite4 := NewGameObject("Sprite")
		sprite := sprite4.AddComponent(NewSprite3(atlas.Texture, uvs2)).(*Sprite)
		sprite.BindAnimations(ind)
		//sprite4.AddComponent(NewSprite(sp))
		sprite.AnimationSpeed = 5
		player := sprite4.AddComponent(NewPlayerController()).(*PlayerController)
		//sprite4.AddComponent(NewRotator())
		ph := sprite4.AddComponent(NewPhysics(false)).(*Physics)
		player.Fire = clone2
		sprite4.Transform().SetParent2(Layer1)
		sprite4.Transform().SetPosition(NewVector2(40, 300))
		sprite4.Transform().SetScale(NewVector2(50, 50))
		ph.Body.SetMass(1)
		ph.Body.IgnoreGravity = false
	//	ph.Body.SetInertia(0)
	*/
	
	floor := NewGameObject("Floor")
	floor.AddComponent(NewSprite(sp))
	phx := floor.AddComponent(NewPhysics(true)).(*Physics)
	floor.Transform().SetParent2(Layer2)
	floor.Transform().SetPosition(NewVector2(100, -20))
	floor.Transform().SetRotation(NewVector3(0, 0, 180))
	floor.Transform().SetScale(NewVector2(1000000, 100))
	phx.Shape.SetFriction(1)
	//phx.Shape.Friction = 1
	_ = phx
	
	floor = NewGameObject("Floor2")
	floor.AddComponent(NewSprite(sp))
	phx = floor.AddComponent(NewPhysics(true)).(*Physics)
	floor.Transform().SetParent2(Layer2)
	floor.Transform().SetPosition(NewVector2(100, -20))
	floor.Transform().SetScale(NewVector2(100, 10000))
	phx.Shape.SetFriction(1)
	//phx.Shape.Friction = 1
	_ = phx
	
	floor = NewGameObject("Floor2")
	floor.AddComponent(NewSprite(sp))
	phx = floor.AddComponent(NewPhysics(true)).(*Physics)
	floor.Transform().SetParent2(Layer2)
	floor.Transform().SetPosition(NewVector2(500, -20))
	floor.Transform().SetScale(NewVector2(100, 10000))
	phx.Shape.SetFriction(1)
	//phx.Shape.Friction = 1
	_ = phx
	
	s.AddGameObject(s.Camera)
	fmt.Println("Scene loaded")
}

func (s *GameScene) SceneBase() *SceneData {
	return s.SceneData
}

func (s *GameScene) New() Scene {
	gs := new(GameScene)
	gs.SceneData = NewScene("GameScene")
	return gs
}
