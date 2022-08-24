package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	accountRoutes "github.com/uchennaemeruche/go-bank-api/account/routes"
	api "github.com/uchennaemeruche/go-bank-api/api/util"
	transferRoutes "github.com/uchennaemeruche/go-bank-api/transfer/routes"
	userRoutes "github.com/uchennaemeruche/go-bank-api/user/routes"

	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
)

type Server struct {
	store  db.Store
	Router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", api.ValidCurrency)
	}

	accountRoutes.Init(router, store)

	transferRoutes.Init(router, store)

	userRoutes.Init(*router, store)

	server.Router = router

	return server
}

func (s *Server) Start(addr string) {
	s.Router.Run(addr)

}
