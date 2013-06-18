package zumbies

import (
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/input"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	//"log"
)

type PlayerController struct {
	engine.BaseComponent
	Player          *Player
	Joint           *chipmunk.PivotJoint
	JointGameObject *engine.GameObject
	WalkSpeed       float32
}

func NewPlayerController(player *Player) *PlayerController {
	return &PlayerController{engine.NewComponent(), player, nil, nil, 200}
}

func (this *PlayerController) Start() {
	j := engine.NewGameObject("Joint")
	j.Transform().SetParent2(GameSceneGeneral.Layer1)
	j.Transform().SetWorldPosition(this.Transform().WorldPosition())
	j.AddComponent(engine.NewPhysics(false))
	j.Physics.Body.SetMass(engine.Inf)
	j.Physics.Body.SetMoment(engine.Inf)
	this.Joint = chipmunk.NewPivotJoint(j.Physics.Body, this.GameObject().Physics.Body)
	engine.Space.AddConstraint(this.Joint)
	j.Physics.Shape.IsSensor = true
	this.Joint.MaxBias = vect.Float(this.WalkSpeed)
	this.Joint.MaxForce = 3000
	this.JointGameObject = j
}

/*
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
*/

func (this *PlayerController) FixedUpdate() {
	this.GameObject().Physics.Shape.Layer = chipmunk.Layer(this.Player.Map.Layer)
	t := this.JointGameObject.Transform()

	var move engine.Vector = this.Transform().WorldPosition()

	if input.KeyDown('W') {
		move.Y += 100
	}
	if input.KeyDown('S') {
		move.Y += -100
	}
	if input.KeyDown('A') {
		move.X += -100
	}
	if input.KeyDown('D') {
		move.X += 100
	}

	t.SetWorldPosition(move)
}

/*
	if newPosition is bad disable walking
	else move to newPosition

*/
