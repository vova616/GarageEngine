package engine

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

type Mouse struct {
	BaseComponent
}

func NewMouse() *Mouse {
	return &Mouse{NewComponent()}
}

func (m *Mouse) OnComponentAdd() {
	gameObject := m.GameObject()
	gameObject.Tag = MouseTag
	gameObject.AddComponent(NewPhysicsShape(false, chipmunk.NewCircle(vect.Vect{0, 0}, 0.5)))
	ph := gameObject.Physics
	ph.Body.SetMass(Inf)
	ph.Body.SetMoment(Inf)
	ph.Body.IgnoreGravity = true
	ph.Shape.IsSensor = true
}

func (m *Mouse) Update() {
	m.Transform().SetPosition(CurrentCamera().MouseLocalPosition())
}

func (m *Mouse) Start() {

}

func (m *Mouse) OnCollisionEnter(arbiter Arbiter) bool {
	return onMouseEnterGameObject(arbiter.GameObjectB(), arbiter)
}

func (m *Mouse) OnCollisionExit(arbiter Arbiter) {
	onMouseExitGameObject(arbiter.GameObjectB(), arbiter)
}
