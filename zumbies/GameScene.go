package zumbies

import (
	"fmt"
	"github.com/vova616/garageEngine/engine"
	"github.com/vova616/garageEngine/engine/components"
	_ "image/jpeg"
	_ "image/png"
	//"gl"  
	"math/rand"
	"strconv"
	"time"
	//"strings"
	//"math"
	//"github.com/vova616/chipmunk"
	//"github.com/vova616/chipmunk/vect"
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
	TileAtlas        *engine.ManagedAtlas
	TileIDs          []engine.ID
	PlayerID         engine.ID
)

func CheckError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func (s *GameScene) LoadTextures() {
	TileAtlas = engine.NewManagedAtlas(1024, 1024)

	e, TileID := TileAtlas.LoadImage("./data/zumbies/tiles/tile.png")
	CheckError(e)
	TileIDs = append(TileIDs, TileID)

	e, TileID = TileAtlas.LoadImage("./data/zumbies/tiles/tile2.png")
	CheckError(e)
	TileIDs = append(TileIDs, TileID)

	e, PlayerID = TileAtlas.LoadImage("./data/zumbies/player.png")
	CheckError(e)

	TileAtlas.BuildAtlas()
	TileAtlas.SetFiltering(engine.Nearest, engine.Nearest)
}

func (s *GameScene) Load() {
	s.LoadTextures()
	engine.SetTitle("Zumbies")
	rand.Seed(time.Now().UnixNano())

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

	cam.Transform().SetScalef(1, 1)

	gui := engine.NewGameObject("GUI")
	gui.Transform().SetParent2(cam)

	Layer1 := engine.NewGameObject("Layer1")
	Layer2 := engine.NewGameObject("Layer2")
	Layer3 := engine.NewGameObject("Layer3")

	s.Layer1 = Layer1
	s.Layer2 = Layer2
	s.Layer3 = Layer3

	mouse := engine.NewGameObject("Mouse")
	mouse.AddComponent(engine.NewMouse())
	mouse.Transform().SetParent2(gui)

	FPSDrawer := engine.NewGameObject("FPS")
	txt := FPSDrawer.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
	fps := FPSDrawer.AddComponent(engine.NewFPS()).(*engine.FPS)
	fps.SetAction(func(fps float64) {
		txt.SetString("FPS: " + strconv.FormatFloat(fps, 'f', 2, 32))
	})
	FPSDrawer.Transform().SetParent2(gui)
	FPSDrawer.Transform().SetPositionf(60, float32(engine.Height)-20)
	FPSDrawer.Transform().SetScalef(20, 20)

	Map := engine.NewGameObject("Map")

	uvs := engine.AnimatedUV{}
	for _, id := range TileIDs {
		uvs = append(uvs, engine.IndexUV(TileAtlas, id))
	}

	Map.AddComponent(NewMap(TileAtlas.Texture, uvs))
	Map.Transform().SetParent2(Layer1)
	Map.Transform().SetPositionf(0, 0)
	Map.Transform().SetScalef(32, 32)

	playerObject := engine.NewGameObject("Player")
	player := NewPlayer()
	playerObject.AddComponent(engine.NewSprite2(TileAtlas.Texture, engine.IndexUV(TileAtlas, PlayerID)))
	playerObject.AddComponent(player)
	playerObject.AddComponent(NewPlayerController(player))
	playerObject.AddComponent(components.NewSmoothFollow(nil, 0, 200))
	playerObject.AddComponent(engine.NewPhysics(false, 1, 1))
	playerObject.Physics.Interpolate = true
	playerObject.Transform().SetScalef(64, 64)

	//SPACCCEEEEE
	engine.Space.Gravity.Y = 0
	engine.Space.Iterations = 10

	s.AddGameObject(cam)
	s.AddGameObject(playerObject)
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
