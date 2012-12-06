package Game

import (
	"github.com/vova616/GarageEngine/Engine"
	//"reflect"
)

type Background struct {
	Engine.BaseComponent
	sprite *Engine.Sprite
}

func NewBackground(sprite *Engine.Sprite) *Background {
	return &Background{BaseComponent: Engine.NewComponent(), sprite: sprite}
}

func (sp *Background) Draw() {
	sp.sprite.Render = true
	sp.sprite.DrawScreen()
	sp.sprite.Render = false
}
