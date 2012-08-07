package Components

import (
	. "../../Engine"
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
	}
}


