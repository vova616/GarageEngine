package SpaceCookies

import (
	. "github.com/vova616/GarageEngine/Engine"
	//"reflect"
)

type Background struct {
	BaseComponent
	sprite *Sprite
}

func NewBackground(sprite *Sprite) *Background {
	return &Background{BaseComponent: NewComponent(), sprite: sprite}
}

func (sp *Background) Draw() {
	sp.sprite.Render = true
	sp.sprite.DrawScreen()
	sp.sprite.Render = false
}
