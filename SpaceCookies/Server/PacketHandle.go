package server

import "strings"
import "log"

func OnWelcomePacket(c *Client, p Packet) {
	welcomePacket := p.(Welcome)

	if strings.ToLower(welcomePacket.Name) == "admin" {
		c.Send(NewLoginError("YOU NO ADMIN!#."))
		log.Println("YOU NO ADMIN!#", c.Socket.RemoteAddr(), c.ID)
		return
	} else if len(welcomePacket.Name) == 0 {
		c.Send(NewLoginError("Empty name, try again."))
		log.Println("Empty name", c.Socket.RemoteAddr(), c.ID)
		return
	}

	nameExists := false
	for _, cc := range MainServer.Clients {
		if cc.Name == welcomePacket.Name {
			nameExists = true
			break
		}
	}

	if nameExists {
		c.Send(NewLoginError("YOU SHALL NOT PASS (this name is already taken)."))
		log.Println("YOU SHALL NOT PASS", c.Socket.RemoteAddr(), c.ID)
	} else {
		c.Name = welcomePacket.Name
		log.Println("Connection in!", c.Socket.RemoteAddr(), c.ID, c.Name)
		PlayerEnterGame(c)
	}
}

func OnPlayerRespawn(c *Client, p Packet) {
	PlayerEnterGame(c)
}

func OnPlayerMove(c *Client, p Packet) {
	movePlayer := p.(PlayerMove)

	c.X, c.Y, c.Rotation = movePlayer.X, movePlayer.Y, movePlayer.Rotation

	for _, client := range MainServer.Clients {
		if c != client {
			client.Send(NewPlayerTransform(c.ID, c.X, c.Y, c.Rotation))
		}
	}
}

func PlayerEnterGame(c *Client) {
	c.X, c.Y = 400, 200
	c.Send(NewEnterGame(c.ID, c.Name))
	for id, client := range MainServer.Clients {
		c.Send(NewSpawnPlayer(NewPlayerTransform(id, client.X, client.Y, client.Rotation), NewPlayerInfo(id, client.Name)))
		if c != client {
			client.Send(NewSpawnPlayer(NewPlayerTransform(c.ID, c.X, c.Y, c.Rotation), NewPlayerInfo(c.ID, c.Name)))
		}
	}
}
