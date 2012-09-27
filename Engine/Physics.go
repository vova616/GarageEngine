package Engine

import (
	c "github.com/vova616/chipmunk"
	. "github.com/vova616/chipmunk/vect"
)

type Physics struct {
	BaseComponent
	Body  *c.Body
	Box   *c.BoxShape
	Shape *c.Shape
}

func NewPhysics(static bool, w, h float32) *Physics {
	var body *c.Body

	box := c.NewBox(Vect{0, 0}, Float(w), Float(h))

	if static {
		body = c.NewBodyStatic()
	} else {
		body = c.NewBody(1, box.Moment(1))
	}

	p := &Physics{NewComponent(), body, box.GetAsBox(), box}
	body.UserData = p

	body.AddShape(box)
	return p
}

func NewPhysics2(static bool, shape *c.Shape) *Physics {
	var body *c.Body
	if static {
		body = c.NewBodyStatic()
	} else {
		body = c.NewBody(1, shape.ShapeClass.Moment(1))
	}

	p := &Physics{NewComponent(), body, shape.GetAsBox(), shape}
	body.UserData = p

	body.AddShape(shape)
	return p
}

func (p *Physics) Start() {
	//p.FixedUpdate()
	pos := p.GameObject().Transform().WorldPosition()
	p.Body.SetAngle(Float(p.GameObject().Transform().WorldRotation().Z) * RadianConst)
	p.Body.SetPosition(Vect{Float(pos.X), Float(pos.Y)})

	if p.GameObject().Sprite != nil {
		p.GameObject().Sprite.UpdateShape()
		p.Body.UpdateShapes()
	}

	//p.Body.UpdateShapes()
	Space.AddBody(p.Body)
}

func (p *Physics) OnComponentBind(binded *GameObject) {
	binded.Physics = p
	p.Body.CallbackHandler = p
}

func (c *Physics) CollisionPreSolve(arbiter *c.Arbiter) bool {
	if c.gameObject == nil {
		return true
	}
	return onCollisionPreSolveGameObject(c.GameObject(), (*Arbiter)(arbiter))
}

func (c *Physics) CollisionEnter(arbiter *c.Arbiter) bool {
	if c.gameObject == nil {
		return true
	}
	return onCollisionEnterGameObject(c.GameObject(), (*Arbiter)(arbiter))
}

func (c *Physics) CollisionExit(arbiter *c.Arbiter) {
	if c.gameObject == nil {
		return
	}
	onCollisionExitGameObject(c.GameObject(), (*Arbiter)(arbiter))
}

func (c *Physics) CollisionPostSolve(arbiter *c.Arbiter) {
	if c.gameObject == nil {
		return
	}
	onCollisionPostSolveGameObject(c.GameObject(), (*Arbiter)(arbiter))
}

func (p *Physics) OnDestroy() {
	p.gameObject = nil
	Space.RemoveBody(p.Body)
}

func (p *Physics) Clone() {
	p.Body = p.Body.Clone()
	p.Box = p.Body.Shapes[0].GetAsBox()
	p.Shape = p.Body.Shapes[0]
	p.Body.UserData = p
	//p.Body.UpdateShapes()
	//p.GameObject().Physics = nil
}
