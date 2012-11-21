package Components

import (
	"github.com/vova616/GarageEngine/Engine"
	//"Engine/Input"
	//"log"
)

type Controller struct {
	Engine.BaseComponent
	Speed float32
}

func NewController() *Controller {
	return &Controller{Engine.NewComponent(), 100}
}

func (sp *Controller) Update() {
	/*
		if Input.KeyDown('A') {
			sp.Transform().Position.X -= sp.Speed*DeltaTime()
		}
		if Input.KeyDown('D') {
			sp.Transform().Position.X += sp.Speed*DeltaTime()
		}
		if Input.KeyDown('W') {
			sp.Transform().Position.Y += sp.Speed*DeltaTime()
		}
		if Input.KeyDown('S') {
			sp.Transform().Position.Y -= sp.Speed*DeltaTime()
		}
		if Input.KeyDown('Q') {
			sp.Transform().Rotation.Z -= sp.Speed*DeltaTime()
		}
		if Input.KeyDown('E') {
			sp.Transform().Rotation.Z += sp.Speed*DeltaTime()
		}
		if Input.KeyDown('Z') {
			sp.Transform().Scale.X += sp.Speed*DeltaTime()
			sp.Transform().Scale.Y += sp.Speed*DeltaTime()
		}
		if Input.KeyDown('X') {
			sp.Transform().Scale.X -= sp.Speed*DeltaTime()
			sp.Transform().Scale.Y -= sp.Speed*DeltaTime()
		}
	*/
}
