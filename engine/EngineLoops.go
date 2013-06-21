package engine

import (
	"github.com/vova616/chipmunk/vect"
)

func iter(objsp *[]*GameObject, f func(*GameObject)) {
	objs := *objsp
	for i := len(objs) - 1; i >= 0; i-- {
		obj := objs[i]
		if obj != nil {
			f(obj)
			//Checks if the objs array has been changed
			if i >= len(*objsp) {
				break
			}
			if obj != objs[i] {
				i++
			}
		}
	}
}

func preStepGameObject(g *GameObject) {
	if g.Physics != nil && g.active && !g.Physics.Body.IsStatic() && g.Physics.started() {
		pos := g.Transform().WorldPosition()
		angle := g.Transform().Angle() * RadianConst

		if g.Physics.Interpolate {
			//Interpolation check: if position/angle has been changed directly and not by the physics engine, change g.Physics.lastPosition/lastAngle
			if vect.Float(pos.X) != g.Physics.interpolatedPosition.X || vect.Float(pos.Y) != g.Physics.interpolatedPosition.Y {
				g.Physics.interpolatedPosition = vect.Vect{vect.Float(pos.X), vect.Float(pos.Y)}
				g.Physics.Body.SetPosition(g.Physics.interpolatedPosition)
			}
			if vect.Float(angle) != g.Physics.interpolatedAngle {
				g.Physics.interpolatedAngle = vect.Float(angle)
				g.Physics.Body.SetAngle(g.Physics.interpolatedAngle)
			}
		} else {
			var pPos vect.Vect
			pPos.X, pPos.Y = vect.Float(pos.X), vect.Float(pos.Y)

			g.Physics.Body.SetAngle(vect.Float(angle))
			g.Physics.Body.SetPosition(pPos)
		}
		g.Physics.lastPosition = g.Physics.Body.Position()
		g.Physics.lastAngle = g.Physics.Body.Angle()
	}
}

func postStepGameObject(g *GameObject) {
	if g.Physics != nil && g.active && !g.Physics.Body.IsStatic() && g.Physics.started() {
		/*
			When parent changes his position/rotation it changes his children position/rotation too but the physics engine thinks its in different position
			so we need to check how much it changed and apply to the new position/rotation so we wont fuck up things too much.

			Note:If position/angle is changed in between preStep and postStep it will be overrided.
		*/
		if CorrectWrongPhysics {
			b := g.Physics.Body
			angle := float32(b.Angle())
			lAngle := float32(g.Physics.lastAngle)
			lAngle += angle - lAngle

			pos := b.Position()
			lPos := g.Physics.lastPosition
			lPos.X += (pos.X - lPos.X)
			lPos.Y += (pos.Y - lPos.Y)

			if g.Physics.Interpolate {
				g.Physics.interpolatedAngle = vect.Float(lAngle)
				g.Physics.interpolatedPosition = lPos
			}

			b.SetPosition(lPos)
			b.SetAngle(g.Physics.interpolatedAngle)

			g.Transform().SetWorldRotationf(lAngle * DegreeConst)
			g.Transform().SetWorldPositionf(float32(lPos.X), float32(lPos.Y))
		} else {
			b := g.Physics.Body
			angle := b.Angle()
			pos := b.Position()

			if g.Physics.Interpolate {
				g.Physics.interpolatedAngle = angle
				g.Physics.interpolatedPosition = pos
			}

			g.Transform().SetWorldRotationf(float32(angle) * DegreeConst)
			g.Transform().SetWorldPositionf(float32(pos.X), float32(pos.Y))
		}
	}
}

func interpolateGameObject(g *GameObject) {
	if g.Physics != nil && g.Physics.Interpolate && g.active && !g.Physics.Body.IsStatic() && g.Physics.started() {
		nextPos := g.Physics.Body.Position()
		currPos := g.Physics.lastPosition

		nextAngle := g.Physics.Body.Angle()
		currAngle := g.Physics.lastAngle

		alpha := vect.Float(fixedTime / stepTime)
		x := currPos.X + ((nextPos.X - currPos.X) * alpha)
		y := currPos.Y + ((nextPos.Y - currPos.Y) * alpha)
		a := currAngle + ((nextAngle - currAngle) * alpha)
		g.Transform().SetWorldPositionf(float32(x), float32(y))
		g.Transform().SetWorldRotationf(float32(a) * DegreeConst)

		g.Physics.interpolatedAngle = a
		g.Physics.interpolatedPosition.X, g.Physics.interpolatedPosition.Y = x, y
	}
}

func drawGameObject(gameObject *GameObject) {
	if !gameObject.active {
		return
	}
	//mat := gameObject.Transform().Matrix()

	//gl.LoadMatrixf(mat.Ptr())

	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].Draw()
		}
	}
}

func startGameObject(gameObject *GameObject) {
	if !gameObject.active {
		return
	}
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if !comps[i].started() {
			comps[i].setStarted(true)
			comps[i].Start()
		}
	}
}

func destoyGameObject(gameObject *GameObject) {
	if gameObject.destoryMark {
		gameObject.destroy()
	}
}

func onCollisionPreSolveGameObject(gameObject *GameObject, arb Arbiter) bool {
	if !gameObject.active {
		return true
	}
	l := len(gameObject.components)
	comps := gameObject.components

	b := true
	for i := l - 1; i >= 0; i-- {
		b = b && comps[i].OnCollisionPreSolve(arb)
	}
	return b
}

func onCollisionPostSolveGameObject(gameObject *GameObject, arb Arbiter) {
	if !gameObject.active {
		return
	}
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].OnCollisionPostSolve(arb)
		}
	}
}

func onCollisionEnterGameObject(gameObject *GameObject, arb Arbiter) bool {
	if gameObject == nil || !gameObject.active {
		return true
	}
	l := len(gameObject.components)
	comps := gameObject.components

	b := true
	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			b = b && comps[i].OnCollisionEnter(arb)
		}
	}
	return b
}

func onCollisionExitGameObject(gameObject *GameObject, arb Arbiter) {
	if gameObject == nil || !gameObject.active {
		return
	}
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].OnCollisionExit(arb)
		}
	}
}

func onMouseEnterGameObject(gameObject *GameObject, arb Arbiter) bool {
	if gameObject == nil || !gameObject.active {
		return true
	}
	l := len(gameObject.components)
	comps := gameObject.components

	b := true
	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			b = b && comps[i].OnMouseEnter(arb)
		}
	}
	return b
}

func onMouseExitGameObject(gameObject *GameObject, arb Arbiter) {
	if gameObject == nil || !gameObject.active {
		return
	}
	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].OnMouseExit(arb)
		}
	}
}

func udpateGameObject(gameObject *GameObject) {
	if !gameObject.active {
		return
	}

	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].Update()
		}
	}
}

func lateudpateGameObject(gameObject *GameObject) {
	if !gameObject.active {
		return
	}

	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].LateUpdate()
		}
	}
}

func fixedUdpateGameObject(gameObject *GameObject) {
	if !gameObject.active || gameObject.Physics == nil {
		return
	}

	l := len(gameObject.components)
	comps := gameObject.components

	for i := l - 1; i >= 0; i-- {
		if comps[i].started() {
			comps[i].FixedUpdate()
		}
	}
}
