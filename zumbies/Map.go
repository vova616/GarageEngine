package zumbies

import (
	"github.com/vova616/garageEngine/engine"
	//"github.com/vova616/garageEngine/engine/input"
	//"log" 
	//"github.com/go-gl/glfw"
	//c "github.com/vova616/chipmunk"
	//. "github.com/vova616/chipmunk/vect"
	"math/rand"
)

type Map struct {
	engine.BaseComponent
	Sprite        *engine.Sprite
	Tiles         []uint32
	Sprites       map[int]*engine.Sprite
	FreeSprites   []*engine.Sprite
	Width, Height int
	TileSize      float32
}

func NewMap(tex *engine.Texture, uv engine.AnimatedUV) *Map {

	return &Map{BaseComponent: engine.NewComponent(),
		Sprite: engine.NewSprite3(tex, uv)}
}

func (m *Map) Start() {

	m.TileSize = m.Transform().WorldScale().Y
	w, h := (float32(engine.Width)/m.TileSize)+3, (float32(engine.Height)/m.TileSize)+3

	tW, tH := 1000, 1000
	m.Tiles = make([]uint32, int(tW*tH))

	m.Sprites = make(map[int]*engine.Sprite)
	m.FreeSprites = make([]*engine.Sprite, 0, int(w*h))

	for i, _ := range m.Tiles {
		m.Tiles[i] = rand.Uint32() % 2 //uint32(len(m.Sprite.UVs))
		//m.Tiles[i] = 1
	}
	/*
		tW, tH = 60, 23
		m.Tiles = []uint32{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	*/
	m.Width = int(tW)
	m.Height = int(tH)
	/*
		for index := uint32(0); index < uint32(cap(m.FreeSprites)); index++ {
			tile := engine.NewGameObject("Tile")
			tile.Transform().SetParent2(m.GameObject())
			tile.Transform().SetScalef(1.005, 1.005)

			tile.Transform().SetParent2(nil)
			if 1 == 1 {
				pos := tile.Transform().Position()
				tile.Transform().SetPositionf(pos.X+0.5, pos.Y+0.5)
				engine.GetScene().SceneBase().AddGameObject(tile)
				tile.SetActive(false)
			}

			var s *engine.Sprite = new(engine.Sprite)
			*s = *m.Sprite
			s.SetAnimationIndex(0)
			s.AnimationSpeed = 0
			tile.GameObject().AddComponent(s)

			m.FreeSprites = append(m.FreeSprites, s)
		}
	*/
}

func (m *Map) LateUpdate() {
	return
	camera := engine.GetScene().SceneBase().Camera
	cameraPos := camera.Transform().WorldPosition()
	mapPos := m.Transform().WorldPosition()

	for i, s := range m.Sprites {
		tPos := s.Transform().WorldPosition()
		if tPos.X > cameraPos.X+float32(engine.Width)+m.TileSize || tPos.X < cameraPos.X-m.TileSize {
			m.FreeSprites = append(m.FreeSprites, s)
			//s.Transform().SetParent(nil)
			s.GameObject().SetActive(false)
			delete(m.Sprites, i)
			continue
		}
		if tPos.Y > cameraPos.Y+float32(engine.Height)+m.TileSize || tPos.Y < cameraPos.Y-m.TileSize {
			m.FreeSprites = append(m.FreeSprites, s)
			//s.Transform().SetParent(nil)
			s.GameObject().SetActive(false)
			delete(m.Sprites, i)
			continue
		}

	}

	width := int(float32(engine.Width)/m.TileSize) + 2
	height := int(float32(engine.Height)/m.TileSize) + 2

	camX := (int(cameraPos.X) / int(m.TileSize)) * int(m.TileSize)
	camY := (int(cameraPos.Y) / int(m.TileSize)) * int(m.TileSize)

	for x := -1; x < width; x++ {
		for y := -1; y < height; y++ {
			var tilePos engine.Vector
			tilePos.X = float32(camX + (x * int(m.TileSize)))
			tilePos.Y = float32(camY + (y * int(m.TileSize)))

			if tilePos.X > cameraPos.X+float32(engine.Width)+m.TileSize || tilePos.X < cameraPos.X-m.TileSize {
				continue
			}
			if tilePos.Y > cameraPos.Y+float32(engine.Height)+m.TileSize || tilePos.Y < cameraPos.Y-m.TileSize {
				continue
			}

			pos := tilePos.X + (tilePos.Y * float32(width))
			s, exists := m.Sprites[int(pos)]
			if exists {
				continue
			}

			tx := (int(tilePos.X-mapPos.X) / int(m.TileSize)) + (m.Width / 2)
			ty := (int(tilePos.Y-mapPos.Y) / int(m.TileSize)) + (m.Height / 2)
			if tx < 0 || ty < 0 || tx >= m.Width || ty >= m.Height {
				continue
			}

			tileIndex := (tx) + (ty * m.Width)
			if tileIndex >= len(m.Tiles) {
				continue
			}

			if len(m.FreeSprites) == 0 {
				panic("Not enough sprites")
			}

			s, m.FreeSprites = m.FreeSprites[len(m.FreeSprites)-1], m.FreeSprites[:len(m.FreeSprites)-1]
			m.Sprites[int(pos)] = s
			s.GameObject().SetActive(true)
			s.Transform().SetWorldPositionf(tilePos.X, tilePos.Y)
			s.Transform().SetScalef(m.TileSize+0.000, m.TileSize+0.000)
			s.SetAnimationIndex(int(m.Tiles[tileIndex]))
		}
	}
}

func (m *Map) Draw() {
	camera := engine.GetScene().SceneBase().Camera
	cameraPos := camera.Transform().WorldPosition()
	mapPos := m.Transform().WorldPosition()

	width := int(float32(engine.Width)/m.TileSize) + 2
	height := int(float32(engine.Height)/m.TileSize) + 2

	camX := (int(cameraPos.X) / int(m.TileSize)) * int(m.TileSize)
	camY := (int(cameraPos.Y) / int(m.TileSize)) * int(m.TileSize)

	for x := -1; x < width; x++ {
		for y := -1; y < height; y++ {
			var tilePos engine.Vector
			tilePos.X = float32(camX + (x * int(m.TileSize)))
			tilePos.Y = float32(camY + (y * int(m.TileSize)))

			if tilePos.X > cameraPos.X+float32(engine.Width)+m.TileSize || tilePos.X < cameraPos.X-m.TileSize {
				continue
			}
			if tilePos.Y > cameraPos.Y+float32(engine.Height)+m.TileSize || tilePos.Y < cameraPos.Y-m.TileSize {
				continue
			}

			tx := (int(tilePos.X-mapPos.X) / int(m.TileSize)) + (m.Width / 2)
			ty := (int(tilePos.Y-mapPos.Y) / int(m.TileSize)) + (m.Height / 2)
			if tx < 0 || ty < 0 || tx >= m.Width || ty >= m.Height {
				continue
			}

			tileIndex := (tx) + (ty * m.Width)
			if tileIndex >= len(m.Tiles) {
				continue
			}

			engine.DrawSprite(m.Sprite.Texture, m.Sprite.UVs[int(m.Tiles[tileIndex])], tilePos, engine.NewVector2(m.TileSize, m.TileSize), 0, engine.AlignCenter, engine.One)
		}
	}
}
