package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type Message struct {
	UserID    string
	Content   string
	Timestamp time.Time
}

type Client struct {
	ID       string
	Messages chan Message
}

type ChatServer struct {
	mu       sync.Mutex
	clients  map[string]*Client
	history  []Message
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[string]*Client),
		history: make([]Message, 0),
	}
}

// Join registers a new client and notifies others
func (s *ChatServer) Join(userID string, reply *[]Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create new client with message channel
	client := &Client{
		ID:       userID,
		Messages: make(chan Message, 100),
	}
	s.clients[userID] = client

	// Broadcast join notification to all OTHER clients
	joinMsg := Message{
		UserID:    "SYSTEM",
		Content:   fmt.Sprintf("User [%s] joined", userID),
		Timestamp: time.Now(),
	}
	s.history = append(s.history, joinMsg)

	for id, c := range s.clients {
		if id != userID {
			select {
			case c.Messages <- joinMsg:
			default:
				// Channel full, skip
			}
		}
	}

	// Return full history to the joining client
	*reply = make([]Message, len(s.history))
	copy(*reply, s.history)

	log.Printf("User [%s] joined. Total clients: %d", userID, len(s.clients))
	return nil
}


// SendMessage broadcasts a message to all other clients (no self-echo)
func (s *ChatServer) SendMessage(args *Message, reply *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	args.Timestamp = time.Now()
	s.history = append(s.history, *args)

	// Broadcast to all OTHER clients (no self-echo)
	for id, c := range s.clients {
		if id != args.UserID {
			select {
			case c.Messages <- *args:
			default:
				// Channel full, skip
			}
		}
	}

	*reply = true
	log.Printf("Message from [%s]: %s", args.UserID, args.Content)
	return nil
}

// GetMessages retrieves pending messages for a client (long-polling style)
func (s *ChatServer) GetMessages(userID string, reply *[]Message) error {
	s.mu.Lock()
	client, exists := s.clients[userID]
	s.mu.Unlock()

	if !exists {
		return fmt.Errorf("client %s not found", userID)
	}

	messages := make([]Message, 0)
	
	// Non-blocking read of all available messages
	for {
		select {
		case msg := <-client.Messages:
			messages = append(messages, msg)
		default:
			*reply = messages
			return nil
		}
	}
}

// Leave removes a client and notifies others
func (s *ChatServer) Leave(userID string, reply *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[userID]; !exists {
		*reply = false
		return nil
	}

	delete(s.clients, userID)

	// Broadcast leave notification
	leaveMsg := Message{
		UserID:    "SYSTEM",
		Content:   fmt.Sprintf("User [%s] left", userID),
		Timestamp: time.Now(),
	}
	s.history = append(s.history, leaveMsg)

	for _, c := range s.clients {
		select {
		case c.Messages <- leaveMsg:
		default:
		}
	}

	*reply = true
	log.Printf("User [%s] left. Total clients: %d", userID, len(s.clients))
	return nil
}

// GetHistory returns full chat history
func (s *ChatServer) GetHistory(userID string, reply *[]Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	*reply = make([]Message, len(s.history))
	copy(*reply, s.history)
	return nil
}

func main() {
	server := NewChatServer()
	rpc.Register(server)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	defer listener.Close()

	log.Println("Chat server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
