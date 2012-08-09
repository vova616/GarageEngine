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
	sp.Transform().SetRotation(sp.Transform().Rotation().Add(NewVector3(0,1,0)))
}