package SpaceCookies

import (
	"github.com/vova616/GarageEngine/Engine"
	//"reflect"
)

type DamageDealer struct {
	Engine.BaseComponent
	Damage float32
}

func NewDamageDealer(dmg float32) *DamageDealer {
	return &DamageDealer{BaseComponent: Engine.NewComponent(), Damage: dmg}
}
