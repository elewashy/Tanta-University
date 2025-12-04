# RPC Chat System with Real-Time Broadcasting

A distributed chat system built with Go using RPC (Remote Procedure Call) for communication between multiple clients and a server. Features real-time message broadcasting using Go concurrency primitives (goroutines, channels, and Mutex).

## Demo Video

[![Demo Video](https://img.shields.io/badge/Watch-Demo%20Video-red?style=for-the-badge&logo=googledrive)](https://drive.google.com/file/d/1nLSSE170dKaQdzgKLsG4hqNf4awR58k0/view?usp=sharing)

> 📹 **[Click here to watch the demo video](https://drive.google.com/file/d/1nLSSE170dKaQdzgKLsG4hqNf4awR58k0/view?usp=sharing)**

---

## Features

- **Real-time Broadcasting**: Messages are instantly broadcast to all connected clients
- **Join Notifications**: All clients are notified when a new user joins ("User [ID] joined")
- **No Self-Echo**: Senders don't receive their own messages back
- **Chat History**: New clients receive full chat history upon joining
- **Concurrent Design**: Uses goroutines and channels for non-blocking message handling
- **Thread-Safe**: Shared client list protected with `sync.Mutex`

## Project Structure

```
RPC_chat_system/
├── server/
│   └── server.go    # RPC server with broadcasting logic
├── client/
│   └── client.go    # Interactive chat client
├── go.mod
└── README.md
```

## Requirements

- Go 1.21 or higher

## Installation & Running

### 1. Clone the repository

```bash
git clone https://github.com/elewashy/Tanta-University.git
cd Tanta-University/Level_3/Semester_1/Distributed_Systems/Assighnment/RPC_chat_system
```

### 2. Build the project

```bash
go build -o server.exe ./server
go build -o client.exe ./client
```

### 3. Run the server

```bash
./server.exe
```

### 4. Run clients (in separate terminals)

```bash
./client.exe Mohamed
./client.exe Ahmed
./client.exe Tamer
```

## Usage

1. Start the server first
2. Connect multiple clients with unique usernames
3. Type messages and press Enter to send
4. Type `quit` to exit

## Architecture

```
┌─────────┐     RPC      ┌─────────────┐     RPC      ┌─────────┐
│ Client1 │◄────────────►│             │◄────────────►│ Client2 │
└─────────┘              │   Server    │              └─────────┘
                         │             │
┌─────────┐     RPC      │  - Mutex    │
│ Client3 │◄────────────►│  - Channels │
└─────────┘              │  - History  │
                         └─────────────┘
```

## Concurrency Model

| Component | Purpose |
|-----------|---------|
| `sync.Mutex` | Protects shared client map and message history |
| `chan Message` | Per-client buffered channel for incoming broadcasts |
| Goroutine (server) | Handles each client connection concurrently |
| Goroutine (client) | Polls for new messages while main thread handles input |

