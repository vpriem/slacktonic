package slacktonic_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vpriem/slacktonic"
)

func TestVerifier(t *testing.T) {
	v := slacktonic.NewVerifier("secret")

	t.Run("should return error on missing headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/slack", nil)

		cmd, err := v.Verify(req)
		require.Equal(t, "", cmd.Command)
		require.EqualError(t, err, "failed to create verifier: missing headers")
	})

	t.Run("should return error on old timestamp", func(t *testing.T) {
		form := url.Values{}
		form.Add("command", "/test")
		body := form.Encode()
		ts := time.Now().AddDate(0, 0, -1).Unix()

		req := httptest.NewRequest(http.MethodPost, "/slack", strings.NewReader(body))
		req.Header.Set("X-Slack-Signature", sign(t, "secret", ts, body))
		req.Header.Set("X-Slack-Request-Timestamp", fmt.Sprintf("%d", ts))

		cmd, err := v.Verify(req)
		require.Equal(t, "", cmd.Command)
		require.EqualError(t, err, "failed to create verifier: timestamp is too old")
	})

	t.Run("should return ErrUnauthorized on invalid signature", func(t *testing.T) {
		form := url.Values{}
		form.Add("command", "/test")
		body := form.Encode()
		ts := time.Now().Unix()

		req := httptest.NewRequest(http.MethodPost, "/slack", strings.NewReader(body))
		req.Header.Set("X-Slack-Signature", sign(t, "wrong-secret", ts, body))
		req.Header.Set("X-Slack-Request-Timestamp", fmt.Sprintf("%d", ts))

		cmd, err := v.Verify(req)
		require.Equal(t, "", cmd.Command)
		require.Error(t, err)
		require.ErrorIs(t, err, slacktonic.ErrUnauthorized)
	})

	t.Run("should success", func(t *testing.T) {
		form := url.Values{}
		form.Add("command", "/test")
		body := form.Encode()
		ts := time.Now().Unix()

		req := httptest.NewRequest(http.MethodPost, "/slack", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-Slack-Signature", sign(t, "secret", ts, body))
		req.Header.Set("X-Slack-Request-Timestamp", fmt.Sprintf("%d", ts))

		cmd, err := v.Verify(req)
		require.Equal(t, "/test", cmd.Command)
		require.NoError(t, err)
	})

}

func sign(t *testing.T, secret string, timestamp int64, body string) string {
	baseString := fmt.Sprintf("v0:%d:%s", timestamp, body)
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(baseString))
	if err != nil {
		t.Fatal(err)
	}
	return "v0=" + hex.EncodeToString(h.Sum(nil))
}
