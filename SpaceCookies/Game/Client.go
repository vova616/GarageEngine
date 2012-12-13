package Game

import (
	//"log"
	"encoding/gob"
	"fmt"
	"github.com/vova616/GarageEngine/Engine"
	"github.com/vova616/GarageEngine/SpaceCookies/Server"
	"log"
	"net"
	"time"
)

var (
	MyClient     *Client = nil
	LoginErrChan chan error
)

const ServerIP = "localhost:123"
const ServerLocalIP = "localhost:123"

type Client struct {
	Engine.BaseComponent
	Socket *net.TCPConn
	Name   string
	ID     Server.ID
	Ship   *ShipController

	Encoder *gob.Encoder
	Decoder *gob.Decoder

	Jobs         chan func()
	GameJobs     chan func()
	Disconnected bool

	lastTransformUpdate        time.Time
	lastX, lastY, lastRotation float32
}

func Connect(name string, errChan *chan error) {
	*errChan = make(chan error)
	/*
		addr, err := net.ResolveTCPAddr("tcp", ServerIP)
		if err != nil {
			*errChan <- err
			return
		}
	*/
	//con, err := net.DialTCP("tcp", nil, addr)
	con, err := net.DialTimeout("tcp", ServerIP, time.Second*4)
	if err != nil {
		con, err = net.DialTimeout("tcp", ServerLocalIP, time.Second*4)
		if err != nil {
			*errChan <- fmt.Errorf("Game Server is down :(.")
			return
		}
	}
	tcpCon := con.(*net.TCPConn)
	MyClient = &Client{BaseComponent: Engine.NewComponent(), Socket: tcpCon, Name: name, Encoder: gob.NewEncoder(tcpCon), Decoder: gob.NewDecoder(tcpCon), Jobs: make(chan func(), 1000), GameJobs: make(chan func(), 1000)}
	go MyClient.Run()
	LoginErrChan = *errChan
}

func (c *Client) Update() {
	b := true
	for b {
		select {
		case job := <-c.GameJobs:
			job()
		default:
			b = false
		}
	}
}

func (c *Client) Send(p Server.Packet) {
	if c.Disconnected {
		return
	}
	c.Encoder.Encode(&p)
}

func (c *Client) LateUpdate() {
	//if time.Since(c.lastTransformUpdate) > time.Second/60 {
	//	c.lastTransformUpdate = time.Now()
	p := c.Transform().WorldPosition()
	r := c.Transform().Angle()
	if c.lastX != p.X || c.lastY != p.Y || c.lastRotation != r {
		c.Jobs <- func() {
			c.Send(Server.NewPlayerMove(Server.NewPlayerTransform(c.ID, p.X, p.Y, r)))
		}
		c.lastX, c.lastY, c.lastRotation = p.X, p.Y, r
	}
	//}
}

func (c *Client) DoJobs() {
	for job := range c.Jobs {
		job()
	}
}

func (c *Client) Run() {
	defer c.OnPanic()
	go c.DoJobs()
	err := MyClient.SendName()
	if err != nil {
		panic(err)
	}

	for {
		var packet Server.Packet
		e := c.Decoder.Decode(&packet)
		if e != nil {
			panic(e)
		}
		c.Jobs <- func() { c.HandlePacket(packet) }
	}
}

func (c *Client) HandlePacket(packet Server.Packet) {
	defer c.OnPanic()
	switch packet.ID() {
	case Server.ID_SpawnPlayer:
		spawnPlayer := packet.(Server.SpawnPlayer)
		c.GameJobs <- func() {
			if spawnPlayer.PlayerInfo.PlayerID == c.ID {
				SpawnMainPlayer(spawnPlayer)
			} else {
				SpawnPlayer(spawnPlayer)
			}
		}
	case Server.ID_EnterGame:
		enterGame := packet.(Server.EnterGame)
		c.ID = enterGame.PlayerID
		c.Name = enterGame.Name
		Engine.LoadScene(GameSceneGeneral)
	case Server.ID_LoginError:
		error := packet.(Server.LoginError)
		LoginErrChan <- fmt.Errorf(error.Error)
		panic(error)
	case Server.ID_PlayerTransform:
		trans := packet.(Server.PlayerTransform)
		c.GameJobs <- func() {
			p, exist := Players[trans.PlayerID]
			if !exist {
				println("player does not exists")
				return
			}
			p.Transform().SetPositionf(trans.X, trans.Y)
			p.Transform().SetRotationf(trans.Rotation)
		}
	}
}

func (c *Client) OnPanic() {
	if x := recover(); x != nil && !c.Disconnected {
		log.Println(c.Name, "Disconnected. Reason:", x)
		c.Disconnected = true
		c.Socket.Close()
		if MyClient == c {
			MyClient = nil
		}
	}
}

func (c *Client) SendName() error {
	p := Server.NewWelcome(c.Name)
	e := c.Encoder.Encode(&p)
	if e != nil {
		return e
	}
	return nil
}
