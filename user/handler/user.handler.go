package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	api "github.com/uchennaemeruche/go-bank-api/api/util"
	"github.com/uchennaemeruche/go-bank-api/user/service"
	"github.com/uchennaemeruche/go-bank-api/util"
)

type UserHandler interface {
	CreateUser(*gin.Context)
	LoginUser(ctx *gin.Context)
	RenewAccessToken(ctx *gin.Context)
	Logout(ctx *gin.Context)
	ToggleBlockSession(ctx *gin.Context)
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

	res := service.NewUserResponse(user)

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

	response, err := h.service.LoginUser(req.Username, req.Password, ctx.Request.UserAgent(), ctx.ClientIP(), h.config.AccessTokenDuration, h.config.RefreshTokenDuration)

	if err != nil {
		errCode := err.(*api.RequestError).Code
		httpCode := http.StatusInternalServerError
		if errCode == 404 {
			httpCode = http.StatusNotFound
		} else if errCode == 401 {
			httpCode = http.StatusUnauthorized
		}
		ctx.JSON(httpCode, api.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *handler) RenewAccessToken(ctx *gin.Context) {
	var req service.RenewAcesTokenReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			ctx.JSON(http.StatusBadRequest, api.FormatValidationErr(validationErrs))
			return
		}
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	response, err := h.service.GetNewAccessToken(req.RefreshToken, h.config.AccessTokenDuration)

	if err != nil {
		errCode := err.(*api.RequestError).Code
		httpCode := http.StatusInternalServerError
		if errCode == 404 {
			httpCode = http.StatusNotFound
		} else if errCode == 401 {
			httpCode = http.StatusUnauthorized
		}
		ctx.JSON(httpCode, api.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *handler) Logout(ctx *gin.Context) {
	var req service.DestroySessionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			ctx.JSON(http.StatusBadRequest, api.FormatValidationErr(validationErrs))
			return
		}
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	isDestoyed, err := h.service.DestroySession(req.RefreshToken)
	if err != nil || !isDestoyed {
		if err != nil {
			errCode := err.(*api.RequestError).Code
			httpCode := http.StatusInternalServerError
			if errCode == 401 {
				httpCode = http.StatusUnauthorized
			}
			ctx.JSON(httpCode, api.ErrorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})

}
func (h *handler) ToggleBlockSession(ctx *gin.Context) {
	var req service.ToggleSessionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			ctx.JSON(http.StatusBadRequest, api.FormatValidationErr(validationErrs))
			return
		}
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	isDestoyed, err := h.service.ToggleBlockSession(req.RefreshToken, req.Status)
	if err != nil || !isDestoyed {
		if err != nil {
			errCode := err.(*api.RequestError).Code
			httpCode := http.StatusInternalServerError
			if errCode == 401 {
				httpCode = http.StatusUnauthorized
			}
			ctx.JSON(httpCode, api.ErrorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User session token blocked"})
}
