package game

import (
	"github.com/vova616/GarageEngine/engine"
	"math/rand"
)

type ResizeScript struct {
	engine.BaseComponent
	MinX, MediumX, MaxX float32
	MinY, MediumY, MaxY float32
	State               int
	Speed               float32
}

func NewResizeScript(minX, mediumX, maxX, minY, mediumY, maxY float32) *ResizeScript {
	return &ResizeScript{BaseComponent: engine.NewComponent(), MinX: minX, MediumX: mediumX, MaxX: maxX, MinY: minY, MediumY: mediumY, MaxY: maxY, Speed: 3}
}

func (m *ResizeScript) SetValues(minX, mediumX, maxX, minY, mediumY, maxY float32) {
	m.MinX = minX
	m.MediumX = mediumX
	m.MaxX = maxX
	m.MinY = minY
	m.MediumY = mediumY
	m.MaxY = maxY

	s := m.Transform().Scale()
	if s.X > m.MaxX {
		s.X = m.MaxX
	}
	if s.Y > m.MaxY {
		s.Y = m.MaxY
	}
	if s.X < m.MinX {
		s.X = m.MinX
	}
	if s.Y < m.MinY {
		s.Y = m.MinY
	}
	m.Transform().SetScale(s)
}

func (m *ResizeScript) Update() {
	delta := float32(engine.DeltaTime())
	if m.State == 0 {
		sx := (rand.Float32() * m.MaxX) + m.MinX
		sy := (rand.Float32() * m.MaxY) + m.MinY
		sx = m.MinX
		sy = m.MinY
		m.Transform().SetScalef(sx, sy)
		m.State = 1
	} else if m.State == 1 {
		deltaX := m.MaxX - m.MinX
		deltaY := m.MaxY - m.MinY
		s := m.Transform().Scale()
		s.X, s.Y = s.X+(deltaX*m.Speed*delta), s.Y+(deltaY*m.Speed*delta)
		m.Transform().SetScale(s)
		if s.X > m.MaxX || s.Y > m.MaxY {
			m.State = 2
		}
	} else if m.State == 2 {
		deltaX := m.MaxX - m.MediumX
		deltaY := m.MaxY - m.MediumY
		s := m.Transform().Scale()
		s.X, s.Y = s.X-(deltaX*m.Speed*delta), s.Y-(deltaY*m.Speed*delta)
		m.Transform().SetScale(s)
		if s.X < m.MediumX || s.Y < m.MediumY {
			m.State = 3
		}
	} else if m.State == 3 {
		deltaX := m.MaxX - m.MediumX
		deltaY := m.MaxY - m.MediumY
		s := m.Transform().Scale()
		s.X, s.Y = s.X+(deltaX*m.Speed*delta), s.Y+(deltaY*m.Speed*delta)
		m.Transform().SetScale(s)
		if s.X > m.MaxX || s.Y > m.MaxY {
			m.State = 2
		}
	}
}
