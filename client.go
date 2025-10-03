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
				log.Printf("websocket closed due to an error: %v", err)
			}
			break
		}

		var chatMessage ChatMessage
		if err := json.Unmarshal(messageBytes, &chatMessage); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		// Set the username and timestamp
		chatMessage.Username = c.username
		chatMessage.Timestamp = time.Now()

		// Save message to database
		if c.hub.db != nil {
			err = c.hub.db.SaveMessage(chatMessage.Username, chatMessage.Content)
			if err != nil {
				log.Printf("error saving message to database: %v", err)
			}
		}

		// Broadcast the message
		messageBytes, err = json.Marshal(chatMessage)
		if err != nil {
			log.Printf("error marshaling message: %v", err)
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
