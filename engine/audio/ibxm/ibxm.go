package ibxm

import (
	"github.com/LaPingvino/GarageEngine/engine/audio"
	"github.com/vova616/ibxmgo"
	"os"
)

var test = audio.AudioClip(&IBXM{})

type IBXM struct {
	ibxm        *ibxmgo.IBXM
	audioBuffer []int32
	format      audio.Format
	length      int
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

	return &IBXM{clip, nil, audio.Stereo16, clip.Length() * 2}, nil
}

/* Returns the length of the buffer required by NextBuffer(). */
func (this *IBXM) BufferLength() int {
	return this.ibxm.AudioBufferLength()
}

func (this *IBXM) SetPosition(pos int) {
	this.ibxm.SetSequencePos(pos)
}

func (this *IBXM) NextBuffer(outputBuf []int16, mono bool) (samples int) {
	if this.audioBuffer == nil {
		this.audioBuffer = make([]int32, this.ibxm.AudioBufferLength())
	}

	n, _ := this.ibxm.GetAudio(this.audioBuffer)
	if mono {
		for i := 0; i < n; i++ {
			x := int((float64(this.audioBuffer[i*2]) + float64(this.audioBuffer[i*2+1])) / 2)
			if x > 32767 {
				x = 32767
			} else if x < -32768 {
				x = -32768
			}
			outputBuf[i] = int16(x)
		}
		//n *= 2
		this.format = audio.Mono16
	} else {
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
		this.format = audio.Stereo16
	}
	return n
}

func (this *IBXM) Clone() (audio.AudioClip, error) {
	ib, e := ibxmgo.NewIBXM(this.ibxm.Module(), this.ibxm.SampleRate())
	return &IBXM{ib, nil, this.format, this.length}, e
}

func (this *IBXM) AudioFormat() audio.Format {
	return this.format
}

func (this *IBXM) Length() int {
	return this.length
}

func (this *IBXM) SampleRate() int {
	return this.ibxm.SampleRate()
}
