package zumbies

import (
	"github.com/vova616/garageEngine/engine"
	"github.com/vova616/garageEngine/engine/input"
	//"log"
)

type PlayerController struct {
	engine.BaseComponent
	Player       *Player
	WalkSpeed    float64
	LastPosition engine.Vector
}

func NewPlayerController(player *Player) *PlayerController {
	return &PlayerController{engine.NewComponent(), player, 100, engine.Zero}
}

func (this *PlayerController) Start() {
	this.LastPosition = this.GameObject().Transform().WorldPosition()
}

func (this *PlayerController) FixedUpdate() {
	//t := this.GameObject().Transform()
	var speed engine.Vector

	if input.KeyDown('A') {
		speed.X -= float32(this.WalkSpeed)
	}
	if input.KeyDown('D') {
		speed.X += float32(this.WalkSpeed)
	}
	if input.KeyDown('S') {
		speed.Y -= float32(this.WalkSpeed)
	}
	if input.KeyDown('W') {
		speed.Y += float32(this.WalkSpeed)
	}

	this.GameObject().Physics.Body.SetVelocity(speed.X, speed.Y)
}

func (this *PlayerController) Update() {
	pos := this.Transform().WorldPosition()
	cmap := this.Player.Map
	Player := this.Player
	myPos := pos

	pos.Y -= 32
	newTile, x, y := cmap.PositionToTile(pos)
	_ = newTile
	if !cmap.IsTileWalkabke(x, y) {
		p := this.LastPosition
		p.Y -= 32
		oldTile, x2, y2 := cmap.PositionToTile(p)

		if oldTile.LayerConnected() {
			//println("Layer Move", Player.Map.Layer)
			movedLayer := false
			if Player.Map.Layer+1 < len(Layers) {
				newLayer := Layers[Player.Map.Layer+1]
				//println("Trying to go up")
				if newLayer.IsTileWalkabke(x, y) {
					movedLayer = true
					Player.Map = newLayer
					return
				}
			}
			if Player.Map.Layer-1 >= 0 && !movedLayer {
				newLayer := Layers[Player.Map.Layer-1]
				//println("Trying to go down")
				if newLayer.IsTileWalkabke(x, y) {
					Player.Map = newLayer
					return
				}
			}
		}

		p.Y += 32
		//println(x, y, x2, y2)

		dx := x - x2
		dy := y - y2
		xw := cmap.IsTileWalkabke(x2+dx, y2)
		yw := cmap.IsTileWalkabke(x2, y2+dy)
		//println(x2, y2, dx, dy, xw, yw)
		//xyw := Map1.IsTileWalkabke(x2+dx, y2+dy)
		if xw {
			p.X = myPos.X
		}
		if yw {
			p.Y = myPos.Y
		}

		this.GameObject().Transform().SetWorldPosition(p)
	}
	this.LastPosition = this.GameObject().Transform().WorldPosition()
}
