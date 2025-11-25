package main

import (
	"chatroom/shared"
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

// ChatServer handles RPC calls for the chatroom
type ChatServer struct {
	messages []shared.Message
	mu       sync.RWMutex
}

// SendMessage adds a message to the chat history and returns all messages
func (cs *ChatServer) SendMessage(args *shared.SendMessageArgs, reply *shared.SendMessageReply) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Add the new message to the history
	cs.messages = append(cs.messages, shared.Message{Sender: args.Sender, Content: args.Content})

	// Return the complete history
	reply.History = make([]shared.Message, len(cs.messages))
	copy(reply.History, cs.messages)
	reply.Success = true

	fmt.Printf("Received message from %s: %s\n", args.Sender, args.Content)
	return nil
}

// GetHistory returns all messages in the chat history
func (cs *ChatServer) GetHistory(args *shared.GetHistoryArgs, reply *shared.GetHistoryReply) error {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	reply.History = make([]shared.Message, len(cs.messages))
	copy(reply.History, cs.messages)

	return nil
}

func main() {
	// Create and register the chat server
	chatServer := new(ChatServer)
	rpc.Register(chatServer)

	// Listen on TCP port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		fmt.Println("Port 8080 is already in use.")
		fmt.Println("Please stop the existing server or choose a different port.")
		return
	}
	defer listener.Close()

	fmt.Println("Chatroom server started on port 8080")
	fmt.Println("Waiting for clients...")

	// Accept connections and serve RPC requests
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
