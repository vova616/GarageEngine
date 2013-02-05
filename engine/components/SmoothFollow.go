package components

import (
	"github.com/vova616/garageEngine/engine"
	"math"
	//"log"
)

type SmoothFollow struct {
	engine.BaseComponent
	Target *engine.GameObject
	Speed  float32
	MaxDis float32
}

func NewSmoothFollow(target *engine.GameObject, speed float32, maxdis float32) *SmoothFollow {
	return &SmoothFollow{engine.NewComponent(), target, speed, maxdis}
}

func (sp *SmoothFollow) Start() {
	if sp.Target == nil {
		sp.Target = sp.GameObject()
	}
}

func (sp *SmoothFollow) LateUpdate() {
	camera := engine.GetScene().SceneBase().Camera
	if camera != nil {
		myPos := engine.Vector{sp.Target.Transform().Position().X - float32(engine.Width/2), sp.Target.Transform().Position().Y - float32(engine.Height/2), 0}
		camPos := camera.Transform().Position()

		if sp.Speed > 0 {
			camPos = engine.Lerp(camPos, myPos, float32(engine.DeltaTime())*sp.Speed)
			disX := camPos.X - myPos.X
			disY := camPos.Y - myPos.Y
			if float32(math.Abs(float64(disX))) > sp.MaxDis {
				if disX < 0 {
					camPos.X = myPos.X - sp.MaxDis
				} else {
					camPos.X = myPos.X + sp.MaxDis
				}
			}
			if float32(math.Abs(float64(disY))) > sp.MaxDis {
				if disY < 0 {
					camPos.Y = myPos.Y - sp.MaxDis
				} else {
					camPos.Y = myPos.Y + sp.MaxDis
				}
			}
		} else {
			camPos = myPos
		}
		camera.Transform().SetPosition(camPos)
	}
}
