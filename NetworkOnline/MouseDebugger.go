package NetworkOnline

import (
	. "github.com/vova616/GarageEngine/Engine"
	in "github.com/vova616/GarageEngine/Engine/Input"
	//"log" 
	"github.com/jteeuwen/glfw"
	c "github.com/vova616/chipmunk"
	. "github.com/vova616/chipmunk/vect"
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

		sprite3 := NewGameObject("Sprite")
		sprite3.AddComponent(NewSprite(cir))
		sprite3.Transform().SetParent2(GameSceneGeneral.Layer2)
		sprite3.Transform().SetWorldPosition(mousePosition)

		sprite3.Transform().SetScale(NewVector2(30, 30))

		phx := sprite3.AddComponent(NewPhysics2(false, c.NewCircle(Vect{0, 0}, 15))).(*Physics)
		phx.Shape.SetFriction(0.5)
		//phx.Shape.Group = 1
		phx.Shape.SetElasticity(0.5)
	}
	if in.MouseDown(glfw.MouseRight) {

		mousePosition := m.Transform().WorldPosition()

		sprite3 := NewGameObject("Sprite")
		sprite3.AddComponent(NewSprite(box))
		sprite3.Transform().SetParent2(GameSceneGeneral.Layer2)
		sprite3.Transform().SetWorldPosition(mousePosition)

		sprite3.Transform().SetScale(NewVector2(30, 30))
		phx := sprite3.AddComponent(NewPhysics(false, 50, 50)).(*Physics)
		phx.Shape.SetFriction(0.5)
		//phx.Shape.Group = 2
		phx.Shape.SetElasticity(0.5)
	}
}
