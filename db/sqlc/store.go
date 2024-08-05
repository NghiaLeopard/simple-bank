package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx (ctx context.Context,arg TransferTxParams) (TransferTxResult,error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context,fn func(*Queries) error) error {
	tx,err := store.db.BeginTx(ctx,nil)

	if err != nil {
		return err
	}

	q := New(tx)

	err1 := fn(q)

	if err1 != nil {
		if rbErr := tx.Rollback();rbErr != nil {
			return fmt.Errorf("tx err: %v,rb err: %v",err,rbErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

func (store *SQLStore) TransferTx (ctx context.Context,arg TransferTxParams) (TransferTxResult,error) {
	var result TransferTxResult

	err := store.execTx(ctx,func (q *Queries) error {
		var err error
		result.Transfer,err = q.CreateTransfer(ctx,CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry,err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry,err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}
	
		result.FromAccount,result.ToAccount,err  = addMoney(ctx,arg.FromAccountID,arg.ToAccountID,arg.Amount,q)

		if err != nil {
			return err
		}

		return nil
	})

	return result,err
}

func addMoney(
	ctx context.Context,
	FromAccountID int64,
	ToAccountID int64,
	Amount int64,
	q *Queries,
)(account1 Account,account2 Account,err error) {
	if FromAccountID > ToAccountID {
		ToAccountID = FromAccountID
		FromAccountID = ToAccountID
		Amount = -Amount
	}

	account1,err = q.UpdateAccount(ctx,UpdateAccountParams{
		ID: FromAccountID,
		Amount: -Amount,
	})

	if err != nil {
		return
	}

	account2,err = q.UpdateAccount(ctx,UpdateAccountParams{
		ID: ToAccountID,
		Amount: Amount,
	})

	if err != nil {
		return
	}

	return 
}