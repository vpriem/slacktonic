package slacktonic_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/require"
	"github.com/vpriem/slacktonic"
)

func TestGetSlashCommand(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		setupContext    func(*gin.Context)
		expectedExists  bool
		expectedCommand slack.SlashCommand
	}{
		{
			name:            "context key does not exists",
			setupContext:    func(c *gin.Context) {},
			expectedExists:  false,
			expectedCommand: slack.SlashCommand{},
		},
		{
			name: "is not a SlashCommand",
			setupContext: func(c *gin.Context) {
				c.Set("slack-slash-command", "not a slash command")
			},
			expectedExists:  false,
			expectedCommand: slack.SlashCommand{},
		},
		{
			name: "is a SlashCommand",
			setupContext: func(c *gin.Context) {
				c.Set("slack-slash-command", slack.SlashCommand{
					Command: "/test",
					Text:    "example text",
				})
			},
			expectedExists: true,
			expectedCommand: slack.SlashCommand{
				Command: "/test",
				Text:    "example text",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			tt.setupContext(c)

			cmd, exists := slacktonic.GetSlashCommand(c)
			require.Equal(t, tt.expectedExists, exists)
			require.Equal(t, tt.expectedCommand, cmd)
		})
	}
}
