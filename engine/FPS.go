package engine

type FPS struct {
	BaseComponent
	FPS, updateInterval, accum, frames, timeleft float64
	do                                           func(float64)
}

func NewFPS() *FPS {
	return &FPS{NewComponent(), 0, 1, 0, 0, 0.5, nil}
}

func (sp *FPS) SetAction(action func(float64)) {
	sp.do = action
}

func (sp *FPS) Update() {
	sp.timeleft -= DeltaTime()
	sp.accum += DeltaTime()
	sp.frames++

	// Interval ended - update GUI text and start new interval
	if sp.timeleft <= 0.0 {
		sp.FPS = (sp.frames / sp.accum)
		InternalFPS = sp.FPS
		sp.timeleft += sp.updateInterval
		sp.accum = 0.0
		sp.frames = 0
		if sp.do != nil {
			sp.do(sp.FPS)
		}
	}
}
