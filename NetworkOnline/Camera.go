package NetworkOnline

import (
	. "github.com/vova616/GarageEngine/Engine"
	//"log"
	"github.com/vova616/GarageEngine/Engine/Input"
)

type CameraController struct {
	BaseComponent
}

func NewCameraCtl() *CameraController {
	return &CameraController{NewComponent()}
}

func (sp *CameraController) Update() {
	t := sp.GameObject().Transform()

	if Input.KeyDown('A') {
		t.Translatef(10, 0, 0)
	}
	if Input.KeyDown('D') {
		t.Translatef(-10, 0, 0)
	}
	if Input.KeyDown('S') {
		t.Translatef(0, 10, 0)
	}
	if Input.KeyDown('W') {
		t.Translatef(0, -10, 0)
	}
	if Input.KeyDown('E') {
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
	if Input.KeyDown('Q') {
		s := t.Scale()
		s.X += 0.05
		s.Y += 0.05
		t.SetScale(s)
	}
}
