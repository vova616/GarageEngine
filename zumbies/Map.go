package zumbies

import (
	"github.com/vova616/garageEngine/engine"
	"github.com/vova616/garageEngine/engine/input"
	//"log" 
	//"fmt"
	//"github.com/go-gl/glfw"
	//c "github.com/vova616/chipmunk"
	//. "github.com/vova616/chipmunk/vect"
	//"math/rand"
)

type Map struct {
	engine.BaseComponent
	Sprite *engine.Sprite

	Tiles []Tile

	Width, Height int
	TileSize      float32

	uvs       []engine.UV
	positions []engine.Vector
	scales    []engine.Vector
	rotations []float32
	alings    []engine.AlignType
	colors    []engine.Vector

	Disco       float32
	DiscoStyle  int
	EnableDisco bool
}

func NewMap(tex *engine.Texture, uv engine.AnimatedUV) *Map {

	return &Map{BaseComponent: engine.NewComponent(),
		Sprite: engine.NewSprite3(tex, uv)}
}

type Tile uint32

const (
	_ = iota
	_ = iota
	_ = iota
	_ = iota
	_ = iota
	_ = iota
	_ = iota

	SideLeft  = Tile(0)
	SideRight = Tile(1 << iota)
	SideUp    = Tile(1 << iota)
	SideDown  = SideRight | SideUp
	SideReset = ^SideDown

	CollisionNone  = Tile(0)
	CollisionLeft  = Tile(1 << iota)
	CollisionRight = Tile(1 << iota)
	CollisionUp    = Tile(1 << iota)
	CollisionDown  = Tile(1 << iota)
	CollisionAll   = CollisionLeft | CollisionRight | CollisionUp | CollisionDown
	CollitionReset = ^CollisionAll
)

func test() {
	println(SideLeft)
	println(SideRight)
	println(SideUp)
	println(SideDown)
	println(SideReset)
	println(Tile.SetSide(1, SideDown))
	println(Tile.SetSide(1, SideDown).Side())
	println(Tile.SetSide(1, SideDown).SetSide(SideLeft))
	println(Tile.SetSide(255, SideDown).Type())
}

func (t Tile) SetSide(side Tile) Tile {
	return Tile((uint32(t) & uint32(SideReset)) | uint32(side))
}

func (t Tile) SetCollision(collision Tile) Tile {
	return Tile((uint32(t) & uint32(CollitionReset)) | uint32(collision))
}

func (t Tile) Collision() Tile {
	if t&CollisionAll == CollisionAll {
		return CollisionAll
	}
	if t&CollisionDown == CollisionDown {
		return CollisionDown
	}
	if t&CollisionUp == CollisionUp {
		return CollisionUp
	}
	if t&CollisionRight == CollisionRight {
		return CollisionRight
	}
	if t&CollisionLeft == CollisionLeft {
		return CollisionLeft
	}
	return CollisionNone
}

func (t Tile) Side() Tile {
	if t&SideDown == SideDown {
		return SideDown
	}
	if t&SideUp == SideUp {
		return SideUp
	}
	if t&SideRight == SideRight {
		return SideRight
	}
	return SideLeft
}

func (t Tile) Angle() float32 {
	if t&SideDown == SideDown {
		return 90
	}
	if t&SideUp == SideUp {
		return -90
	}
	if t&SideRight == SideRight {
		return 180
	}
	return 0
}

func (t Tile) Type() Tile {
	return Tile(byte(t))
}

func (m *Map) Start() {

	m.TileSize = m.Transform().WorldScale().Y
	w, h := (float32(engine.Width)/m.TileSize)+3, (float32(engine.Height)/m.TileSize)+3

	//tW, tH := 1000, 1000
	//m.Tiles = make([]Tile, int(tW*tH))

	m.uvs = make([]engine.UV, int(w*h))
	m.positions = make([]engine.Vector, int(w*h))
	m.scales = make([]engine.Vector, int(w*h))
	m.rotations = make([]float32, int(w*h))
	m.colors = make([]engine.Vector, int(w*h))
	m.alings = make([]engine.AlignType, int(w*h))

	for i, _ := range m.alings {
		m.alings[i] = engine.AlignCenter
	}
	for i, _ := range m.colors {
		m.colors[i] = engine.One
	}

	//for i, _ := range m.Tiles {
	//m.Tiles[i] = rand.Uint32() % 2 //uint32(len(m.Sprite.UVs))
	//	m.Tiles[i] = 1
	//}

	//m.Width = int(tW)
	//m.Height = int(tH)
}

func (m *Map) Update() {
	if input.KeyPress('X') {
		m.EnableDisco = !m.EnableDisco
	}
}

func (m *Map) GetTile(x, y int) (tile Tile, exists bool) {
	if x >= m.Width || y >= m.Height || x < 0 || y < 0 {
		return 0, false
	}
	return m.Tiles[x+(y*m.Width)], true
}

func (m *Map) PositionToTile(worldPosition engine.Vector) (tile Tile, x, y int) {

	mapPos := m.Transform().WorldPosition()
	//Calculate center of map
	mapPos.X -= (m.TileSize * float32(m.Width) / 2)
	mapPos.Y -= (m.TileSize * float32(m.Height) / 2)

	//World to Local position
	worldPosition.X -= mapPos.X
	worldPosition.Y -= mapPos.Y

	//Turn upside down
	worldPosition.Y = (float32(m.Height) * m.TileSize) - worldPosition.Y

	//Calculate tile position in array
	xx := (worldPosition.X / m.TileSize)
	yy := (worldPosition.Y / m.TileSize)

	//Check if outside of map
	if xx < 0 || yy < 0 || xx >= float32(m.Width) || yy >= float32(m.Height) {
		return 0, -1, -1
	}

	//return the tile
	x, y = int(xx), int(yy)
	return m.Tiles[x+(y*m.Width)], int(x), int(y)
}

func (m *Map) GetTilePos(x, y int) (pos engine.Vector, exists bool) {
	if x >= m.Width || y >= m.Height || x < 0 || y < 0 {
		return engine.Zero, false
	}
	mapPos := m.Transform().WorldPosition()
	//Calculate center of map
	mapPos.X -= (m.TileSize * float32(m.Width) / 2)
	mapPos.Y -= (m.TileSize * float32(m.Height) / 2)

	//Sprites are drawn from center, lets reposition them to the corner
	mapPos.X += (m.TileSize / 2)
	mapPos.Y -= (m.TileSize / 2)

	//Calculate local position and turn upside down
	pos.X = float32(x) * m.TileSize
	pos.Y = (float32(m.Height) - float32(y)) * m.TileSize

	//Turn local to world
	pos.X += mapPos.X
	pos.Y += mapPos.Y

	return pos, true
}

func (m *Map) Draw() {
	camera := engine.GetScene().SceneBase().Camera
	cameraPos := camera.Transform().WorldPosition()

	//calculate max tiles on screen
	width := int(float32(engine.Width)/(m.TileSize-1)) + 2
	height := int(float32(engine.Height)/(m.TileSize-1)) + 2

	index := 0
	for x := -1; x < width; x++ {
		for y := -1; y < height; y++ {
			var tilePos engine.Vector

			//Try to find tiles on screen
			tilePos.X = cameraPos.X + float32((x * int(m.TileSize-1)))
			tilePos.Y = cameraPos.Y + float32((y * int(m.TileSize-1)))

			//Get tile
			tileType, tx, ty := m.PositionToTile(tilePos)

			//Check if visible/exists
			if tileType != 0 {
				//Add to draw list
				tile := Tile(tileType - 1)
				m.uvs[index] = m.Sprite.UVs[tile.Type()]
				pos, e := m.GetTilePos(tx, ty)
				if !e {
					panic("Does not exists")
				}
				//Pixel fix
				//pos.X += 0.75
				//pos.Y += 0.75
				m.positions[index] = pos
				m.rotations[index] = tile.Angle() + m.Disco
				m.scales[index] = engine.NewVector2(m.TileSize+float32(int(m.Disco)%52), m.TileSize+float32(int(m.Disco)%52))
				index++
				//engine.DrawSprite(m.Sprite.Texture, m.Sprite.UVs[int(m.Tiles[tileIndex])], tilePos, engine.NewVector2(m.TileSize, m.TileSize), 0, engine.AlignCenter, engine.One)
			}
		}
	}

	if m.EnableDisco {
		if m.DiscoStyle == 1 {
			m.Disco++
		} else {
			m.Disco--
		}
		if m.Disco > 50 {
			m.DiscoStyle = 0
		}
		if m.Disco < 0 {
			m.DiscoStyle = 1
		}
	}

	if index > 0 {
		engine.DrawSprites(m.Sprite.Texture, m.uvs[:index], m.positions[:index], m.scales[:index], m.rotations[:index], m.alings[:index], m.colors[:index])
	}

}
