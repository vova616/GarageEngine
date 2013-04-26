package audio

import (
	//"errors"
	//"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/go-openal/openal"
)

type Format int

const (
	Mono8    = iota
	Mono16   = iota
	Stereo8  = iota
	Stereo16 = iota
)

func (f Format) AlFormat() openal.Format {
	switch f {
	case Mono8:
		return openal.FormatMono8
	case Mono16:
		return openal.FormatMono16
	case Stereo8:
		return openal.FormatStereo8
	case Stereo16:
		return openal.FormatStereo16
	}
	panic("Unkown format")
}

type AudioClip interface {
	/*
		Buffer size in int16 required to pass to NextBuffer
	*/
	BufferLength() int
	NextBuffer([]int16) int
	Clone() (AudioClip, error)
	/*
		Samples in this clip
	*/
	Length() int
	SampleRate() int
	AudioFormat() Format
}
