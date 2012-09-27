package SpaceCookies

import (
	. "github.com/vova616/GarageEngine/Engine"
	in "github.com/vova616/GarageEngine/Engine/Input"
	//"log" 
	"github.com/jteeuwen/glfw"
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
	if in.MouseDown(glfw.MouseLeft) {

		mousePosition := m.Transform().WorldPosition()

		c := cookie.Clone()
		//c.Tag = CookieTag
		c.Transform().SetParent2(GameSceneGeneral.Layer2)
		size := 25 + rand.Float32()*100
		c.Transform().SetPosition(mousePosition)
		c.Transform().SetScalef(size, size, 1)
	}
	if in.MouseDown(glfw.MouseRight) {

		mousePosition := m.Transform().WorldPosition()

		sprite3 := NewGameObject("Sprite")
		sprite3.AddComponent(NewSprite(box))
		sprite3.Transform().SetParent2(GameSceneGeneral.Layer2)
		sprite3.Transform().SetWorldPosition(mousePosition)
		sprite3.Tag = CookieTag
		sprite3.Transform().SetScale(NewVector2(30, 30))
		phx := sprite3.AddComponent(NewPhysics(false, 50, 50)).(*Physics)
		phx.Shape.SetFriction(0.5)
		//phx.Shape.Group = 2
		phx.Shape.SetElasticity(0.5)
	}
}
