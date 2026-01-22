// Package claude provides usage tracking for Claude Code by parsing local session files.
package claude

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SessionData holds parsed data from a single session.
type SessionData struct {
	SessionID       string
	StartTime       time.Time
	EndTime         time.Time
	DurationHours   float64
	PromptCount     int
	SonnetResponses int
	OpusResponses   int
	Project         string
}

// Message represents a single message from the JSONL file.
type Message struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	IsMeta    bool   `json:"isMeta"`
	UserType  string `json:"userType"`
	Message   struct {
		Role    string      `json:"role"`
		Model   string      `json:"model"`
		Content interface{} `json:"content"`
	} `json:"message"`
}

// ParseJSONLFile parses a single JSONL session file.
func ParseJSONLFile(path string) (*SessionData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	session := &SessionData{
		SessionID: filepath.Base(path),
		Project:   filepath.Base(filepath.Dir(path)),
	}

	var timestamps []time.Time
	scanner := bufio.NewScanner(file)

	// Increase buffer size for long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var msg Message
		if err := json.Unmarshal(line, &msg); err != nil {
			continue
		}

		// Parse timestamp
		if msg.Timestamp != "" {
			if ts, err := parseTimestamp(msg.Timestamp); err == nil {
				timestamps = append(timestamps, ts)
			}
		}

		// Count user prompts (excluding meta messages and commands)
		if msg.Type == "user" && msg.Message.Role == "user" && !msg.IsMeta && msg.UserType == "external" {
			if !isCommandMessage(msg.Message.Content) {
				session.PromptCount++
			}
		}

		// Count model responses
		if msg.Type == "assistant" {
			model := strings.ToLower(msg.Message.Model)
			if strings.Contains(model, "opus") {
				session.OpusResponses++
			} else if strings.Contains(model, "sonnet") {
				session.SonnetResponses++
			}
		}
	}

	// Calculate session duration
	if len(timestamps) > 0 {
		session.StartTime = timestamps[0]
		session.EndTime = timestamps[0]
		for _, ts := range timestamps {
			if ts.Before(session.StartTime) {
				session.StartTime = ts
			}
			if ts.After(session.EndTime) {
				session.EndTime = ts
			}
		}
		session.DurationHours = session.EndTime.Sub(session.StartTime).Hours()
	}

	return session, scanner.Err()
}

// parseTimestamp parses an ISO timestamp string.
func parseTimestamp(ts string) (time.Time, error) {
	// Try multiple formats
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, ts); err == nil {
			return t, nil
		}
	}

	// Strip milliseconds and try again
	if idx := strings.Index(ts, "."); idx > 0 {
		if zIdx := strings.Index(ts[idx:], "Z"); zIdx > 0 {
			clean := ts[:idx] + "Z"
			return time.Parse("2006-01-02T15:04:05Z", clean)
		}
	}

	return time.Time{}, nil
}

// isCommandMessage checks if message content is a local command.
func isCommandMessage(content interface{}) bool {
	switch c := content.(type) {
	case string:
		return strings.Contains(c, "<command-name>") || strings.Contains(c, "<local-command-stdout>")
	case []interface{}:
		for _, item := range c {
			if m, ok := item.(map[string]interface{}); ok {
				if text, ok := m["text"].(string); ok {
					if strings.Contains(text, "<command-name>") || strings.Contains(text, "<local-command-stdout>") {
						return true
					}
				}
			}
		}
	}
	return false
}

// GetClaudeProjectsDir returns the path to Claude projects directory.
func GetClaudeProjectsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "projects")
}

// FindAllSessions finds all JSONL session files in the Claude projects directory.
func FindAllSessions() ([]string, error) {
	projectsDir := GetClaudeProjectsDir()

	if _, err := os.Stat(projectsDir); os.IsNotExist(err) {
		return nil, nil // No projects directory yet
	}

	var sessions []string
	err := filepath.Walk(projectsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !info.IsDir() && strings.HasSuffix(path, ".jsonl") {
			sessions = append(sessions, path)
		}
		return nil
	})

	return sessions, err
}
