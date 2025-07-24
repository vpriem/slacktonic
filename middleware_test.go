package slacktonic_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vpriem/slacktonic"
	"go.uber.org/mock/gomock"
)

func TestMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	v := NewMockVerifier(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(slacktonic.Middleware("secret",
		slacktonic.WithVerifier(v),
		slacktonic.WithLogger(nil),
	))
	r.POST("/slack", func(c *gin.Context) {
		cmd, ok := slacktonic.GetSlashCommand(c)
		if ok {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"command": cmd.Command})
		} else {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
		}
	})

	t.Run("should return 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/slack", nil)

		v.EXPECT().Verify(req).Return(slack.SlashCommand{}, slacktonic.ErrUnauthorized)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		require.JSONEq(t, `{"error":"unauthorized"}`, w.Body.String())
	})

	t.Run("should return 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/slack", nil)

		v.EXPECT().Verify(req).Return(slack.SlashCommand{}, assert.AnError)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.JSONEq(t, `{"error":"invalid signature"}`, w.Body.String())
	})

	t.Run("should succeed", func(t *testing.T) {
		form := url.Values{}
		form.Add("command", "/test")
		body := form.Encode()

		req := httptest.NewRequest(http.MethodPost, "/slack", strings.NewReader(body))

		cmd := slack.SlashCommand{Command: "/test"}
		v.EXPECT().Verify(req).Return(cmd, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, `{"command":"/test"}`, w.Body.String())
	})
}

func TestMiddlewareWithExpectedCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	v := NewMockVerifier(ctrl)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(slacktonic.Middleware("secret",
		slacktonic.WithCommand("/test"),
		slacktonic.WithVerifier(v),
	))
	r.POST("/slack", func(c *gin.Context) {
		cmd, ok := slacktonic.GetSlashCommand(c)
		if ok {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"command": cmd.Command})
		} else {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
		}
	})

	t.Run("should accept command", func(t *testing.T) {
		form := url.Values{}
		form.Add("command", "/test")
		body := form.Encode()

		req := httptest.NewRequest(http.MethodPost, "/slack", strings.NewReader(body))

		cmd := slack.SlashCommand{Command: "/test"}
		v.EXPECT().Verify(req).Return(cmd, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, `{"command":"/test"}`, w.Body.String())
	})

	t.Run("should reject non-matching command", func(t *testing.T) {
		form := url.Values{}
		form.Add("command", "/nope")
		body := form.Encode()

		req := httptest.NewRequest(http.MethodPost, "/slack", strings.NewReader(body))

		cmd := slack.SlashCommand{Command: "/foo"}
		v.EXPECT().Verify(req).Return(cmd, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.JSONEq(t, `{"error":"command mismatch"}`, w.Body.String())
	})
}
