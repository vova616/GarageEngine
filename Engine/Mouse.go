package Engine

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

func (m *Mouse) OnComponentBind(gameObject *GameObject) {
	gameObject.Tag = MouseTag
}

func (m *Mouse) Update() {
	m.Transform().SetPosition(mainScene.SceneBase().Camera.MouseLocalPosition())
}

func (m *Mouse) Start() {
	m.GameObject().AddComponent(NewPhysics2(false, chipmunk.NewCircle(vect.Vect{0, 0}, 0.5)))
	ph := m.GameObject().Physics
	ph.Body.IgnoreGravity = true
	ph.Shape.IsSensor = true
}

func (m *Mouse) OnCollisionEnter(arbiter Arbiter) bool {
	return onMouseEnterGameObject(arbiter.GameObjectB(), arbiter)
}

func (m *Mouse) OnCollisionExit(arbiter Arbiter) {
	onMouseExitGameObject(arbiter.GameObjectB(), arbiter)
}
