package components

import (
	"github.com/vova616/GarageEngine/engine"
	//"Engine/input"
	//"log"
)

type Controller struct {
	engine.BaseComponent
	Speed float32
}

func NewController() *Controller {
	return &Controller{engine.NewComponent(), 100}
}

func (sp *Controller) Update() {
	/*
		if input.KeyDown('A') {
			sp.Transform().Position.X -= sp.Speed*DeltaTime()
		}
		if input.KeyDown('D') {
			sp.Transform().Position.X += sp.Speed*DeltaTime()
		}
		if input.KeyDown('W') {
			sp.Transform().Position.Y += sp.Speed*DeltaTime()
		}
		if input.KeyDown('S') {
			sp.Transform().Position.Y -= sp.Speed*DeltaTime()
		}
		if input.KeyDown('Q') {
			sp.Transform().Rotation.Z -= sp.Speed*DeltaTime()
		}
		if input.KeyDown('E') {
			sp.Transform().Rotation.Z += sp.Speed*DeltaTime()
		}
		if input.KeyDown('Z') {
			sp.Transform().Scale.X += sp.Speed*DeltaTime()
			sp.Transform().Scale.Y += sp.Speed*DeltaTime()
		}
		if input.KeyDown('X') {
			sp.Transform().Scale.X -= sp.Speed*DeltaTime()
			sp.Transform().Scale.Y -= sp.Speed*DeltaTime()
		}
	*/
}
