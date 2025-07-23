package middlewares

import (
	"errors"
	"log/slog"
	"net/http"

	"study_buddy/internal/model"
	"study_buddy/pkg/errlist"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			statusCode := getStatusCode(err.Err)

			last := err.Err
			for errors.Unwrap(last) != nil {
				last = errors.Unwrap(last)
			}

			log.Error("ERR: ", err.Error())

			c.AbortWithStatusJSON(statusCode, model.ErrorResponse{
				Error: last.Error(),
			})
		}
	}
}

func getStatusCode(err error) int {
	// map custom errors to appropriate HTTP status codes
	switch {
	case errors.Is(err, errlist.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, errlist.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, errlist.ErrUnauthorized),
		errors.Is(err, errlist.ErrUserNotFound),
		errors.Is(err, errlist.ErrInvalidPassword),
		errors.Is(err, errlist.ErrUserIsNotVerified):
		return http.StatusUnauthorized
	case errors.Is(err, errlist.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, errlist.ErrAlreadyExists),
		errors.Is(err, errlist.ErrUserExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
