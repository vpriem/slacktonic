package slacktonic

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/slack-go/slack"
)

//go:generate go run go.uber.org/mock/mockgen -package=slacktonic_test -destination=verifier_mock_test.go . Verifier

// Verifier is an interface for verifying Slack requests and parsing slash commands.
// It allows for custom verification logic or mocking during tests.
type Verifier interface {
	// Verify verifies the authenticity of a Slack request and parses the slash command.
	// It returns the parsed slash command and nil if verification succeeds,
	// or nil and an error if verification fails.
	Verify(*http.Request) (slack.SlashCommand, error)
}

// SlackVerifier is the default implementation of the Verifier interface.
// It verifies Slack request signatures using the provided signing secret.
type SlackVerifier struct {
	secret string
}

// NewVerifier creates a new SlackVerifier with the provided signing secret.
func NewVerifier(secret string) *SlackVerifier {
	return &SlackVerifier{secret}
}

// Verify verifies the authenticity of a Slack request using the signing secret
// and parses the slash command from the request.
// It returns the parsed slash command and nil if verification succeeds,
// or nil and an error if verification fails.
func (s *SlackVerifier) Verify(req *http.Request) (slack.SlashCommand, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return slack.SlashCommand{}, fmt.Errorf("failed to read body: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	verifier, err := slack.NewSecretsVerifier(req.Header, s.secret)
	if err != nil {
		return slack.SlashCommand{}, fmt.Errorf("failed to create verifier: %w", err)
	}
	if _, err = verifier.Write(bodyBytes); err != nil {
		return slack.SlashCommand{}, fmt.Errorf("failed to write body to verifier: %w", err)
	}
	if err = verifier.Ensure(); err != nil {
		return slack.SlashCommand{}, fmt.Errorf("%w: %v", ErrUnauthorized, err)
	}

	cmd, err := slack.SlashCommandParse(req)
	if err != nil {
		return slack.SlashCommand{}, fmt.Errorf("failed to parse slash command: %w", err)
	}
	return cmd, nil
}
