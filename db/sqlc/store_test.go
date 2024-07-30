package db

import (
	"context"
	"testing"

	"github.com/NghiaLeopard/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: utils.RandomBalance(),
	}

	n := 10

	// existed :=  make(map[int64]bool)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0 ; i < n ; i++ {
		if(i % 2 == 1) {
			arg = TransferTxParams{
				FromAccountID: account2.ID,
				ToAccountID: account1.ID,
				Amount: arg.Amount,
			}
		}
		go func() {
			result,err := store.TransferTx(context.Background(),arg)

			errs <- err
			results <- result
		}()
	}

	for i := 0 ; i < n ; i++ {
		err := <- errs
		require.NoError(t,err)

		result := <- results

		require.NotEmpty(t,result)

		// require.NoError(t,err)
		// 	require.NotEmpty(t,result)
		
		// 	require.Equal(t,result.Transfer.FromAccountID,arg.FromAccountID)
		// 	require.Equal(t,result.Transfer.ToAccountID,arg.ToAccountID)
		// 	require.Equal(t,result.Transfer.Amount,arg.Amount)
		// 	require.NotZero(t,result.Transfer.ID)
		// 	require.NotZero(t,result.Transfer.CreatedAt)
	
		// 	transferResult,err := store.GetTransfer(context.Background(),result.Transfer.ID)
		// 	require.NoError(t,err)
		// 	require.NotEmpty(t,transferResult)
	
		// 	// FromEntry
		// 	require.NotEmpty(t,result.FromEntry)
		// 	require.Equal(t,result.FromEntry.AccountID,arg.FromAccountID)
		// 	require.Equal(t,result.FromEntry.Amount,-arg.Amount)
	
		// 	entryResult,err := store.GetEntry(context.Background(),result.FromEntry.ID)
		// 	require.NoError(t,err)
		// 	require.NotEmpty(t,entryResult)
		// 	require.NotZero(t,entryResult.ID)
		// 	require.NotZero(t,entryResult.CreatedAt)
	
		// 	// ToEntry
		// 	require.NotEmpty(t,result.ToEntry)
		// 	require.Equal(t,result.ToEntry.AccountID,arg.ToAccountID)
		// 	require.Equal(t,result.ToEntry.Amount,arg.Amount)
		// 	ToEntryResult,err := store.GetEntry(context.Background(),result.ToEntry.ID)
		// 	require.NoError(t,err)
		// 	require.NotEmpty(t,ToEntryResult)
		// 	require.NotZero(t,ToEntryResult.ID)
		// 	require.NotZero(t,ToEntryResult.CreatedAt)


		// 	// check account
		// 	require.NotEmpty(t,result.FromAccount)
		// 	require.Equal(t,result.FromAccount.ID,account1.ID)

		// 	require.NotEmpty(t,result.ToAccount)
		// 	require.Equal(t,result.ToAccount.ID,account2.ID)

		// 	// check balance
		// 	diff1 := account1.Balance - result.FromAccount.Balance
		// 	diff2 :=  result.ToAccount.Balance - account2.Balance 

		// 	require.Equal(t,diff1,diff2)
		// 	require.True(t,diff1 > 0)
		// 	require.True(t,diff1 % arg.Amount == 0)

		// 	k := diff1 / arg.Amount
		// 	require.NotContains(t,existed,k)
		// 	existed[k] = true
	}

	// check final balance account
	updateAccount1,err := store.GetAccount(context.Background(),account1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,updateAccount1)
	require.Equal(t,account1.Balance ,updateAccount1.Balance)


	updateAccount2,err := store.GetAccount(context.Background(),account2.ID)
	require.NoError(t,err)
	require.NotEmpty(t,updateAccount2)
	require.Equal(t,account2.Balance ,updateAccount2.Balance)

	t.Log(updateAccount1.Balance,updateAccount2.Balance)
	t.Log(account1.Balance,account2.Balance)

}
