package networkOnline

import (
	"github.com/vova616/garageEngine/engine"
	//"log"
	"github.com/vova616/garageEngine/engine/input"
)

type CameraController struct {
	engine.BaseComponent
	Speed float64
}

func NewCameraCtl(speed float64) *CameraController {
	return &CameraController{engine.NewComponent(), speed}
}

func (sp *CameraController) Update() {
	t := sp.GameObject().Transform()

	if input.KeyDown('A') {
		t.Translatef(float32(-sp.Speed*engine.DeltaTime()), 0)
	}
	if input.KeyDown('D') {
		t.Translatef(float32(sp.Speed*engine.DeltaTime()), 0)
	}
	if input.KeyDown('S') {
		t.Translatef(0, float32(-sp.Speed*engine.DeltaTime()))
	}
	if input.KeyDown('W') {
		t.Translatef(0, float32(sp.Speed*engine.DeltaTime()))
	}
	if input.KeyDown('E') {
		s := t.Scale()
		s.X -= 0.05
		s.Y -= 0.05
		if s.X <= 0.2 {
			s.X = 0.2
		}
		if s.Y <= 0.2 {
			s.Y = 0.2
		}
		t.SetScale(s)
	}
	if input.KeyDown('Q') {
		s := t.Scale()
		s.X += 0.05
		s.Y += 0.05
		t.SetScale(s)
	}
}
