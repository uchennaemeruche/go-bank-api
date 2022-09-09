package userRoutes

import (
	"github.com/gin-gonic/gin"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
	"github.com/uchennaemeruche/go-bank-api/token"
	"github.com/uchennaemeruche/go-bank-api/user/handler"
	"github.com/uchennaemeruche/go-bank-api/user/service"
	"github.com/uchennaemeruche/go-bank-api/util"
)

func Init(router gin.Engine, store db.Store, tokenMaker token.Maker, config util.Config) {
	r := router.Group("/users")

	userService := service.NewUserService(store, tokenMaker)
	handler := handler.NewUserHandler(userService, config)

	r.POST("", handler.CreateUser)
	r.POST("/login", handler.LoginUser)
	r.POST("/session/renew", handler.RenewAccessToken)
	r.POST("/session/destroy", handler.Logout)
	r.POST("/session/block", handler.ToggleBlockSession)
}
