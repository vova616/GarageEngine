package Components

import (
	. "github.com/vova616/GarageEngine/Engine"
	//"fmt"
	)

type GUI struct {
	BaseComponent
}

func NewGUI() *GUI {
	return &GUI{NewComponent()}
} 

func (m *GUI) Update() {
	parent := m.Transform().Parent()
	if parent != nil {
		m.Transform().SetPosition(parent.Position().Mul(MinusOne))
		//m.Transform().SetScale(NewVector3(1,1,1))
		//fmt.Println(m.Transform().Position())
	}
}