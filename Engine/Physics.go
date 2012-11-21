package Engine

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

type Physics struct {
	BaseComponent
	Body  *chipmunk.Body
	Box   *chipmunk.BoxShape
	Shape *chipmunk.Shape
}

func NewPhysics(static bool, w, h float32) *Physics {
	var body *chipmunk.Body

	box := chipmunk.NewBox(vect.Vect{0, 0}, vect.Float(w), vect.Float(h))

	if static {
		body = chipmunk.NewBodyStatic()
	} else {
		body = chipmunk.NewBody(1, box.Moment(1))
	}

	p := &Physics{NewComponent(), body, box.GetAsBox(), box}
	body.UserData = p

	body.AddShape(box)
	return p
}

func NewPhysics2(static bool, shape *chipmunk.Shape) *Physics {
	var body *chipmunk.Body
	if static {
		body = chipmunk.NewBodyStatic()
	} else {
		body = chipmunk.NewBody(1, shape.ShapeClass.Moment(1))
	}

	p := &Physics{NewComponent(), body, shape.GetAsBox(), shape}
	body.UserData = p

	body.AddShape(shape)
	return p
}

func (p *Physics) Start() {
	pos := p.GameObject().Transform().WorldPosition()
	p.Body.SetAngle(vect.Float(p.GameObject().Transform().WorldRotation().Z) * RadianConst)
	p.Body.SetPosition(vect.Vect{vect.Float(pos.X), vect.Float(pos.Y)})

	if p.GameObject().Sprite != nil {
		p.GameObject().Sprite.UpdateShape()
		p.Body.UpdateShapes()
	}

	//p.Body.UpdateShapes()
	Space.AddBody(p.Body)
}

func (p *Physics) OnComponentBind(gobj *GameObject) {
	gobj.Physics = p
	p.Body.CallbackHandler = p
}

func (p *Physics) CollisionPreSolve(arbiter *chipmunk.Arbiter) bool {
	if p.gameObject == nil {
		return true
	}
	return onCollisionPreSolveGameObject(p.GameObject(), (*Arbiter)(arbiter))
}

func (p *Physics) CollisionEnter(arbiter *chipmunk.Arbiter) bool {
	if p.gameObject == nil {
		return true
	}
	return onCollisionEnterGameObject(p.GameObject(), (*Arbiter)(arbiter))
}

func (p *Physics) CollisionExit(arbiter *chipmunk.Arbiter) {
	if p.gameObject == nil {
		return
	}
	onCollisionExitGameObject(p.GameObject(), (*Arbiter)(arbiter))
}

func (p *Physics) CollisionPostSolve(arbiter *chipmunk.Arbiter) {
	if p.gameObject == nil {
		return
	}
	onCollisionPostSolveGameObject(p.GameObject(), (*Arbiter)(arbiter))
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
