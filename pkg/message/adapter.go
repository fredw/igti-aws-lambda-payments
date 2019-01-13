package message

// Adapter represents an adapter to handle the messages
type Adapter interface {
	GetMessages() (Messages, error)
	Delete(id *string) error
	MoveToDLQ(m Message) error
}
