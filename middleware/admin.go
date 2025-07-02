package middleware

import (
	log "log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/startrek92/kube-admission-webhook/auth"
	"github.com/startrek92/kube-admission-webhook/config"
	"github.com/startrek92/kube-admission-webhook/utils"
)

const UserContextKey = "user"

func AdminAuthMiddleware(conf config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authServerConfig := conf.AuthServer
		log.Info("auth server config", "value", authServerConfig)
		verify_token_url := authServerConfig.URL + "/" + authServerConfig.VerifyTokenUrl
		authenticator := &auth.BearerAuthenticator{
			AuthServerURL: verify_token_url,
			HttpClient:    *utils.NewHTTPSClient(true, 3*time.Second),
			Timeout:       2 * time.Second,
		}

		user, err := authenticator.AuthenticateWithContext(ctx, c)
		if err != nil {
			log.Warn("Authentication failed", "path", c.FullPath(), "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Set(UserContextKey, user)
		log.Info("Authenticated", "path", c.FullPath(), "user", user.Email, "role", user.Role)

		c.Next()
	}
}
