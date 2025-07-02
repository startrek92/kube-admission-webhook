package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	log "log/slog"

	"github.com/gin-gonic/gin"
	"github.com/startrek92/kube-admission-webhook/utils"
)

type BearerAuthenticator struct {
	AuthServerURL string
	HttpClient    utils.HTTPSClient
	Timeout       time.Duration
}

func (b *BearerAuthenticator) Authenticate(c *gin.Context) (*AuthenticatedUserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout)
	defer cancel()
	return b.AuthenticateWithContext(ctx, c)
}

func (b *BearerAuthenticator) AuthenticateWithContext(ctx context.Context, c *gin.Context) (*AuthenticatedUserInfo, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		log.Warn("Authorization header missing", "path", c.FullPath())
		return nil, errors.New("missing Authorization header")
	}

	log.Info("Calling auth server", "url", b.AuthServerURL, "path", c.FullPath())

	req, err := http.NewRequestWithContext(ctx, "GET", b.AuthServerURL, nil)
	if err != nil {
		log.Error("Failed to create auth server request", "error", err)
		return nil, err
	}
	req.Header.Set("Authorization", token)

	resp, err := b.HttpClient.Do(req)
	if err != nil {
		log.Error("Auth server request failed", "error", err)
		return nil, errors.New("auth server request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("Auth server denied token", "status", resp.StatusCode)
		return nil, errors.New("unauthorized")
	}

	user := &AuthenticatedUserInfo{
		ID:    "1234",
		Email: "admin@example.com",
		Role:  "admin",
	}

	log.Info("Token validated", "user", user.Email, "role", user.Role)
	return user, nil
}
