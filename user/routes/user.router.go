package userRoutes

import (
	"github.com/gin-gonic/gin"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
	"github.com/uchennaemeruche/go-bank-api/user/handler"
	"github.com/uchennaemeruche/go-bank-api/user/service"
)

func Init(router gin.Engine, store db.Store) {
	r := router.Group("/users")

	userService := service.NewUserService(store)
	handler := handler.NewUserHandler(userService)

	r.POST("", handler.CreateUser)
}
