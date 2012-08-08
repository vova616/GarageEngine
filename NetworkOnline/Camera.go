package NetworkOnline

import (
	. "../Engine"
	//"log"
	"../Engine/Input"
)

type Camera struct {
	BaseComponent
}

func NewCamera() *Camera {
	return &Camera{NewComponent()}
}

func (sp *Camera) Update() {
	if Input.KeyDown('A') {	
		sp.GameObject().Transform().Translate2(5,0,0)
	}
	if Input.KeyDown('D') {	
		sp.GameObject().Transform().Translate2(-5,0,0)
	}
	if Input.KeyDown('S') {	
		sp.GameObject().Transform().Translate2(0,5,0)
	}
	if Input.KeyDown('W') {	
		sp.GameObject().Transform().Translate2(0,-5,0)
	}
}