package NetworkOnline

import (
	. "github.com/vova616/GarageEngine/Engine"
	//"log"
)

type Rotator struct {
	BaseComponent
}

func NewRotator() *Rotator {
	return &Rotator{NewComponent()}
}

func (sp *Rotator) Update() {
	//log.Panicln("Rotate")
	rot := sp.Transform().Rotation()
	sp.Transform().SetRotation(rot.Add(NewVector3(0, 1, 0)))
}
