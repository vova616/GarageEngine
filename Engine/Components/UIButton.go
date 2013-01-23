package Components

import (
	//"gl"
	"github.com/vova616/GarageEngine/Engine"
	"github.com/vova616/GarageEngine/Engine/Input"
)

type UIButton struct {
	Engine.BaseComponent
	onPressCallback func()
	onHoverCallback func(bool)
	mouseOn         bool
}

func NewUIButton(ClickCallback func(), HoverCallback func(bool)) *UIButton {
	return &UIButton{Engine.NewComponent(), ClickCallback, HoverCallback, false}
}

func (btn *UIButton) Update() {
	if btn.mouseOn {
		if btn.onPressCallback != nil && Input.MousePress(Input.Mouse1) {
			btn.onPressCallback()
		}
	}
}

func (btn *UIButton) OnMouseEnter(arbiter Engine.Arbiter) bool {
	btn.mouseOn = true
	if btn.onHoverCallback != nil {
		btn.onHoverCallback(true)
	}
	return true
}

func (btn *UIButton) OnMouseExit(arbiter Engine.Arbiter) {
	btn.mouseOn = false
	if btn.onHoverCallback != nil {
		btn.onHoverCallback(false)
	}
}
