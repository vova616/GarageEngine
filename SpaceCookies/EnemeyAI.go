package SpaceCookies

import (
	//"fmt"
	. "github.com/vova616/GarageEngine/Engine"
	"math/rand"
)

type EnemeyAI struct {
	BaseComponent
	Target *GameObject
}

func NewEnemeyAI(target *GameObject) *EnemeyAI {
	return &EnemeyAI{BaseComponent: NewComponent(), Target: target}
}

func (ai *EnemeyAI) Start() {
	if ai.Target == nil {
		ai.Target = Player
	}
	StartCoroutine(func() {
		for ai.GameObject() != nil {
			Wait(rand.Float32() * 5)
			if ai.GameObject() == nil {
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
}

func (ai *EnemeyAI) Update() {

}
