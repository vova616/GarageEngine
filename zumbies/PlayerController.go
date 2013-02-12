package zumbies

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"github.com/vova616/garageEngine/engine"
	"github.com/vova616/garageEngine/engine/input"
	//"log"
)

type PlayerController struct {
	engine.BaseComponent
	Player    *Player
	WalkSpeed float64
}

func NewPlayerController(player *Player) *PlayerController {
	return &PlayerController{engine.NewComponent(), player, 100}
}

func (this *PlayerController) Start() {
	this.GameObject().Physics.Body.UpdatePositionFunc = func(body *chipmunk.Body, dt vect.Float) {

		v := body.Velocity()
		newPos := vect.Add(body.Position(), vect.Mult(vect.Add(v, body.VBias()), dt))
		cmap := this.Player.Map
		p := engine.Vector{float32(newPos.X), float32(newPos.Y - 22), 0}
		xc, yc := cmap.GetCollisions(p, 25, 10)
		if xc != nil {

			oldPos := body.Position()
			worldPosition := engine.Vector{float32(oldPos.X), float32(oldPos.Y - 22), 0}

			_, x1, y1 := cmap.PositionToTile(worldPosition.Add(engine.Vector{25 / 2, 10 / 2, 0}))
			_, x2, y2 := cmap.PositionToTile(worldPosition.Add(engine.Vector{-25 / 2, -10 / 2, 0}))
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			if y1 > y2 {
				y1, y2 = y2, y1
			}

			for i := 0; i < len(xc); i++ {
				x, y := xc[i], yc[i]
				if !(y >= y1 && y2 >= y) && (x >= x1 && x2 >= x) {
					v.Y = 0
				}
				if !(x >= x1 && x2 >= x) && (y >= y1 && y2 >= y) {
					v.X = 0
				}
			}

			newPos = vect.Add(body.Position(), vect.Mult(vect.Add(v, body.VBias()), dt))
		}

		body.SetPosition(newPos)

		body.SetAngle(body.Angle() + vect.Float(body.AngularVelocity()+body.WBias())*dt)

		body.SetVBias(vect.Vector_Zero)
		body.SetWBias(0)
	}
}

func (this *PlayerController) FixedUpdate() {
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

/*
	if newPosition is bad disable walking
	else move to newPosition

*/
