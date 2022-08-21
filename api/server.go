package api

import (
	"github.com/gin-gonic/gin"
	accountRoutes "github.com/uchennaemeruche/go-bank-api/account/routes"
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

	accountRoutes.Init(router, server.store)

	server.router = router

	return server
}

func (s *Server) Start(addr string) {
	s.router.Run(addr)
}
