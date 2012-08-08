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
	t := sp.GameObject().Transform()
	
	if Input.KeyDown('A') {	
		t.Translate2(5,0,0)
	}
	if Input.KeyDown('D') {	
		t.Translate2(-5,0,0)
	}
	if Input.KeyDown('S') {	
		t.Translate2(0,5,0)
	}
	if Input.KeyDown('W') {	
		t.Translate2(0,-5,0)
	}
}