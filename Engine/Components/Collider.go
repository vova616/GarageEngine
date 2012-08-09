package Components

import (
	"image"
	. "github.com/vova616/GarageEngine/Engine"
	)

type Collider struct {
	BaseComponent
	Rect *image.Rectangle
}

func NewCollider() *Collider {
	return &Collider{NewComponent(),nil}
} 