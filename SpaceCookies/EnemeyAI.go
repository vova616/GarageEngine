package SpaceCookies

import (
	//"fmt"
	"github.com/vova616/GarageEngine/Engine"
	"math/rand"
)

type EnemeyType int

const (
	Enemey_Boss = iota
	Enemey_Cookie
)

type EnemeyAI struct {
	Engine.BaseComponent
	Target *Engine.GameObject
	Type   EnemeyType
}

func NewEnemeyAI(target *Engine.GameObject, typ EnemeyType) *EnemeyAI {
	return &EnemeyAI{BaseComponent: Engine.NewComponent(), Target: target, Type: typ}
}

func (ai *EnemeyAI) Start() {
	if ai.Target == nil {
		ai.Target = Player
	}

	isPlayerClose := func(distance float32) func() Engine.Command {
		return func() Engine.Command {
			if ai.GameObject() == nil || ai.Target.GameObject() == nil {
				return Engine.Close
			}
			myPos := ai.Transform().WorldPosition()
			targetPos := ai.Target.Transform().WorldPosition()
			if targetPos.Distance(myPos) < distance {
				return Engine.Continue
			}
			return Engine.Yield
		}
	}

	prepareForAttack := func() Engine.Command {
		if ai.GameObject() == nil || ai.Target.GameObject() == nil {
			return Engine.Close
		}

		ai.GameObject().Physics.Body.SetTorque(7000)
		return Engine.Continue
	}

	attack := func() Engine.Command {
		if ai.GameObject() == nil || ai.Target.GameObject() == nil {
			return Engine.Close
		}
		myPos := ai.Transform().WorldPosition()
		targetPos := ai.Target.Transform().WorldPosition()

		dir := targetPos.Sub(myPos)
		dir.Normalize()

		rnd := rand.Float32() * 0.5
		if rand.Float32() > 0.5 {
			rnd = -rnd
		}

		ai.GameObject().Physics.Body.AddForce((dir.X+rnd)*(200*rand.Float32()+50), (dir.Y+rnd)*(200*rand.Float32()+50))
		return Engine.Continue
	}

	randomMove := func() Engine.Command {
		if ai.GameObject() == nil || ai.Target.GameObject() == nil {
			return Engine.Close
		}

		myPos := ai.Transform().WorldPosition()
		targetPos := ai.Target.Transform().WorldPosition()
		if targetPos.Distance(myPos) < 500 {

			if rand.Float32() > 0.5 {

				dir := targetPos.Sub(myPos)
				dir.Normalize()

				rnd := rand.Float32() * 0.5
				if rand.Float32() > 0.5 {
					rnd = -rnd
				}

				ai.GameObject().Physics.Body.AddForce((-dir.X+rnd)*50, (-dir.Y+rnd)*50)
			} else {
				dir := targetPos.Sub(myPos)
				dir.Normalize()

				rnd := rand.Float32() * 0.5
				if rand.Float32() > 0.5 {
					rnd = -rnd
				}

				ai.GameObject().Physics.Body.AddForce((dir.X+rnd)*200, (dir.Y+rnd)*200)
			}
		}

		return Engine.Continue
	}

	sendCookies := func() Engine.Command {
		if ai.GameObject() == nil || ai.Target.GameObject() == nil {
			return Engine.Close
		}
		myPos := ai.Transform().WorldPosition()
		targetPos := ai.Target.Transform().WorldPosition()

		dir := targetPos.Sub(myPos)
		dir.Normalize()

		rnd := rand.Float32() * 0.2
		if rand.Float32() > 0.5 {
			rnd = -rnd
		}

		c := cookie.Clone()
		//c.Tag = CookieTag
		c.Transform().SetParent2(GameSceneGeneral.Layer2)
		size := 50 + rand.Float32()*100
		c.Transform().SetScalef(size, size)

		s := ai.Transform().WorldScale()
		s = s.Add(c.Transform().WorldScale())
		s = s.Mul2(0.5)
		p := myPos.Add(dir.Mul(s))

		c.Transform().SetPosition(p)
		c.GameObject().Physics.Body.AddForce((dir.X+rnd)*250, (dir.Y+rnd)*250)

		return Engine.Continue
	}

	appear := func() Engine.Command {
		if ai.GameObject() == nil {
			return Engine.Close
		}

		ai.Transform().SetPositionf(1500, 1500)

		return Engine.Continue
	}

	prepareForNextAttack := func() Engine.Command {
		if ai.GameObject() == nil || ai.Target.GameObject() == nil {
			return Engine.Close
		}

		ai.GameObject().Physics.Body.SetTorque(-10)
		ai.GameObject().Physics.Body.SetAngularVelocity(0)

		return Engine.Restart
	}

	co := false
	if co {
		Engine.StartCoroutine(func() {
			for {
				Engine.CoSleep(5)
				if !ai.GameObject().IsValid() {
					return
				}
				myPos := ai.Transform().WorldPosition()
				targetPos := ai.Target.Transform().WorldPosition()
				if targetPos.Distance(myPos) < 600 {
					dir := targetPos.Sub(myPos)
					dir.Normalize()

					rnd := rand.Float32() * 0.5
					if rand.Float32() > 0.5 {
						rnd = -rnd
					}

					ai.GameObject().Physics.Body.AddForce((dir.X+rnd)*50000, (dir.Y+rnd)*50000)
				}
			}

		})
	} else {
		if ai.Type == Enemey_Cookie {
			Engine.StartBehavior(Engine.SleepRand(5), isPlayerClose(600), prepareForAttack, Engine.Sleep(1.5), attack, Engine.WaitContinue(prepareForNextAttack, nil, 1.5))
		} else {
			Engine.StartBehavior(Engine.Sleep(120), appear, Engine.Sequence(Engine.SleepRand(0.5), isPlayerClose(800), randomMove, sendCookies))
		}
	}

}

func (ai *EnemeyAI) Update() {

}

func (sp *EnemeyAI) OnHit(enemey *Engine.GameObject, damager *DamageDealer) {

}

func (sp *EnemeyAI) OnDie(byTimer bool) {

	sxps := 4
	size := float32(0.5)
	if sp.Type == Enemey_Boss {
		sxps = 10
		size = 3
		Wall.Destroy()
		queenDead = true
	} else {
		CreatePowerUp(sp.Transform().WorldPosition())
	}

	for i := 0; i < sxps; i++ {
		n := Explosion.Clone()
		n.Transform().SetParent2(GameSceneGeneral.Layer1)
		n.Transform().SetWorldPosition(sp.Transform().WorldPosition())
		s := n.Transform().Scale()
		n.Transform().SetScale(s.Mul2((rand.Float32() * 3) + size))
		n.AddComponent(Engine.NewPhysics(false, 1, 1))

		n.Transform().SetRotationf(rand.Float32() * 360)
		rot := n.Transform().Rotation2D()
		n.Physics.Body.SetVelocity(-rot.X*15, -rot.Y*15)

		n.Physics.Body.SetMass(1)
		n.Physics.Shape.Group = 1
		n.Physics.Shape.IsSensor = true
	}

	sp.GameObject().Destroy()
}
