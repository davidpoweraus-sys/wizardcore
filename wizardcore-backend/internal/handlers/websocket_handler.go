package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	gorillawebsocket "github.com/gorilla/websocket"
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/middleware"
	internalws "github.com/yourusername/wizardcore-backend/internal/websocket"
)

var upgrader = gorillawebsocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, you should validate the origin
		return true
	},
}

type WebSocketHandler struct {
	hub *internalws.Hub
}

func NewWebSocketHandler(hub *internalws.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// ServeWebSocket handles WebSocket connections
func (h *WebSocketHandler) ServeWebSocket(c *gin.Context) {
	// Get user ID from authentication middleware
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	client := internalws.NewClient(h.hub, conn, userID)
	h.hub.Register <- client

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}

// JoinMatchRoom handles a request to join a match room
func (h *WebSocketHandler) JoinMatchRoom(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	matchIDStr := c.Param("match_id")
	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	// In a real implementation, we would verify that the user is a participant of the match
	// For now, we'll just allow joining

	// We need to find the client for this user and add to room
	// Since we don't have a mapping from userID to client, we'll need to store it elsewhere.
	// For simplicity, we'll just broadcast a message that the user joined.
	// This is a placeholder; actual implementation would require more logic.

	// Dummy usage to avoid compiler warnings
	_ = userID
	_ = matchID

	c.JSON(http.StatusOK, gin.H{"message": "Join request accepted"})
}