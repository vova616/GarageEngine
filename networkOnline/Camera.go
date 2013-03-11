package networkOnline

import (
	"github.com/vova616/GarageEngine/engine"
	//"log"
	"github.com/vova616/GarageEngine/engine/input"
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
		engine.CurrentCamera().SetSize(engine.CurrentCamera().Size() + float32(engine.DeltaTime()))
	}
	if input.KeyDown('Q') {
		engine.CurrentCamera().SetSize(engine.CurrentCamera().Size() - float32(engine.DeltaTime()))
	}

	if input.KeyDown('Z') {
		a := t.Angle()
		t.SetWorldRotationf(a + float32(engine.DeltaTime())*10)
	}
	if input.KeyDown('X') {
		a := t.Angle()
		t.SetWorldRotationf(a - float32(engine.DeltaTime())*10)
	}
}
