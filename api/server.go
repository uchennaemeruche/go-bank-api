package api

import (
	"github.com/gin-gonic/gin"
	"github.com/uchennaemeruche/go-bank-api/account/handler"
	"github.com/uchennaemeruche/go-bank-api/account/service"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: &store,
	}

	router := gin.Default()

	accountService := service.NewAccountService(server.store)

	handler := handler.NewAccountHandler(accountService)

	router.GET("/accounts/:id", handler.GetAccount)
	router.GET("/accounts", handler.CreateAccount)
	router.POST("/accounts", handler.CreateAccount)

	server.router = router

	return server
}

func (s *Server) Start(addr string) {
	s.router.Run(addr)
}
