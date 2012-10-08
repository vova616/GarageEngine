package SpaceCookies

import (
	//"fmt"
	engine "github.com/vova616/GarageEngine/Engine"
	"math/rand"
)

type EnemeyAI struct {
	engine.BaseComponent
	Target *engine.GameObject
}

func NewEnemeyAI(target *engine.GameObject) *EnemeyAI {
	return &EnemeyAI{BaseComponent: engine.NewComponent(), Target: target}
}

func (ai *EnemeyAI) Start() {
	if ai.Target == nil {
		ai.Target = Player
	}

	isPlayerClose := func() engine.Command {
		if ai.GameObject() == nil || ai.Target.GameObject() == nil {
			return engine.Close
		}
		myPos := ai.Transform().WorldPosition()
		targetPos := ai.Target.Transform().WorldPosition()
		if targetPos.Distance(myPos) < 600 {
			return engine.Continue
		}
		return engine.Yield
	}

	prepareForAttack := func() engine.Command {
		if ai.GameObject() == nil || ai.Target.GameObject() == nil {
			return engine.Close
		}

		ai.GameObject().Physics.Body.SetTorque(7000)
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

		ai.GameObject().Physics.Body.AddForce((dir.X+rnd)*50000, (dir.Y+rnd)*50000)
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
		engine.StartBehavior(engine.SleepRand(5), isPlayerClose, prepareForAttack, engine.Sleep(1.5), attack, engine.WaitContinue(prepareForNextAttack, nil, 1.5))
	}

}

func (ai *EnemeyAI) Update() {

}

func (sp *EnemeyAI) OnHit(enemey *engine.GameObject, damager *DamageDealer) {

}

func (sp *EnemeyAI) OnDie() {
	for i := 0; i < 10; i++ {
		n := Explosion.Clone()
		n.Transform().SetParent2(GameSceneGeneral.Layer1)
		n.Transform().SetWorldPosition(sp.Transform().WorldPosition())
		s := n.Transform().Scale()
		n.Transform().SetScale(s.Mul2((rand.Float32() * 3) + 0.5))
		n.AddComponent(engine.NewPhysics(false, 1, 1))

		n.Transform().SetRotationf(0, 0, rand.Float32()*360)
		rot := n.Transform().Rotation2D()
		n.Physics.Body.SetVelocity(-rot.X*15, -rot.Y*15)

		n.Physics.Body.SetMass(1)
		n.Physics.Shape.Group = 1
		n.Physics.Shape.IsSensor = true
	}

	CreatePowerUp(sp.Transform().WorldPosition())

	sp.GameObject().Destroy()
}
