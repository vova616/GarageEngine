package zumbies

import (
	"github.com/vova616/garageEngine/engine"
	//"log"
	"github.com/vova616/garageEngine/engine/input"
)

type PlayerController struct {
	engine.BaseComponent
	Player    *Player
	WalkSpeed float64
}

func NewPlayerController(player *Player) *PlayerController {
	return &PlayerController{engine.NewComponent(), player, 100}
}

func (this *PlayerController) Update() {
	//t := this.GameObject().Transform()
	var speed engine.Vector

	if input.KeyDown('A') {
		speed.X -= float32(this.WalkSpeed)
	}
	if input.KeyDown('D') {
		speed.X += float32(this.WalkSpeed)
	}
	if input.KeyDown('S') {
		speed.Y -= float32(this.WalkSpeed)
	}
	if input.KeyDown('W') {
		speed.Y += float32(this.WalkSpeed)
	}

	this.GameObject().Physics.Body.SetVelocity(speed.X, speed.Y)
}
