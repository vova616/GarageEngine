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
		if ai.GameObject() == nil {
			return engine.Close
		}
		myPos := ai.Transform().WorldPosition()
		targetPos := ai.Target.Transform().WorldPosition()
		if targetPos.Distance(myPos) < 600 {
			return engine.Continue
		}
		return engine.Yield
	}

	attack := func() engine.Command {
		if ai.GameObject() == nil {
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
		return engine.Restart
	}

	co := false
	if co {
		engine.StartCoroutine(func() {
			for {
				engine.Wait(5)
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
		engine.StartBehavior(engine.SleepRand(5), isPlayerClose, attack)
	}

}

func (ai *EnemeyAI) Update() {

}
