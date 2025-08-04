# slacktonic

[![Go Reference](https://pkg.go.dev/badge/github.com/vpriem/slacktonic.svg)](https://pkg.go.dev/github.com/vpriem/slacktonic)
[![Go Report Card](https://goreportcard.com/badge/github.com/vpriem/slacktonic)](https://goreportcard.com/report/github.com/vpriem/slacktonic)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go middleware for handling Slack slash commands in [Gin](https://github.com/gin-gonic/gin) applications.

## Features

- Verifies Slack request signatures using your signing secret
- Parses slash commands from requests
- Makes slash commands available in your Gin handlers

## Installation

```bash
go get github.com/vpriem/slacktonic
```

## Usage

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vpriem/slacktonic"
)

func main() {
	r := gin.Default()

	// Add the middleware to verify Slack requests
	r.Use(slacktonic.Middleware("your-slack-signing-secret"))

	// Handle slash commands
	r.POST("/slack", func(c *gin.Context) {
		// Get the slash command from the context
		cmd, ok := slacktonic.GetSlashCommand(c)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No slash command found"})
			return
		}

		// Handle the command
		switch cmd.Command {
		case "/hello":
			c.JSON(http.StatusOK, gin.H{"text": "Hello, " + cmd.UserName + "!"})
		default:
			c.JSON(http.StatusOK, gin.H{"text": "Unknown command: " + cmd.Command})
		}
	})

	r.Run(":8080")
}
```

## Configuration Options

The middleware can be configured using option functions passed to the `Middleware` function:

### Custom Logging

By default, the middleware uses a no-op logger that discards all log messages. You can provide your own logger that implements the `Logger` interface using the `WithLogger` option function:

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
r.Use(slacktonic.Middleware("your-slack-signing-secret", slacktonic.WithLogger(logger)))
```

The logger interface is compatible with Go's standard `log/slog` package.

### Custom Verifier

By default, the middleware uses a `SlackVerifier` that verifies Slack request signatures using the provided signing secret. You can provide your own verifier that implements the `Verifier` interface using the `WithVerifier` option function:

```go
// For testing purposes
mockVerifier := &MockVerifier{} // implements the Verifier interface
r.Use(slacktonic.Middleware("", slacktonic.WithVerifier(mockVerifier)))
```

This option is primarily intended for testing your integration with Slack. It allows you to mock the verification process during tests, making it possible to test your handlers without needing real Slack credentials or making actual API calls.

### Command Validation

You can specify an expected command to validate incoming slash commands. If the received command doesn't match the expected one, the middleware will return a 400 Bad Request response:

```go
// Only accept the /hello command
r.Use(slacktonic.Middleware("your-slack-signing-secret", slacktonic.WithCommand("/hello")))
```

This is useful when you want to ensure that only specific commands are processed by your handler.

## License

See [LICENSE](LICENSE) file.
