package SpaceCookies

import (
	"github.com/vova616/GarageEngine/Engine"
	"github.com/vova616/GarageEngine/Engine/Input"
	//"log" 
	//"github.com/jteeuwen/glfw"
	//c "github.com/vova616/chipmunk"
	//. "github.com/vova616/chipmunk/vect"
	"math/rand"
)

type MouseDebugger struct {
	Engine.BaseComponent
}

func NewMouseDebugger() *MouseDebugger {
	return &MouseDebugger{Engine.NewComponent()}
}

func (m *MouseDebugger) Update() {
	if Input.MouseDown(Input.MouseMiddle) {
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
	if Input.MouseDown(Input.MouseRight) {

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

	if Input.KeyPress('R') {
		Engine.LoadScene(GameSceneGeneral)
	}
	if queenDead {
		if Input.KeyPress(Input.KeyF1) {
			PowerUpShip(HP)
		}
		if Input.KeyPress(Input.KeyF2) {
			PowerUpShip(Damage)
		}
		if Input.KeyPress(Input.KeyF3) {
			PowerUpShip(Range)
		}
		if Input.KeyPress(Input.KeyF4) {
			PowerUpShip(Speed)
		}
	}
}
