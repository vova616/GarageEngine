package Components

import (
	"github.com/vova616/GarageEngine/Engine"
	//"log"
)

type SmoothFollow struct {
	Engine.BaseComponent
	Target *Engine.GameObject
}

func NewSmoothFollow(target *Engine.GameObject) *SmoothFollow {
	return &SmoothFollow{Engine.NewComponent(), target}
}

func (sp *SmoothFollow) Start() {
	if sp.Target == nil {
		sp.Target = sp.GameObject()
	}
}

func (sp *SmoothFollow) LateUpdate() {
	camera := Engine.GetScene().SceneBase().Camera
	if camera != nil {
		myPos := Engine.Vector{sp.Target.Transform().Position().X - float32(Engine.Width/2), sp.Target.Transform().Position().Y - float32(Engine.Height/2), 0}
		camPos := camera.Transform().Position()
		myPos = Engine.Lerp(camPos, myPos, Engine.DeltaTime()*3)
		camera.Transform().SetPosition(myPos)
	}
}
