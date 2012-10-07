package NetworkOnline

import (
	"fmt"
	. "github.com/vova616/GarageEngine/Engine"
	. "github.com/vova616/GarageEngine/Engine/Components"
	_ "image/jpeg"
	_ "image/png"
	//"gl"  
	"strconv"
	//"time" 
	//"strings"
	//"math"
	c "github.com/vova616/chipmunk"
	. "github.com/vova616/chipmunk/vect"
	//"image"
)

type GameScene struct {
	*SceneData
	Layer1 *GameObject
	Layer2 *GameObject
	Layer3 *GameObject
}

var (
	GameSceneGeneral *GameScene
	cir              *Texture
	box              *Texture
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

	cam.Transform().SetScale(NewVector3(1, 1, 1))

	gui := NewGameObject("GUI")

	Layer1 := NewGameObject("Layer1")
	Layer2 := NewGameObject("Layer2")
	Layer3 := NewGameObject("Layer3")

	s.Layer1 = Layer1
	s.Layer2 = Layer2
	s.Layer3 = Layer3

	mouse := NewGameObject("Mouse")
	mouse.AddComponent(NewMouseDebugger())
	mouse.AddComponent(NewMouse())
	mouse.Transform().SetParent2(cam)

	FPSDrawer := NewGameObject("FPS")
	txt := FPSDrawer.AddComponent(NewUIText(ArialFont2, "")).(*UIText)
	fps := FPSDrawer.AddComponent(NewFPS()).(*FPS)
	fps.SetAction(func(fps float32) {
		txt.SetString("FPS: " + strconv.FormatFloat(float64(fps), 'f', 2, 32))
	})
	FPSDrawer.Transform().SetParent2(cam)
	FPSDrawer.Transform().SetPosition(NewVector2(60, float32(Height)-20))
	FPSDrawer.Transform().SetScale(NewVector2(20, 20))

	//SPACCCEEEEE
	Space.Gravity.Y = -300
	Space.Iterations = 30

	atlas := NewManagedAtlas(512, 512)
	e := atlas.AddGroup("./data/fire")
	if e != nil {
		fmt.Println(e)
	}
	e = atlas.AddGroup("./data/Charecter")
	if e != nil {
		fmt.Println(e)
	}
	atlas.AddImage(LoadImageQuiet("./data/rect.png"), 333)
	atlas.AddImage(LoadImageQuiet("./data/circle.png"), 222)

	atlas.BuildAtlas()

	uvsFire, indFire := AnimatedGroupUVs(atlas, "fire")
	_ = uvsFire
	_ = indFire

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

	box, _ = LoadTexture("./data/rect.png")
	cir, _ = LoadTexture("./data/circle.png")

	for i := 0; i < 0; i++ {
		sprite3 := NewGameObject("Sprite" + fmt.Sprint(i))
		sprite3.AddComponent(NewSprite2(atlas.Texture, IndexUV(atlas, 333)))
		sprite3.Transform().SetParent2(Layer2)
		sprite3.Transform().SetRotation(NewVector3(0, 0, 180))
		sprite3.Transform().SetPosition(NewVector2(160, 120+float32(i*31)))
		sprite3.Transform().SetScale(NewVector2(30, 30))

		phx := sprite3.AddComponent(NewPhysics(false, 30, 30)).(*Physics)
		phx.Shape.SetFriction(1)
		phx.Shape.SetElasticity(0.0)
		phx.Body.SetMass(1)
	}

	for i := 0; i < 2000; i++ {
		sprite3 := NewGameObject("Sprite" + fmt.Sprint(i))
		sprite3.AddComponent(NewSprite2(atlas.Texture, IndexUV(atlas, 222)))
		sprite3.Transform().SetParent2(Layer2)
		//i*31 - 220
		//+(float32(i%4))*25
		//+float32(i*30)
		sprite3.Transform().SetPosition(NewVector2(200+float32(i%4)*25, float32(i*30)+120))
		sprite3.Transform().SetScale(NewVector2(30, 30))
		phx := sprite3.AddComponent(NewPhysics2(false, c.NewCircle(Vect{0, 0}, 15))).(*Physics)
		phx.Body.SetMass(10)
		phx.Body.SetMoment(phx.Shape.ShapeClass.Moment(10))

		phx.Shape.SetFriction(0.8)
		phx.Shape.SetElasticity(0.8)
	}

	floor := NewGameObject("Floor")
	floor.AddComponent(NewSprite(box))
	phx := floor.AddComponent(NewPhysics(true, 1000000, 100)).(*Physics)
	floor.Transform().SetParent2(Layer2)
	floor.Transform().SetPosition(NewVector2(100, -20))
	floor.Transform().SetScale(NewVector2(10000, 100))
	phx.Shape.SetFriction(1)
	phx.Shape.SetElasticity(1)

	floor2 := floor.Clone()
	floor2.Transform().SetParent2(Layer2)
	floor2.Transform().SetPosition(NewVector2(800, -20))
	floor2.Transform().SetRotationf(0, 0, 90)

	floor3 := floor2.Clone()
	floor3.Transform().SetParent2(Layer2)
	floor3.Transform().SetPosition(NewVector2(0, -20))

	uvs2, ind := AnimatedGroupUVs(atlas, "stand", "walk")
	_ = uvs2
	_ = ind

	sprite4 := NewGameObject("Sprite")
	sprite := sprite4.AddComponent(NewSprite3(atlas.Texture, uvs2)).(*Sprite)
	sprite.BindAnimations(ind)
	//sprite4.AddComponent(NewSprite(sp))
	sprite.AnimationSpeed = 5
	player := sprite4.AddComponent(NewPlayerController()).(*PlayerController)
	//sprite4.AddComponent(NewRotator())
	ph := sprite4.AddComponent(NewPhysics(false, 100, 100)).(*Physics)
	player.Fire = clone2
	sprite4.Transform().SetParent2(Layer1)
	sprite4.Transform().SetPositionf(900, 100, 0)
	sprite4.Transform().SetScalef(100, 100, 0)
	ph.Body.SetMass(100)
	ph.Body.IgnoreGravity = false
	ph.Body.SetMoment(Inf)
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
	shadowShader := NewGameObject("Shadow")
	sCam := NewCamera()
	shadowShader.AddComponent(sCam)
	sShadow := NewShadowShader(s.Camera)
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

func (s *GameScene) SceneBase() *SceneData {
	return s.SceneData
}

func (s *GameScene) New() Scene {
	gs := new(GameScene)
	gs.SceneData = NewScene("GameScene")
	return gs
}
