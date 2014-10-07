package wav

import (
	"encoding/binary"
	"fmt"
	"github.com/LaPingvino/GarageEngine/engine/audio"
	"io"
	"os"
)

var test = audio.AudioClip(&Wav{})

type Wav struct {
	Format
	Data        []int16
	index       int
	audioFormat audio.Format
}

func (this *Wav) BufferLength() int {
	return int(this.Format.SampleRate) / 2
}

func (this *Wav) NextBuffer(buff []int16, mono bool) int {
	if mono && this.audioFormat != audio.Mono16 {
		panic("does not support stereo to mono conversion")
	}

	n := copy(buff, this.Data[this.index:])
	this.index += n
	if this.index >= len(this.Data) {
		this.index = 0
	}
	return n
}

func (this *Wav) Clone() (audio.AudioClip, error) {
	return &Wav{this.Format, this.Data, 0, this.audioFormat}, nil
}

func (this *Wav) SetPosition(pos int) {
	this.index = pos
}

func (this *Wav) Length() int {
	return len(this.Data)
}

func (this *Wav) SampleRate() int {
	return int(this.Format.SampleRate)
}
func (this *Wav) AudioFormat() audio.Format {
	return this.audioFormat
}

type Format struct {
	FormatTag     int16
	Channels      int16
	SampleRate    int32
	AvgBytes      int32
	BlockAlign    int16
	BitsPerSample int16
}

type Format2 struct {
	Format
	SizeOfExtension int16
}

type Format3 struct {
	Format2
	ValidBitsPerSample int16
	ChannelMask        int32
	SubFormat          [16]byte
}

func NewClip(path string) (*Wav, error) {
	mr, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	clip, e := ReadWavHeader(mr)
	if e != nil {
		return nil, e
	}
	return clip, nil
}

func ReadWavHeader(reader io.Reader) (*Wav, error) {
	var buff [4]byte
	reader.Read(buff[:4])

	if string(buff[:4]) != "RIFF" {
		return nil, fmt.Errorf("Not a WAV file.\n")
	}

	var size int32
	binary.Read(reader, binary.LittleEndian, &size)

	reader.Read(buff[:4])

	if string(buff[:4]) != "WAVE" {
		return nil, fmt.Errorf("Not a WAV file.\n")
	}

	reader.Read(buff[:4])

	if string(buff[:4]) != "fmt " {
		return nil, fmt.Errorf("Not a WAV file.\n")
	}

	binary.Read(reader, binary.LittleEndian, &size)

	var format Format

	switch size {
	case 16:
		binary.Read(reader, binary.LittleEndian, &format)
	case 18:
		var f2 Format2
		binary.Read(reader, binary.LittleEndian, &f2)
		format = f2.Format
	case 40:
		var f3 Format3
		binary.Read(reader, binary.LittleEndian, &f3)
		format = f3.Format
	}

	reader.Read(buff[:4])

	if string(buff[:4]) != "data" {
		return nil, fmt.Errorf("Not supported WAV file.\n")
	}

	binary.Read(reader, binary.LittleEndian, &size)

	wavData := make([]byte, size)
	n, e := reader.Read(wavData)
	if e != nil {
		return nil, fmt.Errorf("Cannot read WAV data.\n")
	}
	if int32(n) != size {
		return nil, fmt.Errorf("WAV data size doesnt match.\n")
	}

	wavInts := make([]int16, size/2)
	for i := 0; i < len(wavInts); i++ {
		wavInts[i] = int16(binary.LittleEndian.Uint16(wavData[i*2:]))
	}

	af := audio.Mono8
	if format.Channels == 1 {
		af = audio.Mono16
	} else {
		af = audio.Stereo16
	}

	return &Wav{format, wavInts, 0, af}, nil
}
