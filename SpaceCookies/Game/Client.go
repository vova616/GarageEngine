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

const ServerIP = "game.vovchik.org:123"
const ServerLocalIP = "localhost:123"

type Client struct {
	Socket *net.TCPConn
	Name   string
	Ship   *ShipController

	Encoder *gob.Encoder
	Decoder *gob.Decoder

	Jobs         chan func()
	Disconnected bool
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
	MyClient = &Client{Socket: tcpCon, Name: name, Encoder: gob.NewEncoder(tcpCon), Decoder: gob.NewDecoder(tcpCon), Jobs: make(chan func(), 1000)}
	go MyClient.Run()
	LoginErrChan = *errChan
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
		p := packet
		c.Jobs <- func() { c.HandlePacket(p) }
	}
}

func (c *Client) HandlePacket(packet Server.Packet) {
	defer c.OnPanic()
	switch packet.ID() {
	case Server.ID_EnterGame:
		Engine.LoadScene(GameSceneGeneral)
	case Server.ID_LoginError:
		error := packet.(Server.LoginError)
		LoginErrChan <- fmt.Errorf(error.Error)
		panic(error)
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
