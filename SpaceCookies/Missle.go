package SpaceCookies

import (
	. "github.com/vova616/GarageEngine/Engine"
	"math/rand"
	//"reflect"
)

type Missle struct {
	BaseComponent
	Speed     float32
	Explosion *GameObject
	exploded  bool
	Damage    float32
}

func NewMissle(speed float32) *Missle {
	return &Missle{BaseComponent: NewComponent(), Speed: speed, Damage: 50}
}

func (ms *Missle) OnComponentBind(gameObject *GameObject) {
	gameObject.Tag = MissleTag
}

func (ms *Missle) OnCollisionEnter(arbiter *Arbiter) bool {
	if ms.exploded {
		return true
	}
	ms.exploded = true
	if arbiter.GameObjectA().Tag == CookieTag || arbiter.GameObjectB().Tag == CookieTag {
		ms.CreateBlow()
		ms.GameObject().Destroy()
	}

	return true
}

func (ms *Missle) CreateBlow() {
	//ms.GameObject().Destroy()
	if ms.Explosion == nil {
		return
	}
	for i := 0; i < 10; i++ {
		n := ms.Explosion.Clone()
		n.Transform().SetParent2(GameSceneGeneral.Layer1)
		n.Transform().SetWorldPosition(ms.Transform().WorldPosition())
		s := n.Transform().Scale()
		n.Transform().SetScale(s.Mul2(rand.Float32() + 0.5))
		n.AddComponent(NewPhysics(false, 1, 1))

		n.Transform().SetRotationf(0, 0, rand.Float32()*360)
		rot := n.Transform().Rotation2D()
		n.Physics.Body.SetVelocity(-rot.X*10, -rot.Y*10)

		n.Physics.Body.SetMass(1)
		n.Physics.Shape.Group = 1
		n.Physics.Shape.IsSensor = true
	}
}
