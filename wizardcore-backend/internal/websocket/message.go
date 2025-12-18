package websocket

import (
	"encoding/json"
)

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// MatchJoin indicates a user wants to join a match
	MatchJoin MessageType = "match_join"
	// MatchLeave indicates a user leaves a match
	MatchLeave MessageType = "match_leave"
	// CodeUpdate indicates a code update from a participant
	CodeUpdate MessageType = "code_update"
	// MatchStart indicates the match has started
	MatchStart MessageType = "match_start"
	// MatchEnd indicates the match has ended
	MatchEnd MessageType = "match_end"
	// Error indicates an error message
	Error MessageType = "error"
	// Ping is used for keep-alive
	Ping MessageType = "ping"
	// Pong is used for keep-alive response
	Pong MessageType = "pong"
)

// Message represents a WebSocket message
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// MatchJoinPayload payload for MatchJoin
type MatchJoinPayload struct {
	MatchID string `json:"match_id"`
	UserID  string `json:"user_id"`
}

// CodeUpdatePayload payload for CodeUpdate
type CodeUpdatePayload struct {
	MatchID string `json:"match_id"`
	UserID  string `json:"user_id"`
	Code    string `json:"code"`
	LanguageID int `json:"language_id"`
}

// MatchStartPayload payload for MatchStart
type MatchStartPayload struct {
	MatchID      string `json:"match_id"`
	ExerciseID   string `json:"exercise_id"`
	TimeLimit    int    `json:"time_limit"`
	Participants []string `json:"participants"`
}

// MatchEndPayload payload for MatchEnd
type MatchEndPayload struct {
	MatchID string `json:"match_id"`
	Results []ParticipantResult `json:"results"`
}

// ParticipantResult result of a participant
type ParticipantResult struct {
	UserID string `json:"user_id"`
	Score  int    `json:"score"`
	Result string `json:"result"` // win, loss, draw
	XP     int    `json:"xp"`
}