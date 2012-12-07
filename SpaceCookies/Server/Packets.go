package Server

import (
	"encoding/gob"
)

type PacketID int16

const (
	ID_Welcome    = PacketID(iota)
	ID_EnterGame  = PacketID(iota)
	ID_LoginError = PacketID(iota)
)

type Packet interface {
	ID() PacketID
}

func init() {
	gob.Register(Welcome{})
	gob.Register(EnterGame{})
	gob.Register(LoginError{})
}

type Welcome struct {
	Name string
}

func (w Welcome) ID() PacketID {
	return ID_Welcome
}

func NewWelcome(name string) Packet {
	return Welcome{name}
}

type EnterGame struct {
}

func (e EnterGame) ID() PacketID {
	return ID_EnterGame
}

func NewEnterGame() Packet {
	return EnterGame{}
}

type LoginError struct {
	Error string
}

func (e LoginError) ID() PacketID {
	return ID_LoginError
}

func NewLoginError(error string) Packet {
	return LoginError{error}
}
