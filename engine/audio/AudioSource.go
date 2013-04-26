package audio

import (
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/go-openal/openal"
)

type AudioSource struct {
	engine.BaseComponent
	source openal.Source
	is2D   bool

	Clip        AudioClip
	buffers     []openal.Buffer
	audioBuffer []int16
}

func NewAudioSource(clip AudioClip) *AudioSource {
	return &AudioSource{engine.NewComponent(), openal.NewSource(), false, clip, nil, nil}
}

func (this *AudioSource) Start() {
	this.buffers = openal.NewBuffers(4)
	this.updateBuffers()
}

func (this *AudioSource) updateBuffers() {
	if this.Clip != nil {
		if this.audioBuffer == nil || len(this.audioBuffer) != this.Clip.BufferLength() {
			this.audioBuffer = make([]int16, this.Clip.BufferLength())
			for i := 0; i < len(this.buffers); i++ {
				n := this.Clip.NextBuffer(this.audioBuffer)
				this.buffers[i].SetDataInt(this.Clip.AudioFormat().AlFormat(), this.audioBuffer[:n], int32(this.Clip.SampleRate()))
			}
			this.source.QueueBuffers(this.buffers)
			this.source.Play()
			return
		}
		gBuffs := int(this.source.BuffersProcessed())
		for i := 0; i < gBuffs; i++ {
			n := this.Clip.NextBuffer(this.audioBuffer)

			b := this.source.UnqueueBuffer()
			b.SetDataInt(this.Clip.AudioFormat().AlFormat(), this.audioBuffer[:n], int32(this.Clip.SampleRate()))
			this.source.QueueBuffer(b)
			if this.source.State() != openal.Playing {
				this.source.Play()
			}
		}
	}
}

func (this *AudioSource) Update() {
	this.UpdateTransform()
	this.updateBuffers()
}

func (this *AudioSource) UpdateTransform() {
	if !this.is2D {
		pos := this.Transform().WorldPosition()
		rot := this.Transform().WorldRotation()
		this.source.SetPosition(pos.X, pos.Y, pos.Z)
		this.source.SetDirection(rot.X, rot.Y, rot.Z)
	}
}

func (this *AudioSource) SetPitch(pitch float32) {
	this.source.SetPitch(pitch)
}

func (this *AudioSource) SetGain(gain float32) {
	this.source.SetGain(gain)
}

func (this *AudioSource) SetLooping(loop bool) {
	this.source.SetLooping(loop)
}

func (this *AudioSource) Set2D(v bool) {
	this.is2D = v
	this.source.SetSourceRelative(v)
	if v {
		this.source.SetPosition(0, 0, 0)
		this.source.SetDirection(0, 0, 0)
		this.source.SetVelocity(0, 0, 0)
	} else {
		this.UpdateTransform()
	}
}
