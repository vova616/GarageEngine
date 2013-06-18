package zumbies

import (
	"fmt"
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/components"
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

	MainPlayer *Player
	Layers     []*Map
	Mouse      *engine.GameObject
)

func CheckError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func (s *GameScene) LoadTextures() {
	TileAtlas = engine.NewManagedAtlas(512, 512)

	e, TileID := TileAtlas.LoadImage("./data/zumbies/tiles/tile.png")
	CheckError(e)
	TileIDs = append(TileIDs, TileID)

	e, TileID = TileAtlas.LoadImage("./data/zumbies/tiles/tile2.png")
	CheckError(e)
	TileIDs = append(TileIDs, TileID)

	e, TileID = TileAtlas.LoadImage("./data/zumbies/tiles/tile3.png")
	CheckError(e)
	TileIDs = append(TileIDs, TileID)

	e, TileID = TileAtlas.LoadImage("./data/zumbies/tiles/tile4.png")
	CheckError(e)
	TileIDs = append(TileIDs, TileID)

	e, TileID = TileAtlas.LoadImage("./data/zumbies/tiles/tile5.png")
	CheckError(e)
	TileIDs = append(TileIDs, TileID)

	e, TileID = TileAtlas.LoadImage("./data/zumbies/tiles/tile6.png")
	CheckError(e)
	TileIDs = append(TileIDs, TileID)

	e, PlayerID = TileAtlas.LoadImage("./data/zumbies/zombie.png")
	CheckError(e)

	TileAtlas.BuildAtlas()
	TileAtlas.SetFiltering(engine.Nearest, engine.Nearest)
}

func (s *GameScene) Load() {
	s.LoadTextures()
	engine.SetTitle("Zumbies")
	rand.Seed(time.Now().UnixNano())

	Layers = make([]*Map, 0, 10)

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
	gui.Transform().SetParent2(cam)

	Layer1 := engine.NewGameObject("Layer1")
	Layer2 := engine.NewGameObject("Layer2")
	Layer3 := engine.NewGameObject("Layer3")

	s.Layer1 = Layer1
	s.Layer2 = Layer2
	s.Layer3 = Layer3

	Mouse = engine.NewGameObject("Mouse")
	Mouse.AddComponent(engine.NewMouse())
	Mouse.Transform().SetParent2(gui)

	FPSDrawer := engine.NewGameObject("FPS")
	txt := FPSDrawer.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
	fps := FPSDrawer.AddComponent(engine.NewFPS()).(*engine.FPS)
	fps.SetAction(func(fps float64) {
		txt.SetString("FPS: " + strconv.FormatFloat(fps, 'f', 2, 32))
	})
	FPSDrawer.Transform().SetParent2(gui)
	FPSDrawer.Transform().SetPositionf(60, float32(engine.Height)-20)
	FPSDrawer.Transform().SetScalef(20, 20)

	{
		Map := engine.NewGameObject("Map")
		uvs := engine.AnimatedUV{}
		for _, id := range TileIDs {
			uvs = append(uvs, engine.IndexUV(TileAtlas, id))
		}
		m := NewMap(TileAtlas.Texture, uvs)
		Map.AddComponent(m)
		Map.Transform().SetPositionf(0, 0)
		Map.Transform().SetScalef(32, 32)

		ca4 := Tile(0).SetCollision(CollisionAll)
		l4 := Tile(0).SetLayerConnection(true)

		m.Tiles = []Tile{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 2, ca4, ca4, ca4, ca4, ca4, ca4, ca4, ca4, ca4, ca4, ca4, ca4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, l4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 2, ca4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 1, 1, 2, ca4, ca4, ca4, ca4, ca4, ca4, ca4, ca4, ca4, ca4, ca4, ca4, ca4, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

		m.Width = int(60)
		m.Height = int(23)

		Layers = append(Layers, m)
	}

	{
		Map := engine.NewGameObject("Map")
		uvs := engine.AnimatedUV{}
		for _, id := range TileIDs {
			uvs = append(uvs, engine.IndexUV(TileAtlas, id))
		}
		m := NewMap(TileAtlas.Texture, uvs)
		Map.AddComponent(m)
		Map.Transform().SetPositionf(0, 0)
		Map.Transform().SetScalef(32, 32)

		u5 := Tile(4).SetSide2(SideUp).SetType2(5)
		d5 := Tile(4).SetSide2(SideDown).SetType2(5)
		l5 := Tile(4).SetSide2(SideLeft).SetType2(5)
		r5 := Tile(4).SetSide2(SideRight).SetType2(5)

		dc5 := Tile(5).SetSide(SideDown).SetCollision(CollisionAll)

		l4 := Tile(4).SetLayerConnection(true)

		u6 := Tile(4).SetSide2(SideUp).SetType2(6)
		d6 := Tile(4).SetSide2(SideDown).SetType2(6)
		l6 := Tile(4).SetSide2(SideLeft).SetType2(6)
		r6 := Tile(4).SetSide2(SideRight).SetType2(6)

		_, _, _, _ = u5, d5, l5, r5

		m.Tiles = []Tile{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l6, u5, u5, u5, u5, u5, u5, u5, u5, u5, u5, u5, u6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, l4, dc5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, l5, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, r5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, d6, d5, d5, d5, d5, d5, d5, d5, d5, d5, d5, d5, r6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

		m.Width = int(60)
		m.Height = int(23)

		Layers = append(Layers, m)

	}

	for i := len(Layers) - 1; i >= 0; i-- {
		Layers[i].Layer = int(1 << uint(i))
		Layers[i].Transform().SetParent2(Layer1)
	}

	playerObject := engine.NewGameObject("Player")
	player := NewPlayer()
	playerObject.AddComponent(engine.NewSprite2(TileAtlas.Texture, engine.IndexUV(TileAtlas, PlayerID)))
	playerObject.AddComponent(player)
	playerObject.AddComponent(NewPlayerController(player))
	playerObject.AddComponent(components.NewSmoothFollow(nil, 0, 200))
	playerObject.AddComponent(engine.NewPhysicsCircle(false))
	playerObject.Physics.Interpolate = true
	playerObject.Physics.Body.SetMoment(engine.Inf)
	playerObject.Transform().SetScalef(64, 64)
	playerObject.Transform().SetWorldPositionf(159.99995, 32)
	playerObject.Transform().SetDepth(1)
	MainPlayer = player

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
