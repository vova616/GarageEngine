package SpaceCookies

import (
	. "github.com/vova616/GarageEngine/Engine"
	//"reflect"
)

type Destoyable struct {
	BaseComponent
	Alive  bool
	HP     float32
	FullHP float32
}

func NewDestoyable(hp float32) *Destoyable {
	return &Destoyable{BaseComponent: NewComponent(), FullHP: hp, Alive: true, HP: hp}
}

func (ds *Destoyable) OnCollisionEnter(arbiter *Arbiter) bool {
	if !ds.Alive {
		return true
	}
	var missle *Missle = nil
	if arbiter.GameObjectA().Tag == MissleTag {
		missle = arbiter.GameObjectA().ComponentTypeOfi(missle).(*Missle)
	} else if arbiter.GameObjectB().Tag == MissleTag {
		missle = arbiter.GameObjectB().ComponentTypeOfi(missle).(*Missle)
	}

	if missle != nil {
		ds.HP -= missle.Damage
	}
	if ds.HP <= 0 {
		ds.Alive = false
		ds.GameObject().Destroy()
	}

	return true
}
