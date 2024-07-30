package api

import (
	"database/sql"
	"log"
	"net/http"

	db "github.com/NghiaLeopard/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Currency: req.Currency,
		Balance: 0,
	}

	account,err := server.store.CreateAccount(ctx,arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}

	account,err := server.store.GetAccount(ctx,req.ID)

	if err != nil {
		if err  == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound,errResponse(err))
			return 
		}
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)
}

type getListRequest struct {
	PageId int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,max=10"`
}

func (server *Server) getListAccount(ctx *gin.Context) {
	var req getListRequest

	err := ctx.ShouldBindQuery(&req)

	log.Println(req.PageSize,req.PageId)


	if err != nil {
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	listAccount, err := server.store.ListAccounts(ctx,arg)


	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,listAccount)
}

type updateAccountRequest struct {
	ID int64 `json:"id" binding:"required"`
	Amount int64 `json:"amount" binding:"required,min=1"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest

	if err := ctx.ShouldBindJSON(&req);err != nil {
		ctx.JSON(http.StatusBadRequest,errResponse(err))
	}

	arg := db.UpdateAccountParams{
		Amount: req.Amount,
		ID: req.ID,
	}

	updateAccount,err := server.store.UpdateAccount(ctx,arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
	}

	ctx.JSON(http.StatusOK,updateAccount)

}