package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
}

func (c *Client) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket connection closed unexpectedly for %s: %v", c.username, err)
			} else {
				log.Printf("WebSocket connection closed normally for %s", c.username)
			}
			break
		}

		var chatMessage ChatMessage
		if err := json.Unmarshal(messageBytes, &chatMessage); err != nil {
			log.Printf("Error unmarshaling message from %s: %v", c.username, err)
			continue
		}

		chatMessage.Username = c.username
		chatMessage.Timestamp = time.Now()

		log.Printf("Message received from %s: %s", c.username, chatMessage.Content)

		if c.hub.db != nil {
			err = c.hub.db.SaveMessage(chatMessage.Username, chatMessage.Content)
			if err != nil {
				log.Printf("Error saving message to database from %s: %v", c.username, err)
			} else {
				log.Printf("Message saved to database from %s", c.username)
			}
		}

		messageBytes, err = json.Marshal(chatMessage)
		if err != nil {
			log.Printf("Error marshaling message from %s: %v", c.username, err)
			continue
		}

		c.hub.broadcast <- messageBytes
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
