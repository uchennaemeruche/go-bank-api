package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	api "github.com/uchennaemeruche/go-bank-api/api/util"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
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

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        UserResponse
}

type UserHandler interface {
	CreateUser(*gin.Context)
	LoginUser(ctx *gin.Context)
}

type handler struct {
	service service.UserService
	config  util.Config
}

func NewUserHandler(service service.UserService, config util.Config) UserHandler {
	return &handler{
		service: service,
		config:  config,
	}
}

func NewUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
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
		// if err.(*api.RequestError).Code == 403 {
		// 	ctx.JSON(http.StatusForbidden, api.ErrorResponse(err))
		// 	return
		// }
		target := &api.RequestError{}
		if errors.As(err, &target) && target.Code == 403 {
			ctx.JSON(http.StatusForbidden, api.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse(err))
		return
	}

	res := NewUserResponse(user)

	ctx.JSON(http.StatusOK, res)
}

type UserLoginReq struct {
	Username          string    `json:"username"`
	Password          string    `json:"password"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (h *handler) LoginUser(ctx *gin.Context) {
	var req service.UserLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			ctx.JSON(http.StatusBadRequest, api.FormatValidationErr(validationErrs))
			return
		}
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	accessTOken, user, err := h.service.LoginUser(req.Username, req.Password, h.config.AccessTokenDuration)

	if err != nil {
		errCode := err.(*api.RequestError).Code
		httpCode := 500
		if errCode == 404 {
			httpCode = http.StatusNotFound
		} else if errCode == 401 {
			httpCode = http.StatusUnauthorized
		} else {
			httpCode = http.StatusInternalServerError
		}
		ctx.JSON(httpCode, api.ErrorResponse(err))
		return
	}

	rsp := LoginResponse{
		AccessToken: accessTOken,
		User:        NewUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
