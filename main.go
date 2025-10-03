package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	hub *Hub
	db  *Database
}

func NewServer(db *Database) *Server {
	hub := newHub(db)
	return &Server{
		hub: hub,
		db:  db,
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("Login attempt from %s - User-Agent: %s", r.RemoteAddr, r.UserAgent())

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		log.Printf("CORS preflight request handled for %s", r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		log.Printf("Invalid method %s for login from %s", r.Method, r.RemoteAddr)
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		log.Printf("Invalid JSON in login request from %s: %v", r.RemoteAddr, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Invalid JSON",
		})
		return
	}

	log.Printf("Login attempt for username: %s from %s", loginReq.Username, r.RemoteAddr)

	user, err := s.db.AuthenticateUser(loginReq.Username, loginReq.Password)
	if err != nil {
		log.Printf("Authentication failed for %s from %s: %v", loginReq.Username, r.RemoteAddr, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	log.Printf("User %s successfully authenticated from %s", user.Username, r.RemoteAddr)
	json.NewEncoder(w).Encode(LoginResponse{
		Success:  true,
		Username: user.Username,
		Message:  "Login successful",
	})
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		log.Printf("WebSocket connection attempt without username from %s", r.RemoteAddr)
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	log.Printf("WebSocket connection attempt for user %s from %s", username, r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed for %s from %s: %v", username, r.RemoteAddr, err)
		return
	}

	log.Printf("WebSocket connection established for user %s from %s", username, r.RemoteAddr)

	client := &Client{
		hub:      s.hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
	}

	client.hub.register <- client

	go client.write()
	go client.read()
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := map[string]any{
		"status":             "ok",
		"clients":            len(s.hub.clients),
		"max_clients":        MaxClients,
		"database_connected": s.db != nil,
	}

	json.NewEncoder(w).Encode(response)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Info: No .env file found (this is normal for Railway deployment)")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is required")
	}

	log.Printf("Connecting to MongoDB...")

	db, err := NewDatabase(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to MongoDB successfully")

	if err := db.CreateDefaultUsers(); err != nil {
		log.Printf("Warning: Error creating default users: %v", err)
	} else {
		log.Println("Default users ready")
	}

	server := NewServer(db)

	go server.hub.run()

	r := mux.NewRouter()
	r.Use(corsMiddleware)

	r.HandleFunc("/login", server.handleLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/ws", server.handleWebSocket)
	r.HandleFunc("/health", server.handleHealth).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Chat server starting on port %s", port)
	log.Printf("WebSocket endpoint: /ws")
	log.Printf("Login endpoint: /login")
	log.Printf("Health check: /health")

	if err := http.ListenAndServe("0.0.0.0:"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
