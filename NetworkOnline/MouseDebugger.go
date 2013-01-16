package NetworkOnline

import (
	"github.com/vova616/GarageEngine/Engine"
	"github.com/vova616/GarageEngine/Engine/Input"
	//"log" 
	"github.com/go-gl/glfw"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

type MouseDebugger struct {
	Engine.BaseComponent
}

func NewMouseDebugger() *MouseDebugger {
	return &MouseDebugger{Engine.NewComponent()}
}

func (m *MouseDebugger) Update() {
	if Input.MouseDown(glfw.MouseLeft) {

		mousePosition := m.Transform().WorldPosition()

		sprite3 := Engine.NewGameObject("Sprite")
		sprite3.AddComponent(Engine.NewSprite(cir))
		sprite3.Transform().SetParent2(GameSceneGeneral.Layer2)
		sprite3.Transform().SetWorldPosition(mousePosition)

		sprite3.Transform().SetScale(Engine.NewVector2(30, 30))

		phx := sprite3.AddComponent(Engine.NewPhysics2(false, chipmunk.NewCircle(vect.Vect{0, 0}, 15))).(*Engine.Physics)
		phx.Shape.SetFriction(0.5)
		//phx.Shape.Group = 1
		phx.Shape.SetElasticity(0.5)
	}
	if Input.MouseDown(glfw.MouseRight) {

		mousePosition := m.Transform().WorldPosition()

		sprite3 := Engine.NewGameObject("Sprite")
		sprite3.AddComponent(Engine.NewSprite(box))
		sprite3.Transform().SetParent2(GameSceneGeneral.Layer2)
		sprite3.Transform().SetWorldPosition(mousePosition)

		sprite3.Transform().SetScale(Engine.NewVector2(30, 30))
		phx := sprite3.AddComponent(Engine.NewPhysics(false, 50, 50)).(*Engine.Physics)
		phx.Shape.SetFriction(0.5)
		//phx.Shape.Group = 2
		phx.Shape.SetElasticity(0.5)
	}
}
