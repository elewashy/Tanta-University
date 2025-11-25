# Simple Chatroom - RPC Application with Docker

A simple chatroom application built with Go RPC (Remote Procedure Call) that allows clients to send messages to a central server and retrieve the complete chat history.

## Docker Hub Image

**Docker Hub Repository**: [elewashy/simple_chatroom](https://hub.docker.com/r/elewashy/simple_chatroom)

Pull the image:
```bash
docker pull elewashy/simple_chatroom:latest
```

## Video Demonstration

**Video Recording**: [https://drive.google.com/file/d/1MHC9E8geaAcvhlFsUIJ_M7IVyKaV6Ma_/view?usp=sharing](https://drive.google.com/file/d/1MHC9E8geaAcvhlFsUIJ_M7IVyKaV6Ma_/view?usp=sharing)

This video demonstrates the chatroom application in action, showing:
- Server startup and client connections
- Multiple clients sending messages
- Real-time chat history updates
- Graceful shutdown handling

## Repository

**GitHub Repository**: [https://github.com/elewashy/Tanta-University](https://github.com/elewashy/Tanta-University)

### Repository Download Steps

1. **Clone the repository using Git:**
   ```bash
   git clone https://github.com/elewashy/Tanta-University.git
   ```

2. **Navigate to the project directory:**
   ```bash
   cd Tanta-University/Level_3/Semester_1/Distributed_Systems/Assighnment/Simple_Chatroom
   ```

3. **Or download manually:**
   - Go to [https://github.com/elewashy/Tanta-University](https://github.com/elewashy/Tanta-University)
   - Click on "Code" button
   - Select "Download ZIP"
   - Extract the ZIP file
   - Navigate to: `Level_3/Semester_1/Distributed_Systems/Assighnment/Simple_Chatroom`

## Overview

This project implements a client-server chatroom using Go's `net/rpc` package. The server maintains a persistent message history, and clients can send messages and retrieve the complete chat history at any time.

## Project Structure

```
Assignment_3/
├── server.go         # RPC server that stores and manages chat messages
├── client.go         # RPC client that connects to server and sends/receives messages
├── go.mod            # Go module file
├── Dockerfile        # Docker configuration for containerizing the server
├── .dockerignore     # Files to exclude from Docker build
├── shared/
│   └── types.go      # Shared types for RPC communication
└── README.md         # This file
```

## Features

### Server (`server.go`)
- Listens on TCP port 8080 for incoming RPC connections
- Stores all messages in a thread-safe list
- Provides two RPC methods:
  - `SendMessage`: Adds a new message to the history and returns the complete history
  - `GetHistory`: Returns all messages in the chat history
- Uses mutex locks to ensure thread-safe access to the message list

### Client (`client.go`)
- Prompts user to enter their name at startup
- Connects to the RPC server on `localhost:8080`
- Continuously reads user input (whole lines, not just tokens)
- Sends messages to the server via RPC with sender name
- Displays updated chat history after each message with sender names
- Handles graceful shutdown on:
  - Typing "exit" command
  - Pressing Ctrl+C
- Fetches and displays initial chat history upon connection
- Includes error handling for server disconnections

## Docker Deployment

### Building the Docker Image

1. **Build the Docker image:**
   ```bash
   docker build -t chatroom-server .
   ```

2. **Tag the image for Docker Hub:**
   ```bash
   docker tag chatroom-server elewashy/simple_chatroom:latest
   ```

### Running with Docker

1. **Run the container:**
   ```bash
   docker run -d -p 8080:8080 --name chatroom elewashy/simple_chatroom:latest
   ```

2. **Check if the container is running:**
   ```bash
   docker ps
   ```

3. **View server logs:**
   ```bash
   docker logs chatroom
   ```

4. **Stop the container:**
   ```bash
   docker stop chatroom
   ```

5. **Remove the container:**
   ```bash
   docker rm chatroom
   ```

### Pushing to Docker Hub

1. **Login to Docker Hub:**
   ```bash
   docker login
   ```

2. **Push the image:**
   ```bash
   docker push elewashy/simple_chatroom:latest
   ```

### Testing the Dockerized Server

Once the container is running, you can connect to it using the client:

```bash
go run client.go
```

**Note**: The client connects to `localhost:8080`, which works when the Docker container is running with port mapping `-p 8080:8080`.

### Docker Multi-Stage Build

The Dockerfile uses a multi-stage build approach:
- **Stage 1 (Builder)**: Compiles the Go application using `golang:1.21-alpine`
- **Stage 2 (Runtime)**: Creates a minimal runtime image using `alpine:latest` with only the compiled binary

This approach significantly reduces the final image size while maintaining full functionality.

## How to Run (Without Docker)

### Prerequisites
- Go 1.21 or higher installed on your system

### Starting the Server

1. Open a terminal/command prompt
2. Navigate to the project directory
3. Run the server:
   ```bash
   go run server.go
   ```

You should see:
```
Chatroom server started on port 8080
Waiting for clients...
```

### Running the Client

1. Open a **new** terminal/command prompt
2. Navigate to the project directory
3. Run the client:
   ```bash
   go run client.go
   ```

You should see:
```
Enter your name: [Your Name]
Welcome, [Your Name]!
Connected to chatroom server!
Type messages and press Enter. Type 'exit' to quit.
----------------------------------------
```

4. Start typing messages and press Enter after each message
5. To exit, type `exit` or press Ctrl+C

### Running Multiple Clients

You can run multiple client instances simultaneously. Each client will:
- See the complete chat history when connecting
- Send messages that are visible to all clients
- Receive updated history after each message

## Technical Details

### Input Handling
The client uses `bufio.Scanner` instead of `fmt.Scan` to read complete lines of input, allowing messages with spaces to be sent correctly.

### Error Handling
- Server connection errors are caught and displayed to the user
- If the server goes down, the client detects the error and exits gracefully
- All RPC calls include error checking

### Thread Safety
The server uses `sync.RWMutex` to ensure thread-safe access to the message list when multiple clients are connected simultaneously.

## RPC Methods

### SendMessage
- **Arguments**: `SendMessageArgs{Sender: string, Content: string}`
- **Returns**: `SendMessageReply{Success: bool, History: []Message}`
- **Description**: Adds a message to the server's history and returns the complete history. Messages include the sender's name.

### GetHistory
- **Arguments**: `GetHistoryArgs{}`
- **Returns**: `GetHistoryReply{History: []Message}`
- **Description**: Retrieves all messages stored on the server

## Notes

- The server must be running before clients can connect
- Messages are stored in memory and will be lost when the server stops
- The server can handle multiple concurrent client connections
- Messages are displayed with sender name and sequential number: `[1] SenderName: message content`
- Each user must enter their name when starting the client
- If port 8080 is already in use, stop the existing server process before starting a new one

