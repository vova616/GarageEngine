package audio

import (
	//"errors"
	//"github.com/LaPingvino/GarageEngine/engine"
	"github.com/vova616/go-openal/openal"
)

type Format int
type DistanceModel uint32

const (
	Mono8    = Format(iota)
	Mono16   = Format(iota)
	Stereo8  = Format(iota)
	Stereo16 = Format(iota)
)

const (
	InverseDistance         = DistanceModel(openal.InverseDistance)
	InverseDistanceClamped  = DistanceModel(openal.InverseDistanceClamped)
	LinearDistance          = DistanceModel(openal.LinearDistance)
	LinearDistanceClamped   = DistanceModel(openal.LinearDistanceClamped)
	ExponentDistance        = DistanceModel(openal.ExponentDistance)
	ExponentDistanceClamped = DistanceModel(openal.ExponentDistanceClamped)
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
	NextBuffer([]int16, bool) int
	Clone() (AudioClip, error)
	/*
		Samples in this clip
	*/
	Length() int
	SampleRate() int
	SetPosition(int)
	AudioFormat() Format
}
