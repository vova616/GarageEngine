package NetworkOnline

import (
	. "github.com/vova616/GarageEngine/Engine"
	. "github.com/vova616/GarageEngine/Engine/Components"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	//"gl"  
	"strconv"
	//"time" 
	//"strings"
	//"math"
	c "github.com/vova616/chipmunk"
	. "github.com/vova616/chipmunk/vect"
)

type GameScene struct {
	*SceneData
	Camera *GameObject
	Layer1 *GameObject
	Layer2 *GameObject
}

var (
	GameSceneGeneral *GameScene
	cir *Texture
	box *Texture
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
	_ = ArialFont2

	GameSceneGeneral = s
	s.Camera = NewGameObject("Camera")
	s.Camera.AddComponent(NewCamera()) 
	gui := NewGameObject("GUI")
	gui.AddComponent(NewGUI())
	//gui.Transform().SetParent2(s.Camera)
	Layer1 := NewGameObject("Layer1")
	Layer2 := NewGameObject("Layer2")
	Layer3 := NewGameObject("Layer3")
	Layer1.Transform().SetParent2(s.Camera)
	Layer2.Transform().SetParent2(s.Camera)
	Layer3.Transform().SetParent2(s.Camera)
	 
	
	//Layer2.Transform().SetScale(NewVector3(0.5,0.5,0.5))
	//s.Camera.Transform().Translate2(100,0,0)

	s.Layer1 = Layer1
	s.Layer2 = Layer2
 
	mouse := NewGameObject("Mouse")
	
	mouse.AddComponent(NewMouse())
	mouse.AddComponent(NewMouseDebugger())
	mouse.Transform().SetParent2(gui)
	

	FPSDrawer := NewGameObject("FPS")
	txt := FPSDrawer.AddComponent(NewUIText(ArialFont2, "")).(*UIText)
	fps := FPSDrawer.AddComponent(NewFPS()).(*FPS)
	fps.SetAction(func(fps float32) {
		txt.SetString("FPS: " + strconv.FormatFloat(float64(fps), 'f', 2, 32))
	})
	FPSDrawer.Transform().SetParent2(gui)
	FPSDrawer.Transform().SetPosition(NewVector2(60, float32(Height)-20))
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
	
	//SPACCCEEEEE
	Space.Gravity.Y = -300
	Space.Iterations = 10
	
	//Space.Gravity.X = 0 

	atlas := NewManagedAtlas(1024, 512)

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
	clone2.Transform().SetPosition(NewVector2(775, 300))
	clone2.Transform().SetScale(NewVector2(58, 58))
	clone2.Transform().SetParent2(Layer1)

	 
	f := clone2.Clone()
	f.Transform().SetPosition(NewVector2(25, 300))
	f.Transform().SetParent2(Layer1)

	box, _ = LoadTexture("./data/rect.png")
	//spc := NewSprite(sp)
	
	cir, _ = LoadTexture("./data/circle.png")
	
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
		sprite3.AddComponent(NewSprite(box))
		sprite3.Transform().SetParent2(Layer2)
		sprite3.Transform().SetRotation(NewVector3(0, 0, 180))
		sprite3.Transform().SetPosition(NewVector2(160, 120+float32(i*31)))
		sprite3.Transform().SetScale(NewVector2(30, 30))
		
		phx := sprite3.AddComponent(NewPhysics(false,30,30)).(*Physics)
		phx.Shape.SetFriction(1)
		phx.Shape.SetElasticity(0.0)
		phx.Body.SetMass(1)
	} 
  
	for i := 0; i > 0; i-- {
		sprite3 := NewGameObject("Sprite" + fmt.Sprint(i))
		sprite3.AddComponent(NewSprite(cir))
		sprite3.Transform().SetParent2(Layer2)
		sprite3.Transform().SetPosition(NewVector2(200+(float32(i%4))*25, 120+float32(i*20)))
		sprite3.Transform().SetScale(NewVector2(30, 30))
		phx := sprite3.AddComponent(NewPhysics2(false, c.NewCircle(Vect{0,0},Float(15)))).(*Physics)
		phx.Shape.SetFriction(0.2)
		phx.Shape.SetElasticity(0.8)
	}
	 
	floor := NewGameObject("Floor")
	floor.AddComponent(NewSprite(box))
	phx := floor.AddComponent(NewPhysics(true, 1000000, 100)).(*Physics)
	floor.Transform().SetParent2(Layer2)
	floor.Transform().SetPosition(NewVector2(100, -20))
	floor.Transform().SetRotation(NewVector3(0, 0, 180))
	floor.Transform().SetScale(NewVector2(1000000, 100))
	phx.Shape.SetFriction(1)
	phx.Shape.SetElasticity(1)
	//phx.Shape.Friction = 1
	_ = phx
	

	floor = NewGameObject("Floor2")
	floor.AddComponent(NewSprite(box))
	phx = floor.AddComponent(NewPhysics(true, 100, 10000)).(*Physics)
	floor.Transform().SetParent2(Layer2)
	floor.Transform().SetPosition(NewVector2(800, -20))
	floor.Transform().SetScale(NewVector2(100, 10000))
	phx.Shape.SetFriction(1)
	phx.Shape.SetElasticity(1)
	//phx.Shape.Friction = 1
	_ = phx 
	
	floor = NewGameObject("Floor2")
	floor.AddComponent(NewSprite(box))
	phx = floor.AddComponent(NewPhysics(true, 100, 10000)).(*Physics)
	floor.Transform().SetParent2(Layer2)
	floor.Transform().SetPosition(NewVector2(0, -20))
	floor.Transform().SetScale(NewVector2(100, 10000))
	phx.Shape.SetFriction(1)
	phx.Shape.SetElasticity(1)
	//phx.Shape.Friction = 1
	_ = phx
	
 
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
	

	
	s.AddGameObject(s.Camera)
	s.AddGameObject(gui)
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
