package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/uchennaemeruche/go-bank-api/api/util"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
	"github.com/uchennaemeruche/go-bank-api/token"
	"github.com/uchennaemeruche/go-bank-api/transfer"
	"github.com/uchennaemeruche/go-bank-api/transfer/service"
)

type TransferHandler interface {
	CreateTransfer(*gin.Context)
}

type handler struct {
	service service.TransferService
}

func NewTransferHandler(service service.TransferService) TransferHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateTransfer(ctx *gin.Context) {
	var input transfer.TransferReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	account, isValid := h.validAcount(ctx, input.FromAccountID, input.Currency)
	if !isValid {
		return
	}

	authPayload := ctx.MustGet(api.AuthPayloadKey).(*token.Payload)
	if authPayload.Username != account.Owner {
		ctx.JSON(http.StatusUnauthorized, api.ErrorResponse(errors.New("drawing account does not belong to the currently logged in user")))
		return
	}

	// _, isValid = h.validAcount(ctx, input.ToAccountID, input.Currency)
	if _, isValid = h.validAcount(ctx, input.ToAccountID, input.Currency); !isValid {
		return
	}

	result, err := h.service.Create(input.FromAccountID, input.ToAccountID, input.Amount, input.Currency)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (h *handler) validAcount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := h.service.ValidAccount(accountId, currency)
	if err != nil {
		if err.(*api.RequestError).Code == 404 {
			ctx.JSON(http.StatusNotFound, api.ErrorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return account, false
	}

	return account, true

}
