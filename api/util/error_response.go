package api

import (
	"github.com/gin-gonic/gin"
)

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

type RequestError struct {
	Err  error
	Code int
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}
