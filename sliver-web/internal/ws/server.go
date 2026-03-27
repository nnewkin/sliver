package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	ID       int64
	Username string
	Role     string
	Conn     *websocket.Conn
	Send     chan []byte
	Server   *Server
}

type Message struct {
	Type      string      `json:"type"`
	Channel   string      `json:"channel,omitempty"`
	Action    string      `json:"action,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

type Server struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	running    bool
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		running:    true,
	}
}

func (s *Server) Start() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()

			msg := Message{
				Type:      "event",
				Channel:   "system",
				Data:      map[string]interface{}{"event": "client_connected", "username": client.Username},
				Timestamp: time.Now(),
			}
			data, _ := json.Marshal(msg)
			s.broadcastToAll(data)

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.Send)
			}
			s.mu.Unlock()

			msg := Message{
				Type:      "event",
				Channel:   "system",
				Data:      map[string]interface{}{"event": "client_disconnected", "username": client.Username},
				Timestamp: time.Now(),
			}
			data, _ := json.Marshal(msg)
			s.broadcastToAll(data)

		case message := <-s.broadcast:
			s.broadcastToAll(message)
		}
	}
}

func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
	for client := range s.clients {
		close(client.Send)
	}
	s.clients = make(map[*Client]bool)
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID int64, username, role string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade failed: %v\n", err)
		return
	}

	client := &Client{
		ID:       userID,
		Username: username,
		Role:     role,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Server:   s,
	}

	s.register <- client

	go client.writePump()
	go client.readPump()
}

func (s *Server) broadcastToAll(message []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for client := range s.clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(s.clients, client)
		}
	}
}

func (s *Server) SendToUser(userID int64, msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for client := range s.clients {
		if client.ID == userID {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}

func (s *Server) BroadcastEvent(channel, event string, data interface{}) {
	msg := Message{
		Type:      "event",
		Channel:   channel,
		Data:      map[string]interface{}{"event": event, "data": data},
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	s.broadcast <- data
}

func (s *Server) GetClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}

func (c *Client) readPump() {
	defer func() {
		c.Server.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		c.handleMessage(&msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(msg *Message) {
	switch msg.Action {
	case "subscribe":
		if channels, ok := msg.Data.([]interface{}); ok {
			for _, ch := range channels {
				if channel, ok := ch.(string); ok {
					c.Server.BroadcastEvent(channel, "subscribed", map[string]interface{}{"channel": channel})
				}
			}
		}

	case "ping":
		resp := Message{
			Type:      "pong",
			Timestamp: time.Now(),
		}
		data, _ := json.Marshal(resp)
		c.Send <- data
	}
}
