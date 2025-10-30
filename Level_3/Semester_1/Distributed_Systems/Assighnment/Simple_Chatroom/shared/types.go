package shared

// Message represents a single chat message
type Message struct {
	Sender  string
	Content string
}

// SendMessageArgs holds the arguments for sending a message
type SendMessageArgs struct {
	Sender  string
	Content string
}

// SendMessageReply holds the reply after sending a message
type SendMessageReply struct {
	Success bool
	History []Message
}

// GetHistoryArgs holds arguments for getting history
type GetHistoryArgs struct {
}

// GetHistoryReply holds the reply with chat history
type GetHistoryReply struct {
	History []Message
}
