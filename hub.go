package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const MaxClients = 4

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	db         *Database
	mutex      sync.RWMutex
}

func newHub(db *Database) *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		db:         db,
		mutex:      sync.RWMutex{},
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			if len(h.clients) >= MaxClients {
				log.Printf("Maximum number of clients reached, rejecting new connection from %s", client.username)
				h.mutex.Unlock()
				close(client.send)
				client.conn.Close()
				continue
			}

			h.clients[client] = true
			clientCount := len(h.clients)
			h.mutex.Unlock()
			log.Printf("Client %s connected. Total clients: %d", client.username, clientCount)

			if h.db != nil {
				messages, err := h.db.GetRecentMessages()
				if err != nil {
					log.Printf("Error getting recent messages for %s: %v", client.username, err)
				} else {
					log.Printf("Sending %d historical messages to %s", len(messages), client.username)
					for _, msg := range messages {
						chatMessage := ChatMessage{
							ID:        primitive.NewObjectID().Hex(),
							Type:      "message",
							Username:  msg.Username,
							Content:   msg.Content,
							Timestamp: msg.Timestamp,
						}
						messageBytes, err := json.Marshal(chatMessage)
						if err != nil {
							log.Printf("Error marshaling historical message for %s: %v", client.username, err)
							continue
						}
						select {
						case client.send <- messageBytes:
						default:
							log.Printf("Warning: Client %s send channel full, dropping historical message", client.username)
							h.mutex.Lock()
							close(client.send)
							delete(h.clients, client)
							h.mutex.Unlock()
							return
						}
					}
				}
			}

			joinMessage := ChatMessage{
				ID:        primitive.NewObjectID().Hex(),
				Type:      "user_joined",
				Username:  client.username,
				Content:   client.username + " joined the chat",
				Timestamp: time.Now().UTC(),
			}
			messageBytes, err := json.Marshal(joinMessage)
			if err != nil {
				log.Printf("Error marshaling join message for %s: %v", client.username, err)
			} else {
				h.broadcastToOthers(client, messageBytes)
				log.Printf("Broadcasted join message for %s", client.username)
			}

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				clientCount := len(h.clients)
				h.mutex.Unlock()
				log.Printf("Client %s disconnected. Total clients: %d", client.username, clientCount)

				leaveMessage := ChatMessage{
					ID:        primitive.NewObjectID().Hex(),
					Type:      "user_left",
					Username:  client.username,
					Content:   client.username + " left the chat",
					Timestamp: time.Now().UTC(),
				}
				messageBytes, err := json.Marshal(leaveMessage)
				if err != nil {
					log.Printf("Error marshaling leave message for %s: %v", client.username, err)
				} else {
					h.broadcastToAll(messageBytes)
					log.Printf("Broadcasted leave message for %s", client.username)
				}
			} else {
				h.mutex.Unlock()
			}

		case message := <-h.broadcast:
			h.broadcastToAll(message)
		}
	}
}

func (h *Hub) broadcastToAll(message []byte) {
	h.mutex.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mutex.RUnlock()

	for _, client := range clients {
		select {
		case client.send <- message:
		default:
			log.Printf("Warning: Client %s send channel full, removing client", client.username)
			h.mutex.Lock()
			close(client.send)
			delete(h.clients, client)
			h.mutex.Unlock()
		}
	}
}

func (h *Hub) broadcastToOthers(sender *Client, message []byte) {
	h.mutex.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		if client != sender {
			clients = append(clients, client)
		}
	}
	h.mutex.RUnlock()

	for _, client := range clients {
		select {
		case client.send <- message:
		default:
			log.Printf("Warning: Client %s send channel full, removing client", client.username)
			h.mutex.Lock()
			close(client.send)
			delete(h.clients, client)
			h.mutex.Unlock()
		}
	}
}
