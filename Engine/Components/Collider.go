package Components

import (
	"image"
	. "../../Engine"
	)

type Collider struct {
	BaseComponent
	Rect *image.Rectangle
}

func NewCollider() *Collider {
	return &Collider{NewComponent(),nil}
} 