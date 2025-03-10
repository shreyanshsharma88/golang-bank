package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/shreyanshsharma88/golang-bank/auth"
	db "github.com/shreyanshsharma88/golang-bank/db/sqlc"
	"github.com/shreyanshsharma88/golang-bank/utils"
)

type Server struct {
	store      *db.Store
	router     *gin.Engine
	tokenMaker auth.Maker
	config     utils.Config
}

func (server *Server) renderRoutes() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRouter.POST("/accounts", server.createAccount)
	authRouter.GET("/accounts/:id", server.getAccount)
	authRouter.GET("/accounts", server.listAccounts)
	authRouter.DELETE("/accounts/:id", server.deleteAccount)
	authRouter.PUT("/accounts/:id", server.updateAccount)

	authRouter.POST("/transfers", server.createTransfers)

	server.router = router

}

func NewServer(config utils.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := auth.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return &Server{}, err
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.renderRoutes()

	return server, nil

}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
