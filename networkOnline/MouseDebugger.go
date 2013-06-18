package networkOnline

import (
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/input"
	//"log"
	"github.com/go-gl/glfw"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

type MouseDebugger struct {
	engine.BaseComponent
}

func NewMouseDebugger() *MouseDebugger {
	return &MouseDebugger{engine.NewComponent()}
}

func (m *MouseDebugger) Update() {
	if input.MouseDown(glfw.MouseLeft) {

		mousePosition := m.Transform().WorldPosition()

		sprite3 := engine.NewGameObject("Sprite")
		sprite3.AddComponent(engine.NewSprite(cir))
		sprite3.Transform().SetParent2(GameSceneGeneral.Layer2)
		sprite3.Transform().SetWorldPosition(mousePosition)

		sprite3.Transform().SetScale(engine.NewVector2(30, 30))

		phx := sprite3.AddComponent(engine.NewPhysicsShape(false, chipmunk.NewCircle(vect.Vect{0, 0}, 15))).(*engine.Physics)
		phx.Shape.SetFriction(0.5)
		//phx.Shape.Group = 1
		phx.Shape.SetElasticity(0.5)
	}
	if input.MouseDown(glfw.MouseRight) {

		mousePosition := m.Transform().WorldPosition()

		sprite3 := engine.NewGameObject("Sprite")
		sprite3.AddComponent(engine.NewSprite(box))
		sprite3.Transform().SetParent2(GameSceneGeneral.Layer2)
		sprite3.Transform().SetWorldPosition(mousePosition)

		sprite3.Transform().SetScale(engine.NewVector2(30, 30))
		phx := sprite3.AddComponent(engine.NewPhysics(false)).(*engine.Physics)
		phx.Shape.SetFriction(0.5)
		//phx.Shape.Group = 2
		phx.Shape.SetElasticity(0.5)
	}
}
