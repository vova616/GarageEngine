package Game

import (
	"github.com/vova616/GarageEngine/Engine"
	//"reflect"
	"time"
)

type Destoyable struct {
	Engine.BaseComponent
	Alive           bool
	HP              float32
	FullHP          float32
	Team            int
	destoyableFuncs DestoyableFuncs

	createTime    time.Time
	aliveDuration time.Duration
	autoDestory   bool
}

func NewDestoyable(hp float32, team int) *Destoyable {
	return &Destoyable{BaseComponent: Engine.NewComponent(), FullHP: hp, Alive: true, HP: hp, Team: team}
}

type DestoyableFuncs interface {
	OnDie(byTimer bool)
	OnHit(*Engine.GameObject, *DamageDealer)
}

func (ds *Destoyable) Start() {
	ds.createTime = time.Now()
	ds.destoyableFuncs, _ = ds.GameObject().ComponentImplements(&ds.destoyableFuncs).(DestoyableFuncs)
}

func (ds *Destoyable) SetDestroyTime(sec float32) {
	ds.autoDestory = true
	ds.aliveDuration = time.Millisecond * time.Duration(1000*sec)
}

func (ds *Destoyable) Update() {
	if ds.autoDestory && ds.GameObject() != nil {
		if time.Now().After(ds.createTime.Add(ds.aliveDuration)) {
			if ds.destoyableFuncs != nil {
				ds.destoyableFuncs.OnDie(true)
			} else {
				ds.GameObject().Destroy()
			}
		}
	}
}

func (ds *Destoyable) OnCollisionEnter(arbiter *Engine.Arbiter) bool {
	if !ds.Alive {
		return true
	}
	var dmg *DamageDealer = nil
	var enemy *Engine.GameObject
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
			ds.destoyableFuncs.OnDie(false)
		} else {
			ds.GameObject().Destroy()
		}
	}

	return true
}
