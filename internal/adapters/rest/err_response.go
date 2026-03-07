package rest

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// errorResponse contains a user-facing message and an optional internal error
// string. The internal error is omitted when empty.
type errorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func newErrorResponse(c *gin.Context, log *slog.Logger, statusCode int, message string, err error) {
	// Log the internal error server-side
	if err != nil {
		log.Error(message, slog.String("error", err.Error()))
	} else {
		log.Error(message)
	}

	// Return both a user-facing message and the internal error (if present).
	// Note: exposing internal error text may leak implementation details; if
	// that's a concern, remove the Error field or set it conditionally.
	errText := ""
	if err != nil {
		errText = err.Error()
	}

	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message, Error: errText})
}
