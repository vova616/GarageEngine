package audio

import (
	//"errors"
	"github.com/LaPingvino/GarageEngine/engine"
	"github.com/vova616/go-openal/openal"
)

var (
	Device               *openal.Device
	Context              *openal.Context
	Listener             *openal.Listener
	currentDistanceModel DistanceModel
)

func init() {
	Listener = nil
	Device = openal.OpenDevice("")
	Context = Device.CreateContext()
	Context.Activate()
	currentDistanceModel = DistanceModel(openal.GetDistanceModel())
}

type AudioListener struct {
	engine.BaseComponent
	listener *openal.Listener
}

func NewAudioListener() *AudioListener {
	if Listener == nil {
		Listener = new(openal.Listener)
		Listener.SetOrientation(0, 0, -1, 0, 1, 0)
	}
	return &AudioListener{engine.NewComponent(), Listener}
}

func (this *AudioListener) Update() {
	pos := this.Transform().WorldPosition()
	this.listener.SetPosition(pos.X, pos.Y, pos.Z)
}
