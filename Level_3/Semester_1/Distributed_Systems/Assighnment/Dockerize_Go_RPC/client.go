package main

import (
	"bufio"
	"chatroom/shared"
	"fmt"
	"net/rpc"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// Connect to the RPC server
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		fmt.Println("Make sure the server is running on port 8080")
		return
	}
	defer client.Close()

	// Prompt for user name
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your name: ")
	if !scanner.Scan() {
		fmt.Println("Error reading name")
		return
	}
	userName := strings.TrimSpace(scanner.Text())
	if userName == "" {
		userName = "Anonymous"
	}

	fmt.Printf("\nWelcome, %s!\n", userName)
	fmt.Println("Connected to chatroom server!")
	fmt.Println("Type messages and press Enter. Type 'exit' to quit.")
	fmt.Println("----------------------------------------")

	// Set up signal handler for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Channel to signal when we should exit
	exitChan := make(chan bool, 1)

	// Goroutine to handle Ctrl+C
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		exitChan <- true
	}()

	// Fetch initial history
	fmt.Println("\nChat History:")
	fmt.Println("----------------------------------------")
	fetchHistory(client)

	// Main loop: read input and send messages
	for {
		fmt.Print("\nEnter message: ")

		// Check if we should exit
		select {
		case <-exitChan:
			return
		default:
		}

		// Read a line of input
		if !scanner.Scan() {
			break
		}

		message := strings.TrimSpace(scanner.Text())

		// Check for exit command
		if message == "exit" {
			fmt.Println("Exiting chatroom...")
			return
		}

		// Skip empty messages
		if message == "" {
			continue
		}

		// Send message to server
		args := shared.SendMessageArgs{Sender: userName, Content: message}
		var reply shared.SendMessageReply

		err := client.Call("ChatServer.SendMessage", &args, &reply)
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			fmt.Println("Server connection may be lost. Please restart the client.")
			return
		}

		// Display updated history
		if reply.Success {
			fmt.Println("\nUpdated Chat History:")
			fmt.Println("----------------------------------------")
			printHistory(reply.History)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}

// fetchHistory retrieves and displays the chat history from the server
func fetchHistory(client *rpc.Client) {
	args := shared.GetHistoryArgs{}
	var reply shared.GetHistoryReply

	err := client.Call("ChatServer.GetHistory", &args, &reply)
	if err != nil {
		fmt.Printf("Error fetching history: %v\n", err)
		return
	}

	printHistory(reply.History)
}

// printHistory displays the chat history
func printHistory(history []shared.Message) {
	if len(history) == 0 {
		fmt.Println("(No messages yet)")
		return
	}

	for i, msg := range history {
		fmt.Printf("[%d] %s: %s\n", i+1, msg.Sender, msg.Content)
	}
}
