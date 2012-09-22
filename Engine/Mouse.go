package Engine

import (
	c "github.com/vova616/chipmunk"
	. "github.com/vova616/chipmunk/vect"
)

type Mouse struct {
	BaseComponent
}

func NewMouse() *Mouse {
	return &Mouse{NewComponent()}
}

func (m *Mouse) Update() {
	m.Transform().SetPosition(mainScene.SceneBase().Camera.MouseLocalPosition())
}

func (m *Mouse) Start() {
	m.GameObject().AddComponent(NewPhysics2(false, c.NewCircle(Vect{0, 0}, Float(0.5))))
	ph := m.GameObject().Physics
	ph.Body.IgnoreGravity = true
	ph.Shape.IsSensor = true
}

func (m *Mouse) OnCollisionEnter(arbiter *Arbiter) bool {
	if m.GameObject().Physics.Body == arbiter.ShapeA.Body {
		return onMouseEnterGameObject(arbiter.GameObjectB(), arbiter)
	} else {
		return onMouseEnterGameObject(arbiter.GameObjectA(), arbiter)
	}

	return true
}

func (m *Mouse) OnCollisionExit(arbiter *Arbiter) {
	if m.GameObject().Physics.Body == arbiter.ShapeA.Body {
		onMouseExitGameObject(arbiter.GameObjectB(), arbiter)
	} else {
		onMouseExitGameObject(arbiter.GameObjectA(), arbiter)
	}
}
