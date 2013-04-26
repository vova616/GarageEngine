package ibxm

import (
	"github.com/vova616/GarageEngine/engine/audio"
	"github.com/vova616/ibxmgo"
	"os"
)

type IBXM struct {
	ibxm        *ibxmgo.IBXM
	audioBuffer []int32
}

func NewClip(path string) (*IBXM, error) {
	mr, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	mod, e := ibxmgo.Decode(mr)
	if e != nil {
		return nil, e
	}

	clip, e := ibxmgo.NewIBXM(mod, 48000)
	if e != nil {
		return nil, e
	}
	return &IBXM{clip, nil}, nil
}

/* Returns the length of the buffer required by NextBuffer(). */
func (this *IBXM) BufferLength() int {
	return (this.ibxm.CalculateTickLen(32, 128000) + 65) * 2
}

func (this *IBXM) NextBuffer(outputBuf []int16) (samples int) {
	if this.audioBuffer == nil {
		this.audioBuffer = make([]int32, this.ibxm.AudioBufferLength())
	}

	n, _ := this.ibxm.GetAudio(this.audioBuffer)
	n *= 2
	for i := 0; i < n; i++ {
		x := this.audioBuffer[i]
		if x > 32767 {
			x = 32767
		} else if x < -32768 {
			x = -32768
		}
		outputBuf[i] = int16(x)
	}
	return n
}

func (this *IBXM) Clone() (audio.AudioClip, error) {
	ib, e := ibxmgo.NewIBXM(this.ibxm.Module(), this.ibxm.SampleRate())
	return &IBXM{ib, nil}, e
}

func (this *IBXM) AudioFormat() audio.Format {
	return audio.Stereo16
}

func (this *IBXM) Length() int {
	return this.ibxm.Length()
}

func (this *IBXM) SampleRate() int {
	return this.ibxm.SampleRate()
}
