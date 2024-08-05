package api

import (
	"fmt"

	db "github.com/NghiaLeopard/simple-bank/db/sqlc"
	"github.com/NghiaLeopard/simple-bank/token"
	"github.com/NghiaLeopard/simple-bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config utils.Config
	store db.Store
	router *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config utils.Config,store db.Store) (*Server,error) {
	tokenMaker,err := token.NewPasetoMaker([]byte(config.SymmetricKey))

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{config: config,store: store,tokenMaker: tokenMaker}
	router := gin.Default()

	if v,ok := binding.Validator.Engine().(*validator.Validate);ok {
		v.RegisterValidation("currency",ValidCurrency)
	}

	router.GET("/accounts",server.getListAccount)
	router.GET("/account/:id",server.getAccount)
	router.POST("/account",server.createAccount)
	router.POST("/user",server.CreateUser)
	router.POST("/user/login",server.LoginUser)
	router.POST("/transfer",server.createTransfer)
	router.PATCH("/account",server.updateAccount)


	server.router = router

	return server,nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}