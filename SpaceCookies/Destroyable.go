package SpaceCookies

import (
	. "github.com/vova616/GarageEngine/Engine"
	//"reflect"
)

type Destoyable struct {
	BaseComponent
	Alive           bool
	HP              float32
	FullHP          float32
	Team            int
	destoyableFuncs DestoyableFuncs
}

func NewDestoyable(hp float32, team int) *Destoyable {
	return &Destoyable{BaseComponent: NewComponent(), FullHP: hp, Alive: true, HP: hp, Team: team}
}

type DestoyableFuncs interface {
	OnDie()
	OnHit(*GameObject, *DamageDealer)
}

func (ds *Destoyable) Start() {
	ds.destoyableFuncs, _ = ds.GameObject().ComponentImplements(&ds.destoyableFuncs).(DestoyableFuncs)
}

func (ds *Destoyable) OnCollisionEnter(arbiter *Arbiter) bool {
	if !ds.Alive {
		return true
	}
	var dmg *DamageDealer = nil
	var enemy *GameObject
	var enemyDestoyable *Destoyable

	if arbiter.GameObjectA() == ds.GameObject() {
		enemy = arbiter.GameObjectB()
	} else {
		enemy = arbiter.GameObjectA()
	}

	dmg, _ = enemy.ComponentTypeOfi(dmg).(*DamageDealer)
	enemyDestoyable, _ = enemy.ComponentTypeOfi(enemyDestoyable).(*Destoyable)

	if enemyDestoyable == nil || enemyDestoyable.Team == ds.Team {
		return true
	}

	if dmg != nil {
		ds.HP -= dmg.Damage
	}
	if ds.destoyableFuncs != nil {
		ds.destoyableFuncs.OnHit(enemy, dmg)
	}

	if ds.HP <= 0 {
		ds.Alive = false
		if ds.destoyableFuncs != nil {
			ds.destoyableFuncs.OnDie()
		} else {
			ds.GameObject().Destroy()
		}
	}

	return true
}
