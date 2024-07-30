package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/NghiaLeopard/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner: utils.RandomOwner(),
		Balance: utils.RandomBalance(),
		Currency: utils.RandomCurrency(),
	}

	account,err := testQuery.CreateAccount(context.Background(),arg)

	require.NoError(t,err)

	require.NotEmpty(t,account)
	require.Equal(t,arg.Owner,account.Owner)
	require.Equal(t,arg.Balance,account.Balance)
	require.Equal(t,arg.Currency,account.Currency)

	require.NotZero(t,account.ID)
	require.NotZero(t,account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T){
	accountRandom := createRandomAccount(t)

	account,err := testQuery.GetAccount(context.Background(),accountRandom.ID)

	require.NoError(t,err)
	require.NotEmpty(t,account)

	require.Equal(t,accountRandom.ID,account.ID)
	require.Equal(t,accountRandom.Balance,account.Balance)
	require.Equal(t,accountRandom.Currency,account.Currency)
	require.Equal(t,accountRandom.Owner,account.Owner)
	require.WithinDuration(t,accountRandom.CreatedAt,account.CreatedAt,time.Second)
}

func TestUpdateAccount(t *testing.T) {
	accountRandom := createRandomAccount(t)
	arg := UpdateAccountParams {
		ID: accountRandom.ID,
		Amount: utils.RandomBalance(),
	}
	
	account,err := testQuery.UpdateAccount(context.Background(),arg)

	require.NoError(t,err)
	require.NotEmpty(t,account)

	require.Equal(t,accountRandom.ID,account.ID)
	require.Equal(t,arg.Amount + accountRandom.Balance,account.Balance)
	require.Equal(t,accountRandom.Currency,account.Currency)
	require.Equal(t,accountRandom.Owner,account.Owner)
	require.WithinDuration(t,accountRandom.CreatedAt,account.CreatedAt,time.Second)

}

func TestDeleteAccount(t *testing.T) {
	accountRandom := createRandomAccount(t)
	
	err := testQuery.DeleteAccount(context.Background(),accountRandom.ID)

	require.NoError(t,err)

	account,err1 := testQuery.GetAccount(context.Background(),accountRandom.ID)

	require.Error(t,err1)
	require.EqualError(t,err1,sql.ErrNoRows.Error())
	require.Empty(t,account)
}

func TestGetListAccounts(t *testing.T) {
	arg:= ListAccountsParams{
		Limit: 5,
		// skip first 5 records and get 5 next records 
		Offset: 5,
	}

	for i := 0; i < 10 ; i++ {
		createRandomAccount(t)
	} 

	listAccount,err := testQuery.ListAccounts(context.Background(),arg)
	require.NoError(t,err)
	require.Equal(t,len(listAccount),5)

	for _,value := range listAccount {
		require.NotEmpty(t,value)
	}
}