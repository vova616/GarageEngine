package NetworkOnline

import (
	"github.com/vova616/GarageEngine/Engine"
	//"log"
)

type Rotator struct {
	Engine.BaseComponent
}

func NewRotator() *Rotator {
	return &Rotator{Engine.NewComponent()}
}

func (sp *Rotator) Update() {
	//log.Panicln("Rotate")
	rot := sp.Transform().Rotation()
	sp.Transform().SetRotation(rot.Add(Engine.NewVector3(0, 1, 0)))
}
