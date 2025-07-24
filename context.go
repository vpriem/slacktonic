// Package slacktonic provides middleware for handling Slack slash commands in Gin applications.
package slacktonic

import (
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
)

// contextKey is the key used to store the slash command in the Gin context.
const contextKey = "slack-slash-command"

// GetSlashCommand retrieves the Slack slash command from the Gin context.
//
// This function should be called after the Middleware has processed the request.
// The Middleware stores the slash command in the context if it was successfully
// parsed from the request.
//
// Parameters:
//   - c: The Gin context from which to retrieve the slash command
//
// Returns:
//   - slack.SlashCommand: The slash command if it exists in the context
//   - bool: true if the slash command was found, false otherwise
func GetSlashCommand(c *gin.Context) (slack.SlashCommand, bool) {
	if c == nil {
		return slack.SlashCommand{}, false
	}

	value, exists := c.Get(contextKey)
	if !exists {
		return slack.SlashCommand{}, false
	}

	cmd, ok := value.(slack.SlashCommand)
	if !ok {
		return slack.SlashCommand{}, false
	}

	return cmd, true
}
