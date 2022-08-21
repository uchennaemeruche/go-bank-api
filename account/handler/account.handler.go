package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uchennaemeruche/go-bank-api/account/entity"
	"github.com/uchennaemeruche/go-bank-api/account/service"
	"github.com/uchennaemeruche/go-bank-api/api/util"
)

type AccountHandler interface {
	GetAccount(*gin.Context)
	CreateAccount(*gin.Context)
	ListAccount(*gin.Context)
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
			ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("string not allowed")))
			return
		}
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	account, err := h.service.GetOne(uri.ID)
	if err != nil {
		if err.(*util.RequestError).Code == 404 {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
func (h *handler) CreateAccount(ctx *gin.Context) {
	var input entity.CreateAccountReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	account, err := h.service.Create(input.Owner, input.Currency)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

func (h *handler) ListAccount(ctx *gin.Context) {
	var input entity.ListAccountReq
	if err := ctx.ShouldBindQuery(&input); err != nil {
		if strings.Contains(err.Error(), "strconv.ParseInt") {
			ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("you passed strings instead of numbers")))
			return
		}
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	accounts, err := h.service.ListAccount(input.PageSize, input.PageId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
