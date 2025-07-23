package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"study_buddy/internal/config"
	"study_buddy/pkg/constants"
	"study_buddy/pkg/errlist"
	"study_buddy/pkg/utils"

	"github.com/gin-gonic/gin"
)

func UserIdentity(hashCfg *config.HashConfig) gin.HandlerFunc {
	return func(g *gin.Context) {
		header := g.GetHeader(constants.AuthorizationHeader)
		headerParts := strings.Split(header, " ")
		if header == "" || len(headerParts) != 2 || headerParts[0] != "Bearer" {
			err := fmt.Errorf("invalid authorization header: %w", errlist.ErrUnauthorized)
			g.AbortWithStatusJSON(http.StatusUnauthorized, utils.SetError(err.Error()))
			return
		}

		if len(headerParts[1]) == 0 {
			err := fmt.Errorf("empty token: %w", errlist.ErrUnauthorized)
			g.AbortWithStatusJSON(http.StatusUnauthorized, utils.SetError(err.Error()))
			return
		}

		claims, err := ParseToken(headerParts[1], hashCfg)
		if err != nil {
			utils.SetAuthorizationToken(g, "")
			g.AbortWithStatusJSON(http.StatusUnauthorized, utils.SetError(err.Error()))
			return
		}

		g.Set(constants.UserID, claims.UserID)
		g.Set(constants.Role, claims.Role)
		g.Set(constants.Name, claims.Name)
		g.Set(constants.Email, claims.Email)
		g.Next()
	}
}
