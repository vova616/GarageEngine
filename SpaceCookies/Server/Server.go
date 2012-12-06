package Server

import (
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
}

func (c *Client) Run() {
	log.Println("Connection in!", c.Socket.LocalAddr(), c.ID)
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
				conn, id, "", Vector{}, Vector{},
			}
			MainServer.Clients[c.ID] = c
			go c.Run()
		}
	}
}
