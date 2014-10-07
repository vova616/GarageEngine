package audio

import (
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/go-openal/openal"
)

type AudioSource struct {
	engine.BaseComponent
	source        openal.Source
	isMono        bool
	distanceModel DistanceModel

	Clip        AudioClip
	buffers     []openal.Buffer
	audioBuffer []int16

	play     bool
	loop     bool
	position int
}

func NewAudioSource(clip AudioClip) *AudioSource {
	as := &AudioSource{engine.NewComponent(), openal.NewSource(), false, InverseDistanceClamped, nil, nil, nil, true, false, 0}
	if clip != nil {
		c, e := clip.Clone()
		if e != nil {
			panic(e)
		}
		as.Clip = c
	}
	as.source.SetMinGain(0)
	as.source.SetMaxGain(1)
	if clip.AudioFormat() == Mono16 || clip.AudioFormat() == Mono8 {
		as.Set2D()
	}
	return as
}

func (this *AudioSource) SetDistanceModel(model DistanceModel) {
	this.distanceModel = model
}

func (this *AudioSource) Start() {
	this.buffers = openal.NewBuffers(4)
	this.updateBuffers()
}

func (this *AudioSource) Pause() {
	this.play = false
}

func (this *AudioSource) Play() {
	if !this.play {
		this.position = 0
	}
	this.play = true
}

func (this *AudioSource) Stop() {
	this.position = 0
	this.play = false
	this.source.Stop()
	gBuffs := int(this.source.BuffersProcessed())
	for i := 0; i < gBuffs; i++ {
		this.source.UnqueueBuffer()
	}
	this.Clip.SetPosition(0)
}

func (this *AudioSource) IsPlaying() bool {
	return this.play
}

func (this *AudioSource) updateBuffers() {
	if !this.play {
		return
	}
	if this.Clip != nil {
		gBuffs := int(this.source.BuffersProcessed())

		for i := 0; i < gBuffs; i++ {
			if !this.loop && this.position >= this.Clip.Length() {
				this.play = false
				break
			}
			n := this.Clip.NextBuffer(this.audioBuffer, this.isMono)
			this.position += n
			b := this.source.UnqueueBuffer()
			b.SetDataInt(this.Clip.AudioFormat().AlFormat(), this.audioBuffer[:n], int32(this.Clip.SampleRate()))
			this.source.QueueBuffer(b)
			state := this.source.State()
			if !this.play {
				if state == openal.Playing {
					this.source.Pause()
				}
			} else {
				if state != openal.Playing {
					this.source.Play()
				}
			}
		}

		if this.audioBuffer == nil || len(this.audioBuffer) != this.Clip.BufferLength() {
			this.audioBuffer = make([]int16, this.Clip.BufferLength())
		}

		queued := int(this.source.BuffersQueued())
		if queued == 0 {
			i := 0
			for ; i < len(this.buffers); i++ {
				if !this.loop && this.position >= this.Clip.Length() {
					break
				}
				n := this.Clip.NextBuffer(this.audioBuffer, this.isMono)
				this.position += n
				this.buffers[i].SetDataInt(this.Clip.AudioFormat().AlFormat(), this.audioBuffer[:n], int32(this.Clip.SampleRate()))
			}
			this.source.QueueBuffers(this.buffers[:i])
			if this.play {
				this.source.Play()
			}
		}
	}
}

func (this *AudioSource) Update() {
	if currentDistanceModel != this.distanceModel {
		openal.SetDistanceModel(openal.GetDistanceModel())
	}
	this.UpdateTransform()
	this.updateBuffers()
}

func (this *AudioSource) UpdateTransform() {
	if this.isMono {
		pos := this.Transform().WorldPosition()
		rot := this.Transform().WorldRotation()
		this.source.SetPosition(pos.X, pos.Y, pos.Z)
		this.source.SetDirection(rot.X, rot.Y, rot.Z)
		//this.source.SetDirection(0, 0, -1)
	}
}

func (this *AudioSource) SetPitch(pitch float32) {
	this.source.SetPitch(pitch)
}

func (this *AudioSource) SetGain(gain float32) {
	this.source.SetGain(gain)
}

func (this *AudioSource) SetMaxDistance(distance float32) {
	this.source.SetMaxDistance(distance)
}

func (this *AudioSource) SetReferenceDistance(distance float32) {
	this.source.SetReferenceDistance(distance)
}

func (this *AudioSource) SetRolloffFactor(factor float32) {
	this.source.SetRolloffFactor(factor)
}

func (this *AudioSource) SetLooping(loop bool) {
	this.loop = loop
}

func (this *AudioSource) Set2D() {
	this.isMono = false
	this.source.SetSourceRelative(true)
	this.source.SetPosition(0, 0, 0)
	this.source.SetDirection(0, 0, 0)
	this.source.SetVelocity(0, 0, 0)
}

func (this *AudioSource) SetMono(v bool) {
	this.isMono = v
	this.source.SetSourceRelative(v)
	if v {
		this.source.SetPosition(0, 0, 0)
		this.source.SetDirection(0, 0, 0)
		this.source.SetVelocity(0, 0, 0)
	} else {
		this.UpdateTransform()
	}
}
