package Components

import (
	"github.com/vova616/GarageEngine/Engine"
	"math"
	//"log"
)

type SmoothFollow struct {
	Engine.BaseComponent
	Target *Engine.GameObject
	Speed  float32
	MaxDis float32
}

func NewSmoothFollow(target *Engine.GameObject, speed float32, maxdis float32) *SmoothFollow {
	return &SmoothFollow{Engine.NewComponent(), target, speed, maxdis}
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
		camNPos := Engine.Lerp(camPos, myPos, Engine.DeltaTime()*sp.Speed)
		disX := camNPos.X - myPos.X
		disY := camNPos.Y - myPos.Y
		if float32(math.Abs(float64(disX))) > sp.MaxDis {
			if disX < 0 {
				camNPos.X = myPos.X - sp.MaxDis
			} else {
				camNPos.X = myPos.X + sp.MaxDis
			}
		}
		if float32(math.Abs(float64(disY))) > sp.MaxDis {
			if disY < 0 {
				camNPos.Y = myPos.Y - sp.MaxDis
			} else {
				camNPos.Y = myPos.Y + sp.MaxDis
			}
		}

		camera.Transform().SetPosition(camNPos)
	}
}
