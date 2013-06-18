package game

import (
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/input"
	//"log"
	//"github.com/go-gl/glfw"
	//c "github.com/vova616/chipmunk"
	//. "github.com/vova616/chipmunk/vect"
	"math/rand"
)

type MouseDebugger struct {
	engine.BaseComponent
}

func NewMouseDebugger() *MouseDebugger {
	return &MouseDebugger{engine.NewComponent()}
}

func (m *MouseDebugger) Update() {
	if input.MouseDown(input.MouseMiddle) {
		if queenDead {
			mousePosition := m.Transform().WorldPosition()

			c := cookie.Clone()
			//c.Tag = CookieTag
			c.Transform().SetParent2(GameSceneGeneral.Layer2)
			size := 25 + rand.Float32()*100
			c.Transform().SetPosition(mousePosition)
			c.Transform().SetScalef(size, size)
		}
	}
	if input.MouseDown(input.MouseRight) {

		mousePosition := m.Transform().WorldPosition()

		b := defender.Clone()
		/*
			phx := b.AddComponent(NewPhysics(false, 50, 50)).(*Physics)
			phx.Shape.SetFriction(0.5)
			//phx.Shape.Group = 2
			phx.Shape.SetElasticity(0.5)
		*/
		b.Transform().SetParent2(GameSceneGeneral.Layer2)
		b.Transform().SetWorldPosition(mousePosition)
		b.Transform().SetScalef(50, 50)

	}

	if input.KeyPress('R') {
		MyClient.SendRespawn()
	}

	if input.MouseWheelDelta != 0 && engine.Debug {
		engine.CurrentCamera().SetSize(engine.CurrentCamera().Size() + float32(-input.MouseWheelDelta))
	}

	if queenDead {
		if input.KeyPress(input.KeyF1) {
			PowerUpShip(HP)
		}
		if input.KeyPress(input.KeyF2) {
			PowerUpShip(Damage)
		}
		if input.KeyPress(input.KeyF3) {
			PowerUpShip(Range)
		}
		if input.KeyPress(input.KeyF4) {
			PowerUpShip(Speed)
		}
	}
}
