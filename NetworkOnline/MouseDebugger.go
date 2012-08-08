package NetworkOnline

import (
	. "../Engine"
	in "../Engine/Input"
	//"log"
	"github.com/jteeuwen/glfw"
	c "chipmunk"
	. "chipmunk/vect"
)

type MouseDebugger struct {
	BaseComponent
}

func NewMouseDebugger() *MouseDebugger {
	return &MouseDebugger{NewComponent()}
}

func (m *MouseDebugger) Update() {
	if in.MouseDown(glfw.MouseLeft) {
		sprite3 := NewGameObject("Sprite")
		sprite3.AddComponent(NewSprite(cir))
		sprite3.Transform().SetParent2(GameSceneGeneral.Layer2)
		sprite3.Transform().SetWorldPosition(m.GameObject().Transform().WorldPosition())
		
		
		
		sprite3.Transform().SetScale(NewVector2(30, 30))
		
		
		
		phx := sprite3.AddComponent(NewPhysics2(false, c.NewCircle(Vect{0,0},Float(15)))).(*Physics)
		phx.Shape.SetFriction(0.5) 
		phx.Shape.SetElasticity(0.5)
	}
	if in.MouseDown(glfw.MouseRight) {
		sprite3 := NewGameObject("Sprite")
		sprite3.AddComponent(NewSprite(box))
		sprite3.Transform().SetParent2(GameSceneGeneral.Layer2)
		sprite3.Transform().SetWorldPosition(m.GameObject().Transform().WorldPosition())
		
		sprite3.Transform().SetScale(NewVector2(30, 30))
		phx := sprite3.AddComponent(NewPhysics(false,50,50)).(*Physics)
		phx.Shape.SetFriction(0.5) 
		phx.Shape.SetElasticity(0.5)
	}
}
