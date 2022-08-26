package transferroutes

import (
	"github.com/gin-gonic/gin"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
	"github.com/uchennaemeruche/go-bank-api/transfer/handler"
	"github.com/uchennaemeruche/go-bank-api/transfer/service"
)

func Init(router *gin.Engine, store db.Store, authMiddleware gin.HandlerFunc) {
	r := router.Group("/transfers").Use(authMiddleware)

	service := service.NewTransferService(store)
	handler := handler.NewTransferHandler(service)

	r.POST("", handler.CreateTransfer)
}
