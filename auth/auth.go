package auth

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Authenticator interface {
	Authenticate(c *gin.Context) (*AuthenticatedUserInfo, error)
	AuthenticateWithContext(ctx context.Context, c *gin.Context) (*AuthenticatedUserInfo, error)
}
