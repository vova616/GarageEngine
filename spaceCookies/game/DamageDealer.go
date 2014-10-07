package game

import (
	"github.com/LaPingvino/GarageEngine/engine"
	//"reflect"
)

type DamageDealer struct {
	engine.BaseComponent
	Damage float32
}

func NewDamageDealer(dmg float32) *DamageDealer {
	return &DamageDealer{BaseComponent: engine.NewComponent(), Damage: dmg}
}
