package interfaces

import "context"

type IChatService interface {
	// SendChat is method to send chat
	SendChat(ctx context.Context, message string) (string, error)
}
