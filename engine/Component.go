package engine

import (
//c "github.com/vova616/chipmunk"
)

type BaseComponent struct {
	hasStarted bool
	gameObject *GameObject
}

func NewComponent() BaseComponent {
	return BaseComponent{}
}

func (c *BaseComponent) onAdd(component Component, gameObject *GameObject) {
	c.gameObject = gameObject
	component.OnComponentAdd()
}

func (c *BaseComponent) GameObject() *GameObject {
	if c.gameObject != nil && c.gameObject.IsValid() == false {
		return nil
	}
	return c.gameObject
}

func (c *BaseComponent) Transform() *Transform {
	g := c.GameObject()
	if g != nil {
		return g.Transform()
	}
	return nil
}

/*
type CollisionCallback interface {
	OnCollisionEnter(arbiter Arbiter) bool
	OnCollisionPreSolve(arbiter Arbiter) bool
	OnCollisionPostSolve(arbiter Arbiter)
	OnCollisionExit(arbiter Arbiter)
}
*/

type Component interface {
	Draw()
	PostDraw()
	Update()
	FixedUpdate()
	Start()
	Awake()
	Clone()
	LateUpdate()

	OnEnable()
	OnDisable()

	OnCollisionEnter(arbiter Arbiter) bool
	OnCollisionPreSolve(arbiter Arbiter) bool
	OnCollisionPostSolve(arbiter Arbiter)
	OnCollisionExit(arbiter Arbiter)

	OnMouseEnter(arbiter Arbiter) bool
	OnMouseExit(arbiter Arbiter)

	OnComponentAdd()
	OnDestroy()
	started() bool
	setStarted(b bool)
	setGameObject(gobj *GameObject)
	onAdd(component Component, gameObject *GameObject)
}

func (c *BaseComponent) started() bool {
	return c.hasStarted
}

func (c *BaseComponent) setStarted(b bool) {
	c.hasStarted = b
}

func (c *BaseComponent) setGameObject(gobj *GameObject) {
	c.gameObject = gobj
}

func (c *BaseComponent) Draw() {

}

func (c *BaseComponent) PostDraw() {

}

func (c *BaseComponent) Update() {

}

func (c *BaseComponent) FixedUpdate() {

}

func (c *BaseComponent) OnEnable() {

}

func (c *BaseComponent) OnDisable() {

}

func (c *BaseComponent) Start() {

}

func (c *BaseComponent) Awake() {

}

func (c *BaseComponent) Clone() {

}

func (c *BaseComponent) LateUpdate() {

}

func (c *BaseComponent) OnCollisionPreSolve(arbiter Arbiter) bool {
	return true
}

func (c *BaseComponent) OnCollisionEnter(arbiter Arbiter) bool {
	return true
}

func (c *BaseComponent) OnCollisionExit(arbiter Arbiter) {

}

func (c *BaseComponent) OnMouseEnter(arbiter Arbiter) bool {
	return true
}

func (c *BaseComponent) OnMouseExit(arbiter Arbiter) {

}

func (c *BaseComponent) OnCollisionPostSolve(arbiter Arbiter) {

}

func (c *BaseComponent) OnComponentAdd() {

}

func (c *BaseComponent) OnDestroy() {

}
