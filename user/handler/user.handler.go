package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	api "github.com/uchennaemeruche/go-bank-api/api/util"
	"github.com/uchennaemeruche/go-bank-api/user/service"
	"github.com/uchennaemeruche/go-bank-api/util"
)

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type UserHandler interface {
	CreateUser(*gin.Context)
}

type handler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateUser(ctx *gin.Context) {
	var req service.CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			ctx.JSON(http.StatusBadRequest, api.FormatValidationErr(validationErrs))
			return
		}
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	user, err := h.service.Create(req.Username, hashedPassword, req.FullName, req.Email)
	if err != nil {
		fmt.Println("ERR HERE:", err)
		if err.(*api.RequestError).Code == 403 {
			ctx.JSON(http.StatusForbidden, api.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	res := UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, res)
}
