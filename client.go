package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

		// Validate message content
		if len(chatMessage.Content) == 0 {
			log.Printf("Empty message content from %s, ignoring", c.username)
			continue
		}

		if len(chatMessage.Content) > 500 {
			log.Printf("Message too long from %s (%d chars), truncating", c.username, len(chatMessage.Content))
			chatMessage.Content = chatMessage.Content[:500]
		}

		// Generate unique message ID and set server timestamp
		messageID := primitive.NewObjectID()
		serverTime := time.Now().UTC()

		chatMessage.ID = messageID.Hex()
		chatMessage.Username = c.username
		chatMessage.Timestamp = serverTime
		chatMessage.Type = "message"

		log.Printf("Message received from %s (ID: %s): %s", c.username, chatMessage.ID, chatMessage.Content)

		// Save to database with consistent timestamp and type
		if c.hub.db != nil {
			err = c.hub.db.SaveMessage(chatMessage.Username, chatMessage.Content, chatMessage.Timestamp, chatMessage.Type)
			if err != nil {
				log.Printf("Error saving message to database from %s: %v", c.username, err)
				// Continue even if DB save fails - don't block real-time chat
			} else {
				log.Printf("Message saved to database from %s (ID: %s)", c.username, chatMessage.ID)
			}
		}

		messageBytes, err = json.Marshal(chatMessage)
		if err != nil {
			log.Printf("Error marshaling message from %s: %v", c.username, err)
			continue
		}

		// Broadcast to all connected clients
		select {
		case c.hub.broadcast <- messageBytes:
			log.Printf("Message broadcasted from %s (ID: %s)", c.username, chatMessage.ID)
		default:
			log.Printf("Warning: Broadcast channel full, message from %s may be dropped", c.username)
		}
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
