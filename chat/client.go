package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	userData map[string]interface{}
	room     *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var m *message
		err := c.socket.ReadJSON(&m)
		if err != nil {
			log.Println("error when read json from browser: ", err.Error())
			return
		}
		m.When = time.Now()
		m.Name = c.userData["name"].(string)
		m.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
		c.room.foward <- m
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			log.Println(err.Error())
			return
		}
	}
}
