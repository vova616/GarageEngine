package SpaceCookies

import (
	. "github.com/vova616/GarageEngine/Engine"
	"math/rand"
)

type ResizeScript struct {
	BaseComponent
	MinX, MaxX float32
	MinY, MaxY float32
	State      int
	Speed      float32
}

func NewResizeScript(minX, maxX, minY, maxY float32) *ResizeScript {
	return &ResizeScript{BaseComponent: NewComponent(), MinX: minX, MaxX: maxX, MinY: minY, MaxY: maxY, Speed: 3}
}

func (m *ResizeScript) Update() {
	if m.State == 0 {
		sx := (rand.Float32() * m.MaxX) + m.MinX
		sy := (rand.Float32() * m.MaxY) + m.MinY
		sx = m.MinX
		sy = m.MinY
		m.Transform().SetScalef(sx, sy, 1)
		m.State = 1
	} else if m.State == 1 {
		deltaX := m.MaxX - m.MinX
		deltaY := m.MaxX - m.MinX
		s := m.Transform().Scale()
		s.X, s.Y = s.X+(deltaX*m.Speed*DeltaTime()), s.Y+(deltaY*m.Speed*DeltaTime())
		m.Transform().SetScale(s)
		if s.X > m.MaxX || s.Y > m.MaxY {
			m.State = 2
		}
	} else if m.State == 2 {
		deltaX := m.MaxX - m.MinX
		deltaY := m.MaxX - m.MinX
		s := m.Transform().Scale()
		s.X, s.Y = s.X-(deltaX*m.Speed*DeltaTime()), s.Y-(deltaY*m.Speed*DeltaTime())
		m.Transform().SetScale(s)
		if s.X < m.MinX || s.Y < m.MinY {
			m.State = 0
		}
	}
}
