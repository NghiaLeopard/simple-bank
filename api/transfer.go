package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/NghiaLeopard/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountId int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}

	if !(server.validAccount(ctx,req.FromAccountId,req.Currency)) {
		return
	}

	if !(server.validAccount(ctx,req.ToAccountId,req.Currency)) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID: req.ToAccountId,
		Amount: req.Amount,
	}

	result, err := server.store.TransferTx(ctx,arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
	}

	ctx.JSON(http.StatusOK,result)
}

func (server *Server) validAccount(ctx *gin.Context,accountID int64,currency string) bool {
	account,err := server.store.GetAccount(ctx,accountID)

	if err != nil {
		if err := sql.ErrNoRows; err != nil {
		ctx.JSON(http.StatusNotFound,errResponse(err))

		return false
		}
		ctx.JSON(http.StatusInternalServerError,errResponse(err))

		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest,errResponse(err))

		return false
	}

	return true
}