package Engine

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

type Component interface {
	Draw()
	PostDraw()
	Update()
	FixedUpdate()
	Start()
	Clone()
	LateUpdate()
	OnCollision(collision Collision)
	OnCollisionEnter(collision Collision)
	OnCollisionExit(collision Collision)
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

func (c *BaseComponent) Clone() {

}

func (c *BaseComponent) LateUpdate() {

}

func (c *BaseComponent) OnCollision(collision Collision) {

}

func (c *BaseComponent) OnCollisionEnter(collision Collision) {

}

func (c *BaseComponent) OnCollisionExit(collision Collision) {

}

func (c *BaseComponent) OnComponentBind(binded *GameObject) {

}

func (c *BaseComponent) Destroy() {

}
