package Engine

import (
	c "github.com/vova616/chipmunk"
	. "github.com/vova616/chipmunk/vect"
)

type Collision struct {
	Data      *c.Arbiter
	ColliderA *GameObject
	ColliderB *GameObject
}

func NewCollision(arbiter *c.Arbiter) Collision {
	a, _ := arbiter.ShapeA.Body.UserData.(*Physics)
	b, _ := arbiter.ShapeB.Body.UserData.(*Physics)

	if a == nil || b == nil {
		panic("dafuq")
	}

	arb := *arbiter

	return Collision{&arb, a.GameObject(), b.GameObject()}
}

type Physics struct {
	BaseComponent
	Body  *c.Body
	Box   *c.BoxShape
	Shape *c.Shape

	lastCollision    *c.Arbiter
	currentCollision *c.Arbiter
}

var (
	x = float64(100)
)

func NewPhysics(static bool, w, h float32) *Physics {
	var body *c.Body

	box := c.NewBox(Vect{0, 0}, Float(w), Float(h))

	if static {
		body = c.NewBodyStatic()
	} else {
		body = c.NewBody(1, box.Moment(1))
	}

	p := &Physics{NewComponent(), body, box.GetAsBox(), box, &c.Arbiter{}, &c.Arbiter{}}
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

	p := &Physics{NewComponent(), body, shape.GetAsBox(), shape, &c.Arbiter{}, &c.Arbiter{}}
	body.UserData = p

	body.AddShape(shape)
	return p
}

func (p *Physics) Start() {
	//p.FixedUpdate()
	pos := p.GameObject().Transform().WorldPosition()
	p.Body.SetAngle(Float(180-p.GameObject().Transform().WorldRotation().Z) * RadianConst)
	p.Body.SetPosition(Vect{Float(pos.X), Float(pos.Y)})
	//p.Body.UpdateShapes()
	Space.AddBody(p.Body)
}

func (p *Physics) OnComponentBind(binded *GameObject) {
	binded.Physics = p
}

func (p *Physics) Clone() {
	p.lastCollision = &c.Arbiter{}
	p.currentCollision = &c.Arbiter{}
	p.Body = p.Body.Clone()
	p.Box = p.Body.Shapes[0].GetAsBox()
	p.Shape = p.Body.Shapes[0]
	p.Body.UserData = p
	//p.Body.UpdateShapes()
	//p.GameObject().Physics = nil
}
