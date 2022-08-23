package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uchennaemeruche/go-bank-api/account/entity"
	"github.com/uchennaemeruche/go-bank-api/account/service"
	api "github.com/uchennaemeruche/go-bank-api/api/util"
)

type AccountHandler interface {
	GetAccount(*gin.Context)
	CreateAccount(*gin.Context)
	ListAccount(*gin.Context)
	UpdateAccount(ctx *gin.Context)
	DeleteAccount(ctx *gin.Context)
}

type handler struct {
	service service.AccountService
}

func NewAccountHandler(service service.AccountService) AccountHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) GetAccount(ctx *gin.Context) {

	var uri entity.GetAccountReq

	if err := ctx.ShouldBindUri(&uri); err != nil {
		if strings.Contains(err.Error(), "strconv.ParseInt") {
			ctx.JSON(http.StatusBadRequest, api.ErrorResponse(errors.New("string not allowed")))
			return
		}
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	account, err := h.service.GetOne(uri.ID)
	if err != nil {
		if err.(*api.RequestError).Code == 404 {
			ctx.JSON(http.StatusNotFound, api.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
func (h *handler) CreateAccount(ctx *gin.Context) {
	var input entity.CreateAccountReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	account, err := h.service.Create(input.Owner, input.Currency, input.AccountType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

func (h *handler) ListAccount(ctx *gin.Context) {
	var input entity.ListAccountReq
	if err := ctx.ShouldBindQuery(&input); err != nil {
		if strings.Contains(err.Error(), "strconv.ParseInt") {
			ctx.JSON(http.StatusBadRequest, api.ErrorResponse(errors.New("you passed strings instead of numbers")))
			return
		}
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	accounts, err := h.service.ListAccount(input.PageSize, input.PageId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (h *handler) UpdateAccount(ctx *gin.Context) {
	var input entity.UpdateAccountReq
	var uri entity.GetAccountReq

	if err := ctx.ShouldBindUri(&uri); err != nil {
		if strings.Contains(err.Error(), "strconv.ParseInt") {
			ctx.JSON(http.StatusBadRequest, api.ErrorResponse(errors.New("string not allowed")))
			return
		}
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		if strings.Contains(err.Error(), "strconv.ParseInt") {
			ctx.JSON(http.StatusBadRequest, api.ErrorResponse(errors.New("string not allowed")))
			return
		}
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	account, err := h.service.UpdateAccount(uri.ID, input.Balance)

	if err != nil {
		if err.(*api.RequestError).Code == 404 {
			ctx.JSON(http.StatusNotFound, api.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (h *handler) DeleteAccount(ctx *gin.Context) {
	var uri entity.GetAccountReq

	if err := ctx.ShouldBindUri(&uri); err != nil {
		if strings.Contains(err.Error(), "strconv.ParseInt") {
			ctx.JSON(http.StatusBadRequest, api.ErrorResponse(errors.New("string not allowed")))
			return
		}
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	err := h.service.DeleteAccount(uri.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Account deleted"})
}
