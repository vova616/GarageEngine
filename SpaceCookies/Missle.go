package SpaceCookies

import (
	"github.com/vova616/GarageEngine/Engine"
	"math/rand"
	//"reflect"
)

type Missle struct {
	Engine.BaseComponent
	Speed     float32
	Explosion *Engine.GameObject
	exploded  bool
}

func NewMissle(speed float32) *Missle {
	return &Missle{BaseComponent: Engine.NewComponent(), Speed: speed}
}

func (ms *Missle) OnComponentBind(gameObject *Engine.GameObject) {
	gameObject.Tag = MissleTag
}

func (ms *Missle) OnHit(enemey *Engine.GameObject, damager *DamageDealer) {

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
			n.AddComponent(Engine.NewPhysics(false, 1, 1))
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
