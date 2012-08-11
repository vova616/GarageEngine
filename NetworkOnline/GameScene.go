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
	Layer1 *GameObject
	Layer2 *GameObject
	Layer3 *GameObject
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
	_ = ArialFont

	GameSceneGeneral = s
	
	
	s.Camera = NewCamera()
	
	cam := NewGameObject("Camera") 
	cam.AddComponent(s.Camera) 
	cam.AddComponent(NewCameraCtl()) 
	
	cam.Transform().SetScale(NewVector3(0.8,0.8,0.8))
	
	gui := NewGameObject("GUI")
	
	Layer1 := NewGameObject("Layer1")
	Layer2 := NewGameObject("Layer2")
	Layer3 := NewGameObject("Layer3")
	 
	s.Layer1 = Layer1
	s.Layer2 = Layer2
	s.Layer3 = Layer3
 
	mouse := NewGameObject("Mouse")
	mouse.AddComponent(NewMouseDebugger())
	

	FPSDrawer := NewGameObject("FPS")
	txt := FPSDrawer.AddComponent(NewUIText(ArialFont2, "")).(*UIText)
	fps := FPSDrawer.AddComponent(NewFPS()).(*FPS)
	fps.SetAction(func(fps float32) {
		txt.SetString("FPS: " + strconv.FormatFloat(float64(fps), 'f', 2, 32))
	}) 
	FPSDrawer.Transform().SetParent2(gui)
	FPSDrawer.Transform().SetPosition(NewVector2(60, float32(Height)-20))
	FPSDrawer.Transform().SetScale(NewVector2(20, 20))
	

	//SPACCCEEEEE
	Space.Gravity.Y = -700
	Space.Iterations = 10
	

	atlas, e := AtlasLoadDirectory("./data/fire")
	if atlas == nil {
		panic(e)
	}
	if e != nil {
		fmt.Println(e)
	}
	atlas.BuildAtlas()
 
	uvsFire, indFire := AnimatedGroupUVs(atlas, "fire")
	_ = uvsFire 
	_ = indFire

 	box, _ = LoadTexture("./data/rect.png")
	//spc := NewSprite(sp)
	
	cir, _ = LoadTexture("./data/circle.png")
	
	
	
	clone2 := NewGameObject("Sprite")
	s2 := clone2.AddComponent(NewSprite3(atlas.Texture, uvsFire)).(*Sprite)
	s2.BindAnimations(indFire)
	s2.AnimationSpeed = 6
	clone2.Transform().SetPosition(NewVector2(775, 300))
	clone2.Transform().SetScale(NewVector2(58, 58))
	clone2.Transform().SetParent2(Layer1)

	 
	f := clone2.Clone()
	f.Transform().SetPosition(NewVector2(25, 300))
	f.Transform().SetParent2(Layer1)

	
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
  
	for i := 0; i < 0; i++ {
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

	atlas, e = AtlasLoadDirectory("./data/Charecter")
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
	
	sprite4 := NewGameObject("Sprite")
	sprite := sprite4.AddComponent(NewSprite3(atlas.Texture, uvs2)).(*Sprite)
	sprite.BindAnimations(ind)
	//sprite4.AddComponent(NewSprite(sp))
	sprite.AnimationSpeed = 5
	player := sprite4.AddComponent(NewPlayerController()).(*PlayerController)
	//sprite4.AddComponent(NewRotator())
	ph := sprite4.AddComponent(NewPhysics(false,100,100)).(*Physics)
	player.Fire = clone2
	sprite4.Transform().SetParent2(Layer1)
	sprite4.Transform().SetPosition(NewVector2(900, 100))
	sprite4.Transform().SetScale(NewVector2(100, 100))
	ph.Body.SetMass(1)
	ph.Body.IgnoreGravity = false
	ph.Body.SetMoment(Inf)
	ph.Shape.SetFriction(0.5)
	ph.Shape.SetElasticity(0)
	
	

	s.AddGameObject(cam)
	s.AddGameObject(mouse)
	s.AddGameObject(gui)
	s.AddGameObject(Layer1)
	s.AddGameObject(Layer2)
	s.AddGameObject(Layer3)
	
	
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
