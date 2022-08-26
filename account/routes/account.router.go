package accountRoutes

import (
	"github.com/gin-gonic/gin"
	"github.com/uchennaemeruche/go-bank-api/account/handler"
	"github.com/uchennaemeruche/go-bank-api/account/service"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
)

func Init(router *gin.Engine, store db.Store, authMiddleware gin.HandlerFunc) {

	r := router.Group("/accounts").Use(authMiddleware)

	accountService := service.NewAccountService(store)

	handler := handler.NewAccountHandler(accountService)
	r.GET("/:id", handler.GetAccount)
	r.GET("", handler.ListAccount)
	r.POST("", handler.CreateAccount)
	r.PUT("/:id", handler.UpdateAccount)
	r.DELETE("/:id", handler.DeleteAccount)
}
