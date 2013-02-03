package server

import (
	"encoding/gob"
	"log"
	"net"
	"sync/atomic"
)

var MainServer *Server

type Job func()
type Server struct {
	Socket  *net.TCPListener
	Clients map[ID]*Client
	Jobs    chan Job
	IDGen   *IDGenerator
}

func (s *Server) Run() {
	for job := range s.Jobs {
		job()
	}
}

type Client struct {
	Socket   *net.TCPConn
	ID       ID
	Name     string
	X, Y     float32
	Rotation float32

	Decoder *gob.Decoder
	Encoder *gob.Encoder

	Disconnected int32
}

func (c *Client) Run() {
	defer c.OnPanic()
	c.Decoder = gob.NewDecoder(c.Socket)
	c.Encoder = gob.NewEncoder(c.Socket)

	for {
		var packet Packet
		e := c.Decoder.Decode(&packet)
		if e != nil {
			panic(e)
		}
		MainServer.Jobs <- func() { c.HandlePacket(packet) }
	}
}

func (c *Client) HandlePacket(p Packet) {
	defer c.OnPanic()
	switch p.ID() {
	case ID_Welcome:
		OnWelcomePacket(c, p)
	case ID_PlayerMove:
		OnPlayerMove(c, p)
	case ID_Respawn:
		OnPlayerRespawn(c, p)
	}
}

func (c *Client) Send(p Packet) {
	if atomic.LoadInt32(&c.Disconnected) == 0 {
		e := c.Encoder.Encode(&p)
		if e != nil {
			panic(e)
		}
	}
}

func (c *Client) OnPanic() {
	if x := recover(); x != nil {
		if atomic.CompareAndSwapInt32(&c.Disconnected, 0, 1) {
			log.Println(c.Name, "Disconnected. Reason:", x)
			MainServer.Jobs <- func() {
				delete(MainServer.Clients, c.ID)
				MainServer.IDGen.PutID(c.ID)
			}
		}
	}
}

func StartServer() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:123")
	if err != nil {
		log.Println(err)
		return
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Server started!")
	//MainServer.IDGen can be not safe because the only place we use it is when we adding/removing clients from the list and we need to do it safe anyway
	MainServer = &Server{ln, make(map[ID]*Client), make(chan Job, 1000), NewIDGenerator(100000, false)}
	go MainServer.Run()

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Println(err)
			break
		}
		MainServer.Jobs <- func() {
			id := MainServer.IDGen.NextID()
			c := &Client{
				Socket: conn, ID: id,
			}
			MainServer.Clients[c.ID] = c
			go c.Run()
		}
	}
}
