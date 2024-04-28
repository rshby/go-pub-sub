package controller

import (
	"github.com/gin-gonic/gin"
	"go-pub-sub/internal/service/interfaces"
	"net/http"
)

type ChatController struct {
	ChatServce interfaces.IChatService
}

func NewChatController(chatService interfaces.IChatService) *ChatController {
	return &ChatController{ChatServce: chatService}
}

// Send is handler to send message from client
func (c *ChatController) Send(ctx *gin.Context) {
	message := ctx.Query("message")
	if message == "" {
		statusCode := http.StatusBadRequest
		ctx.JSON(statusCode, gin.H{
			"message": "message params required",
		})
		return
	}

	// call service to pucblish message to google pub sub
	messageId, err := c.ChatServce.SendChat(ctx, message)
	if err != nil {
		statusCode := http.StatusInternalServerError
		ctx.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	// success
	statusCode := http.StatusOK
	ctx.JSON(statusCode, gin.H{
		"message":    "success send chat message",
		"message_id": messageId,
	})
}
