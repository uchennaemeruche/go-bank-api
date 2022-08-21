package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
	"github.com/uchennaemeruche/go-bank-api/handler"
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

	handler := handler.NewAccountHandler()

	router.GET("/accounts", handler.GetAccount)

	server.router = router

	return server
}

func (s *Server) Start(addr string) {
	s.router.Run(addr)
}
