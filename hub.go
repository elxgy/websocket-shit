package main

import (
	"encoding/json"
	"log"
)

const MaxClients = 4

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	db         *Database
}

func newHub(db *Database) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		db:         db,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if len(h.clients) >= MaxClients {
				log.Printf("Maximum number of clients reached, rejecting new connection")
				close(client.send)
				client.conn.Close()
				continue
			}

			h.clients[client] = true
			log.Printf("Client %s connected. Total clients: %d", client.username, len(h.clients))

			if h.db != nil {
				messages, err := h.db.GetRecentMessages()
				if err != nil {
					log.Printf("Error getting recent messages: %v", err)
				} else {
					for _, msg := range messages {
						chatMessage := ChatMessage{
							Type:      "message",
							Username:  msg.Username,
							Content:   msg.Content,
							Timestamp: msg.Timestamp,
						}
						messageBytes, err := json.Marshal(chatMessage)
						if err != nil {
							log.Printf("Error marshaling message: %v", err)
							continue
						}
						select {
						case client.send <- messageBytes:
						default:
							close(client.send)
							delete(h.clients, client)
							return
						}
					}
				}
			}

			joinMessage := ChatMessage{
				Type:     "user_joined",
				Username: client.username,
				Content:  client.username + " joined the chat",
			}
			messageBytes, _ := json.Marshal(joinMessage)
			h.broadcastToOthers(client, messageBytes)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client %s disconnected. Total clients: %d", client.username, len(h.clients))

				leaveMessage := ChatMessage{
					Type:     "user_left",
					Username: client.username,
					Content:  client.username + " left the chat",
				}
				messageBytes, _ := json.Marshal(leaveMessage)
				h.broadcastToAll(messageBytes)
			}

		case message := <-h.broadcast:
			h.broadcastToAll(message)
		}
	}
}

func (h *Hub) broadcastToAll(message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) broadcastToOthers(sender *Client, message []byte) {
	for client := range h.clients {
		if client != sender {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}
