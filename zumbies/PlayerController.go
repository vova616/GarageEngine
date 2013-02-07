package zumbies

import (
	"github.com/vova616/garageEngine/engine"
	"github.com/vova616/garageEngine/engine/input"
	//"log"
)

type PlayerController struct {
	engine.BaseComponent
	Player       *Player
	WalkSpeed    float64
	LastPosition engine.Vector
}

func NewPlayerController(player *Player) *PlayerController {
	return &PlayerController{engine.NewComponent(), player, 100, engine.Zero}
}

func (this *PlayerController) Start() {
	this.LastPosition = this.GameObject().Transform().WorldPosition()
}

func (this *PlayerController) FixedUpdate() {
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

	tpos := this.Transform().WorldPosition()
	pos := tpos
	tpos.Y -= 32
	_, x, y := Map1.PositionToTile(tpos)

	t, exists := Map1.GetTile(x, y)
	if !exists || t.Collision() != CollisionNone {

		dir := pos.Sub(this.LastPosition)
		if dir.X != 0 {
			speed.X = 0
		}
		if dir.Y != 0 {
			speed.Y = 0
		}

		this.GameObject().Transform().SetWorldPosition(this.LastPosition)
	}

	this.GameObject().Physics.Body.SetVelocity(speed.X, speed.Y)
	this.LastPosition = this.GameObject().Transform().WorldPosition()
}

func (this *PlayerController) Update() {
	pos := this.Transform().WorldPosition()
	pos.Y -= 32
	_, x, y := Map1.PositionToTile(pos)
	t, exists := Map1.GetTile(x, y)
	if !exists || t.Collision() != CollisionNone {
		this.GameObject().Transform().SetWorldPosition(this.LastPosition)
	}
}
