package SpaceCookies

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
	//c "github.com/vova616/chipmunk"
	//. "github.com/vova616/chipmunk/vect"
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

	cam.Transform().SetScale(NewVector3(1, 1, 1))

	gui := NewGameObject("GUI")

	Layer1 := NewGameObject("Layer1")
	Layer2 := NewGameObject("Layer2")
	Layer3 := NewGameObject("Layer3")

	s.Layer1 = Layer1
	s.Layer2 = Layer2
	s.Layer3 = Layer3

	mouse := NewGameObject("Mouse")
	mouse.AddComponent(NewMouse())
	mouse.AddComponent(NewMouseDebugger())
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
	Space.Gravity.Y = 0
	Space.Iterations = 10

	const SpaceShip = 333
	const Missle = 334

	atlas2, e := AtlasFromSheet("./data/SpaceCookies/Explosion.png", 128, 128, 6*8)
	if e != nil {
		fmt.Println(e)
	}
	atlas2.BuildAtlas()

	atlas := NewManagedAtlas(2048, 2048)
	atlas.AddImage(LoadImageQuiet("./data/SpaceCookies/Ship1.png"), SpaceShip)
	atlas.AddImage(LoadImageQuiet("./data/SpaceCookies/missile_MIRV.png"), Missle)
	e = atlas.AddGroupSheet("./data/SpaceCookies/Explosion.png", 128, 128, 6*8)
	if e != nil {
		fmt.Println(e)
	}

	atlas.BuildAtlas()

	box, _ = LoadTexture("./data/rect.png")
	cir, _ = LoadTexture("./data/circle.png")

	ship := NewGameObject("Ship")
	ship.AddComponent(NewSprite2(atlas.Texture, IndexUV(atlas, SpaceShip)))
	shipController := ship.AddComponent(NewShipController()).(*ShipController)
	ship.Transform().SetParent2(Layer2)
	ship.Transform().SetPosition(NewVector2(400, 400))
	ship.Transform().SetScale(NewVector2(100, 100))

	uvs, ind := AnimatedGroupUVs(atlas2, "Explosion")
	Explosion := NewGameObject("Explosion")
	Explosion.AddComponent(NewSprite3(atlas2.Texture, uvs))
	Explosion.Sprite.BindAnimations(ind)
	Explosion.Sprite.AnimationSpeed = 20
	Explosion.Sprite.AnimationEndCallback = func(sprite *Sprite) {
		sprite.GameObject().Destroy()
	}
	Explosion.Transform().SetScale(NewVector2(30, 30))

	missle := NewGameObject("Missle")
	missle.AddComponent(NewSprite2(atlas.Texture, IndexUV(atlas, Missle)))
	missle.AddComponent(NewPhysics(false, 10, 10))
	missle.Transform().SetScale(NewVector2(10, 10))
	m := NewMissle(30000)
	missle.AddComponent(m)
	shipController.Missle = m
	m.Explosion = Explosion

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
