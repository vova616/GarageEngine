package game

import (
	"github.com/vova616/GarageEngine/engine"
	//"reflect"
)

type Background struct {
	engine.BaseComponent
	sprite *engine.Sprite
}

func NewBackground(sprite *engine.Sprite) *Background {
	return &Background{BaseComponent: engine.NewComponent(), sprite: sprite}
}

func (sp *Background) Draw() {
	sp.sprite.Render = true
	sp.sprite.DrawScreen()
	sp.sprite.Render = false
}
