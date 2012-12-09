package Server

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
		c.Send(NewEnterGame())
	}
}
