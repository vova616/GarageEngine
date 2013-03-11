package networkOnline

import (
	"github.com/vova616/GarageEngine/engine"
	//"log"
)

type Rotator struct {
	engine.BaseComponent
}

func NewRotator() *Rotator {
	return &Rotator{engine.NewComponent()}
}

func (sp *Rotator) Update() {
	//log.Panicln("Rotate")
	rot := sp.Transform().Rotation()
	sp.Transform().SetRotation(rot.Add(engine.NewVector3(0, 1, 0)))
}
