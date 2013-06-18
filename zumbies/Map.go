package zumbies

import (
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/input"
	//"log"
	//"fmt"
	//"github.com/go-gl/glfw"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	//"math/rand"
	"sort"
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
	alings    []engine.Align
	colors    []engine.Color

	Disco       float32
	DiscoStyle  int
	EnableDisco bool

	Layer int
}

func NewMap(tex *engine.Texture, uv engine.AnimatedUV) *Map {

	return &Map{BaseComponent: engine.NewComponent(),
		Sprite: engine.NewSprite3(tex, uv)}
}

type Tile int64

const (
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	_               = iota
	SideRight       = Tile(1 << iota)
	SideUp          = Tile(1 << iota)
	SideRight2      = Tile(1 << iota)
	SideUp2         = Tile(1 << iota)
	CollisionLeft   = Tile(1 << iota)
	CollisionRight  = Tile(1 << iota)
	CollisionUp     = Tile(1 << iota)
	CollisionDown   = Tile(1 << iota)
	LayerConnection = Tile(1 << iota)

	SideLeft  = Tile(0)
	SideDown  = SideRight | SideUp
	SideReset = ^SideDown

	SideLeft2  = Tile(0)
	SideDown2  = SideRight2 | SideUp2
	SideReset2 = ^SideDown2

	CollisionNone  = Tile(0)
	CollisionAll   = CollisionLeft | CollisionRight | CollisionUp | CollisionDown
	CollitionReset = ^CollisionAll

	LayerConnectionReset = ^LayerConnection
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

	println(Tile.SetType2(255, 15).Type())
	println(Tile.SetType2(12, 255).Type2())
	println(Tile(4).SetSide2(SideRight).SetType2(6).Side2() == SideRight, SideRight, SideRight2)

	println((SideRight << 2) == SideRight2)
}

func (t Tile) SetSide(side Tile) Tile {
	return (t & SideReset) | side
}

func (t Tile) SetSide2(side Tile) Tile {
	if side < SideRight2 {
		side = side << 2
	}
	return (t & SideReset2) | side
}

func (t Tile) SetCollision(collision Tile) Tile {
	return (t & CollitionReset) | collision
}

func (t Tile) SetLayerConnection(on bool) Tile {
	if on {
		return t | LayerConnection
	}
	return t & LayerConnectionReset
}

func (t Tile) LayerConnected() bool {
	return t&LayerConnection == LayerConnection
}

func (t Tile) Collision() Tile {
	if t&CollisionAll == CollisionAll || t == 0 {
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

func (t Tile) Side2() Tile {
	if t&SideDown2 == SideDown2 {
		return SideDown
	}
	if t&SideUp2 == SideUp2 {
		return SideUp
	}
	if t&SideRight2 == SideRight2 {
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

func (t Tile) Angle2() float32 {
	if t&SideDown2 == SideDown2 {
		return 90
	}
	if t&SideUp2 == SideUp2 {
		return -90
	}
	if t&SideRight2 == SideRight2 {
		return 180
	}
	return 0
}

func (t Tile) Type() Tile {
	return Tile(byte(t))
}

func (t Tile) Type2() Tile {
	return Tile(uint16(t) >> 8)
}

func (t Tile) SetType(typ byte) Tile {
	return (t & ^0xff) | Tile(typ)
}

func (t Tile) SetType2(typ2 byte) Tile {
	return (t & ^0xff00) | (Tile(typ2) << 8) | Tile(byte(t))
}

func (m *Map) Start() {

	m.TileSize = m.Transform().WorldScale().Y
	w, h := (float32(engine.Width)/m.TileSize)+3, (float32(engine.Height)/m.TileSize)+3

	//tW, tH := 1000, 1000
	//m.Tiles = make([]Tile, int(tW*tH))

	m.uvs = make([]engine.UV, int(w*h*2))
	m.positions = make([]engine.Vector, int(w*h*2))
	m.scales = make([]engine.Vector, int(w*h*2))
	m.rotations = make([]float32, int(w*h*2))
	m.colors = make([]engine.Color, int(w*h*2))
	m.alings = make([]engine.Align, int(w*h*2))

	for i, _ := range m.alings {
		m.alings[i] = engine.AlignCenter
	}
	for i, _ := range m.colors {
		m.colors[i] = engine.Color_White
	}

	m.GenerateCollision()

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

func (m *Map) IsTileWalkabke(x, y int) bool {
	if x >= m.Width || y >= m.Height || x < 0 || y < 0 {
		return false
	}
	t := m.Tiles[x+(y*m.Width)]
	if t == 0 {
		return false
	}
	return t.Collision() == CollisionNone
}

func (m *Map) CheckCollision(worldPosition engine.Vector, width, height float32) bool {
	_, x1, y1 := m.PositionToTile(worldPosition.Add(engine.Vector{width / 2, height / 2, 0}))
	_, x2, y2 := m.PositionToTile(worldPosition.Add(engine.Vector{-width / 2, -height / 2, 0}))
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	for ; y1 <= y2; y1++ {
		for x := x1; x <= x2; x++ {
			if !m.IsTileWalkabke(x, y1) {
				return true
			}
		}
	}

	return false
}

func (m *Map) GetCollisions(worldPosition engine.Vector, width, height float32) (x, y []int) {
	return m.GetCollisions2(worldPosition, width, height, nil, nil)
}

func (m *Map) GetCollisions2(worldPosition engine.Vector, width, height float32, xCollisions, yCollisions []int) (x, y []int) {
	_, x1, y1 := m.PositionToTile(worldPosition.Add(engine.Vector{width / 2, height / 2, 0}))
	_, x2, y2 := m.PositionToTile(worldPosition.Add(engine.Vector{-width / 2, -height / 2, 0}))
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for ; y1 <= y2; y1++ {
		for x := x1; x <= x2; x++ {
			if !m.IsTileWalkabke(x, y1) {
				if xCollisions == nil {
					xCollisions = make([]int, 0, 4)
				}
				if yCollisions == nil {
					yCollisions = make([]int, 0, 4)
				}
				xCollisions = append(xCollisions, x)
				yCollisions = append(yCollisions, y1)
			}
		}
	}

	return xCollisions, yCollisions
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
	if xx < 0 {
		x = -1
	} else {
		x = int(xx)
	}
	if yy < 0 {
		y = -1
	} else {
		y = int(yy)
	}

	//Check if outside of map
	if xx < 0 || yy < 0 || xx >= float32(m.Width) || yy >= float32(m.Height) {
		return 0, x, y
	}

	//return the tile
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

	//Sprites are drawn from center, lets reposition them to the center from the corner
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

func (m *Map) GenerateCollision2() {
	tilesy := make(map[int][]int)
	tilesx := make(map[int][]int)
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if t, e := m.GetTile(x, y); e {
				if t.Collision() != CollisionNone {
					tilesy[y] = append(tilesy[y], x)
					tilesx[x] = append(tilesx[x], y)
				}
			}
		}
	}

	var shapes []*chipmunk.Shape

	centerx := vect.Float(m.TileSize * float32(m.Width) / 2)
	centery := vect.Float(m.TileSize * float32(m.Height) / 2)

	println("Map layer", m.Layer)
	for y, xarr := range tilesy {
		minx := xarr[0]
		maxx := xarr[0]
		sort.Ints(xarr)

		for i := 1; i < len(xarr); i++ {
			x := xarr[i]
			if maxx+1 == x {
				maxx = x
			} else {
				if maxx-minx > 1 {
					shapes = append(shapes, chipmunk.NewSegment(
						vect.Vect{vect.Float(float32(minx)*m.TileSize) - centerx, -vect.Float(float32(y)*m.TileSize) + centery},
						vect.Vect{vect.Float(float32(maxx)*m.TileSize) - centerx, -vect.Float(float32(y)*m.TileSize) + centery},
						1))
					println(y, minx, maxx)
				}
				minx = x
				maxx = x
			}
		}
		if maxx-minx > 1 {
			shapes = append(shapes, chipmunk.NewSegment(
				vect.Vect{vect.Float(float32(minx)*m.TileSize) - centerx, -vect.Float(float32(y)*m.TileSize) + centery},
				vect.Vect{vect.Float(float32(maxx)*m.TileSize) - centerx, -vect.Float(float32(y)*m.TileSize) + centery},
				1))
			println(y, minx, maxx)
		}

	}

	for x, yarr := range tilesx {
		miny := yarr[0]
		maxy := yarr[0]
		sort.Ints(yarr)
		println(len(yarr))
		for i := 1; i < len(yarr); i++ {
			y := yarr[i]
			if maxy+1 == y {
				maxy = y
			} else {
				if maxy-miny > 1 {
					shapes = append(shapes, chipmunk.NewSegment(
						vect.Vect{vect.Float(float32(x)*m.TileSize) - centerx, -vect.Float(float32(miny)*m.TileSize) + centery},
						vect.Vect{vect.Float(float32(x)*m.TileSize) - centerx, -vect.Float(float32(maxy)*m.TileSize) + centery},
						1))
					println(x, miny, maxy)
				}
				miny = y
				maxy = y
			}
		}
		if maxy-miny > 1 {
			shapes = append(shapes, chipmunk.NewSegment(
				vect.Vect{vect.Float(float32(x)*m.TileSize) - centerx, -vect.Float(float32(miny)*m.TileSize) + centery},
				vect.Vect{vect.Float(float32(x)*m.TileSize) - centerx, -vect.Float(float32(maxy)*m.TileSize) + centery},
				1))
			println(x, miny, maxy)
		}

	}
	println("Map end")
	m.GameObject().AddComponent(engine.NewPhysicsShapes(true, shapes))

}

func (m *Map) GenerateCollision() {

	var shapes []*chipmunk.Shape

	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if t, e := m.GetTile(x, y); e {
				if t.Collision() != CollisionNone {
					p, _ := m.GetTilePos(x, y)
					shape := chipmunk.NewBox(vect.Vect{vect.Float(p.X), vect.Float(p.Y)},
						vect.Float(m.TileSize), vect.Float(m.TileSize))

					shape.Layer = chipmunk.Layer(m.Layer)
					shape.SetFriction(0)
					shapes = append(shapes, shape)
				}
			}
		}
	}

	m.GameObject().AddComponent(engine.NewPhysicsShapes(true, shapes))

}

func (m *Map) Draw() {
	if m.Layer < MainPlayer.Map.Layer {
		return
	}
	camera := engine.GetScene().SceneBase().Camera
	cameraPos := camera.Transform().WorldPosition()
	w, h := camera.ScreenSize()
	cameraPos.X -= w / 2
	cameraPos.Y -= h / 2

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
			if tileType.Type() != 0 {
				//Add to draw list
				m.uvs[index] = m.Sprite.UVs[tileType.Type()-1]
				pos, e := m.GetTilePos(tx, ty)
				if !e {
					panic("Does not exists")
				}
				m.positions[index] = pos
				m.rotations[index] = tileType.Angle() + m.Disco
				m.scales[index] = engine.NewVector2(m.TileSize+float32(int(m.Disco)%52), m.TileSize+float32(int(m.Disco)%52))
				index++
				//engine.DrawSprite(m.Sprite.Texture, m.Sprite.UVs[int(m.Tiles[tileIndex])], tilePos, engine.NewVector2(m.TileSize, m.TileSize), 0, engine.AlignCenter, engine.One)
			}
			if tileType.Type2() != 0 {
				//Add to draw list
				m.uvs[index] = m.Sprite.UVs[tileType.Type2()-1]
				pos, e := m.GetTilePos(tx, ty)
				if !e {
					panic("Does not exists")
				}
				m.positions[index] = pos
				m.rotations[index] = tileType.Angle2() + m.Disco
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
