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

	MainServer = &Server{ln, make(map[ID]*Client), make(chan Job, 1000), NewIDGenerator(1000000)}
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
			MainServer.Clients[id] = c
			go c.Run()
		}
	}
}
