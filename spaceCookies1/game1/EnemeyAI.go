package game

import (
	//"fmt"
	"github.com/vova616/garageEngine/engine"
	"math/rand"
)

type EnemeyType int

const (
	Enemey_Boss = iota
	Enemey_Cookie
)

type EnemeyAI struct {
	engine.BaseComponent
	Target *engine.GameObject
	Type   EnemeyType
}

func NewEnemeyAI(target *engine.GameObject, typ EnemeyType) *EnemeyAI {
	return &EnemeyAI{BaseComponent: engine.NewComponent(), Target: target, Type: typ}
}

func (ai *EnemeyAI) Start() {
	if ai.Target == nil {
		ai.Target = Player
	}

	isPlayerClose := func(distance float32) func() engine.Command {
		return func() engine.Command {
			if ai.GameObject() == nil || ai.Target.GameObject() == nil {
				return engine.Close
			}
			myPos := ai.Transform().WorldPosition()
			targetPos := ai.Target.Transform().WorldPosition()
			if targetPos.Distance(myPos) < distance {
				return engine.Continue
			}
			return engine.Yield
		}
	}

	prepareForAttack := func() engine.Command {
		if ai.GameObject() == nil || ai.Target.GameObject() == nil {
			return engine.Close
		}

		ai.GameObject().Physics.Body.SetTorque(10000)
		return engine.Continue
	}

	attack := func() engine.Command {
		if ai.GameObject() == nil || ai.Target.GameObject() == nil {
			return engine.Close
		}
		myPos := ai.Transform().WorldPosition()
		targetPos := ai.Target.Transform().WorldPosition()

		dir := targetPos.Sub(myPos)
		dir.Normalize()

		rnd := rand.Float32() * 0.5
		if rand.Float32() > 0.5 {
			rnd = -rnd
		}

		attackSpeed := float32(70000)
		minAttackSpeed := float32(20000)

		attackSpeed -= minAttackSpeed

		ai.GameObject().Physics.Body.AddForce((dir.X+rnd)*((attackSpeed*rand.Float32())+minAttackSpeed), (dir.Y+rnd)*((attackSpeed*rand.Float32())+minAttackSpeed))
		return engine.Continue
	}

	randomMove := func() engine.Command {
		if ai.GameObject() == nil || ai.Target.GameObject() == nil {
			return engine.Close
		}
		attackSpeed := float32(40000)
		moveSpeed := float32(20000)
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

				ai.GameObject().Physics.Body.AddForce((-dir.X+rnd)*moveSpeed, (-dir.Y+rnd)*moveSpeed)
			} else {
				dir := targetPos.Sub(myPos)
				dir.Normalize()

				rnd := rand.Float32() * 0.5
				if rand.Float32() > 0.5 {
					rnd = -rnd
				}

				ai.GameObject().Physics.Body.AddForce((dir.X+rnd)*attackSpeed, (dir.Y+rnd)*attackSpeed)
			}
		}

		return engine.Continue
	}

	sendCookies := func() engine.Command {
		if ai.GameObject() == nil || ai.Target.GameObject() == nil {
			return engine.Close
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

		attackSpeed := float32(70000)

		c.GameObject().Physics.Body.AddForce((dir.X+rnd)*attackSpeed, (dir.Y+rnd)*attackSpeed)

		return engine.Continue
	}

	appear := func() engine.Command {
		if ai.GameObject() == nil {
			return engine.Close
		}

		ai.Transform().SetPositionf(1500, 1500)

		return engine.Continue
	}

	prepareForNextAttack := func() engine.Command {
		if ai.GameObject() == nil || ai.Target.GameObject() == nil {
			return engine.Close
		}

		ai.GameObject().Physics.Body.SetTorque(-10)
		ai.GameObject().Physics.Body.SetAngularVelocity(0)

		return engine.Restart
	}

	co := false
	if co {
		engine.StartCoroutine(func() {
			for {
				engine.CoSleep(5)
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
			engine.StartBehavior(engine.SleepRand(5), isPlayerClose(600), prepareForAttack, engine.Sleep(1.5), attack, engine.WaitContinue(prepareForNextAttack, nil, 1.5))
		} else {
			engine.StartBehavior(engine.Sleep(120), appear, engine.Sequence(engine.SleepRand(0.5), isPlayerClose(800), randomMove, sendCookies))
		}
	}

}

func (ai *EnemeyAI) Update() {

}

func (sp *EnemeyAI) OnHit(enemey *engine.GameObject, damager *DamageDealer) {

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
		n.AddComponent(engine.NewPhysics(false, 1, 1))

		n.Transform().SetRotationf(rand.Float32() * 360)
		rot := n.Transform().Direction()
		n.Physics.Body.SetVelocity(-rot.X*15, -rot.Y*15)

		n.Physics.Body.SetMass(1)
		n.Physics.Shape.Group = 1
		n.Physics.Shape.IsSensor = true
	}

	sp.GameObject().Destroy()
}
