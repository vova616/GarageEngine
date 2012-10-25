package SpaceCookies

import (
	. "github.com/vova616/GarageEngine/Engine"
	"github.com/vova616/GarageEngine/Engine/Input"
	//"log" 
	//"github.com/jteeuwen/glfw"
	//c "github.com/vova616/chipmunk"
	//. "github.com/vova616/chipmunk/vect"
	"math/rand"
)

type MouseDebugger struct {
	BaseComponent
}

func NewMouseDebugger() *MouseDebugger {
	return &MouseDebugger{NewComponent()}
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
			c.Transform().SetScalef(size, size, 1)
		}
	}
	if Input.MouseDown(Input.MouseRight) {

		mousePosition := m.Transform().WorldPosition()

		sprite3 := NewGameObject("Sprite")
		ds := NewDestoyable(30, 3)
		ds.SetDestroyTime(5)
		sprite3.AddComponent(ds)
		sprite3.AddComponent(NewSprite(box))
		sprite3.Transform().SetParent2(GameSceneGeneral.Layer2)
		sprite3.Transform().SetWorldPosition(mousePosition)
		sprite3.Tag = CookieTag
		sprite3.Transform().SetScale(NewVector2(50, 50))
		phx := sprite3.AddComponent(NewPhysics(false, 50, 50)).(*Physics)
		phx.Shape.SetFriction(0.5)
		//phx.Shape.Group = 2
		phx.Shape.SetElasticity(0.5)
	}

	if Input.KeyPress('R') {
		LoadScene(GameSceneGeneral)
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
