package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	accountRoutes "github.com/uchennaemeruche/go-bank-api/account/routes"
	api "github.com/uchennaemeruche/go-bank-api/api/util"
	"github.com/uchennaemeruche/go-bank-api/token"
	transferRoutes "github.com/uchennaemeruche/go-bank-api/transfer/routes"
	userRoutes "github.com/uchennaemeruche/go-bank-api/user/routes"
	"github.com/uchennaemeruche/go-bank-api/util"

	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
)

type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	Router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create a token maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", api.ValidCurrency)
	}

	router.Use()

	accountRoutes.Init(router, store, authMiddleware(tokenMaker))

	transferRoutes.Init(router, store, authMiddleware(tokenMaker))

	userRoutes.Init(*router, store, tokenMaker, config)

	server.Router = router

	return server, nil
}

func (s *Server) Start(addr string) {
	s.Router.Run(addr)
}
