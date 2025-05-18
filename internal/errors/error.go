package errors

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func New(status int, message string) *ApiError {
	return &ApiError{Status: status, Message: message}
}

func Respond(c *gin.Context, err error) {
	var apiErr *ApiError
	if errors.As(err, &apiErr) {
		c.JSON(apiErr.Status, apiErr)
	}
}
