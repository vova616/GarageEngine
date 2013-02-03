package networkOnline

import (
	"github.com/vova616/garageEngine/engine"
	//"log"
	"github.com/vova616/garageEngine/engine/input"
)

type CameraController struct {
	engine.BaseComponent
}

func NewCameraCtl() *CameraController {
	return &CameraController{engine.NewComponent()}
}

func (sp *CameraController) Update() {
	t := sp.GameObject().Transform()

	if input.KeyDown('A') {
		t.Translatef(-10, 0)
	}
	if input.KeyDown('D') {
		t.Translatef(10, 0)
	}
	if input.KeyDown('S') {
		t.Translatef(0, -10)
	}
	if input.KeyDown('W') {
		t.Translatef(0, 10)
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
