package SpaceCookies

import (
	. "github.com/vova616/GarageEngine/Engine"
	//"Engine/Components"
	//"github.com/jteeuwen/glfw"
	"github.com/vova616/GarageEngine/Engine/Input"
	//"log"
	c "github.com/vova616/chipmunk"
	. "github.com/vova616/chipmunk/vect"
	//"fmt"
	//"time"

)

type ShipController struct {
	BaseComponent
	Speed           float32
	RotationSpeed   float32
	Missle          *Missle
	MisslesPosition []Vector
}

func NewShipController() *ShipController {
	return &ShipController{NewComponent(), 400000, 100, nil, []Vector{{-28, 10, 0},
		{28, 10, 0}}}
}

func (sp *ShipController) OnComponentBind(binded *GameObject) {
	sp.GameObject().AddComponent(NewPhysics2(false, c.NewCircle(Vect{0, 0}, 15)))
}

func (sp *ShipController) Start() {
	ph := sp.GameObject().Physics
	ph.Body.SetMass(50)
	ph.Shape.Group = 1
	//sp.Physics.Shape.Friction = 0.5
}

func (sp *ShipController) Shoot() {
	if sp.Missle != nil {
		s := sp.Transform().Direction2D(Up)
		a := sp.Transform().Rotation()
		//scale := sp.Transform().Scale()
		for _, pos := range sp.MisslesPosition {
			p := sp.Transform().WorldPosition()
			_ = a
			m := Identity()
			//m.Scale(scale.X, scale.Y, scale.Z)
			m.Translate(pos.X, pos.Y, pos.Z)
			m.Rotate(a.Z, 0, 0, -1)
			m.Translate(p.X, p.Y, p.Z)
			p = m.Translation()

			nfire := sp.Missle.GameObject().Clone()
			nfire.Transform().SetParent2(GameSceneGeneral.Layer3)
			nfire.Transform().SetWorldPosition(p)
			nfire.Physics.Body.IgnoreGravity = true
			nfire.Physics.Body.SetMass(0.1)
			nfire.Tag = MissleTag
			nfire.Physics.Body.AddForce(s.X*3000, s.Y*3000)

			nfire.Physics.Shape.Group = 1
			nfire.Physics.Body.SetMoment(Inf)
			nfire.Transform().SetRotation(sp.Transform().Rotation())
		}
	}
}

func (sp *ShipController) Update() {
	r := sp.Transform().Rotation()
	r2 := sp.Transform().Direction2D(Up)
	ph := sp.GameObject().Physics
	rx, ry := r2.X, r2.Y
	rx, ry = rx*DeltaTime(), ry*DeltaTime()

	if Input.KeyDown('W') {
		ph.Body.AddForce(sp.Speed*rx, sp.Speed*ry)
	}

	if Input.KeyDown('S') {
		ph.Body.AddForce(-sp.Speed*rx, -sp.Speed*ry)
	}

	if Input.KeyDown('D') {
		ph.Body.SetAngularVelocity(0)
		ph.Body.SetTorque(0)
		sp.Transform().SetRotationf(0, 0, r.Z-sp.RotationSpeed*DeltaTime())
	}
	if Input.KeyDown('A') {
		ph.Body.SetAngularVelocity(0)
		ph.Body.SetTorque(0)
		sp.Transform().SetRotationf(0, 0, r.Z+sp.RotationSpeed*DeltaTime())
	}

	if Input.KeyPress('F') {
		sp.Shoot()
	}

	if Input.KeyPress('P') {

		EnablePhysics = !EnablePhysics
	}
}

func (sp *ShipController) LateUpdate() {
	GameSceneGeneral.SceneData.Camera.Transform().SetPosition(NewVector3(sp.Transform().Position().X-float32(Width/2), sp.Transform().Position().Y-float32(Height/2), 0))
	

}
