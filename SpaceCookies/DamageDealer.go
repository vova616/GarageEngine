package SpaceCookies

import (
	. "github.com/vova616/GarageEngine/Engine"
	//"reflect"
)

type DamageDealer struct {
	BaseComponent
	Damage float32
}

func NewDamageDealer(dmg float32) *DamageDealer {
	return &DamageDealer{BaseComponent: NewComponent(), Damage: dmg}
}
