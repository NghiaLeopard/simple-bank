package api

import (
	"database/sql"
	"net/http"

	db "github.com/NghiaLeopard/simple-bank/db/sqlc"
	"github.com/NghiaLeopard/simple-bank/utils"
	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username       string `json:"username" binding:"required,alphanum"`
	Password 	   string `json:"password" binding:"required,min=6"`
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest


	if err := ctx.ShouldBindJSON(&req);err != nil {
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}

	hashedPassword,err := utils.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	}

	user,err := server.store.CreateUser(ctx,arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,user)
}

type LoginUserRequest struct {
	Username       string `json:"username" binding:"required,alphanum"`
	Password 	   string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	Username string
	Token string
}

func (server *Server) LoginUser(ctx *gin.Context) {
	req := &LoginUserRequest{}

	err := ctx.ShouldBindJSON(req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}

	user,err := server.store.GetUser(ctx,req.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound,errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}

	err = utils.CheckPassword(req.Password,user.HashedPassword)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized,errResponse(err))
		return
	}

	token,payload,err := server.tokenMaker.CreateTokenPaseto(req.Username,server.config.Duration)

	if payload == nil || token == "" {
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}

	rsp := LoginUserResponse{
		Username: payload.Username,
		Token: token,
	}

	ctx.JSON(http.StatusOK,rsp)
}
