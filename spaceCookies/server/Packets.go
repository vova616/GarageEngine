package server

import (
	"encoding/gob"
)

type PacketID int16

const (
	ID_Welcome    = PacketID(iota)
	ID_EnterGame  = PacketID(iota)
	ID_LoginError = PacketID(iota)

	ID_SpawnPlayer     = PacketID(iota)
	ID_PlayerInfo      = PacketID(iota)
	ID_PlayerTransform = PacketID(iota)
	ID_RemovePlayer    = PacketID(iota)
	ID_PlayerMove      = PacketID(iota)

	ID_Respawn = PacketID(iota)
)

type Packet interface {
	ID() PacketID
}

func init() {
	gob.Register(Welcome{})
	gob.Register(EnterGame{})
	gob.Register(LoginError{})

	gob.Register(SpawnPlayer{})
	gob.Register(PlayerInfo{})
	gob.Register(PlayerTransform{})
	gob.Register(RemovePlayer{})
	gob.Register(PlayerMove{})

	gob.Register(PlayerRespawn{})
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
	PlayerID ID
	Name     string
}

func (e EnterGame) ID() PacketID {
	return ID_EnterGame
}

func NewEnterGame(id ID, name string) Packet {
	return EnterGame{id, name}
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

type SpawnPlayer struct {
	PlayerTransform
	PlayerInfo
}

func (s SpawnPlayer) ID() PacketID {
	return ID_SpawnPlayer
}

func NewSpawnPlayer(playerTransform PlayerTransform, playerInfo PlayerInfo) Packet {
	return SpawnPlayer{playerTransform, playerInfo}
}

type PlayerTransform struct {
	PlayerID ID
	X, Y     float32
	Rotation float32
}

func (s PlayerTransform) ID() PacketID {
	return ID_PlayerTransform
}

func NewPlayerTransform(playerID ID, X, Y, Rotation float32) PlayerTransform {
	return PlayerTransform{playerID, X, Y, Rotation}
}

type PlayerInfo struct {
	PlayerID ID
	Name     string
}

func (s PlayerInfo) ID() PacketID {
	return ID_PlayerInfo
}

func NewPlayerInfo(playerID ID, name string) PlayerInfo {
	return PlayerInfo{playerID, name}
}

type RemovePlayer struct {
	PlayerID ID
}

func (s RemovePlayer) ID() PacketID {
	return ID_RemovePlayer
}

func NewRemovePlayer(playerID ID) Packet {
	return RemovePlayer{playerID}
}

type PlayerMove struct {
	PlayerTransform
}

func (s PlayerMove) ID() PacketID {
	return ID_PlayerMove
}

func NewPlayerMove(transfrom PlayerTransform) Packet {
	return PlayerMove{transfrom}
}

type PlayerRespawn struct {
}

func (s PlayerRespawn) ID() PacketID {
	return ID_Respawn
}

func NewPlayerRespawn() Packet {
	return PlayerRespawn{}
}
