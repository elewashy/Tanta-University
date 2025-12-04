package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"time"
)

type Message struct {
	UserID    string
	Content   string
	Timestamp time.Time
}

type ChatClient struct {
	client *rpc.Client
	userID string
	done   chan struct{}
}

func NewChatClient(serverAddr, userID string) (*ChatClient, error) {
	client, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}

	return &ChatClient{
		client: client,
		userID: userID,
		done:   make(chan struct{}),
	}, nil
}

func (c *ChatClient) Join() ([]Message, error) {
	var history []Message
	err := c.client.Call("ChatServer.Join", c.userID, &history)
	return history, err
}

func (c *ChatClient) SendMessage(content string) error {
	msg := &Message{
		UserID:  c.userID,
		Content: content,
	}
	var reply bool
	return c.client.Call("ChatServer.SendMessage", msg, &reply)
}

func (c *ChatClient) GetMessages() ([]Message, error) {
	var messages []Message
	err := c.client.Call("ChatServer.GetMessages", c.userID, &messages)
	return messages, err
}

func (c *ChatClient) Leave() error {
	var reply bool
	return c.client.Call("ChatServer.Leave", c.userID, &reply)
}

func (c *ChatClient) Close() {
	close(c.done)
	c.client.Close()
}

// receiveMessages polls for new messages using a goroutine
func (c *ChatClient) receiveMessages() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			messages, err := c.GetMessages()
			if err != nil {
				continue
			}
			for _, msg := range messages {
				if msg.UserID == "SYSTEM" {
					fmt.Printf("\n[%s] %s\n> ",
						msg.Timestamp.Format("15:04:05"), msg.Content)
				} else {
					fmt.Printf("\n[%s] %s: %s\n> ",
						msg.Timestamp.Format("15:04:05"), msg.UserID, msg.Content)
				}
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: client <userID> [serverAddr]")
		os.Exit(1)
	}

	userID := os.Args[1]
	serverAddr := "localhost:8080"
	if len(os.Args) >= 3 {
		serverAddr = os.Args[2]
	}

	client, err := NewChatClient(serverAddr, userID)
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer client.Close()

	// Join and get history
	history, err := client.Join()
	if err != nil {
		log.Fatal("Join error:", err)
	}

	fmt.Printf("Connected as [%s]\n", userID)
	fmt.Println("--- Chat History ---")
	for _, msg := range history {
		if msg.UserID == "SYSTEM" {
			fmt.Printf("[%s] %s\n", msg.Timestamp.Format("15:04:05"), msg.Content)
		} else {
			fmt.Printf("[%s] %s: %s\n", msg.Timestamp.Format("15:04:05"), msg.UserID, msg.Content)
		}
	}
	fmt.Println("--- End History ---")
	fmt.Println("Type messages and press Enter. Type 'quit' to exit.")

	// Start goroutine to receive messages
	go client.receiveMessages()

	// Read and send messages from stdin
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "quit" {
			client.Leave()
			fmt.Println("Goodbye!")
			return
		}
		if text != "" {
			if err := client.SendMessage(text); err != nil {
				fmt.Println("Send error:", err)
			}
		}
		fmt.Print("> ")
	}
}
