package Components

import (
	"github.com/vova616/GarageEngine/Engine"
	"image"
)

type Collider struct {
	Engine.BaseComponent
	Rect *image.Rectangle
}

func NewCollider() *Collider {
	return &Collider{Engine.NewComponent(), nil}
}
