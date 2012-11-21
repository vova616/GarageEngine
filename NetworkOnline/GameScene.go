package NetworkOnline

import (
	"fmt"
	"github.com/vova616/GarageEngine/Engine"
	"github.com/vova616/GarageEngine/Engine/Components"
	_ "image/jpeg"
	_ "image/png"
	//"gl"  
	"strconv"
	//"time" 
	//"strings"
	//"math"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	//"image"
)

type GameScene struct {
	*Engine.SceneData
	Layer1 *Engine.GameObject
	Layer2 *Engine.GameObject
	Layer3 *Engine.GameObject
}

var (
	GameSceneGeneral *GameScene
	cir              *Engine.Texture
	box              *Engine.Texture
)

func (s *GameScene) Load() {
	ArialFont, err := Engine.NewFont("./data/Fonts/arial.ttf", 48)
	if err != nil {
		panic(err)
	}

	ArialFont2, err := Engine.NewFont("./data/Fonts/arial.ttf", 24)
	if err != nil {
		panic(err)
	}
	_ = ArialFont2
	_ = ArialFont

	GameSceneGeneral = s

	s.Camera = Engine.NewCamera()

	cam := Engine.NewGameObject("Camera")
	cam.AddComponent(s.Camera)
	cam.AddComponent(NewCameraCtl())

	cam.Transform().SetScalef(1, 1, 1)

	gui := Engine.NewGameObject("GUI")

	Layer1 := Engine.NewGameObject("Layer1")
	Layer2 := Engine.NewGameObject("Layer2")
	Layer3 := Engine.NewGameObject("Layer3")

	s.Layer1 = Layer1
	s.Layer2 = Layer2
	s.Layer3 = Layer3

	mouse := Engine.NewGameObject("Mouse")
	mouse.AddComponent(NewMouseDebugger())
	mouse.AddComponent(Engine.NewMouse())
	mouse.Transform().SetParent2(cam)

	FPSDrawer := Engine.NewGameObject("FPS")
	txt := FPSDrawer.AddComponent(Components.NewUIText(ArialFont2, "")).(*Components.UIText)
	fps := FPSDrawer.AddComponent(Engine.NewFPS()).(*Engine.FPS)
	fps.SetAction(func(fps float32) {
		txt.SetString("FPS: " + strconv.FormatFloat(float64(fps), 'f', 2, 32))
	})
	FPSDrawer.Transform().SetParent2(cam)
	FPSDrawer.Transform().SetPositionf(60, float32(Engine.Height)-20, 1)
	FPSDrawer.Transform().SetScalef(20, 20, 1)

	//SPACCCEEEEE
	Engine.Space.Gravity.Y = -300
	Engine.Space.Iterations = 30

	atlas := Engine.NewManagedAtlas(512, 512)
	e := atlas.LoadGroup("./data/fire")
	if e != nil {
		fmt.Println(e)
	}
	e = atlas.LoadGroup("./data/Charecter")
	if e != nil {
		fmt.Println(e)
	}
	atlas.LoadImage("./data/rect.png", 333)
	atlas.LoadImage("./data/circle.png", 222)

	atlas.BuildAtlas()

	uvsFire, indFire := Engine.AnimatedGroupUVs(atlas, "fire")
	_ = uvsFire
	_ = indFire

	clone2 := Engine.NewGameObject("Sprite")
	s2 := clone2.AddComponent(Engine.NewSprite3(atlas.Texture, uvsFire)).(*Engine.Sprite)
	s2.BindAnimations(indFire)
	s2.AnimationSpeed = 6
	clone2.Transform().SetPositionf(775, 300, 1)
	clone2.Transform().SetScalef(58, 58, 1)
	clone2.Transform().SetParent2(Layer1)

	f := clone2.Clone()
	f.Transform().SetPositionf(25, 300, 1)
	f.Transform().SetParent2(Layer1)

	box, _ = Engine.LoadTexture("./data/rect.png")
	cir, _ = Engine.LoadTexture("./data/circle.png")

	for i := 0; i < 0; i++ {
		sprite3 := Engine.NewGameObject("Sprite" + fmt.Sprint(i))
		sprite3.AddComponent(Engine.NewSprite2(atlas.Texture, Engine.IndexUV(atlas, 333)))
		sprite3.Transform().SetParent2(Layer2)
		sprite3.Transform().SetRotation(Engine.NewVector3(0, 0, 180))
		sprite3.Transform().SetPosition(Engine.NewVector2(160, 120+float32(i*31)))
		sprite3.Transform().SetScale(Engine.NewVector2(30, 30))

		phx := sprite3.AddComponent(Engine.NewPhysics(false, 30, 30)).(*Engine.Physics)
		phx.Shape.SetFriction(1)
		phx.Shape.SetElasticity(0.0)
		phx.Body.SetMass(1)
	}

	for i := 0; i < 2000; i++ {
		sprite3 := Engine.NewGameObject("Sprite" + fmt.Sprint(i))
		sprite3.AddComponent(Engine.NewSprite2(atlas.Texture, Engine.IndexUV(atlas, 222)))
		sprite3.Transform().SetParent2(Layer2)
		//i*31 - 220
		//+(float32(i%4))*25
		//+float32(i*30)
		sprite3.Transform().SetPosition(Engine.NewVector2(200+float32(i%4)*25, float32(i*30)+120))
		sprite3.Transform().SetScale(Engine.NewVector2(30, 30))
		phx := sprite3.AddComponent(Engine.NewPhysics2(false, chipmunk.NewCircle(vect.Vect{0, 0}, 15))).(*Engine.Physics)
		phx.Body.SetMass(10)
		phx.Body.SetMoment(phx.Shape.ShapeClass.Moment(10))

		phx.Shape.SetFriction(0.8)
		phx.Shape.SetElasticity(0.8)
	}

	floor := Engine.NewGameObject("Floor")
	floor.AddComponent(Engine.NewSprite(box))
	phx := floor.AddComponent(Engine.NewPhysics(true, 1000000, 100)).(*Engine.Physics)
	floor.Transform().SetParent2(Layer2)
	floor.Transform().SetPosition(Engine.NewVector2(100, -20))
	floor.Transform().SetScale(Engine.NewVector2(10000, 100))
	phx.Shape.SetFriction(1)
	phx.Shape.SetElasticity(1)

	floor2 := floor.Clone()
	floor2.Transform().SetParent2(Layer2)
	floor2.Transform().SetPosition(Engine.NewVector2(800, -20))
	floor2.Transform().SetRotationf(0, 0, 90)

	floor3 := floor2.Clone()
	floor3.Transform().SetParent2(Layer2)
	floor3.Transform().SetPosition(Engine.NewVector2(0, -20))

	uvs2, ind := Engine.AnimatedGroupUVs(atlas, "stand", "walk")
	_ = uvs2
	_ = ind

	sprite4 := Engine.NewGameObject("Sprite")
	sprite := sprite4.AddComponent(Engine.NewSprite3(atlas.Texture, uvs2)).(*Engine.Sprite)
	sprite.BindAnimations(ind)
	//sprite4.AddComponent(NewSprite(sp))
	sprite.AnimationSpeed = 5
	player := sprite4.AddComponent(NewPlayerController()).(*PlayerController)
	//sprite4.AddComponent(NewRotator())
	ph := sprite4.AddComponent(Engine.NewPhysics(false, 100, 100)).(*Engine.Physics)
	player.Fire = clone2
	sprite4.Transform().SetParent2(Layer1)
	sprite4.Transform().SetPositionf(900, 100, 0)
	sprite4.Transform().SetScalef(100, 100, 0)
	ph.Body.SetMass(100)
	ph.Body.IgnoreGravity = false
	ph.Body.SetMoment(Engine.Inf)
	ph.Shape.SetFriction(1)
	ph.Shape.SetElasticity(0)
	sprite.Border = true
	sprite.BorderSize = 0

	/*
		floor = NewGameObject("Box")
		bbBox := NewSprite(box)
		floor.AddComponent(bbBox)
		phx = floor.AddComponent(NewPhysics(true, 500, 500)).(*Physics)
		floor.Transform().SetParent2(Layer2)
		floor.Transform().SetPosition(NewVector2(900, 200))
		floor.Transform().SetScale(NewVector2(500, 500))
		phx.Shape.SetFriction(1)
		phx.Shape.SetElasticity(1)
		phx.Shape.IsSensor = true
		//phx.Shape.Friction = 1
		_ = phx
	*/
	shadowShader := Engine.NewGameObject("Shadow")
	sCam := Engine.NewCamera()
	shadowShader.AddComponent(sCam)
	sShadow := Engine.NewShadowShader(s.Camera)
	cam.AddComponent(sShadow)

	//sShadow.Sprite = bbBox

	s.AddGameObject(cam)
	s.AddGameObject(gui)
	s.AddGameObject(Layer1)
	s.AddGameObject(Layer2)
	s.AddGameObject(Layer3)
	//s.AddGameObject(shadowShader)

	fmt.Println("Scene loaded")
}

func (s *GameScene) SceneBase() *Engine.SceneData {
	return s.SceneData
}

func (s *GameScene) New() Engine.Scene {
	gs := new(GameScene)
	gs.SceneData = Engine.NewScene("GameScene")
	return gs
}
