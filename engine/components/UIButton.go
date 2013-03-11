package components

import (
	//"gl"
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/input"
)

type UIButton struct {
	engine.BaseComponent
	onPressCallback func()
	onHoverCallback func(bool)
	mouseOn         bool
}

func NewUIButton(ClickCallback func(), HoverCallback func(bool)) *UIButton {
	return &UIButton{engine.NewComponent(), ClickCallback, HoverCallback, false}
}

func (btn *UIButton) Update() {
	if btn.mouseOn {
		if btn.onPressCallback != nil && input.MousePress(input.Mouse1) {
			btn.onPressCallback()
		}
	}
}

func (btn *UIButton) OnMouseEnter(arbiter engine.Arbiter) bool {
	btn.mouseOn = true
	if btn.onHoverCallback != nil {
		btn.onHoverCallback(true)
	}
	return true
}

func (btn *UIButton) OnMouseExit(arbiter engine.Arbiter) {
	btn.mouseOn = false
	if btn.onHoverCallback != nil {
		btn.onHoverCallback(false)
	}
}
