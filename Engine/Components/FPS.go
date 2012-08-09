package Components

import (
	. "GarageEngine/Engine"
	//. "Engine/Components"
	//"fmt"
	
) 

type FPS struct {
	BaseComponent
	FPS, updateInterval, accum,frames,timeleft float32
	do func(float32)
}

func NewFPS() *FPS {
	return &FPS{NewComponent(),0,0.5,0,0,0.5, nil}
}

func (sp *FPS) SetAction(action func(float32)) {
	sp.do = action
} 

func (sp *FPS) Update() {
    sp.timeleft -= DeltaTime()
    sp.accum += DeltaTime()
    sp.frames++
    
    // Interval ended - update GUI text and start new interval
    if sp.timeleft <= 0.0  {
        sp.FPS = (sp.frames/sp.accum)
        sp.timeleft = sp.updateInterval
        sp.accum = 0.0
        sp.frames = 0
        if sp.do != nil {
        	sp.do(sp.FPS)
        }
    }
}
