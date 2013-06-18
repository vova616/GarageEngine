package networkOnline

import (
	"fmt"
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/bt"
	"github.com/vova616/GarageEngine/engine/components"
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
	*engine.SceneData
	Layer1 *engine.GameObject
	Layer2 *engine.GameObject
	Layer3 *engine.GameObject
}

var (
	GameSceneGeneral *GameScene
	cir              *engine.Texture
	box              *engine.Texture
)

func (s *GameScene) Load() {
	ArialFont, err := engine.NewFont("./data/Fonts/arial.ttf", 48)
	if err != nil {
		panic(err)
	}

	ArialFont2, err := engine.NewFont("./data/Fonts/arial.ttf", 24)
	if err != nil {
		panic(err)
	}
	_ = ArialFont2
	_ = ArialFont

	GameSceneGeneral = s

	s.Camera = engine.NewCamera()

	cam := engine.NewGameObject("Camera")
	cam.AddComponent(s.Camera)
	cam.AddComponent(NewCameraCtl(100))

	cam.Transform().SetScalef(1, 1)

	gui := engine.NewGameObject("GUI")

	Layer1 := engine.NewGameObject("Layer1")
	Layer2 := engine.NewGameObject("Layer2")
	Layer3 := engine.NewGameObject("Layer3")

	s.Layer1 = Layer1
	s.Layer2 = Layer2
	s.Layer3 = Layer3

	mouse := engine.NewGameObject("Mouse")
	mouse.AddComponent(NewMouseDebugger())
	mouse.AddComponent(engine.NewMouse())
	mouse.Transform().SetParent2(cam)

	FPSDrawer := engine.NewGameObject("FPS")
	txt := FPSDrawer.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
	fps := FPSDrawer.AddComponent(engine.NewFPS()).(*engine.FPS)
	txt.SetAlign(engine.AlignLeft)
	fps.SetAction(func(fps float64) {
		txt.SetString("FPS: " + strconv.FormatFloat(fps, 'f', 2, 32))
	})
	FPSDrawer.Transform().SetParent2(cam)
	FPSDrawer.Transform().SetPositionf(-float32(engine.Width)/2+10, +float32(engine.Height)/2-20)
	FPSDrawer.Transform().SetScalef(20, 20)

	{
		FPSDrawer := engine.NewGameObject("Counter")
		txt := FPSDrawer.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
		txt.SetAlign(engine.AlignLeft)
		bt.Start(func() bt.Command {
			txt.SetString(fmt.Sprintf("Bodies: %d", len(engine.Space.Bodies)))
			return bt.Restart
		})
		FPSDrawer.Transform().SetParent2(cam)
		FPSDrawer.Transform().SetPositionf(-float32(engine.Width)/2+10, +float32(engine.Height)/2-45)
		FPSDrawer.Transform().SetScalef(20, 20)
	}

	//SPACCCEEEEE
	engine.Space.Gravity.Y = -300
	engine.Space.Iterations = 10

	atlas := engine.NewManagedAtlas(512, 512)
	e := atlas.LoadGroup("./data/fire")
	if e != nil {
		fmt.Println(e)
	}
	e = atlas.LoadGroup("./data/Charecter")
	if e != nil {
		fmt.Println(e)
	}
	err, rectID := atlas.LoadImage("./data/rect.png")
	if err != nil {
		panic(err)
	}
	err, circleID := atlas.LoadImage("./data/circle.png")
	if err != nil {
		panic(err)
	}

	atlas.BuildAtlas()

	atlas.BuildMipmaps()
	atlas.SetFiltering(engine.MipMapLinearNearest, engine.Nearest)

	uvsFire, indFire := engine.AnimatedGroupUVs(atlas, "fire")
	_ = uvsFire
	_ = indFire

	clone2 := engine.NewGameObject("Sprite")
	s2 := clone2.AddComponent(engine.NewSprite3(atlas.Texture, uvsFire)).(*engine.Sprite)
	s2.BindAnimations(indFire)
	s2.AnimationSpeed = 6
	clone2.Transform().SetPositionf(775, 300)
	clone2.Transform().SetScalef(58, 58)
	clone2.Transform().SetParent2(Layer1)

	f := clone2.Clone()
	f.Transform().SetPositionf(25, 300)
	f.Transform().SetParent2(Layer1)

	box, _ = engine.LoadTexture("./data/rect.png")
	cir, _ = engine.LoadTexture("./data/circle.png")
	cir.BuildMipmaps()
	cir.SetFiltering(engine.MipMapLinearNearest, engine.Nearest)

	ball := engine.NewGameObject("Ball")
	ball.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, circleID)))
	ball.Transform().SetScalef(30, 30)
	ball.AddComponent(engine.NewPhysicsShape(false, chipmunk.NewCircle(vect.Vect{0, 0}, 15)))
	ball.Physics.Body.SetMass(10)
	ball.Physics.Body.SetMoment(ball.Physics.Shape.Moment(10))
	ball.Physics.Shape.SetFriction(0.8)
	ball.Physics.Shape.SetElasticity(0.8)

	for i := 0; i < 0; i++ {
		sprite3 := engine.NewGameObject("Sprite" + fmt.Sprint(i))
		sprite3.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, rectID)))
		sprite3.Transform().SetParent2(Layer2)
		sprite3.Transform().SetRotationf(180)
		sprite3.Transform().SetPositionf(160, 120+float32(i*31))
		sprite3.Transform().SetScalef(30, 30)

		phx := sprite3.AddComponent(engine.NewPhysics(false)).(*engine.Physics)
		phx.Shape.SetFriction(1)
		phx.Shape.SetElasticity(0.0)
		phx.Body.SetMass(1)
	}

	for i := 0; i < 200; i++ {
		sprite3 := ball.Clone()
		sprite3.Transform().SetParent2(Layer2)
		sprite3.Transform().SetPositionf(200+float32(i%4)*25, float32(i*30)+120)
	}

	floor := engine.NewGameObject("Floor")
	floor.AddComponent(engine.NewSprite(box))
	floor.AddComponent(engine.NewPhysics(true))
	floor.Transform().SetParent2(Layer2)
	floor.Transform().SetPositionf(100, -20)
	floor.Transform().SetScalef(10000, 100)
	floor.Physics.Shape.SetFriction(1)
	floor.Physics.Shape.SetElasticity(1)

	floor2 := floor.Clone()
	floor2.Transform().SetParent2(Layer2)
	floor2.Transform().SetPositionf(800, -20)
	floor2.Transform().SetRotationf(90)

	floor3 := floor2.Clone()
	floor3.Transform().SetParent2(Layer2)
	floor3.Transform().SetPositionf(0, -20)

	uvs2, ind := engine.AnimatedGroupUVs(atlas, "stand", "walk")
	_ = uvs2
	_ = ind

	sprite4 := engine.NewGameObject("Sprite")
	sprite := sprite4.AddComponent(engine.NewSprite3(atlas.Texture, uvs2)).(*engine.Sprite)
	sprite.BindAnimations(ind)
	//sprite4.AddComponent(NewSprite(sp))
	sprite.AnimationSpeed = 5
	player := sprite4.AddComponent(NewPlayerController()).(*PlayerController)
	//sprite4.AddComponent(NewRotator())
	ph := sprite4.AddComponent(engine.NewPhysics(false)).(*engine.Physics)
	player.Fire = clone2
	sprite4.Transform().SetParent2(Layer1)
	sprite4.Transform().SetPositionf(900, 100)
	sprite4.Transform().SetScalef(100, 100)
	ph.Body.SetMass(100)
	ph.Body.IgnoreGravity = false
	ph.Body.SetMoment(engine.Inf)
	ph.Shape.SetFriction(1)
	ph.Shape.SetElasticity(0)

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

	//sShadow.Sprite = bbBox
	//
	{
		sprite3 := ball.Clone()
		sprite3.Transform().SetParent2(Layer2)
		sprite3.Transform().SetPositionf(200, 120)

		joint := chipmunk.NewPivotJoint(mouse.Physics.Body, sprite3.Physics.Body)
		joint.MaxForce = 5000000
		//mouseJoint->errorBias = cpfpow(1.0f - 0.15f, 60.0f);
		engine.Space.AddConstraint(joint)
	}

	cam.Transform().SetWorldPosition(player.Transform().WorldPosition())

	s.AddGameObject(cam)
	s.AddGameObject(gui)
	s.AddGameObject(Layer1)
	s.AddGameObject(Layer2)
	s.AddGameObject(Layer3)
	//s.AddGameObject(shadowShader)

	fmt.Println("Scene loaded")
}

func (s *GameScene) New() engine.Scene {
	gs := new(GameScene)
	gs.SceneData = engine.NewScene("GameScene")
	return gs
}
