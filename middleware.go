// Package slacktonic provides middleware for handling Slack slash commands in Gin applications.
//
// It verifies the authenticity of incoming Slack requests using the signing secret
// and parses slash commands, making them available in the Gin context.
//
// Example usage:
//
//	r := gin.New()
//	r.Use(slacktonic.Middleware("your-slack-signing-secret"))
//	r.POST("/slack", func(c *gin.Context) {
//		cmd, ok := slacktonic.GetSlashCommand(c)
//		if ok {
//			// Handle the slash command
//		}
//	})
package slacktonic

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Options configures the behavior of the middleware.
// It is configured through option functions passed to the Middleware function.
type Options struct {
	Logger   Logger
	Command  string
	Verifier Verifier
}

// Middleware returns a Gin middleware function that verifies Slack requests
// and parses slash commands.
//
// The middleware:
// 1. Verifies the request using the provided Verifier (or creates a SlackVerifier with the provided secret)
// 2. Parses the slash command from the request
// 3. Validates the command against the expected command (if specified using WithCommand)
// 4. Stores the slash command in the Gin context for later retrieval
//
// Parameters:
//   - secret: The Slack signing secret used to verify request authenticity
//   - optFns: Optional configuration functions for the middleware (WithCommand, WithLogger, WithVerifier)
//
// Returns a Gin HandlerFunc that can be used with the Gin router.
func Middleware(secret string, optFns ...Option) gin.HandlerFunc {
	opts := &Options{
		Logger:   &noopLogger{},
		Verifier: NewVerifier(secret),
	}
	for _, optFn := range optFns {
		optFn(opts)
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		cmd, err := opts.Verifier.Verify(c.Request)
		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				opts.Logger.WarnContext(ctx, err.Error(), "error", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}

			opts.Logger.WarnContext(ctx, err.Error(), "error", errors.Unwrap(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}

		if opts.Command != "" && opts.Command != cmd.Command {
			opts.Logger.WarnContext(ctx, "command mismatch", "expected", opts.Command, "received", cmd.Command)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "command mismatch"})
			return
		}

		opts.Logger.InfoContext(ctx, "received slack slash command ", "command", cmd)
		c.Set(contextKey, cmd)
		c.Next()
	}
}

type Option func(*Options)

// WithCommand specifies the expected slash command.
// If provided, the middleware will check if the received command matches this value.
// If the command doesn't match, a 400 Bad Request response will be returned.
func WithCommand(command string) Option {
	return func(opt *Options) {
		opt.Command = command
	}
}

// WithLogger sets a logger for the middleware.
// Logger is used to log errors that occur during request processing.
// If not provided, a no-op logger will be used.
func WithLogger(l Logger) Option {
	return func(opt *Options) {
		if l != nil {
			opt.Logger = l
		}
	}
}

// WithVerifier sets a custom verifier for the middleware.
// Verifier is used to verify Slack requests and parse slash commands.
// If not provided, a SlackVerifier will be created using the provided secret.
func WithVerifier(v Verifier) Option {
	return func(opt *Options) {
		if v != nil {
			opt.Verifier = v
		}
	}
}
