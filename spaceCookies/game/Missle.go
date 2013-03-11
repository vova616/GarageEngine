package game

import (
	"github.com/vova616/GarageEngine/engine"
	"math/rand"
	//"reflect"
)

type Missle struct {
	engine.BaseComponent
	Speed     float32
	Explosion *engine.GameObject
	exploded  bool
}

func NewMissle(speed float32) *Missle {
	return &Missle{BaseComponent: engine.NewComponent(), Speed: speed}
}

func (ms *Missle) OnComponentAdd() {
	ms.GameObject().Tag = MissleTag
}

func (ms *Missle) OnHit(enemey *engine.GameObject, damager *DamageDealer) {

}

func (ms *Missle) OnDie(byTimer bool) {
	if ms.Explosion == nil {
		ms.GameObject().Destroy()
		return
	}
	if ms.GameObject() == nil {
		return
	}
	if !byTimer {
		for i := 0; i < 2; i++ {
			n := ms.Explosion.Clone()
			n.Transform().SetParent2(GameSceneGeneral.Layer1)
			n.Transform().SetWorldPosition(ms.Transform().WorldPosition())
			s := n.Transform().Scale()
			n.Transform().SetScale(s.Mul2(rand.Float32() + 0.5))
			n.AddComponent(engine.NewPhysics(false))
			n.Transform().SetRotationf(rand.Float32() * 360)
			rot := n.Transform().Direction()
			n.Physics.Body.SetVelocity(-rot.X*10, -rot.Y*10)

			n.Physics.Body.SetMass(1)
			n.Physics.Shape.Group = 1
			n.Physics.Shape.IsSensor = true
		}
	}
	ms.GameObject().Destroy()
}
