package engine

import (
	"github.com/vova616/chipmunk"
)

type Arbiter struct {
	*chipmunk.Arbiter
	gameObjectA *GameObject //Always the owner
	gameObjectB *GameObject //Always the other target
	Swapped     bool
}

func newArbiter(arb *chipmunk.Arbiter, owner *GameObject) Arbiter {
	newArb := Arbiter{Arbiter: arb}
	pa, ba := arb.BodyA.CallbackHandler.(*Physics)
	if ba {
		newArb.gameObjectA = pa.GameObject()
	} else {
		panic("CallbackHandler is not Physics")
	}
	pb, bb := arb.BodyB.CallbackHandler.(*Physics)
	if bb {
		newArb.gameObjectB = pb.GameObject()
	} else {
		panic("CallbackHandler is not Physics")
	}

	//Find owner
	if newArb.gameObjectA != owner {
		newArb.gameObjectA, newArb.gameObjectB = newArb.gameObjectB, newArb.gameObjectA
		newArb.Swapped = true
	}
	return newArb
}

func (arbiter *Arbiter) GameObjectA() *GameObject {
	return arbiter.gameObjectA
}

func (arbiter *Arbiter) GameObjectB() *GameObject {
	return arbiter.gameObjectB
}

/*
func (arbiter *Arbiter) ShapeA() *chipmunk.Shape {
	if arbiter.Swapped {
		return arbiter.Arbiter.ShapeB
	}
	return arbiter.Arbiter.ShapeA
}

func (arbiter *Arbiter) ShapeB() *chipmunk.Shape {
	if arbiter.Swapped {
		return arbiter.Arbiter.ShapeA
	}
	return arbiter.Arbiter.ShapeB
}
*/
func (arbiter *Arbiter) ShapeB() *GameObject {
	return arbiter.gameObjectB
}

func (arbiter *Arbiter) Normal(contact *chipmunk.Contact) Vector {
	normal := contact.Normal()
	if arbiter.Swapped {
		return NewVector2(float32(-normal.X), float32(-normal.Y))
	}
	return NewVector2(float32(normal.X), float32(normal.Y))
}
