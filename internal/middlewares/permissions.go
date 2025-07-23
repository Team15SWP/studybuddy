package middlewares

import (
	"log/slog"
	"net/http"

	"study_buddy/pkg/constants"
	"study_buddy/pkg/errlist"
	"study_buddy/pkg/utils"

	"github.com/gin-gonic/gin"
)

var RolePermissions = map[int32]map[string][]string{
	constants.Admin: {
		constants.SyllabusEntity: {constants.Create, constants.Get, constants.Update, constants.Delete},
		constants.TaskEntity:     {constants.Create, constants.Get, constants.Update, constants.Delete},
	},
	constants.User: {
		constants.SyllabusEntity: {constants.Get},
		constants.TaskEntity:     {constants.Create, constants.Get, constants.Update, constants.Delete},
	},
}

func PermissionHandler(entity, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get(constants.Role)
		if !ok {
			_ = c.Error(errlist.ErrForbidden)
			return
		}
		permissions := RolePermissions[role.(int32)][entity]
		for _, permission := range permissions {
			if permission == action {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, utils.SetError(errlist.ErrForbidden.Error()))

		var handlerName interface{}
		handlerName, ex := c.Get(constants.HandlerName)
		if !ex {
			handlerName = ""
		}

		slog.Info("permission middleware", constants.HandlerName, handlerName, "ERR", errlist.ErrForbidden)
	}
}
