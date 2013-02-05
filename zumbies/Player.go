package zumbies

import (
	"github.com/vova616/garageEngine/engine"
	//"log"
	//"github.com/vova616/garageEngine/engine/input"
)

type Player struct {
	engine.BaseComponent
}

func NewPlayer() *Player {
	return &Player{engine.NewComponent()}
}

func (sp *Player) Update() {

}
