package Server

import (
	"encoding/gob"
	"log"
	"net"
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

type Vector struct {
	X, Y, Z float32
}

type Client struct {
	Socket   *net.TCPConn
	ID       ID
	Name     string
	Position Vector
	Rotation Vector

	Decoder *gob.Decoder
	Encoder *gob.Encoder
}

func (c *Client) Run() {
	defer c.OnExit()
	c.Decoder = gob.NewDecoder(c.Socket)
	c.Encoder = gob.NewEncoder(c.Socket)

	var packet Packet
	for {
		e := c.Decoder.Decode(&packet)
		if e != nil {
			panic(e)
		}
		switch packet.ID() {
		case ID_Welcome:
			MainServer.Jobs <- func() { OnWelcomePacket(c, packet) }
		}
	}
}

func (c *Client) Send(p Packet) {
	e := c.Encoder.Encode(&p)
	if e != nil {
		log.Println(e)
	}
}

func (c *Client) OnExit() {
	if x := recover(); x != nil {
		log.Println(c.Name, "Disconnected. Reason:", x)
	}
	MainServer.Jobs <- func() {
		delete(MainServer.Clients, c.ID)
		MainServer.IDGen.PutID(c.ID)
	}
}

func StartServer() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:123")
	if err != nil {
		panic(err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
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
