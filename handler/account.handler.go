package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	GetAccount(*gin.Context)
	CreateAccount()
	// ListAccount()
}

type handler struct {
	// service string
}

func NewAccountHandler() AccountHandler {
	return &handler{
		// service: service,
	}
}

func (h *handler) GetAccount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Done")
}
func (h *handler) CreateAccount() {

}
