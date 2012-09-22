package Engine

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
	component.OnComponentBind(gameObject)
}

func (c *BaseComponent) GameObject() *GameObject {
	return c.gameObject
}

func (c *BaseComponent) Transform() *Transform {
	return c.gameObject.Transform()
}

/*
type CollisionCallback interface {
	OnCollisionEnter(arbiter *Arbiter) bool
	OnCollisionPreSolve(arbiter *Arbiter) bool
	OnCollisionPostSolve(arbiter *Arbiter)
	OnCollisionExit(arbiter *Arbiter)
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

	OnCollisionEnter(arbiter *Arbiter) bool
	OnCollisionPreSolve(arbiter *Arbiter) bool
	OnCollisionPostSolve(arbiter *Arbiter)
	OnCollisionExit(arbiter *Arbiter)

	OnMouseEnter(arbiter *Arbiter) bool
	OnMouseExit(arbiter *Arbiter)

	OnComponentBind(binded *GameObject)
	Destroy()
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

func (c *BaseComponent) Start() {

}

func (c *BaseComponent) Awake() {

}

func (c *BaseComponent) Clone() {

}

func (c *BaseComponent) LateUpdate() {

}

func (c *BaseComponent) OnCollisionPreSolve(arbiter *Arbiter) bool {
	return true
}

func (c *BaseComponent) OnCollisionEnter(arbiter *Arbiter) bool {
	return true
}

func (c *BaseComponent) OnCollisionExit(arbiter *Arbiter) {

}

func (c *BaseComponent) OnMouseEnter(arbiter *Arbiter) bool {
	return true
}

func (c *BaseComponent) OnMouseExit(arbiter *Arbiter) {

}

func (c *BaseComponent) OnCollisionPostSolve(arbiter *Arbiter) {

}

func (c *BaseComponent) OnComponentBind(binded *GameObject) {

}

func (c *BaseComponent) Destroy() {

}
